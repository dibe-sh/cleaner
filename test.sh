#!/bin/bash

# Create Setup Test Environment
echo "Setting up Test Environment"
./setup_test_env.sh

# Run the main script in dry-run mode
echo "Running Test Script in Dry-Run Mode"
go run main.go -root=test_project -dry-run -verbose=true

# Run the main script in actual mode (optional if you want to test real deletion)
echo "Running Test Script in Actual Mode"
go run main.go -root=test_project -verbose=true

# Clean up the test environment
echo "Cleaning up Test Environment"
rm -rf test_project
rm -f cleaned_source.txt

# Done
echo "Test Completed"