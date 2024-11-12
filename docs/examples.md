---
title: "Examples"
toc: true
toc_label: "Table of Contents"
toc_icon: "align-left"
---

# Examples

These examples show how to use the Cleaner with different configurations and options.

## Example 1: Basic Directory Cleanup

To clean a `project_directory` while keeping detailed logs of actions taken:

```bash
go run main.go -root=project_directory -verbose=true -save-log=true -log-format=text
```

This command will:

 1. Traverse project_directory
 2. Delete files and folders according to the rules in cleaner_config.json
 3. Save a log of actions in cleaned_source.txt

## Example 2: Dry-Run for Safe Testing

Run a dry-run on the test_project directory to see what would be deleted:

```bash
go run main.go -root=test_project -dry-run -verbose=true
```

This command will:

 1. Traverse test_project
 2. Log the files and folders that would be deleted
 3. Not delete any files or folders

## Example 3: Custom Configuration

Suppose you want to remove only .log files and the dist directory, while excluding README.md:

1. Edit cleaner_config.json:

```json
{
    "directories_to_remove": ["dist"],
    "file_extensions_to_remove": [".log"],
    "exclude_files": ["README.md"],
    "matchRegex": false
}
```

2. Run the Cleaner with the custom configuration:

```bash
go run main.go -root=project_directory -verbose=true
```

This command will clean up dist and all .log files in project_directory, excluding README.md.

These examples cover common use cases. Customize the cleaner_config.json file to fit your exact requirements.
