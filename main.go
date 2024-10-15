package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Config struct {
	DirectoriesToRemove    []string `json:"directories_to_remove"`
	FileExtensionsToRemove []string `json:"file_extensions_to_remove"`
	ExcludeDirectories     []string `json:"exclude_directories"`
	ExcludeFiles           []string `json:"exclude_files"`
}

func main() {
	// Parse command-line arguments
	rootDir := flag.String("root", ".", "Root directory to scan")
	configFile := flag.String("config", "cleaner_config.json", "Path to the JSON configuration file")
	verbose := flag.Bool("verbose", true, "Enable verbose console output")
	saveLog := flag.Bool("save-log", true, "Enable saving logs to file")
	dryRun := flag.Bool("dry-run", false, "Simulate the cleaning process without deleting files")
	logFormat := flag.String("log-format", "text", "Log format: 'text' or 'json'")
	flag.Parse()

	// Check if the root directory exists
	if _, err := os.Stat(*rootDir); os.IsNotExist(err) {
		log.Fatalf("Root directory does not exist: %s\n", *rootDir)
	}

	// Load configuration
	config, err := loadConfig(*configFile)
	if err != nil {
		fmt.Printf("Warning: %v\n", err)
		fmt.Println("Using default configuration.")
		config = getDefaultConfig()
	}

	// Initialize logger
	var logger *log.Logger
	if *saveLog {
		logFile := "cleaned_source.txt"
		f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}
		defer f.Close()
		logger = log.New(f, "", 0) // Use custom log format
	}

	// Create a WaitGroup and a semaphore for concurrency control
	var wg sync.WaitGroup
	maxGoroutines := runtime.NumCPU()
	semaphore := make(chan struct{}, maxGoroutines)
	var logMutex sync.Mutex
	var consoleMutex sync.Mutex

	// Start processing the root directory
	wg.Add(1)
	go walkDir(*rootDir, config, *verbose, *dryRun, *logFormat, &wg, semaphore, logger, &logMutex, &consoleMutex)

	// Wait for all goroutines to finish
	wg.Wait()
	fmt.Println("Cleaning process completed.")
}

func getDefaultConfig() *Config {
	return &Config{
		DirectoriesToRemove: []string{
			"node_modules",
			"dist",
			"build",
			"bin",
			".next",
			".turbo",
			".idea",
			".cache",
		},
		FileExtensionsToRemove: []string{".DS_Store", "__debug_bin"},
		ExcludeDirectories:     []string{".git", ".svn"},
		ExcludeFiles:           []string{},
	}
}

func loadConfig(configPath string) (*Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("unable to open config file: %v", err)
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("unable to decode config JSON: %v", err)
	}

	return &config, nil
}

