#!/bin/bash

# Create the root test directory
mkdir -p test_project

# Navigate into the test directory
cd test_project

# Create directories
mkdir -p .git bin build dist node_modules/module src

# Create files in .git (excluded directory)
echo "git config" > .git/config

# Create files in bin (should be removed with regex and non-regex patterns)
echo "binary data" > bin/program.exe
echo "library data" > bin/helper.dll
echo "temporary data" > bin/temp_file.log  # Should be removed if matched by regex `temp.*\.log`
echo "debug info" > bin/debug_info.txt     # Should be removed if matched by regex `debug_.*\.txt`

# Create files in build (should be removed)
echo "object file data" > build/output.o

# Create files in dist (should be removed)
echo "app js data" > dist/app.js

# Create files in node_modules (should be removed)
echo "module code" > node_modules/module/index.js

# Create source files (should not be removed)
echo "package main" > src/main.go
echo "package main" > src/utils.go
touch src/.DS_Store  # Should be removed

# Create root-level files
echo "ENV variables" > .env  # Should not be removed
echo "Project documentation" > README.md  # Should not be removed
touch .DS_Store       # Should be removed
touch __debug_bin     # Should be removed
touch ___debug_bin_1100     # Should be removed
touch app__debug_bin     # Should be removed
touch tempfile.log    # Should be removed if matched by regex `temp.*\.log`
touch debug_output.txt # Should be removed if matched by regex `debug_.*\.txt`

# Navigate back to the parent directory
cd ..
