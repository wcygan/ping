#!/bin/bash

# Get the current date components
Year=$(date +%Y)
Month=$(date +%m)
Day=$(date +%d)

# Create the directory structure
mkdir -p "notes/$Year/$Month/$Day"

# Get the local time in the desired format including seconds and timezone in military time
LocalTime=$(date +"%H:%M:%S-%Z")

# Construct the file name
FileName="$LocalTime.md"

# Full path to the new file
FilePath="notes/$Year/$Month/$Day/$FileName"

# Create the file
touch "$FilePath"

# Open the file in VSCode
code "$FilePath"