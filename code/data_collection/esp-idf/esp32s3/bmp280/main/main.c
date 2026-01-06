// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

#include <stdio.h>
#include <string.h>
#include <time.h>
#include <sys/time.h>
#include "freertos/FreeRTOS.h"
#include "freertos/task.h"
#include "freertos/event_groups.h"
#include "esp_system.h"
#include "esp_wifi.h"
#include "esp_event.h"
#include "esp_log.h"
#include "esp_sntp.h"
#include "nvs_flash.h"
#include "driver/i2c.h"
#include "esp_http_client.h"

// WiFi Configuration - UPDATE THESE
#define WIFI_SSID "your_wifi_ssid"
#define WIFI_PASS "your_wifi_password"

// API Configuration - UPDATE THIS
#define API_URL "http://192.168.1.60:8080/log"

// I2C Configuration for ESP32-S3-DevKitC-1
#define I2C_MASTER_SCL_IO 9
#define I2C_MASTER_SDA_IO 8
#define I2C_MASTER_NUM I2C_NUM_0
#define I2C_MASTER_FREQ_HZ 100000

// BMP280 Configuration
#define BMP280_ADDR 0x77
#define BMP280_REG_TEMP_XLSB 0xFC
#define BMP280_REG_TEMP_LSB 0xFB
#define BMP280_REG_TEMP_MSB 0xFA
#define BMP280_REG_CTRL_MEAS 0xF4
#define BMP280_REG_CONFIG 0xF5
#define BMP280_REG_ID 0xD0

// Calibration registers
#define BMP280_REG_DIG_T1 0x88

static const char *TAG = "BMP280_LOGGER";
static EventGroupHandle_t wifi_event_group;
const int WIFI_CONNECTED_BIT = BIT0;

// BMP280 calibration data
typedef struct {
    uint16_t dig_T1;
    int16_t dig_T2;
    int16_t dig_T3;
} bmp280_calib_data;

static bmp280_calib_data calib;

// I2C initialization
static esp_err_t i2c_master_init(void) {
    i2c_config_t conf = {
        .mode = I2C_MODE_MASTER,
        .sda_io_num = I2C_MASTER_SDA_IO,
        .scl_io_num = I2C_MASTER_SCL_IO,
        .sda_pullup_en = GPIO_PULLUP_ENABLE,
        .scl_pullup_en = GPIO_PULLUP_ENABLE,
        .master.clk_speed = I2C_MASTER_FREQ_HZ,
    };
    esp_err_t err = i2c_param_config(I2C_MASTER_NUM, &conf);
    if (err != ESP_OK) return err;
    return i2c_driver_install(I2C_MASTER_NUM, conf.mode, 0, 0, 0);
}

// I2C read register
static esp_err_t bmp280_read_reg(uint8_t reg, uint8_t *data, size_t len) {
    return i2c_master_write_read_device(I2C_MASTER_NUM, BMP280_ADDR,
                                        &reg, 1, data, len,
                                        pdMS_TO_TICKS(1000));
}

// I2C write register
static esp_err_t bmp280_write_reg(uint8_t reg, uint8_t data) {
    uint8_t write_buf[2] = {reg, data};
    return i2c_master_write_to_device(I2C_MASTER_NUM, BMP280_ADDR,
                                      write_buf, 2, pdMS_TO_TICKS(1000));
}

// Initialize BMP280 sensor
static esp_err_t bmp280_init(void) {
    uint8_t chip_id;
    esp_err_t err = bmp280_read_reg(BMP280_REG_ID, &chip_id, 1);
    if (err != ESP_OK) {
        ESP_LOGE(TAG, "Failed to read BMP280 chip ID");
        return err;
    }
    ESP_LOGI(TAG, "BMP280 Chip ID: 0x%02X", chip_id);

    // Read calibration data
    uint8_t calib_data[6];
    err = bmp280_read_reg(BMP280_REG_DIG_T1, calib_data, 6);
    if (err != ESP_OK) {
        ESP_LOGE(TAG, "Failed to read calibration data");
        return err;
    }

    calib.dig_T1 = (calib_data[1] << 8) | calib_data[0];
    calib.dig_T2 = (calib_data[3] << 8) | calib_data[2];
    calib.dig_T3 = (calib_data[5] << 8) | calib_data[4];

    // Configure BMP280: normal mode, temp oversampling x2
    err = bmp280_write_reg(BMP280_REG_CTRL_MEAS, 0x4F);
    if (err != ESP_OK) {
        ESP_LOGE(TAG, "Failed to configure BMP280");
        return err;
    }

    vTaskDelay(pdMS_TO_TICKS(100));
    ESP_LOGI(TAG, "BMP280 initialized successfully");
    return ESP_OK;
}

