#!/bin/bash

# all_language_files.sh
#
# DESCRIPTION:
#   Concatenates all source files from multiple programming languages into
#   separate text files, organized by project. Projects are identified by
#   their standard manifest/configuration files at the project root.
#
# SUPPORTED LANGUAGES:
#   - Python:  identified by requirements.txt  -> all_python_files.txt
#   - Java:    identified by build.gradle      -> all_java_files.txt
#   - Go:      identified by go.mod            -> all_go_files.txt
#   - Rust:    identified by Cargo.toml        -> all_rust_files.txt
#   - Swift:   identified by Package.swift     -> all_swift_files.txt
#
# OUTPUT FORMAT:
#   Each output file contains projects separated by headers, with each source
#   file preceded by its full path. Files are sorted alphabetically within
#   each project.
#
# EXCLUSIONS:
#   - Hidden directories (starting with .)
#   - Swift .build directories (dependency checkouts)
#
# SPDX-License-Identifier: GPL-3.0-or-later
# Copyright (C) 2025-2026 ggeoffre, LLC
#

set -euo pipefail

# -----------------------------------------------------------------------------
# Usage function - displays help information
# -----------------------------------------------------------------------------
usage() {
    cat << EOF
Usage: $(basename "$0") [OPTIONS] [SEARCH_DIR]

Concatenate source files from multiple languages into separate text files,
organized by project.

ARGUMENTS:
    SEARCH_DIR          Directory to search for projects (default: current directory)

OPTIONS:
    -h, --help          Show this help message and exit
    -v, --version       Show version information
    -l, --list          List supported languages and their identifiers
    -o, --output-dir    Directory for output files (default: current directory)
    -q, --quiet         Suppress progress messages

SUPPORTED LANGUAGES:
    Language    Identifier          Output File
    --------    ----------          -----------
    Python      requirements.txt    all_python_files.txt
    Java        build.gradle        all_java_files.txt
    Go          go.mod              all_go_files.txt
    Rust        Cargo.toml          all_rust_files.txt
    Swift       Package.swift       all_swift_files.txt

EXAMPLES:
    $(basename "$0")                     # Search current directory
    $(basename "$0") /path/to/projects   # Search specific directory
    $(basename "$0") -o ./output .       # Output files to ./output directory
    $(basename "$0") -q .                # Run quietly without progress messages

OUTPUT FORMAT:
    Each output file contains:
    ========================================
    PROJECT: <project_name>
    PATH: <full_path_to_project>
    ========================================

    ----------------------------------------
    FILE: <full_path_to_source_file>
    ----------------------------------------
    <file contents>

NOTES:
    - Hidden directories (starting with .) are excluded
    - Swift .build directories are excluded to avoid dependency checkouts
    - Files within each project are sorted alphabetically by path

EOF
}

# -----------------------------------------------------------------------------
# Version function
# -----------------------------------------------------------------------------
version() {
    echo "$(basename "$0") version 1.0.0"
}

# -----------------------------------------------------------------------------
# List supported languages
# -----------------------------------------------------------------------------
list_languages() {
    cat << EOF
Supported Languages:

Language    Project Identifier    File Extension    Output File
--------    ------------------    --------------    -----------
Python      requirements.txt      .py               all_python_files.txt
Java        build.gradle          .java             all_java_files.txt
Go          go.mod                .go               all_go_files.txt
Rust        Cargo.toml            .rs               all_rust_files.txt
Swift       Package.swift         .swift            all_swift_files.txt

EOF
}

# -----------------------------------------------------------------------------
# Parse command line arguments
# -----------------------------------------------------------------------------
SEARCH_DIR="."
OUTPUT_DIR="."
QUIET=false

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            usage
            exit 0
            ;;
        -v|--version)
            version
            exit 0
            ;;
        -l|--list)
            list_languages
            exit 0
            ;;
        -o|--output-dir)
            OUTPUT_DIR="$2"
            shift 2
            ;;
        -q|--quiet)
            QUIET=true
            shift
            ;;
        -*)
            echo "Error: Unknown option: $1" >&2
            echo "Use -h or --help for usage information." >&2
            exit 1
            ;;
        *)
            SEARCH_DIR="$1"
            shift
            ;;
    esac
done

# -----------------------------------------------------------------------------
# Validate inputs
# -----------------------------------------------------------------------------

# Check if search directory exists
if [[ ! -d "$SEARCH_DIR" ]]; then
    echo "Error: Search directory does not exist: $SEARCH_DIR" >&2
    exit 1
fi

# Create output directory if it doesn't exist
if [[ ! -d "$OUTPUT_DIR" ]]; then
    mkdir -p "$OUTPUT_DIR"
fi

# -----------------------------------------------------------------------------
# Helper function to print progress messages
# -----------------------------------------------------------------------------
log() {
    if [[ "$QUIET" == false ]]; then
        echo "$1"
    fi
}

