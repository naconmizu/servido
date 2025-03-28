package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"tempo/crates"
)




func ClearServer() error {
	// Get database files
	dbFiles := crates.Grider("list")
	
	// Convert slice to map for O(1) lookups
	dbFileMap := make(map[string]struct{}, len(dbFiles))
	for _, file := range dbFiles {
		dbFileMap[file] = struct{}{}
	}
	
	// Read directory contents
	dirFiles, err := os.ReadDir("./")
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}
	
	// Process each file
	removedCount := 0
	for _, file := range dirFiles {
		if file.IsDir() {
			continue // Skip directories
		}
		
		filename := file.Name()
		// Check if file exists in database files using map lookup
		if _, exists := dbFileMap[filename]; exists {
			// Construct full path
			fullPath := filepath.Join("./", filename)
			
			// Remove file and handle potential errors
			if err := os.Remove(fullPath); err != nil {
				log.Printf("Warning: could not remove file %s: %v", filename, err)
			} else {
				removedCount++
			}
		}
	}
	
	log.Printf("ClearServer: removed %d files", removedCount)
	return nil
}