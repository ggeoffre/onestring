# SPDX-License-Identifier: GPL-3.0-or-later
# Copyright (C) 2025-2026 ggeoffre, LLC

cmake_minimum_required(VERSION 3.13)

# Include Pico SDK
include($ENV{PICO_SDK_PATH}/external/pico_sdk_import.cmake)

# Project name
project(bmp280_wifi_logger C CXX ASM)

set(CMAKE_C_STANDARD 11)
set(CMAKE_CXX_STANDARD 17)

# Initialize the Pico SDK
pico_sdk_init()

# Verify we're building for a Pico W variant (with WiFi)
if(NOT PICO_BOARD MATCHES "pico.*_w")
    message(FATAL_ERROR "This project requires a Pico W board (pico_w or pico2_w). Use: cmake -DPICO_BOARD=pico2_w -DPICO_PLATFORM=rp2350 ..")
endif()

# Add executable
add_executable(bmp280_wifi_logger
    main.c
)

# Link required libraries
# NOTE: Order matters! pico_cyw43_arch must come before pico_stdlib
target_link_libraries(bmp280_wifi_logger
    pico_cyw43_arch_lwip_threadsafe_background
    pico_stdlib
    hardware_i2c
    hardware_gpio
)

# Add include directory for lwipopts.h
target_include_directories(bmp280_wifi_logger PRIVATE
    ${CMAKE_CURRENT_LIST_DIR}
)

# Enable USB output, disable UART output
pico_enable_stdio_usb(bmp280_wifi_logger 1)
pico_enable_stdio_uart(bmp280_wifi_logger 0)

# Create map/bin/hex/uf2 files
pico_add_extra_outputs(bmp280_wifi_logger)