// Read temperature from BMP280
static float bmp280_read_temperature(void) {
    uint8_t data[3];
    esp_err_t err = bmp280_read_reg(BMP280_REG_TEMP_MSB, data, 3);
    if (err != ESP_OK) {
        ESP_LOGE(TAG, "Failed to read temperature");
        return -999.0;
    }

    int32_t adc_T = ((int32_t)data[0] << 12) | ((int32_t)data[1] << 4) | ((int32_t)data[2] >> 4);

    // Temperature compensation formula from BMP280 datasheet
    int32_t var1, var2, T;
    var1 = ((((adc_T >> 3) - ((int32_t)calib.dig_T1 << 1))) * ((int32_t)calib.dig_T2)) >> 11;
    var2 = (((((adc_T >> 4) - ((int32_t)calib.dig_T1)) *
              ((adc_T >> 4) - ((int32_t)calib.dig_T1))) >> 12) *
            ((int32_t)calib.dig_T3)) >> 14;
    T = var1 + var2;
    float temperature = (T * 5 + 128) >> 8;
    return temperature / 100.0;
}

// WiFi event handler
static void wifi_event_handler(void *arg, esp_event_base_t event_base,
                                int32_t event_id, void *event_data) {
    if (event_base == WIFI_EVENT && event_id == WIFI_EVENT_STA_START) {
        esp_wifi_connect();
    } else if (event_base == WIFI_EVENT && event_id == WIFI_EVENT_STA_DISCONNECTED) {
        ESP_LOGI(TAG, "WiFi disconnected, retrying...");
        esp_wifi_connect();
    } else if (event_base == IP_EVENT && event_id == IP_EVENT_STA_GOT_IP) {
        ip_event_got_ip_t *event = (ip_event_got_ip_t *)event_data;
        ESP_LOGI(TAG, "Got IP: " IPSTR, IP2STR(&event->ip_info.ip));
        xEventGroupSetBits(wifi_event_group, WIFI_CONNECTED_BIT);
    }
}

// Initialize WiFi
static void wifi_init(void) {
    wifi_event_group = xEventGroupCreate();

    ESP_ERROR_CHECK(esp_netif_init());
    ESP_ERROR_CHECK(esp_event_loop_create_default());
    esp_netif_create_default_wifi_sta();

    wifi_init_config_t cfg = WIFI_INIT_CONFIG_DEFAULT();
    ESP_ERROR_CHECK(esp_wifi_init(&cfg));

    ESP_ERROR_CHECK(esp_event_handler_register(WIFI_EVENT, ESP_EVENT_ANY_ID,
                                                &wifi_event_handler, NULL));
    ESP_ERROR_CHECK(esp_event_handler_register(IP_EVENT, IP_EVENT_STA_GOT_IP,
                                                &wifi_event_handler, NULL));

    wifi_config_t wifi_config = {
        .sta = {
            .ssid = WIFI_SSID,
            .password = WIFI_PASS,
        },
    };

    ESP_ERROR_CHECK(esp_wifi_set_mode(WIFI_MODE_STA));
    ESP_ERROR_CHECK(esp_wifi_set_config(WIFI_IF_STA, &wifi_config));
    ESP_ERROR_CHECK(esp_wifi_start());

    ESP_LOGI(TAG, "WiFi initialization completed");
}

// SNTP time sync callback
void time_sync_notification_cb(struct timeval *tv) {
    ESP_LOGI(TAG, "Time synchronized");
}

