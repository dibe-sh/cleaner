package main

import (
	"testing"
)

func TestShouldRemoveDir(t *testing.T) {
	dirsToRemove := []string{"node_modules", "dist", "build"}

	tests := []struct {
		dirName string
		want    bool
	}{
		{"node_modules", true},
		{"build", true},
		{"src", false},
	}

	for _, tt := range tests {
		got := shouldRemoveDir(tt.dirName, dirsToRemove)
		if got != tt.want {
			t.Errorf("shouldRemoveDir(%q) = %v; want %v", tt.dirName, got, tt.want)
		}
	}
}

func TestShouldRemoveFile(t *testing.T) {
	tests := []struct {
		fileName   string
		extensions []string
		matchRegex bool
		want       bool
	}{
		// Non-regex cases
		{"program.exe", []string{".exe", ".dll", ".tmp"}, false, true},
		{"tempfile.tmp", []string{".exe", ".dll", ".tmp"}, false, true},
		{"document.txt", []string{".exe", ".dll", ".tmp"}, false, false},

		// Regex cases
		{"tempfile.log", []string{`temp.*\.log`}, true, true},
		{"logfile.log", []string{`temp.*\.log`}, true, false},
		{"debug_info.txt", []string{`debug_.*\.txt`}, true, true},
		{"info.txt", []string{`debug_.*\.txt`}, true, false},
	}

	for _, tt := range tests {
		got := shouldRemoveFile(tt.fileName, tt.extensions, tt.matchRegex)
		if got != tt.want {
			t.Errorf("shouldRemoveFile(%q, %v, %v) = %v; want %v", tt.fileName, tt.extensions, tt.matchRegex, got, tt.want)
		}
	}
}

func TestIsExcluded(t *testing.T) {
	excludeList := []string{".git", ".svn", ".env"}

	tests := []struct {
		name string
		want bool
	}{
		{".git", true},
		{".env", true},
		{"README.md", false},
	}

	for _, tt := range tests {
		got := isExcluded(tt.name, excludeList)
		if got != tt.want {
			t.Errorf("isExcluded(%q) = %v; want %v", tt.name, got, tt.want)
		}
	}
}

func TestLoadConfig(t *testing.T) {
	config := getDefaultConfig()

	if len(config.DirectoriesToRemove) == 0 {
		t.Error("Expected default DirectoriesToRemove to be non-empty")
	}
	if len(config.FileExtensionsToRemove) == 0 {
		t.Error("Expected default FileExtensionsToRemove to be non-empty")
	}
}
