# SPDX-License-Identifier: GPL-3.0-or-later
# Copyright (C) 2025 ggeoffre, LLC

▄▄   ▄▄ ▄▄  ▄▄▄▄ ▄▄▄▄   ▄▄▄  ▄▄▄▄  ▄▄ ▄▄ ▄▄▄▄▄▄ ▄▄ ▄▄  ▄▄▄  ▄▄  ▄▄
██▀▄▀██ ██ ██▀▀▀ ██▄█▄ ██▀██ ██▄█▀ ▀███▀   ██   ██▄██ ██▀██ ███▄██
██   ██ ██ ▀████ ██ ██ ▀███▀ ██      █     ██   ██ ██ ▀███▀ ██ ▀██

For Micropython development, I recommend using a Raspberry Pi 5 with at least 8GB RAM, a USB-C power supply, and a M.2 NVME drive as your development environment.  Given that you will be messing around wiht current and wires, one doesn't want to short out their USB port of their more expensive computer or laptop.  It is also not possible to use a Docker container for Circuitpython development, as you must be able to connect to various development boards using serial communications over USB.  Circuitpython is easire to get a development environment up and running as it does not need to compile your code.  It also has a huge library of modules that can be used to interact with various sensors and devices.

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
       USB-A to micro-USB Cable
       STEMMA QT / Qwiic JST SH 4-Pin Cable

    2. Recomended/Required Software:

        sudo apt update; sudo apt full-upgrade
        sudo apt-get install -y python3 python3-pip git tio minicom nano vim
        sudo shutdown -r now

    3. Check to see if your microcontroller board is mounted as a drive vai teh USB cable (below assumes Pico 2 W board):

        Enter Bootloader: Press and hold the BOOTSEL button on your Pico W, then plug it into your computer via a micro USB cable.
        echo; lsblk | grep -E "sd[a-zA-Z][0-9]"; echo; ls -laR /media/*
        cat /media/$USER/RP2350/INFO_UF2.TXT

    4. Download and install Micropython on you microcontroller board from micropython.org:

        wget -P /media/$USER/RP2350/ https://micropython.org/resources/firmware/RPI_PICO2_W-20251209-v1.27.0.uf2

    5. Use mpremote to install Micropython code and libraries on your microcontroller boards (github.com/micropython/micropython/tree/master/tools/mpremote):

        pip install mpremote

II. Use the following prompts to generate the code:

    A. Write a MicropythonPython program that talks to a BMP280 sensor connected to a Raspberry Pi Pico 2 W microcontroller board over i2c and read the current temperature.
    B. Update the code so it connects to WiFi. Use Micropython's TOML method of keeping the network id and passwrod secure.
    C. Now, add code to get the exact current time and date from the internet in integer form.
    D. Use the temperature from the BMP280, and the timestamp to POST the data to an API running on the network using the Onestring JSON format (below).
    E. Create a settings.toml file to store the api endpoint, network id and password securely.

    The String: {"recorded": 1768570200,"location": "den","sensor": "bmp280","measurement": "temperature","units": "C","value": 22.3}

    NOTE: The goal is to use Micropython to read temperature data off of a bmp280 sensor connected to a Raspberry Pi Pico 2 W microcontroller board over i2c, then post the collected temperature data to an API running on the network)

III. Run your python Code

    wget https://github.com/dafvid/micropython-bmp280/raw/refs/heads/master/bmp280.py
    mpremote cp bmp280.py :lib
    mpremote fs mkdir :lib
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

  Use the same technique for other I2C sensors like the BME280 or SHT31
  Use the same technique for other development boards like the ESP32-C6 or Arduino Nano
  Use the same technique for other platforms like Micropython, Arduino, ESP-IDF and PICO-SDK