// Initialize SNTP
static void obtain_time(void) {
    ESP_LOGI(TAG, "Initializing SNTP");
    esp_sntp_setoperatingmode(SNTP_OPMODE_POLL);
    esp_sntp_setservername(0, "pool.ntp.org");
    esp_sntp_set_time_sync_notification_cb(time_sync_notification_cb);
    esp_sntp_init();

    // Wait for time to be set
    time_t now = 0;
    struct tm timeinfo = {0};
    int retry = 0;
    const int retry_count = 10;

    while (timeinfo.tm_year < (2016 - 1900) && ++retry < retry_count) {
        ESP_LOGI(TAG, "Waiting for system time to be set... (%d/%d)", retry, retry_count);
        vTaskDelay(pdMS_TO_TICKS(2000));
        time(&now);
        localtime_r(&now, &timeinfo);
    }

    if (timeinfo.tm_year < (2016 - 1900)) {
        ESP_LOGE(TAG, "Failed to obtain time");
    } else {
        ESP_LOGI(TAG, "Time obtained successfully");
    }
}

// HTTP POST handler
esp_err_t http_event_handler(esp_http_client_event_t *evt) {
    switch (evt->event_id) {
        case HTTP_EVENT_ON_DATA:
            ESP_LOGI(TAG, "HTTP Response: %.*s", evt->data_len, (char *)evt->data);
            break;
        default:
            break;
    }
    return ESP_OK;
}

// POST data to API
static void post_data_to_api(float temperature, time_t timestamp) {
    char post_data[256];
    snprintf(post_data, sizeof(post_data),
             "{\"recorded\":%ld,\"location\":\"den\",\"sensor\":\"bmp280\","
             "\"measurement\":\"temperature\",\"units\":\"C\",\"value\":%.1f}",
             (long)timestamp, temperature);

    ESP_LOGI(TAG, "Posting data: %s", post_data);

    esp_http_client_config_t config = {
        .url = API_URL,
        .event_handler = http_event_handler,
        .method = HTTP_METHOD_POST,
    };

    esp_http_client_handle_t client = esp_http_client_init(&config);
    esp_http_client_set_header(client, "Content-Type", "application/json");
    esp_http_client_set_post_field(client, post_data, strlen(post_data));

    esp_err_t err = esp_http_client_perform(client);
    if (err == ESP_OK) {
        ESP_LOGI(TAG, "HTTP POST Status = %d",
                 esp_http_client_get_status_code(client));
    } else {
        ESP_LOGE(TAG, "HTTP POST request failed: %s", esp_err_to_name(err));
    }

    esp_http_client_cleanup(client);
}

void app_main(void) {
    ESP_LOGI(TAG, "Starting BMP280 Temperature Logger");

    // Initialize NVS
    esp_err_t ret = nvs_flash_init();
    if (ret == ESP_ERR_NVS_NO_FREE_PAGES || ret == ESP_ERR_NVS_NEW_VERSION_FOUND) {
        ESP_ERROR_CHECK(nvs_flash_erase());
        ret = nvs_flash_init();
    }
    ESP_ERROR_CHECK(ret);

    // Initialize I2C
    ESP_ERROR_CHECK(i2c_master_init());
    ESP_LOGI(TAG, "I2C initialized");

    // Initialize BMP280
    ESP_ERROR_CHECK(bmp280_init());

    // Initialize WiFi
    wifi_init();

    // Wait for WiFi connection
    xEventGroupWaitBits(wifi_event_group, WIFI_CONNECTED_BIT,
                        false, true, portMAX_DELAY);
    ESP_LOGI(TAG, "Connected to WiFi");

    // Obtain time from NTP
    obtain_time();

    // Main loop: read temperature and post to API every 60 seconds
    while (1) {
        float temperature = bmp280_read_temperature();
        time_t now;
        time(&now);

        ESP_LOGI(TAG, "Temperature: %.2f °C, Timestamp: %ld", temperature, (long)now);

        if (temperature != -999.0) {
            post_data_to_api(temperature, now);
        }

        vTaskDelay(pdMS_TO_TICKS(60000)); // Wait 60 seconds
    }
}
