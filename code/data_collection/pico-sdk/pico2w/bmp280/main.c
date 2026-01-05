// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

#include <stdio.h>
#include <string.h>
#include "pico/stdlib.h"
#include "pico/cyw43_arch.h"
#include "hardware/i2c.h"
#include "lwip/pbuf.h"
#include "lwip/udp.h"
#include "lwip/tcp.h"
#include "lwip/dns.h"

// WiFi credentials - UPDATE THESE
#define WIFI_SSID "YOUR_WIFI_SSID"
#define WIFI_PASSWORD "YOUR_WIFI_PASSWORD"

// API endpoint - UPDATE THIS
#define API_HOST "192.168.1.100"  // Use IP address or implement DNS
#define API_PORT 80
#define API_PATH "/api/temperature"

// I2C Configuration
#define I2C_PORT i2c0
#define I2C_SDA 4
#define I2C_SCL 5
#define BMP280_ADDR 0x76  // or 0x77 if SDO is pulled high

// BMP280 Registers
#define BMP280_REG_TEMP_XLSB 0xFC
#define BMP280_REG_TEMP_LSB 0xFB
#define BMP280_REG_TEMP_MSB 0xFA
#define BMP280_REG_CONTROL 0xF4
#define BMP280_REG_CONFIG 0xF5
#define BMP280_REG_ID 0xD0
#define BMP280_REG_RESET 0xE0
#define BMP280_REG_STATUS 0xF3
#define BMP280_REG_CALIB_START 0x88
#define BMP280_REG_CALIB_END 0xA1

// Calibration data structure
typedef struct {
    uint16_t dig_T1;
    int16_t dig_T2;
    int16_t dig_T3;
} bmp280_calib_t;

bmp280_calib_t calib;
int32_t t_fine;

// I2C helper functions
static int bmp280_write_byte(uint8_t reg, uint8_t value) {
    uint8_t buf[2] = {reg, value};
    return i2c_write_blocking(I2C_PORT, BMP280_ADDR, buf, 2, false);
}

static int bmp280_read_bytes(uint8_t reg, uint8_t *buf, uint16_t len) {
    int ret = i2c_write_blocking(I2C_PORT, BMP280_ADDR, &reg, 1, true);
    if (ret < 0) return ret;
    return i2c_read_blocking(I2C_PORT, BMP280_ADDR, buf, len, false);
}

// Initialize BMP280 sensor
bool bmp280_init() {
    printf("Initializing BMP280...\n");

    // Check chip ID
    uint8_t chip_id;
    if (bmp280_read_bytes(BMP280_REG_ID, &chip_id, 1) < 0) {
        printf("Error: Failed to read BMP280 chip ID\n");
        return false;
    }

    printf("BMP280 Chip ID: 0x%02X ", chip_id);
    if (chip_id == 0x58) {
        printf("(BMP280)\n");
    } else if (chip_id == 0x60) {
        printf("(BME280)\n");
    } else {
        printf("(Unknown - expected 0x58 or 0x60)\n");
        printf("Check I2C address (try 0x77 if using 0x76)\n");
        return false;
    }

    // Soft reset
    bmp280_write_byte(BMP280_REG_RESET, 0xB6);
    sleep_ms(10);

    // Read calibration data (6 bytes for temperature)
    uint8_t calib_data[6];
    if (bmp280_read_bytes(BMP280_REG_CALIB_START, calib_data, 6) < 0) {
        printf("Error: Failed to read calibration data\n");
        return false;
    }

    calib.dig_T1 = (calib_data[1] << 8) | calib_data[0];
    calib.dig_T2 = (int16_t)((calib_data[3] << 8) | calib_data[2]);
    calib.dig_T3 = (int16_t)((calib_data[5] << 8) | calib_data[4]);

    printf("Calibration: T1=%u, T2=%d, T3=%d\n",
           calib.dig_T1, calib.dig_T2, calib.dig_T3);

    // Configure sensor
    // Config: standby 1000ms, filter off
    bmp280_write_byte(BMP280_REG_CONFIG, 0x00);

    // Control: osrs_t=1 (x1), osrs_p=0 (skip pressure), mode=normal
    bmp280_write_byte(BMP280_REG_CONTROL, 0x27);

    sleep_ms(100);

    printf("BMP280 initialized successfully\n");
    return true;
}

