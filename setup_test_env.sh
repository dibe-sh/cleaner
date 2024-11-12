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
dd if=/dev/zero of=bin/program.exe bs=1K count=5   # 5 KB
dd if=/dev/zero of=bin/helper.dll bs=1K count=3     # 3 KB
dd if=/dev/zero of=bin/temp_file.log bs=1K count=1  # 1 KB (matches regex `temp.*\.log`)
dd if=/dev/zero of=bin/debug_info.txt bs=1K count=2 # 2 KB (matches regex `debug_.*\.txt`)

# Create files in build (should be removed)
dd if=/dev/zero of=build/output.o bs=1K count=4     # 4 KB

# Create files in dist (should be removed)
dd if=/dev/zero of=dist/app.js bs=1K count=6        # 6 KB

# Create files in node_modules (should be removed)
dd if=/dev/zero of=node_modules/module/index.js bs=1K count=2  # 2 KB

# Create source files (should not be removed)
echo "package main" > src/main.go
echo "package main" > src/utils.go
dd if=/dev/zero of=src/.DS_Store bs=1K count=1      # 1 KB (should be removed)

# Create root-level files
echo "ENV variables" > .env                         # Should not be removed
echo "Project documentation" > README.md            # Should not be removed
dd if=/dev/zero of=.DS_Store bs=1K count=1          # 1 KB (should be removed)
dd if=/dev/zero of=__debug_bin bs=1K count=1        # 1 KB (should be removed)
dd if=/dev/zero of=___debug_bin_1100 bs=1K count=1  # 1 KB (should be removed)
dd if=/dev/zero of=app__debug_bin bs=1K count=1     # 1 KB (should be removed)
dd if=/dev/zero of=tempfile.log bs=1K count=2       # 2 KB (matches regex `temp.*\.log`)
dd if=/dev/zero of=debug_output.txt bs=1K count=3   # 3 KB (matches regex `debug_.*\.txt`)

# Navigate back to the parent directory
cd ..