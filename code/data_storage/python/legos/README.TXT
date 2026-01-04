# SPDX-License-Identifier: GPL-3.0-or-later
# Copyright (C) 2025-2026 ggeoffre, LLC


Task Two - LEGOS
+++

To execute this code, run the following commands:

    1. Start the Mongo Database in a Docker container (from a command prompt within Docker Desktop):

        docker run --name mongo-container -p 27017:27017 -d mongo:8.2.2

    2. Start the Python Runtime in a Docker container (from a command prompt within Docker Desktop):

        docker run --name python-container -p 8080:8080 -v ~/Workspace/:/usr/src/python -d python:3.13.5 sleep infinity

    3. Connect to the Python Runtime (from a command prompt within VS Code):

        docker exec -it -w /usr/src/python python-container /bin/bash

    4. Update Pip and Install the Requirements (from within the Python Runtime Docker container you just connected to):

        pip install --upgrade pip; pip install -r requirements.txt

    5. Run the code (from within the Python Runtime Docker container you just connected to):

        export DATA_HOSTNAME="YOUR_LOACAL_MACHINE_IP_ADDRESS"
        export DATA_ACCESS=mongo; python app.py

        -or-

        export DATA_HOSTNAME="YOUR_LOACAL_MACHINE_IP_ADDRESS"
        export DATA_ACCESS=mongo; python app.py runserver 0.0.0.0:8080


To test this code, run the following commands:

   1. Right click on the script file you want to run and select "Open in Terminal".

   2. Run the following command from the command prompt (depends on which command is supported):

        sh ./test_all_apis.sh CURL

        -or-

        sh ./test_all_apis.sh WGET
