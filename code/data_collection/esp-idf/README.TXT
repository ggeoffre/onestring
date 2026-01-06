# SPDX-License-Identifier: GPL-3.0-or-later
# Copyright (C) 2025-2026 ggeoffre, LLC

▄▄▄▄▄  ▄▄▄▄ ▄▄▄▄      ▄▄ ▄▄▄▄  ▄▄▄▄▄
██▄▄  ███▄▄ ██▄█▀ ▄▄▄ ██ ██▀██ ██▄▄
██▄▄▄ ▄▄██▀ ██        ██ ████▀ ██

ESP-IDF is similar to PICO-SDK in that they are both C/C++ development environments.  Together with Arduino, all three use C/C++ as their development language.  But that is where the similarities end.  The instalation process is different for each, and the command line argumetns are different.  Not to mention teh fact that the entire chip architecture and board orginizartion is different as well.  The software developmetn lifecycle is similar enought that we can still follow the same basic processes.

{
    "recorded" : 1768570200,
    "location" : "den",
    "sensor" : "bmp280",
    "measurement" : "temperature",
    "units" : "C",
    "value" : 22.3
}

I. ESP-IDF Setup
+++

    1. Recomended/Required Software:

       Raspberry Pi 5 8GB 0r 16GB w/ SD Card
       ESP32-S3-DevKitC-1 v1.1
       Adafruit BMP280 I2C Sensor (w/ Stemma/QT)
       USB-A to micro-USB Cable
       STEMMA QT / Qwiic JST SH 4-Pin Cable

    2. Recomended/Required Software:

        sudo apt-get update; sudo apt full-upgrade
        sudo apt-get install -y git tio minicom nano vim cmake python3 python3-pip python3-venv wget flex bison gpref ninja-build ccache libffi-dev libssl-dev dfu-util libusb-1.0-0
        sudo shutdown -r now

    3. Install the ESP-IDF SDK

        ESP-IDF SDK
        https://docs.espressif.com/projects/esp-idf/en/stable/esp32/get-started/index.html
        https://github.com/espressif/esp-idf

        mkdir -p ~/esp
        cd ~/esp
        git clone -b v5.5.2 --recursive https://github.com/espressif/esp-idf.git

        cd ~/esp/esp-idf
        ./install.sh all

        find ~/esp -maxdepth 1 -type d

            ./esp/esp-idf
            ./esp/esp-idf/tools
            ./esp/esp-idf/.git
            ./esp/esp-idf/docs
            ./esp/esp-idf/.gitlab
            ./esp/esp-idf/.github
            ./esp/esp-idf/components
            ./esp/esp-idf/examples

        . $HOME/esp/esp-idf/export.sh
        idf.py --version

    4. Install esptool:

        python -m pip install --upgrade esptool --break-system-packages
        esptool -p /dev/ttyACM0 chip-id
        esptool -p /dev/ttyACM0 erase-flash


II. Use the following prompts to generate the code:

    A. Write ESP-IDF code that talks to a BMP280 sensor connected to a ESP32-S3-DevKitC-1 v1.1 microcontroller board over i2c and read the current temperature.
    B. Update the code so it connects to WiFi.
    C. Now, add code to get the exact current time and date from the internet in integer form.
    D. Use the temperature from the BMP280, and the timestamp to POST the data to an API running on the network using the Onestring JSON format (below).
    E. Write this in C/C++ using the ESP-IDF framework and generate the CMakeLists.txt file necessary to build this code

    The String: {"recorded": 1768570200,"location": "den","sensor": "bmp280","measurement": "temperature","units": "C","value": 22.3}

    NOTE: The goal is to use ESP-IDF code to read temperature data off of a bmp280 sensor connected to a ESP32-S3-DevKitC-1 v1.1 microcontroller board over i2c, then post the collected temperature data to an API running on the network.


III. Build, deploy/run and monitor your Pico-SDK Code

    idf.py set-target esp32s3
    idf.py build
    idf.py -p /dev/ttyACM0 flash monitor


IV. Once fully operational, copy each of the four functions individually, and add it to the following AI prompt...

  "Using the following source code, generate an AI prompt that would produce this code:"

    compared to...

  "Using the following source code to generate an AI prompt, using simple english, that would produce this code:"

    NOTE: Cut/paste the code you generated above that sucessfully ran.


V. Now take all of the original prompts from II above, and add it to the following AI prompt...

  "How could each of these prompts be written in a manner to produce 'better' code:"

    NOTE: use any instruction to replace 'better' (like cleaner, efficient, optimized, readable, concise, performant, scalable, maintainable, secure, etc.)


BEYOND BMP280 & ESP-IDF ON ESP32-S3:

    Use the same technique for other I2C sensors like the BME280/P581, SHT4x/3x, DHT11/DHT22, or HDC3022/AHT20
    Use the same technique for other development boards like the ESP32-C6/S3 or Arduino Nano
    Use the same technique for other platforms like Micropython, Circuitpython, Arduino and PICO-SDK
