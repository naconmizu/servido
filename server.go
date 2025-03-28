package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"tempo/crates"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Server initializes the Gin server with improved implementation
func Server() {

	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file:", err)
	}
	// Set up server logging
	logFile, err := os.OpenFile("server.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	defer logFile.Close()

	// Configure custom logger
	logger := log.New(io.MultiWriter(os.Stdout, logFile), "[TEMPO] ", log.LstdFlags)

	// Set Gin to release mode for production
	gin.SetMode(gin.ReleaseMode)

	// Create a new Gin router with logging and recovery middleware
	c := gin.New()
	c.Use(gin.Recovery())
	c.Use(gin.LoggerWithWriter(io.MultiWriter(os.Stdout, logFile)))

	// Initialize file queue
	arquivos := NewQueue()

	// Define upload directory
	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		logger.Fatal("Failed to create upload directory:", err)
	}

	// Create server context that can be used for graceful shutdown
	serverCtx, _ := context.WithCancel(context.Background())

	// Admin authentication middleware
	adminAuth := func(ctx *gin.Context) {
		adminPassword := os.Getenv("ADMPASSWORD")
		if adminPassword == "" {
			logger.Println("Warning: ADMPASSWORD environment variable is not set")
			adminPassword = "iQuietDownIfItsWhatYouWant" // Fallback for development
		}

		if providedPassword := ctx.Param("admin"); providedPassword != adminPassword {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}

	// GET route to list server files (admin only)
	c.GET("/listserver/:admin", adminAuth, func(ctx *gin.Context) {
		files := arquivos.ToSlice()
		ctx.String(http.StatusOK, strings.Join(files, "\n"))
		logger.Printf("Admin listed %d server files", len(files))
	})

	c.GET("/ipserver/:admin", adminAuth, func(ctx *gin.Context) {
		ip, err := GetLocalIP()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get IP"})
			return
		}
		ctx.String(http.StatusOK, ip)
	})
		
		



	// GET route to list database files (admin only)
	c.GET("/listdatabase/:admin", adminAuth, func(ctx *gin.Context) {
		files := crates.Grider("list")
		ctx.String(http.StatusOK, strings.Join(files, "\n"))
		logger.Printf("Admin listed %d database files", len(files))
		arquivos.PrintQueue()
	})

	// POST route for simple file upload (stored only on server)
	c.POST("/", func(ctx *gin.Context) {
		file, header, err := ctx.Request.FormFile("file")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to receive file"})
			return
		}
		defer file.Close()

		filename := filepath.Base(header.Filename) // Sanitize filename
		filePath := filepath.Join(uploadDir, filename)

		out, err := os.Create(filePath)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file on server"})
			return
		}
		defer out.Close()

		written, err := io.Copy(out, file)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}

		arquivos.Enqueue(filename)
		logger.Printf("File uploaded: %s (%d bytes)", filename, written)

		arquivos.PrintQueue()
		ctx.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("File %s uploaded successfully!", filename),
			"size":    written,
		})
	})

	// POST route for file upload to both server and database
	c.POST("/up", func(ctx *gin.Context) {
		file, header, err := ctx.Request.FormFile("file")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to receive file"})
			return
		}
		defer file.Close()

		filename := filepath.Base(header.Filename) // Sanitize filename
		filePath := filepath.Join(uploadDir, filename)

		out, err := os.Create(filePath)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file on server"})
			return
		}
		defer out.Close()

		written, err := io.Copy(out, file)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}

		// Add to queue and database
		arquivos.Enqueue(filename)
		result := crates.Grider("up", filename)

		logger.Printf("File uploaded to server and database: %s (%d bytes)", filename, written)
		arquivos.PrintQueue()

		ctx.JSON(http.StatusOK, gin.H{
			"message":         "File uploaded successfully to server and database",
			"size":            written,
			"database_result": result,
		})
	})

	// DELETE route for file removal
	c.DELETE("/:filename", func(ctx *gin.Context) {
		filename := filepath.Base(ctx.Param("filename")) // Sanitize filename
		filePath := filepath.Join(uploadDir, filename)

		files := arquivos.ToSlice()
		fileExists := slices.Contains(files, filename)

		if fileExists {
			if err := os.Remove(filePath); err != nil {
				if os.IsNotExist(err) {
					logger.Printf("File %s not found on filesystem", filename)
				} else {
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove file"})
					logger.Printf("Error removing file %s: %v", filename, err)
					return
				}
			}
		}

		// Remove from queue regardless if file exists (clean up queue)
		newQueue := NewQueue()
		for _, file := range files {
			if file != filename {
				newQueue.Enqueue(file)
			}
		}
		arquivos = newQueue

		// Remove from database
		_ = crates.Grider("delete", filename)

		logger.Printf("File removed: %s (existed: %v)", filename, fileExists)
		ctx.String(http.StatusOK, "File successfully removed from server and database")
	})

	// GET route to download file from database
	c.GET("/down/:filename", func(ctx *gin.Context) {
		filename := filepath.Base(ctx.Param("filename")) // Sanitize filename

		// Check if file is in database
		dbFiles := crates.Grider("list")
		if !slices.Contains(dbFiles, filename) {
			logger.Printf("File not found in database: %s", filename)
			ctx.JSON(http.StatusNotFound, gin.H{"error": "File not found in database"})
			return
		}

		_ = crates.Grider("down", filename)

		outputPath := "./uploads/" + filename

		if _, err := os.Stat(outputPath); os.IsNotExist(err) {
			logger.Printf("File check error: %v", err)
			currentDir, _ := os.Getwd()
			logger.Printf("Current directory: %s", currentDir)
			fileInfo, _ := os.Stat(outputPath)
			// fileInfo, _ := os.Stat(filename)
			var permissions string
			if fileInfo != nil {
				permissions = fileInfo.Mode().String()
			} else {
				permissions = "unknown"
			}
			logger.Printf("File permissions: %s", permissions)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve file from database"})
			return
		}

		// Send file to client
		ctx.File(outputPath)
		logger.Printf("File downloaded from database: %s", filename)

		// Clean up file after sending
		go func() {
			// Give time for the file to be sent
			time.Sleep(20 * time.Second)
			filepath := "./uploads/" + filename
			if err := os.Remove(filepath); err != nil {
				logger.Printf("Error removing temporary file %s: %v", filename, err)
			} else {
				logger.Printf("Removed temporary file: %s", filename)
			}
		}()
	})

	// Health check endpoint
	c.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	})

	// Start auto-cleanup goroutine with context for graceful shutdown
	go func() {
		ticker := time.NewTicker(30 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-serverCtx.Done():
				logger.Println("Auto-cleanup routine shutting down")
				return

			case <-ticker.C:

				arquivos = FromSliceToQueue(crates.Grider("list"))

				if arquivos.IsEmpty() {
					logger.Println("Auto-cleanup: Queue is empty")
					continue
				}

				// Dequeue and process file
				filename, ok := arquivos.Dequeue()
				if !ok {
					logger.Println("Auto-cleanup: Failed to dequeue file")
					continue
				}

				logger.Printf("Auto-cleanup: Processing file %s", filename)

				// Create DELETE request
				req, err := http.NewRequest(http.MethodDelete,
					fmt.Sprintf("http://localhost:8080/%s", filename), nil)
				if err != nil {
					logger.Printf("Auto-cleanup: Error creating DELETE request: %v", err)
					// Re-queue the file if we couldn't process it
					arquivos.Enqueue(filename)
					continue
				}

				// Send DELETE request
				client := &http.Client{Timeout: 10 * time.Second}
				resp, err := client.Do(req)
				if err != nil {
					logger.Printf("Auto-cleanup: Error making DELETE request: %v", err)
					// Re-queue the file if we couldn't process it
					arquivos.Enqueue(filename)
					continue
				}

				// Read and close response body
				body, _ := io.ReadAll(resp.Body)
				resp.Body.Close()

				// Clean up server files
				if err := ClearServer(); err != nil {
					logger.Printf("Auto-cleanup: Error in ClearServer: %v", err)
				}

				logger.Printf("Auto-cleanup: File %s removed, response: %s",
					filename, strings.TrimSpace(string(body)))
			}
		}
	}()

	// Start server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: c,
	}

	// Run server in goroutine so we can handle shutdown gracefully
	go func() {
		logger.Println("Server starting on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Server failed: %v", err)
		}
	}()

	// Listen for interrupt signal
	// Note: In a real implementation, you'd use a signal channel here
	// This is just a placeholder for the concept
	select {
	case <-serverCtx.Done():
		logger.Println("Shutting down server...")

		// Create shutdown context with timeout
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Start separate goroutine to warn about shutdown taking too long
		go func() {
			select {
			case <-time.After(7 * time.Second):
				logger.Println("Warning: Server shutdown is taking longer than expected")
			case <-shutdownCtx.Done():
				// Context finished normally, no need to warn
				return
			}
		}()

		// Close any open resources before shutting down HTTP server
		// e.g. database, cache connections, etc.
		// db.Close()

		// Attempt graceful shutdown
		if err := srv.Shutdown(shutdownCtx); err != nil {
			logger.Printf("Server shutdown error: %v", err)

			// If graceful shutdown fails, force close
			if err := srv.Close(); err != nil {
				logger.Fatalf("Server forced close failed: %v", err)
			}
			logger.Println("Server closed forcefully")
		} else {
			logger.Println("Server shut down gracefully")
		}

		// Notify any monitoring systems that the server is down
		// notifySystemMonitor("server-down")
	}
}
