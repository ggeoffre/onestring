# SPDX-License-Identifier: GPL-3.0-or-later
# Copyright (C) 2025 ggeoffre, LLC


Task Two - API
+++

To execute this code, run the following commands:

    1. Start the Python Runtime in a Docker container (from a command prompt within Docker Desktop):

        docker run --name python-container -p 8080:8080 -v ~/Workspace/:/usr/src/python -d python:3.13.5 sleep infinity

    2. Connect to the Python Runtime (from a command prompt within VS Code):

        docker exec -it -w /usr/src/python python-container /bin/bash

    3. Update Pip and Install the Requirements (from within the Python Runtime Docker container you just connected to):

        pip install --upgrade pip; pip install -r requirements.txt
        pip install -r requirements.txt --root-user-action ignore (necessary for postgres)
        NOTE: this is within a docker container that can be recreated, so you can ignore the root user action.

    4. Run the code (from within the Python Runtime Docker container you just connected to):

        python app.py


To test this code, run the following commands:

   1. Right click on the script file you want to run and select "Open in Terminal".

   2. Run the following command from the command prompt (depends on which command is supported):

        sh ./test_all_apis.sh CURL

        -or-

        sh ./test_all_apis.sh WGET
