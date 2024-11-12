---
title: "Configuration Guide"
toc: true
toc_label: "Table of Contents"
toc_icon: "align-left"
---

# Configuration Guide

The Cleaner is highly configurable. You can control which directories and files are deleted by editing the `cleaner_config.json` file.

Here’s a breakdown of each configuration option:

- **`directories_to_remove`**: Specifies directory names that the app should delete if found within the target directory.
- **`file_extensions_to_remove`**: Lists file extensions to be deleted (e.g., `.log`, `.tmp`). Supports regex patterns if `matchRegex` is set to true.
- **`exclude_directories`**: Directories that should not be deleted, even if they match the `directories_to_remove` criteria.
- **`exclude_files`**: Files that should not be deleted, even if they match the `file_extensions_to_remove` criteria.
- **`matchRegex`**: When true, enables regex pattern matching for `file_extensions_to_remove`.

Example configuration file (`cleaner_config.json`):

```json
{
    "directories_to_remove": ["node_modules", "dist", "build"],
    "file_extensions_to_remove": [".log", ".tmp"],
    "exclude_directories": [".git", ".svn"],
    "exclude_files": ["README.md", "config.json"],
    "matchRegex": true
}

Explanation of Regex Use

If matchRegex is enabled, file_extensions_to_remove can contain regex patterns. For example:
 • .temp will match any file that starts with “.temp” just like using `%like%` in SQL.

By customizing this configuration, you can tailor the Cleaner to specific project needs.
