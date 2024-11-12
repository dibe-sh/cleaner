package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
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
	MatchRegex             bool     `json:"matchRegex"`
}

func main() {
	rootDir := flag.String("root", ".", "Root directory to scan")
	configFile := flag.String("config", "cleaner_config.json", "Path to the JSON configuration file")
	verbose := flag.Bool("verbose", true, "Enable verbose console output")
	saveLog := flag.Bool("save-log", true, "Enable saving logs to file")
	dryRun := flag.Bool("dry-run", false, "Simulate the cleaning process without deleting files")
	logFormat := flag.String("log-format", "text", "Log format: 'text' or 'json'")
	flag.Parse()

	if _, err := os.Stat(*rootDir); os.IsNotExist(err) {
		log.Fatalf("Root directory does not exist: %s\n", *rootDir)
	}

	config, err := loadConfig(*configFile)
	if err != nil {
		fmt.Printf("Warning: %v\n", err)
		fmt.Println("Using default configuration.")
		config = getDefaultConfig()
	}

	var logger *log.Logger
	if *saveLog {
		logFile := "cleaned_source.txt"
		f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}
		defer f.Close()
		logger = log.New(f, "", 0)
	}

	var wg sync.WaitGroup
	maxGoroutines := runtime.NumCPU()
	semaphore := make(chan struct{}, maxGoroutines)
	var logMutex sync.Mutex
	var consoleMutex sync.Mutex

	var totalSize int64
	var totalSizeMutex sync.Mutex
	var paths []string

	wg.Add(1)
	go walkDir(*rootDir, config, *verbose, *dryRun, *logFormat, &wg, semaphore, logger, &logMutex, &consoleMutex, &totalSize, &totalSizeMutex, &paths)

	wg.Wait()

	if *dryRun {
		fmt.Printf("Total cleanable size: %.2f MB\n", float64(totalSize)/(1024*1024))
		fmt.Println("Paths that would be removed:", paths)
	} else {
		fmt.Println("Cleaning process completed.")
	}
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
		MatchRegex:             true,
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

func walkDir(
	dir string,
	config *Config,
	verbose, dryRun bool,
	logFormat string,
	wg *sync.WaitGroup,
	semaphore chan struct{},
	logger *log.Logger,
	logMutex, consoleMutex *sync.Mutex,
	totalSize *int64,
	totalSizeMutex *sync.Mutex,
	paths *[]string, // New parameter to collect paths during dry-run
) {
	defer wg.Done()

	semaphore <- struct{}{}
	defer func() { <-semaphore }()

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
			if isExcluded(entryName, config.ExcludeDirectories) {
				continue
			}

			if shouldRemoveDir(entryName, config.DirectoriesToRemove) {
				dirSize := calculateSize(entryPath)
				totalSizeMutex.Lock()
				*totalSize += dirSize
				totalSizeMutex.Unlock()

				if dryRun {
					if verbose {
						consoleMutex.Lock()
						fmt.Printf("[Dry Run] Would remove directory: %s (Size: %.2f MB)\n", entryPath, float64(dirSize)/(1024*1024))
						consoleMutex.Unlock()
					}
					*paths = append(*paths, entryPath)
					if logger != nil {
						logEntry(logger, logMutex, logFormat, "[Dry Run] Would remove directory", entryPath)
					}
				} else {
					err := os.RemoveAll(entryPath)
					if err != nil {
						consoleMutex.Lock()
						fmt.Fprintf(os.Stderr, "Failed to remove directory %s: %v\n", entryPath, err)
						consoleMutex.Unlock()
					} else if logger != nil {
						logEntry(logger, logMutex, logFormat, "Removed directory", entryPath)
					}
				}
			} else {
				wg.Add(1)
				go walkDir(entryPath, config, verbose, dryRun, logFormat, wg, semaphore, logger, logMutex, consoleMutex, totalSize, totalSizeMutex, paths)
			}
		} else {
			if isExcluded(entryName, config.ExcludeFiles) {
				continue
			}

			if shouldRemoveFile(entryName, config.FileExtensionsToRemove, config.MatchRegex) {
				fileSize := calculateSize(entryPath)
				totalSizeMutex.Lock()
				*totalSize += fileSize
				totalSizeMutex.Unlock()

				if dryRun {
					if verbose {
						consoleMutex.Lock()
						fmt.Printf("[Dry Run] Would remove file: %s (Size: %.2f KB)\n", entryPath, float64(fileSize)/1024)
						consoleMutex.Unlock()
					}
					*paths = append(*paths, entryPath)
					if logger != nil {
						logEntry(logger, logMutex, logFormat, "[Dry Run] Would remove file", entryPath)
					}
				} else {
					err := os.Remove(entryPath)
					if err != nil {
						consoleMutex.Lock()
						fmt.Fprintf(os.Stderr, "Failed to remove file %s: %v\n", entryPath, err)
						consoleMutex.Unlock()
					} else if logger != nil {
						logEntry(logger, logMutex, logFormat, "Removed file", entryPath)
					}
				}
			}
		}
	}
}

func calculateSize(path string) int64 {
	var size int64
	filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err == nil && info != nil {
			size += info.Size()
		}
		return nil
	})
	return size
}

func shouldRemoveDir(dirName string, dirsToRemove []string) bool {
	for _, name := range dirsToRemove {
		if dirName == name {
			return true
		}
	}
	return false
}

func shouldRemoveFile(fileName string, extensions []string, matchRegex bool) bool {
	if matchRegex {
		for _, pattern := range extensions {
			matched, err := regexp.MatchString(pattern, fileName)
			if err != nil {
				fmt.Printf("Error in regex pattern %s: %v\n", pattern, err)
				continue
			}
			if matched {
				return true
			}
		}
	} else {
		for _, ext := range extensions {
			if strings.HasSuffix(fileName, ext) {
				return true
			}
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
		logMessage := fmt.Sprintf("%s - %s: %s", timestamp, action, path)
		logger.Println(logMessage)
	}
}
