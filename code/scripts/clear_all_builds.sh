#!/bin/bash

# Clear all Build Artifacts
#
# This script is designed to clean up specific build directories and files 
# from a given directory. It allows the user to specify a directory to clean, 
# or defaults to the current directory if none is provided. The script 
# identifies and displays the sizes of certain build-related directories, 
# prompts the user for confirmation before deleting them, and also finds 
# and removes executable files and specific unwanted files like 
# "package-lock.json", "Cargo.lock", and ".DS_Store". The script includes 
# usage instructions and error handling for non-existent directories.
#
# Copyright 2024-6 ggeoffre, LLC
# SPDX-License-Identifier: BSD-3-Clause
# 

# Display usage message if -h or --help is passed
if [[ "$1" == "-h" || "$1" == "--help" ]]; then
    echo "Usage: $0 [DIR]"
    echo
    echo "This script is designed to clean up specific build directories and files from the specified directory."
    echo "If no directory is provided, it defaults to the current working directory."
    echo
    echo "Options:"
    echo "  -h, --help       Display this help message."
    echo "  DIR              The directory to clean. If not provided, the script will operate in the current directory."
    echo
    echo "Description:"
    echo "This script identifies and removes common build-related directories such as .vscode, .idea, and .gradle,"
    echo "as well as unwanted files like package-lock.json, Cargo.lock, and .DS_Store."
    echo "It prompts for confirmation before performing any deletions to prevent accidental data loss."
    echo
    echo "Examples:"
    echo "  $0                   # Clean the current directory."
    echo "  $0 /path/to/dir     # Clean the specified directory."
    exit 0
fi

# Use the directory path provided as an argument, or default to the current directory
DIR=${1:-.}

# Check if the directory exists
if [[ ! -d "$DIR" ]]; then
    echo "Error: directory '$DIR' does not exist" >&2
    exit 2
fi

# Use the directory path provided as an argument, or default to the current directory
DIR=${1:-.}

# Check if the directory exists
if [[ ! -d "$DIR" ]]; then
  echo "Error: directory '$DIR' does not exist" >&2
  exit 2
fi

echo "Cleaning build directories and files in: $DIR\n"


# BUILD DIRECTORIES AND FILES CLEANUP
# Find all matching build directories and display their sizes
find "$DIR" -type d \( -name ".vscode" -o -name ".idea" -o -name ".gradle" -o -name "__pycache__" -o -name ".venv" -o -name "bin" -o -name ".git" -o -name "build" -o -name "target" -o -name ".build" \) -exec du -sh {} +
if find "$DIR" -type d \( -name ".vscode" -o -name ".idea" -o -name ".gradle" -o -name "__pycache__" -o -name ".venv" -o -name "bin" -o -name ".git" -o -name "build" -o -name "target" -o -name ".build" \) | grep -q .; then
    find "$DIR" -type d \( -name ".vscode" -o -name ".idea" -o -name ".gradle" -o -name "__pycache__" -o -name ".venv" -o -name "bin" -o -name ".git" -o -name "build" -o -name "target" -o -name ".build" \) -exec du -sk {} + | awk '{sum += $1} END {print "Total size: " sum/1024/1024 " GB"}'
fi

# Prompt user for confirmation before deletion
if find "$DIR" -type d \( -name ".vscode" -o -name ".idea" -o -name ".gradle" -o -name "__pycache__" -o -name ".venv" -o -name "bin" -o -name ".git" -o -name "build" -o -name "target" -o -name ".build" \) | grep -q .; then
    read -r -p "Are you sure you want to remove all matching directories under $DIR? [Y/N] " answer
    case "$answer" in
        [Yy])
            find "$DIR" -type d \( -name ".vscode" -o -name ".idea" -o -name ".gradle" -o -name "__pycache__" -o -name ".venv" -o -name "bin" -o -name ".git" -o -name "build" -o -name "target" -o -name ".build" \) -exec rm -rf {} +
            ;;
        *)
            echo "Aborted."
            ;;
    esac
else
    echo "No matching directories found."
fi

# EXECUTABLE FILES AND SPECIFIC FILES CLEANUP
# Find all executable files (non-text) and display their paths
find "$DIR" -type f -exec file {} + | grep "executable" | grep -Ev "text" | cut -d':' -f1

# Prompt user for confirmation before deletion
if find "$DIR" -type f -exec file {} + | grep "executable" | grep -Ev "text" | cut -d':' -f1 | grep -q .; then
    read -r -p "Are you sure you want to remove all? [Y/N] " answer
    case "$answer" in
        [Yy])
            for filename in $(find "$DIR" -type f -exec file {} + | grep "executable" | grep -Ev "text" | cut -d':' -f1); do
                rm "$filename"
            done
            ;;        
        *)
            echo "Aborted."
            ;;
    esac
else
    echo "No executable files found."
fi

# SPECIFIC UNWANTED FILES CLEANUP
# Find all matching files and display their paths
find "$DIR" \( -name "package-lock.json" -o -name "Cargo.lock" -o -name ".DS_Store" \)

# Prompt user for confirmation before deletion
if find "$DIR" \( -name "package-lock.json" -o -name "Cargo.lock" -o -name ".DS_Store" \) | grep -q .; then
    read -r -p "Are you sure you want to remove all matching files? [Y/N] " answer
    case "$answer" in
        [Yy])
            find "$DIR" \( -name "package-lock.json" -o -name "Cargo.lock" -o -name ".DS_Store" \) -delete
            ;;
        *)
            echo "Aborted."
            ;;
    esac
else
    echo "No matching files found."
fi
