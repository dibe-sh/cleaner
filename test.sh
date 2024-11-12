#!/bin/bash

# Create Setup Test Environment
echo "Setting up Test Environment"
./setup_test_env.sh

# Run the main script
echo "Running Test Script"
go run main.go test_project

# Clean up the test environment
echo "Cleaning up Test Environment"
rm -rf test_project
rm cleaned_source.txt

# Done
echo "Test Completed"