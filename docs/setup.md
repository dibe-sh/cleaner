# Setup Instructions

Follow these steps to set up the Cleaner App on your local machine.

## 1. Clone the Repository

First, clone the repository to your local machine:

```bash
git clone https://github.com/username/repository.git
cd repository
```

## 2. Install Dependencies

Make sure you have Go installed (version 1.16 or later is recommended). Run the following command to install dependencies:

```bash
go mod tidy
```

## 3. Set Up Test Environment (Optional)

For testing purposes, you can set up a test environment with sample files and directories using the provided script:

```bash
./setup_test_env.sh
```

This script will create a test_project directory with sample files and directories matching the app’s criteria for deletion. It’s useful for testing the app in dry-run mode.

You’re now ready to use the Cleaner App. Head over to the [Usage Guide](./usage.md) to learn how to run the app.

## Download the Executable

If you prefer to download the executable directly, follow these steps:

1. Download the appropriate executable for your operating system from the [Releases](https://github.com/dibe-sh/cleaner/releases) page.
2. Make the executable file executable (e.g., `chmod +x cleaner` on Unix-based systems).
3. Run the executable from the command line to start cleaning up directories.
4. Refer to the [Usage Guide](./usage.md) for more details on running the app with different options.
5. Enjoy a cleaner workspace with the Cleaner App!
