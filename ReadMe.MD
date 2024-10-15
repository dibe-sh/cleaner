# Cleaner

[![CI](https://github.com/dibe-sh/cleaner/actions/workflows/ci.yml/badge.svg)](https://github.com/dibe-sh/cleaner/actions/workflows/ci.yml)

Cleaner is a cross-platform command-line tool written in Go that recursively scans directories to remove specified files and directories. It offers features like concurrency for faster execution, dry-run mode, exclusion patterns, custom logging formats, and more. Cleaner is configurable via a JSON file, allowing for flexible and powerful cleaning operations.

## Features

- Recursive Scanning: Scans directories to any depth.
- Configurable Cleaning: Remove directories and files based on a JSON configuration.
- Concurrency: Utilizes multiple CPU cores for faster execution.
- Dry-Run Mode: Simulate the cleaning process without deleting files.
- Exclusion Patterns: Exclude specific directories and files from being removed.
- Custom Logging Formats: Choose between text or JSON log formats.
- Verbose and Silent Modes: Control console output with the --verbose flag.
- No External Dependencies: Built as static binaries, requiring no dependencies on the target platforms.

## Table of Contents

- [Cleaner](#cleaner)
  - [Features](#features)
  - [Table of Contents](#table-of-contents)
  - [Installation](#installation)
    - [Download Pre-built Binaries](#download-pre-built-binaries)
    - [Build from Source](#build-from-source)
  - [Usage](#usage)
    - [Basic Usage](#basic-usage)
  - [Help](#help)
  - [Configuration](#configuration)
    - [Default Configuration (cleaner\_config.json)](#default-configuration-cleaner_configjson)
    - [Custom Configuration](#custom-configuration)
  - [Command-Line Flags](#command-line-flags)
  - [Examples](#examples)
  - [Contributing](#contributing)
  - [License](#license)

## Installation

### Download Pre-built Binaries

Pre-built binaries are available for multiple platforms:

- Windows (amd64)
- macOS (amd64, arm64)
- Linux (amd64, arm64)
   Download the appropriate binary for your platform from the [Releases page](https://github.com/dibe-sh/cleaner/releases). The binaries are self-contained and do not require any external dependencies.

Example:

- For macOS on Apple Silicon (M1/M2):

Download cleaner-darwin-arm64 and make it executable:

```bash
chmod +x cleaner-darwin-arm64
```

- For Windows:

Download cleaner-windows-amd64.exe and run it directly.

### Build from Source

If you prefer to build Cleaner from source, ensure you have Go installed (version 1.16 or later).

Steps:

Clone the Repository:

```bash
git clone https://github.com/dibe-sh/cleaner.git
cd cleaner
```

Build the Executable:

```bash
go build -ldflags="-s -w" -o cleaner
```

The -ldflags="-s -w" flags reduce the binary size by stripping debug information.
Make the Executable (if necessary):

```bash
chmod +x cleaner
```

## Usage

### Basic Usage

```bash
./cleaner --root /path/to/start
```

This command starts the cleaning process from the specified root directory using the default configuration.

## Help

```bash
./cleaner --help
```

Displays the list of available command-line flags and their descriptions.

## Configuration

Cleaner uses a JSON configuration file to specify which directories and files to remove or exclude.

### Default Configuration (cleaner_config.json)

```json
{
    "directories_to_remove": [
        "node_modules",
        "dist",
        "build",
        "bin",
        ".next",
        ".turbo",
        ".idea",
        ".cache"
    ],
    "file_extensions_to_remove": [".DS_Store", "__debug_bin"],
    "exclude_directories": [".git", ".svn"],
    "exclude_files": []
}
```

- directories_to_remove: List of directory names to remove.
- file_extensions_to_remove: List of file extensions to remove.
- exclude_directories: List of directory names to exclude from removal.
- exclude_files: List of file names to exclude from removal.

### Custom Configuration

You can specify a custom configuration file using the --config flag:

```bash
./cleaner --root /path/to/start --config /path/to/custom_config.json
```

## Command-Line Flags

- --root: Root directory to start scanning (default is current directory .).
- --config: Path to the JSON configuration file (default is cleaner_config.json).
- --verbose: Enable verbose console output (true or false, default is true).
- --save-log: Enable saving logs to file (true or false, default is true).
- --dry-run: Simulate the cleaning process without deleting files (default is false).
- --log-format: Log format, either 'text' or 'json' (default is 'text').

## Examples

- Run with Default Settings

```bash
./cleaner --root /path/to/start
```

- Dry-Run Mode
   Simulate the cleaning process without deleting any files:

```bash
./cleaner --root /path/to/start --dry-run
```

- Silent Mode without Saving Logs
   Run without console output and without saving logs:

```bash
./cleaner --root /path/to/start --verbose=false --save-log=false
```

Custom Configuration and Log Format
Use a custom configuration file and save logs in JSON format:

```bash
./cleaner --root /path/to/start --config custom_config.json --log-format json
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request on GitHub.

Steps to Contribute

- Fork the repository.
- Create a new branch (git checkout -b feature/YourFeature).
- Commit your changes (git commit -am 'Add some feature').
- Push to the branch (git push origin feature/YourFeature).
- Open a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.
