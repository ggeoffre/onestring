# SPDX-License-Identifier: GPL-3.0-or-later
# Copyright (C) 2025-2026 ggeoffre, LLC

▄▄▄▄  ▄▄  ▄▄▄▄  ▄▄▄       ▄▄▄▄ ▄▄▄▄  ▄▄ ▄▄
██▄█▀ ██ ██▀▀▀ ██▀██ ▄▄▄ ███▄▄ ██▀██ ██▄█▀
██    ██ ▀████ ▀███▀     ▄▄██▀ ████▀ ██ ██

For Pico-SDK development, I recommend using a Raspberry Pi 5 with at least 8GB RAM, a USB-C power supply, and a M.2 NVME drive as your development environment.  Given that you will be messing around wiht current and wires, one doesn't want to short out their USB port of their more expensive computer or laptop.  It is also not possible to use a Docker container for Pico-SDK development, as you must be able to connect to various development boards using serial communications over USB.  Pico-SDK is one of the more difficult development environment up and running as it does not need to compile your code.  While the Pico-SDK does have a robust set of online documentation, it has one of the smallest library of modules that can be used to interact with various sensors and devices, as it only works with Raspberry Pi Pico microcontroller chipsets.

{
    "recorded" : 1768570200,
    "location" : "den",
    "sensor" : "bmp280",
    "measurement" : "temperature",
    "units" : "C",
    "value" : 22.3
}

I. Raspberry Pi Setup
+++

    1. Recomended/Required Software:

       Raspberry Pi 5 8GB 0r 16GB w/ SD Card
       Raspberry Pi Pico 2 W (w/ header pins)
       Adafruit PiCowBell for Pico (w/ Stemma/QT)
       Adafruit BMP280 I2C Sensor (w/ Stemma/QT)
       USB-A to micro-USB Cable
       STEMMA QT / Qwiic JST SH 4-Pin Cable

    2. Recomended/Required Software:

        sudo apt-get update; sudo apt full-upgrade
        sudo apt-get install -y git tio minicom nano vim cmake python3 build-essential gcc-arm-none-eabi libnewlib-arm-none-eabi g++ libstdc++-arm-none-eabi-newlib
        sudo shutdown -r now

    3. Install the Raspberry Pi Pico SDK

        Raspberry Pi Pico SDK
        https://www.raspberrypi.com/documentation/microcontrollers/c_sdk.html
        https://github.com/raspberrypi/pico-sdk
        https://github.com/raspberrypi/pico-setup

        cd {your-pico-sdk-path}
        wget https://raw.githubusercontent.com/raspberrypi/pico-setup/master/pico_setup.sh
        chmod +x pico_setup.sh
        ./pico_setup.sh
        find ./pico -maxdepth 1 -type d

            ./pico
            ./pico/picotool
            ./pico/openocd
            ./pico/pico-playground
            ./pico/pico-examples
            ./pico/pico-extras
            ./pico/debugprobe
            ./pico/pico-sdk

    4. Check to see if your microcontroller board is mounted as a drive vai the USB cable (below assumes Pico 2 W board):

        Enter Bootloader: Press and hold the BOOTSEL button on your Pico W, then plug it into your computer via a micro USB cable.
        echo; lsblk | grep -E "sd[a-zA-Z][0-9]"; echo; ls -laR /media/*
        cat /media/$USER/RP2350/INFO_UF2.TXT

        picotool info -a -F

    5. Flash "nuke" the Raspberry Pi Pico (to remove prior runtimes and code)

        Resetting Flash memory
        https://www.raspberrypi.com/documentation/microcontrollers/pico-series.html#resetting-flash-memory

        wget -P /media/$USER/RP2350/ https://datasheets.raspberrypi.com/soft/flash_nuke.uf2
        picotool info -a -F


II. Use the following prompts to generate the code:

    A. Write a Pico-SDK code that talks to a BMP280 sensor connected to a Raspberry Pi Pico 2 W microcontroller board over i2c and read the current temperature.
    B. Update the code so it connects to WiFi.
    C. Now, add code to get the exact current time and date from the internet in integer form.
    D. Use the temperature from the BMP280, and the timestamp to POST the data to an API running on the network using the Onestring JSON format (below).
    E. Write this in C/C++ using the Pico-SDK framework and generate the CMakeLists.txt file necessary to build this code

    The String: {"recorded": 1768570200,"location": "den","sensor": "bmp280","measurement": "temperature","units": "C","value": 22.3}

    NOTE: The goal is to use Pico-SDK code to read temperature data off of a bmp280 sensor connected to a Raspberry Pi Pico 2 W microcontroller board over i2c, then post the collected temperature data to an API running on the network)


III. Build, deploy/run and monitor your Pico-SDK Code

    Build
    ~~~
    cd {your-path-to-pico-sdk-work-directory}
    mkdir build; cd build
    cmake -DPICO_BOARD=pico2_w -DPICO_PLATFORM=rp2350 ..
    make -j4

    Deploy
    ~~~
    picotool load bmp280_wifi_logger.uf2 -f
    picotool reboot -a

    Monitor
    ~~~
    find /sys/bus/usb/devices/usb*/ -name dev 2> /dev/null | grep tty 2> /dev/null
    tio /dev/ttyACM0

    NOTE: Pico SDK WiFi projects require `lwipopts.h` in project root, and CMake library link order matters! The WiFi library must be linked first so its include paths take precedence and the lwIP configuration is properly resolved.


IV. Once fully operational, copy each of the four functions individually, and add it to the following AI prompt...

  "Using the following source code, generate an AI prompt that would produce this code:"

    compared to...

  "Using the following source code to generate an AI prompt, using simple english, that would produce this code:"

    NOTE: Cute/paste the code you generated above that sucessfully ran.


V. Now take all of the original prompts from II above, and add it to the following AI prompt...

  "How could each of these prompts be written in a manner to produce 'better' code:"

    NOTE: use any instruction to replace 'better' (like cleaner, efficient, optimized, readable, concise, performant, scalable, maintainable, secure, etc.)


BEYOND BMP280 & PICO-SDK ON PICO2W:

    Use the same technique for other I2C sensors like the BME280/P581, SHT4x/3x, DHT11/DHT22, or HDC3022/AHT20
    Use the same technique for other development boards like the ESP32-C6/S3 or Arduino Nano
    Use the same technique for other platforms like Micropython, Circuitpython, Arduino and PICO-SDK
