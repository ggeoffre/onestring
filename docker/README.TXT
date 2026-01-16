// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2025-2026 ggeoffre, LLC

██████╗  ██████╗  ██████╗██╗  ██╗███████╗██████╗
██╔══██╗██╔═══██╗██╔════╝██║ ██╔╝██╔════╝██╔══██╗
██║  ██║██║   ██║██║     █████╔╝ █████╗  ██████╔╝
██║  ██║██║   ██║██║     ██╔═██╗ ██╔══╝  ██╔══██╗
██████╔╝╚██████╔╝╚██████╗██║  ██╗███████╗██║  ██║
╚═════╝  ╚═════╝  ╚═════╝╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝

For the series of tasks associated with this project, Docker is being used to cleanly run the necessary databases and runtimes required for both development and deployment.  The minimum requirements are to have both Docker and VS Code installed on your machine.  Everyone understands why VC Code is required, but wy Docker?

It is the easiest way to get each of the five databases up and running locally on your machine.  Cassandra, Mongo, PostgreSQL, MySQL, and Redis.  You can even be specific as to the exact version of each database you want to use.  This is a huge help as it ensures that we all are managing dependencies in a consistent manner.


DATABASES
+++
The docker command for starting the databases quite simple.  It is broken into just a fiew parts:

  docker run -d --name mongo -p 27017:27017 mongo:8.2.2

    docker run :: means run the container in the background
    -d :: means run detached from the terminal
    --name mongo :: means give the container the name mongo
    -p 27017:27017 :: means map port 27017 on the host to port 27017 on the container
    mongo:8.2.2 :: means use the mongo image version 8.2.2

Here are the run commands for all of the databases:

    docker run -d --name cassandra -p 9042:9042 cassandra:5.0.6
    docker run -d --name mongo -p 27017:27017 mongo:8.2.2
    docker run -d --name postgres -p 5432:5432 postgres:18.1
    docker run -d --name mysql -p 3306:3306 mysql:9.5.0
    docker run -d --name redis -p 6379:6379 redis:8.4.0


RUNTIMES
+++
This one can be a little more confusing, as we will be running the language runtime "inside" the container, yet mounting to a directory "outside" the filesystem.  This is done to ensure that the container has access to the necessary files and dependencies, while also allowing for easy access to the container's filesystem from the host machine.

    docker run -d --name python-container -p 8000:8000 -v /path/to/project:/app python:3.9

        docker run :: means tun the container in the background
        -d :: means run detached from the terminal
        --name python-container :: means give the container the name python-container
        -p 8000:8000 :: means map port 8000 on the host to port 8000 on the container
        -v [/path/to/project]:/app :: means mount the local directory [/path/to/project] on the host to the directory /app on the container
        python:3.9 :: means use the python image version 3.9

    Here are the run commands for all of the runtimes:

        docker run -d --name go-container -p 8080:8080 -v [/path/to/project]:/app golang:1.25.4
        docker run -d --name java-container -p 8080:8080 -v [/path/to/project]:/app gradle:9.2.1-jdk21
        docker run -d --name python-container -p 8080:8080 -v [/path/to/project]:/app python:3.13.5
        docker run -d --name rust-container -p 8080:8080 -v [/path/to/project]:/app rust:1.91.1
        docker run -d --name swift-container -p 8080:8080 -v [/path/to/project]:/app swift:6.2.1

    NOTE: be very careful to update the [/path/to/project] to match your local directory path.  THis is the part the seems to confuse a lot of people.


CONNECTING TO THE RUNTIME
+++
In order to compile and run your code inside of the container, you need to have access to the command prompt inside your container.  And Docker has a way of doing jsut that.

    docker exec -it -w /usr/src/python python-container /bin/bash

        docker exec means execute the following command inside of the container
        -it means interactive mode
        -w /usr/src/python means change the working directory to /usr/src/python
        python-container means the name of the container
        /bin/bash means execute the bash shell inside the container

    Here are the run commands for connecting to all of the runtimes:

        docker exec -it -w /usr/src/python python-container /bin/bash
        docker exec -it -w /usr/src/java java-container /bin/bash
        docker exec -it -w /usr/src/rust rust-container /bin/bash
        docker exec -it -w /usr/src/swift swift-container /bin/bash
        docker exec -it -w /usr/src/go go-container /bin/bash

    IMPORTANT:  The ip address of the localhost runnign inside each container is different than the localhost of your computer itself.  To find the ip address of your computer, you must be running a command prompt on your local computer itself.

    IP ADDRESS:  Use ipconfig on Windows or ifconfig on Linux/MacOS to find the ip address of your computer.


Links
+++

Docker Hub - listing of all publicll available docker images
https://hub.docker.com/search

Getting Started with Docker
https://docs.docker.com/get-started/

Docker's Official YouTube Channel
https://www.youtube.com/user/dockerrun