# =============================================================================
# PYTHON FILES
# =============================================================================
# Project identifier: requirements.txt at project root
# File extension: .py
# Excludes: hidden directories (.venv, .tox, etc.)
# =============================================================================

PYTHON_OUTPUT="$OUTPUT_DIR/all_python_files.txt"

# Initialize empty output file
> "$PYTHON_OUTPUT"

log "Processing Python projects..."

# Find all project roots by locating requirements.txt files
# -path '*/.*' -prune: skip hidden directories
# -o -name "requirements.txt": find requirements.txt files
# -type f -print: only regular files, print results
find "$SEARCH_DIR" -path '*/.*' -prune -o -name "requirements.txt" -type f -print | while read -r req_file; do
    # Extract project directory and name from the requirements.txt path
    project_dir=$(dirname "$req_file")
    project_name=$(basename "$project_dir")

    # Write project header to output file
    echo "========================================" >> "$PYTHON_OUTPUT"
    echo "PROJECT: $project_name" >> "$PYTHON_OUTPUT"
    echo "PATH: $project_dir" >> "$PYTHON_OUTPUT"
    echo "========================================" >> "$PYTHON_OUTPUT"
    echo "" >> "$PYTHON_OUTPUT"

    # Find all .py files within this project, sorted alphabetically
    find "$project_dir" -path '*/.*' -prune -o -name "*.py" -type f -print | sort | while read -r py_file; do
        # Write file header with full path
        echo "----------------------------------------" >> "$PYTHON_OUTPUT"
        echo "FILE: $py_file" >> "$PYTHON_OUTPUT"
        echo "----------------------------------------" >> "$PYTHON_OUTPUT"
        # Append file contents
        cat "$py_file" >> "$PYTHON_OUTPUT"
        echo "" >> "$PYTHON_OUTPUT"
        echo "" >> "$PYTHON_OUTPUT"
    done

    echo "" >> "$PYTHON_OUTPUT"
done

log "Python output written to $PYTHON_OUTPUT"

# =============================================================================
# JAVA FILES
# =============================================================================
# Project identifier: build.gradle at project root
# File extension: .java
# Excludes: hidden directories
# =============================================================================

JAVA_OUTPUT="$OUTPUT_DIR/all_java_files.txt"

# Initialize empty output file
> "$JAVA_OUTPUT"

log "Processing Java projects..."

# Find all project roots by locating build.gradle files
find "$SEARCH_DIR" -path '*/.*' -prune -o -name "build.gradle" -type f -print | while read -r gradle_file; do
    # Extract project directory and name
    project_dir=$(dirname "$gradle_file")
    project_name=$(basename "$project_dir")

    # Write project header
    echo "========================================" >> "$JAVA_OUTPUT"
    echo "PROJECT: $project_name" >> "$JAVA_OUTPUT"
    echo "PATH: $project_dir" >> "$JAVA_OUTPUT"
    echo "========================================" >> "$JAVA_OUTPUT"
    echo "" >> "$JAVA_OUTPUT"

    # Find all .java files within this project, sorted alphabetically
    find "$project_dir" -path '*/.*' -prune -o -name "*.java" -type f -print | sort | while read -r java_file; do
        # Write file header with full path
        echo "----------------------------------------" >> "$JAVA_OUTPUT"
        echo "FILE: $java_file" >> "$JAVA_OUTPUT"
        echo "----------------------------------------" >> "$JAVA_OUTPUT"
        # Append file contents
        cat "$java_file" >> "$JAVA_OUTPUT"
        echo "" >> "$JAVA_OUTPUT"
        echo "" >> "$JAVA_OUTPUT"
    done

    echo "" >> "$JAVA_OUTPUT"
done

log "Java output written to $JAVA_OUTPUT"

# =============================================================================
# GO FILES
# =============================================================================
# Project identifier: go.mod at project root
# File extension: .go
# Excludes: hidden directories
# =============================================================================

GO_OUTPUT="$OUTPUT_DIR/all_go_files.txt"

# Initialize empty output file
> "$GO_OUTPUT"

log "Processing Go projects..."

# Find all project roots by locating go.mod files
find "$SEARCH_DIR" -path '*/.*' -prune -o -name "go.mod" -type f -print | while read -r mod_file; do
    # Extract project directory and name
    project_dir=$(dirname "$mod_file")
    project_name=$(basename "$project_dir")

    # Write project header
    echo "========================================" >> "$GO_OUTPUT"
    echo "PROJECT: $project_name" >> "$GO_OUTPUT"
    echo "PATH: $project_dir" >> "$GO_OUTPUT"
    echo "========================================" >> "$GO_OUTPUT"
    echo "" >> "$GO_OUTPUT"

    # Find all .go files within this project, sorted alphabetically
    find "$project_dir" -path '*/.*' -prune -o -name "*.go" -type f -print | sort | while read -r go_file; do
        # Write file header with full path
        echo "----------------------------------------" >> "$GO_OUTPUT"
        echo "FILE: $go_file" >> "$GO_OUTPUT"
        echo "----------------------------------------" >> "$GO_OUTPUT"
        # Append file contents
        cat "$go_file" >> "$GO_OUTPUT"
        echo "" >> "$GO_OUTPUT"
        echo "" >> "$GO_OUTPUT"
    done

    echo "" >> "$GO_OUTPUT"