// Read temperature from BMP280
float bmp280_read_temperature() {
    uint8_t data[3];

    if (bmp280_read_bytes(BMP280_REG_TEMP_MSB, data, 3) < 0) {
        printf("Error reading temperature\n");
        return -999.0f;
    }

    int32_t adc_T = ((uint32_t)data[0] << 12) | ((uint32_t)data[1] << 4) | ((data[2] >> 4) & 0x0F);

    // Compensation formula from BMP280 datasheet
    int32_t var1 = ((((adc_T >> 3) - ((int32_t)calib.dig_T1 << 1))) * ((int32_t)calib.dig_T2)) >> 11;
    int32_t var2 = (((((adc_T >> 4) - ((int32_t)calib.dig_T1)) *
                      ((adc_T >> 4) - ((int32_t)calib.dig_T1))) >> 12) *
                    ((int32_t)calib.dig_T3)) >> 14;

    t_fine = var1 + var2;
    int32_t T = (t_fine * 5 + 128) >> 8;

    return (float)T / 100.0f;
}

// NTP time result
typedef struct {
    bool received;
    uint32_t timestamp;
} ntp_result_t;

ntp_result_t ntp_result = {false, 0};

// NTP UDP receive callback
void ntp_recv(void *arg, struct udp_pcb *pcb, struct pbuf *p, const ip_addr_t *addr, u16_t port) {
    if (p != NULL && p->len >= 48) {
        uint8_t *data = (uint8_t *)p->payload;
        // Extract timestamp from bytes 40-43 (Transmit Timestamp)
        uint32_t ntp_time = ((uint32_t)data[40] << 24) |
                           ((uint32_t)data[41] << 16) |
                           ((uint32_t)data[42] << 8) |
                           data[43];

        // Convert NTP timestamp (seconds since 1900) to Unix timestamp (seconds since 1970)
        ntp_result.timestamp = ntp_time - 2208988800UL;
        ntp_result.received = true;

        printf("NTP time received: %lu (", ntp_result.timestamp);
        time_t t = ntp_result.timestamp;
        printf("%s)\n", ctime(&t));
    }
    if (p != NULL) {
        pbuf_free(p);
    }
}

// Get current time from NTP server
uint32_t get_ntp_time() {
    printf("Getting NTP time...\n");

    struct udp_pcb *pcb = udp_new();
    if (!pcb) {
        printf("Failed to create UDP PCB\n");
        return 0;
    }

    udp_recv(pcb, ntp_recv, NULL);

    // Use pool.ntp.org - 216.239.35.0 is one of Google's public NTP servers
    ip_addr_t ntp_server;
    IP4_ADDR(&ntp_server, 216, 239, 35, 0);

    // Build NTP request packet
    uint8_t ntp_packet[48];
    memset(ntp_packet, 0, 48);
    ntp_packet[0] = 0x1B; // LI=0, VN=3, Mode=3 (client)

    struct pbuf *p = pbuf_alloc(PBUF_TRANSPORT, 48, PBUF_RAM);
    if (!p) {
        printf("Failed to allocate pbuf\n");
        udp_remove(pcb);
        return 0;
    }

    memcpy(p->payload, ntp_packet, 48);
    err_t err = udp_sendto(pcb, p, &ntp_server, 123);
    pbuf_free(p);

    if (err != ERR_OK) {
        printf("Failed to send NTP request: %d\n", err);
        udp_remove(pcb);
        return 0;
    }

    // Wait for response (up to 5 seconds)
    ntp_result.received = false;
    for (int i = 0; i < 50 && !ntp_result.received; i++) {
        cyw43_arch_poll();
        sleep_ms(100);
    }

    udp_remove(pcb);

    if (!ntp_result.received) {
        printf("NTP timeout - no response\n");
    }

    return ntp_result.received ? ntp_result.timestamp : 0;
}

