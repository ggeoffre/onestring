# SPDX-License-Identifier: GPL-3.0-or-later
# Copyright (C) 2025 ggeoffre, LLC

 ▄▄▄  ▄▄▄▄  ▄▄▄▄  ▄▄ ▄▄ ▄▄ ▄▄  ▄▄  ▄▄▄       ▄▄▄▄ ▄▄    ▄▄
██▀██ ██▄█▄ ██▀██ ██ ██ ██ ███▄██ ██▀██ ▄▄▄ ██▀▀▀ ██    ██
██▀██ ██ ██ ████▀ ▀███▀ ██ ██ ▀██ ▀███▀     ▀████ ██▄▄▄ ██

For ANdroid-CLI development, I recommend using a Raspberry Pi 5 with at least 8GB RAM, a USB-C power supply, and a M.2 NVME drive as your development environment.  Given that you will be messing around wiht current and wires, one doesn't want to short out their USB port of their more expensive computer or laptop.  It is also not possible to use a Docker container for Circuitpython development, as you must be able to connect to various development boards using serial communications over USB.  Android-CLI is a little more dificult than the python based platforms for microcontrollers, given that you must compile your code before installing in on a board.  However, it does have a huge list of supported boards and an equally largelibrary of modules that can be used to interact with various sensors and devices.

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

        sudo apt update; sudo apt full-upgrade
        sudo apt-get install -y git tio minicom nano vim
        sudo shutdown -r now

        ARDUINO-CLI (https://docs.arduino.cc/arduino-cli/installation/)
        ~~~
        cd {your-path-to-arduino-bin}
        curl -fsSL https://raw.githubusercontent.com/arduino/arduino-cli/master/install.sh | sh
        nano ~/.bashrc
          export PATH="{your-path-to-arduino-bin}:$PATH" <= add to last line of .bashrc
        source ~/.bashrc
        echo $PATH
        arduino-cli version
        arduino-cli config dump

    3. Check to see if your microcontroller board is mounted as a drive vai the USB cable (below assumes Pico 2 W board):

        Enter Bootloader: Press and hold the BOOTSEL button on your Pico W, then plug it into your computer via a micro USB cable.
        echo; lsblk | grep -E "sd[a-zA-Z][0-9]"; echo; ls -laR /media/*
        cat /media/$USER/RP2350/INFO_UF2.TXT

    4. Download and install the Arduino Core support for your board:

        Unofficial-list-of-3rd-party-boards-support-urls.md
        https://github.com/ianhi/detect-arduino/blob/master/Unofficial-list-of-3rd-party-boards-support-urls.md

        arduino-cli config add board_manager.additional_urls https://github.com/earlephilhower/arduino-pico/releases/download/global/package_rp2040_index.json
        arduino-cli core update-index
        arduino-cli board search pico2
        arduino-cli core install rp2040:rp2040
        arduino-cli board listall pico2
        arduino-cli board list

    5. Check for and install the required libraries for this project (github.com/micropython/micropython/tree/master/tools/mpremote):

        arduino-cli config set directories.user {your-path-to-arduino-work-directory}
        arduino-cli config dump
        arduino-cli lib install "Adafruit BMP280 Library" "ArduinoJson" "NTPClient"
        arduino-cli lib list


II. Use the following prompts to generate the code:

    A. Write a Arduino sketch that talks to a BMP280 sensor connected to a Raspberry Pi Pico 2 W microcontroller board over i2c and read the current temperature.
    B. Update the sketch so it connects to WiFi.
    C. Now, add code to the sketch to get the exact current time and date from the internet in integer form.
    D. Use the temperature from the BMP280, and the timestamp to POST the data to an API running on the network using the Onestring JSON format (below).

    The String: {"recorded": 1768570200,"location": "den","sensor": "bmp280","measurement": "temperature","units": "C","value": 22.3}

    NOTE: The goal is to use this Arduino sketch to read temperature data off of a bmp280 sensor connected to a Raspberry Pi Pico 2 W microcontroller board over i2c, then post the collected temperature data to an API running on the network)


III. Build, deploy and run your Arduino Sketch

    cd {your-path-to-arduino-work-directory}
    arduino-cli sketch new onestring
    nano onestring.ino
    arduino-cli compile -vu -b rp2040:rp2040:rpipico2w ~/iot/work/arduino/onestring -p /dev/ttyACM0
    mpremote cp main.py :main.py

    find /sys/bus/usb/devices/usb*/ -name dev 2> /dev/null | grep tty 2> /dev/null
    tio /dev/ttyACM0


IV. Once fully operational, copy each of the four functions individually, and add it to the following AI prompt...

  "Using the following source code, generate an AI prompt that would produce this code:"

    compared to...

  "Using the following source code to generate an AI prompt, using simple english, that would produce this code:"

    NOTE: Cute/paste the code you generated above that sucessfully ran.


V. Now take all of the original prompts from II above, and add it to the following AI prompt...

  "How could each of these prompts be written in a manner to produce 'better' code:"

    NOTE: use any instruction to replace 'better' (like cleaner, efficient, optimized, readable, concise, performant, scalable, maintainable, secure, etc.)


BEYOND BMP280 & CIRCUITPYTHON ON PICO2W:

    Use the same technique for other I2C sensors like the BME280/P581, SHT4x/3x, DHT11/DHT22, or HDC3022/AHT20
    Use the same technique for other development boards like the ESP32-C6/S3 or Arduino Nano
    Use the same technique for other platforms like Micropython, Circuitpython, ESP-IDF and PICO-SDK