done

log "Go output written to $GO_OUTPUT"

# =============================================================================
# RUST FILES
# =============================================================================
# Project identifier: Cargo.toml at project root
# File extension: .rs
# Excludes: hidden directories
# =============================================================================

RUST_OUTPUT="$OUTPUT_DIR/all_rust_files.txt"

# Initialize empty output file
> "$RUST_OUTPUT"

log "Processing Rust projects..."

# Find all project roots by locating Cargo.toml files
find "$SEARCH_DIR" -path '*/.*' -prune -o -name "Cargo.toml" -type f -print | while read -r cargo_file; do
    # Extract project directory and name
    project_dir=$(dirname "$cargo_file")
    project_name=$(basename "$project_dir")

    # Write project header
    echo "========================================" >> "$RUST_OUTPUT"
    echo "PROJECT: $project_name" >> "$RUST_OUTPUT"
    echo "PATH: $project_dir" >> "$RUST_OUTPUT"
    echo "========================================" >> "$RUST_OUTPUT"
    echo "" >> "$RUST_OUTPUT"

    # Find all .rs files within this project, sorted alphabetically
    find "$project_dir" -path '*/.*' -prune -o -name "*.rs" -type f -print | sort | while read -r rs_file; do
        # Write file header with full path
        echo "----------------------------------------" >> "$RUST_OUTPUT"
        echo "FILE: $rs_file" >> "$RUST_OUTPUT"
        echo "----------------------------------------" >> "$RUST_OUTPUT"
        # Append file contents
        cat "$rs_file" >> "$RUST_OUTPUT"
        echo "" >> "$RUST_OUTPUT"
        echo "" >> "$RUST_OUTPUT"
    done

    echo "" >> "$RUST_OUTPUT"
done

log "Rust output written to $RUST_OUTPUT"

# =============================================================================
# SWIFT FILES
# =============================================================================
# Project identifier: Package.swift at project root
# File extension: .swift
# Excludes: hidden directories AND .build directories (Swift package manager
#           stores downloaded dependencies in .build/checkouts which contain
#           their own Package.swift files)
# =============================================================================

SWIFT_OUTPUT="$OUTPUT_DIR/all_swift_files.txt"

# Initialize empty output file
> "$SWIFT_OUTPUT"

log "Processing Swift projects..."

# Find all project roots by locating Package.swift files
# Additional exclusion: .build directories to avoid dependency checkouts
find "$SEARCH_DIR" \( -path '*/.*' -o -path '*/.build' \) -prune -o -name "Package.swift" -type f -print | while read -r package_file; do
    # Extract project directory and name
    project_dir=$(dirname "$package_file")
    project_name=$(basename "$project_dir")

    # Write project header
    echo "========================================" >> "$SWIFT_OUTPUT"
    echo "PROJECT: $project_name" >> "$SWIFT_OUTPUT"
    echo "PATH: $project_dir" >> "$SWIFT_OUTPUT"
    echo "========================================" >> "$SWIFT_OUTPUT"
    echo "" >> "$SWIFT_OUTPUT"

    # Find all .swift files within this project, sorted alphabetically
    # Exclude .build directories to avoid compiled dependencies
    find "$project_dir" \( -path '*/.*' -o -path '*/.build' \) -prune -o -name "*.swift" -type f -print | sort | while read -r swift_file; do
        # Write file header with full path
        echo "----------------------------------------" >> "$SWIFT_OUTPUT"
        echo "FILE: $swift_file" >> "$SWIFT_OUTPUT"
        echo "----------------------------------------" >> "$SWIFT_OUTPUT"
        # Append file contents
        cat "$swift_file" >> "$SWIFT_OUTPUT"
        echo "" >> "$SWIFT_OUTPUT"
        echo "" >> "$SWIFT_OUTPUT"
    done

    echo "" >> "$SWIFT_OUTPUT"
done

log "Swift output written to $SWIFT_OUTPUT"

# =============================================================================
# Completion summary
# =============================================================================
log ""
log "Processing complete. Output files:"
log "  - $PYTHON_OUTPUT"
log "  - $JAVA_OUTPUT"
log "  - $GO_OUTPUT"
log "  - $RUST_OUTPUT"
log "  - $SWIFT_OUTPUT"
