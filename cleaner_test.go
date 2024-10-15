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
	extensionsToRemove := []string{".exe", ".dll", ".tmp"}

	tests := []struct {
		fileName string
		want     bool
	}{
		{"program.exe", true},
		{"tempfile.tmp", true},
		{"document.txt", false},
	}

	for _, tt := range tests {
		got := shouldRemoveFile(tt.fileName, extensionsToRemove)
		if got != tt.want {
			t.Errorf("shouldRemoveFile(%q) = %v; want %v", tt.fileName, got, tt.want)
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
	// Since we're not using files, we'll test getDefaultConfig instead
	config := getDefaultConfig()

	if len(config.DirectoriesToRemove) == 0 {
		t.Error("Expected default DirectoriesToRemove to be non-empty")
	}
	if len(config.FileExtensionsToRemove) == 0 {
		t.Error("Expected default FileExtensionsToRemove to be non-empty")
	}
}