// TCP connection state
typedef struct {
    struct tcp_pcb *pcb;
    bool complete;
    bool success;
    char *request;
    uint16_t request_len;
    uint16_t sent;
} tcp_state_t;

// TCP error callback
static void tcp_error_cb(void *arg, err_t err) {
    tcp_state_t *state = (tcp_state_t *)arg;
    // ERR_ABRT (-13) after successful response is normal - server closed connection
    if (err == ERR_ABRT && state->success) {
        // This is fine - we already got our response
        state->complete = true;
    } else {
        printf("TCP error: %d\n", err);
        state->complete = true;
        state->success = false;
    }
}

// TCP sent callback
static err_t tcp_sent_cb(void *arg, struct tcp_pcb *pcb, u16_t len) {
    tcp_state_t *state = (tcp_state_t *)arg;
    state->sent += len;

    if (state->sent >= state->request_len) {
        printf("Request sent completely\n");
    }

    return ERR_OK;
}

// TCP receive callback
static err_t tcp_recv_cb(void *arg, struct tcp_pcb *pcb, struct pbuf *p, err_t err) {
    tcp_state_t *state = (tcp_state_t *)arg;

    if (p == NULL) {
        // Connection closed by server - this is normal after HTTP response
        if (state->success) {
            printf("Connection closed by server (normal)\n");
        } else {
            printf("Connection closed by server\n");
        }
        state->complete = true;
        tcp_close(pcb);
        return ERR_OK;
    }

    if (err == ERR_OK) {
        // Check HTTP response
        char *response = (char *)p->payload;
        if (strstr(response, "200 OK") != NULL || strstr(response, "201") != NULL) {
            printf("HTTP POST successful!\n");
            state->success = true;
        } else {
            printf("HTTP response:\n%.*s\n", p->len > 200 ? 200 : p->len, response);
        }

        tcp_recved(pcb, p->tot_len);
        pbuf_free(p);

        // Close our side of the connection
        tcp_close(pcb);
        state->complete = true;
    }

    return ERR_OK;
}

// TCP connected callback
static err_t tcp_connected_cb(void *arg, struct tcp_pcb *pcb, err_t err) {
    tcp_state_t *state = (tcp_state_t *)arg;

    if (err != ERR_OK) {
        printf("TCP connection failed: %d\n", err);
        state->complete = true;
        state->success = false;
        return err;
    }

    printf("TCP connected, sending request...\n");

    // Send HTTP request
    err_t write_err = tcp_write(pcb, state->request, state->request_len, TCP_WRITE_FLAG_COPY);
    if (write_err != ERR_OK) {
        printf("TCP write failed: %d\n", write_err);
        state->complete = true;
        state->success = false;
        return write_err;
    }

    tcp_output(pcb);
    return ERR_OK;
}

