# SPDX-License-Identifier: GPL-3.0-or-later
# Copyright (C) 2025 ggeoffre, LLC


Task Three - Access
+++

To execute this code, run the following commands:

    1. Start the Python Runtime in a Docker container (from a command prompt within Docker Desktop):

        docker run --name python-container -p 8080:8080 -v ~/Workspace/:/usr/src/python -d python:3.13.5 sleep infinity

    2. Connect to the Python Runtime (from a command prompt within VS Code):

        docker exec -it -w /usr/src/python python-container /bin/bash

    3. Run the code (from within the Python Runtime Docker container you just connected to):

        python main.py