func walkDir(dir string, config *Config, verbose, dryRun bool, logFormat string, wg *sync.WaitGroup, semaphore chan struct{}, logger *log.Logger, logMutex, consoleMutex *sync.Mutex) {
	defer wg.Done()

	// Acquire a semaphore slot
	semaphore <- struct{}{}
	defer func() { <-semaphore }()

	// Read directory entries
	entries, err := os.ReadDir(dir)
	if err != nil {
		consoleMutex.Lock()
		fmt.Fprintf(os.Stderr, "Failed to read directory %s: %v\n", dir, err)
		consoleMutex.Unlock()
		return
	}

	for _, entry := range entries {
		entryName := entry.Name()
		entryPath := filepath.Join(dir, entryName)

		if entry.IsDir() {
			// Check for exclusion
			if isExcluded(entryName, config.ExcludeDirectories) {
				continue
			}

			// Check if the directory should be removed
			if shouldRemoveDir(entryName, config.DirectoriesToRemove) {
				// Log detection
				if verbose {
					consoleMutex.Lock()
					fmt.Printf("Detected directory to remove: %s\n", entryPath)
					consoleMutex.Unlock()
				}

				// Log starting cleaning
				if verbose {
					consoleMutex.Lock()
					fmt.Printf("Cleaning directory: %s\n", entryPath)
					consoleMutex.Unlock()
				}

				if !dryRun {
					// Remove the directory
					err := os.RemoveAll(entryPath)
					if err != nil {
						consoleMutex.Lock()
						fmt.Fprintf(os.Stderr, "Failed to remove directory %s: %v\n", entryPath, err)
						consoleMutex.Unlock()
					} else {
						// Log completion
						if verbose {
							consoleMutex.Lock()
							fmt.Printf("Completed cleaning directory: %s\n", entryPath)
							consoleMutex.Unlock()
						}

						// Log the path of the removed directory
						if logger != nil {
							logEntry(logger, logMutex, logFormat, "Removed directory", entryPath)
						}
					}
				} else {
					// Dry-run mode: simulate removal
					if verbose {
						consoleMutex.Lock()
						fmt.Printf("[Dry Run] Would remove directory: %s\n", entryPath)
						consoleMutex.Unlock()
					}
					if logger != nil {
						logEntry(logger, logMutex, logFormat, "[Dry Run] Would remove directory", entryPath)
					}
				}
			} else {
				// Recursively process subdirectories
				wg.Add(1)
				go walkDir(entryPath, config, verbose, dryRun, logFormat, wg, semaphore, logger, logMutex, consoleMutex)
			}
		} else {
			// Check for exclusion
			if isExcluded(entryName, config.ExcludeFiles) {
				continue
			}

			// Check if the file should be removed
			if shouldRemoveFile(entryName, config.FileExtensionsToRemove) {
				// Log detection
				if verbose {
					consoleMutex.Lock()
					fmt.Printf("Detected file to remove: %s\n", entryPath)
					consoleMutex.Unlock()
				}

				// Log starting cleaning
				if verbose {
					consoleMutex.Lock()
					fmt.Printf("Cleaning file: %s\n", entryPath)
					consoleMutex.Unlock()
				}

				if !dryRun {
					// Remove the file
					err := os.Remove(entryPath)
					if err != nil {
						consoleMutex.Lock()
						fmt.Fprintf(os.Stderr, "Failed to remove file %s: %v\n", entryPath, err)
						consoleMutex.Unlock()
					} else {
						// Log completion
						if verbose {
							consoleMutex.Lock()
							fmt.Printf("Completed cleaning file: %s\n", entryPath)
							consoleMutex.Unlock()
						}

						// Log the path of the removed file
						if logger != nil {
							logEntry(logger, logMutex, logFormat, "Removed file", entryPath)
						}
					}
				} else {
					// Dry-run mode: simulate removal
					if verbose {
						consoleMutex.Lock()
						fmt.Printf("[Dry Run] Would remove file: %s\n", entryPath)
						consoleMutex.Unlock()
					}
					if logger != nil {
						logEntry(logger, logMutex, logFormat, "[Dry Run] Would remove file", entryPath)
					}
				}
			}
		}
	}
}

func shouldRemoveDir(dirName string, dirsToRemove []string) bool {
	for _, name := range dirsToRemove {
		if dirName == name {
			return true
		}
	}
	return false
}

func shouldRemoveFile(fileName string, extensions []string) bool {
	for _, ext := range extensions {
		if strings.HasSuffix(fileName, ext) {
			return true
		}
	}
	return false
}

func isExcluded(name string, excludeList []string) bool {
	for _, exclude := range excludeList {
		if name == exclude {
			return true
		}
	}
	return false
}

func logEntry(logger *log.Logger, logMutex *sync.Mutex, format, action, path string) {
	logMutex.Lock()
	defer logMutex.Unlock()

	timestamp := time.Now().Format(time.RFC3339)
	if format == "json" {
		logMessage := fmt.Sprintf(`{"timestamp":"%s","action":"%s","path":"%s"}`, timestamp, action, path)
		logger.Println(logMessage)
	} else {
		// Default to text format
		logMessage := fmt.Sprintf("%s - %s: %s", timestamp, action, path)
		logger.Println(logMessage)
	}
}