// Post data to API
bool post_to_api(uint32_t timestamp, float temperature) {
    // Build JSON payload
    char json[256];
    snprintf(json, sizeof(json),
             "{\"recorded\":%lu,\"location\":\"den\",\"sensor\":\"bmp280\","
             "\"measurement\":\"temperature\",\"units\":\"C\",\"value\":%.1f}",
             timestamp, temperature);

    // Build HTTP request
    char request[512];
    int request_len = snprintf(request, sizeof(request),
                               "POST %s HTTP/1.1\r\n"
                               "Host: %s\r\n"
                               "Content-Type: application/json\r\n"
                               "Content-Length: %d\r\n"
                               "Connection: close\r\n"
                               "\r\n"
                               "%s",
                               API_PATH, API_HOST, strlen(json), json);

    printf("\nPosting to http://%s:%d%s\n", API_HOST, API_PORT, API_PATH);
    printf("JSON: %s\n", json);

    // Parse IP address
    ip_addr_t server_ip;
    if (!ipaddr_aton(API_HOST, &server_ip)) {
        printf("Invalid IP address: %s\n", API_HOST);
        return false;
    }

    // Create TCP connection
    struct tcp_pcb *pcb = tcp_new();
    if (!pcb) {
        printf("Failed to create TCP PCB\n");
        return false;
    }

    tcp_state_t state = {
        .pcb = pcb,
        .complete = false,
        .success = false,
        .request = request,
        .request_len = request_len,
        .sent = 0
    };

    tcp_arg(pcb, &state);
    tcp_err(pcb, tcp_error_cb);
    tcp_recv(pcb, tcp_recv_cb);
    tcp_sent(pcb, tcp_sent_cb);

    err_t err = tcp_connect(pcb, &server_ip, API_PORT, tcp_connected_cb);
    if (err != ERR_OK) {
        printf("TCP connect failed: %d\n", err);
        tcp_close(pcb);
        return false;
    }

    // Wait for completion (up to 10 seconds)
    for (int i = 0; i < 100 && !state.complete; i++) {
        cyw43_arch_poll();
        sleep_ms(100);
    }

    if (!state.complete) {
        printf("HTTP POST timeout\n");
        tcp_abort(pcb);
        return false;
    }

    return state.success;
}

int main() {
    stdio_init_all();
    sleep_ms(3000); // Wait for USB serial connection

    printf("\n");
    printf("========================================\n");
    printf("  BMP280 WiFi Temperature Logger\n");
    printf("========================================\n\n");

    // Initialize I2C
    printf("Initializing I2C...\n");
    i2c_init(I2C_PORT, 100 * 1000); // 100 kHz
    gpio_set_function(I2C_SDA, GPIO_FUNC_I2C);
    gpio_set_function(I2C_SCL, GPIO_FUNC_I2C);
    gpio_pull_up(I2C_SDA);
    gpio_pull_up(I2C_SCL);
    printf("I2C initialized (SDA=GPIO%d, SCL=GPIO%d)\n\n", I2C_SDA, I2C_SCL);

    // Initialize BMP280
    if (!bmp280_init()) {
        printf("\nFATAL: Failed to initialize BMP280!\n");
        printf("Check connections and I2C address (0x76 or 0x77)\n");
        return 1;
    }
    printf("\n");

    // Initialize WiFi
    printf("Initializing WiFi...\n");
    if (cyw43_arch_init()) {
        printf("FATAL: Failed to initialize WiFi chip!\n");
        return 1;
    }

    cyw43_arch_enable_sta_mode();
    printf("Connecting to '%s'...\n", WIFI_SSID);

    if (cyw43_arch_wifi_connect_timeout_ms(WIFI_SSID, WIFI_PASSWORD,
                                            CYW43_AUTH_WPA2_AES_PSK, 30000)) {
        printf("FATAL: Failed to connect to WiFi!\n");
        printf("Check SSID and password\n");
        return 1;
    }

    printf("WiFi connected!\n");
    printf("IP Address: %s\n\n", ip4addr_ntoa(netif_ip4_addr(netif_list)));

    // Main loop
    int reading_count = 0;
    while (true) {
        printf("========================================\n");
        printf("Reading #%d\n", ++reading_count);
        printf("========================================\n");

        // Read temperature
        float temperature = bmp280_read_temperature();
        if (temperature < -100) {
            printf("Failed to read temperature, skipping...\n\n");
            sleep_ms(60000);
            continue;
        }
        printf("Temperature: %.2f°C\n", temperature);

        // Get current time from NTP
        uint32_t timestamp = get_ntp_time();
        if (timestamp == 0) {
            printf("Failed to get NTP time, skipping POST...\n\n");
            sleep_ms(60000);
            continue;
        }

        // Post to API
        if (post_to_api(timestamp, temperature)) {
            printf("✓ Data posted successfully\n");
        } else {
            printf("✗ Failed to post data\n");
        }

        // Wait 5 minutes before next reading
        printf("\nWaiting 61 seconds for next reading...\n\n");
        sleep_ms(61000);
    }

    cyw43_arch_deinit();
    return 0;
}
