---
title: "Usage Guide"
layout: default
nav_order: 3  
---

# Usage Guide

This guide provides details on how to run the Cleaner, including examples for basic use and specific options.

## Basic Usage

To clean up a directory using the Cleaner, use the `-root` option to specify the target directory:

```bash
go run main.go -root=<directory> -verbose=true
```

Replace `<directory>` with the path to the directory you want to clean up. The `-verbose` flag enables detailed logging for each action taken by the app.

This command will:
 • Traverse the specified directory
 • Delete files and directories that match the cleanup criteria in the configuration
 • Provide detailed logs if -verbose=true

## Dry-Run Mode

To simulate the cleanup process without deleting any files, use the `-dry-run` option:

```bash
go run main.go -root=<directory> -dry-run=true
```

This command will:
 • Traverse the specified directory
 • Log the files and directories that would be deleted
 • Not delete any files or directories

## Logging Output

The Cleaner provides detailed logging to track the cleanup process. By default, logs are displayed in the console. To save logs to a file, use the `-log-file` option:

```bash
go run main.go -root=<directory> -save-log=true -log-format=text
```

Options for log-format:
 • text: Simple text format
 • json: JSON format, useful for programmatic analysis

## Options

The Cleaner provides several options to customize the cleanup process. Here are some common options:

- **`-root`**: Specifies the root directory to clean up.
- **`-verbose`**: Enables detailed logging for each action.
- **`-dry-run`**: Simulates the cleanup process without deleting any files.
- **`-config`**: Specifies a custom configuration file (default is `cleaner_config.json`).
- **`-help`**: Displays help information about available options.
