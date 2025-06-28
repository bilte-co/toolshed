package cli

import (
	"fmt"
	"log/slog"
	"math/rand"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ServeCmd represents the serve command
type ServeCmd struct {
	Port int    `short:"p" help:"Port to listen on (default: random available port)"`
	Dir  string `short:"d" help:"Directory to serve (default: current directory)"`
}

func (cmd *ServeCmd) Run(ctx *CLIContext) error {
	// Set defaults
	if cmd.Dir == "" {
		var err error
		cmd.Dir, err = os.Getwd()
		if err != nil {
			ctx.Logger.Error("Failed to get current directory", "error", err)
			return fmt.Errorf("failed to get current directory: %w", err)
		}
	}

	// Clean and validate directory path
	cmd.Dir = filepath.Clean(cmd.Dir)

	// Check if directory exists and is accessible
	info, err := os.Stat(cmd.Dir)
	if err != nil {
		ctx.Logger.Error("Directory not accessible", "dir", cmd.Dir, "error", err)
		return fmt.Errorf("directory not accessible: %w", err)
	}
	if !info.IsDir() {
		ctx.Logger.Error("Path is not a directory", "path", cmd.Dir)
		return fmt.Errorf("path is not a directory: %s", cmd.Dir)
	}

	// Get an available port
	port, err := cmd.getPort()
	if err != nil {
		ctx.Logger.Error("Failed to get available port", "error", err)
		return fmt.Errorf("failed to get available port: %w", err)
	}

	// Create a custom file server with security
	fs := &secureFileSystem{http.Dir(cmd.Dir)}
	handler := &loggingHandler{
		handler: http.FileServer(fs),
		logger:  ctx.Logger,
	}

	// Setup server
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	server := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Log startup information
	absDir, _ := filepath.Abs(cmd.Dir)
	url := fmt.Sprintf("http://127.0.0.1:%d", port)
	
	ctx.Logger.Info("Starting HTTP server",
		"port", port,
		"directory", absDir,
		"url", url,
	)

	fmt.Printf("Serving %s on port %d\n", absDir, port)
	fmt.Printf("Server running at %s\n", url)
	fmt.Println("Press Ctrl+C to stop")

	// Start server
	return server.ListenAndServe()
}

// getPort returns the specified port or finds an available random port between 4000-8999
func (cmd *ServeCmd) getPort() (int, error) {
	if cmd.Port != 0 {
		// Check if specified port is available
		if err := cmd.checkPortAvailable(cmd.Port); err != nil {
			return 0, fmt.Errorf("specified port %d is not available: %w", cmd.Port, err)
		}
		return cmd.Port, nil
	}

	// Find a random available port in range 4000-8999
	const minPort = 4000
	const maxPort = 8999
	const maxAttempts = 100

	for attempt := 0; attempt < maxAttempts; attempt++ {
		// Generate random port in range
		port := minPort + rand.Intn(maxPort-minPort+1)
		
		if err := cmd.checkPortAvailable(port); err == nil {
			return port, nil
		}
	}

	return 0, fmt.Errorf("failed to find available port in range %d-%d after %d attempts", minPort, maxPort, maxAttempts)
}

// checkPortAvailable checks if a port is available for binding
func (cmd *ServeCmd) checkPortAvailable(port int) error {
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	listener.Close()
	return nil
}

// secureFileSystem wraps http.Dir to prevent directory traversal
type secureFileSystem struct {
	fs http.FileSystem
}

func (sfs *secureFileSystem) Open(name string) (http.File, error) {
	// Clean the path to prevent directory traversal
	name = filepath.Clean(name)
	
	// Ensure path doesn't contain .. or other traversal attempts
	if strings.Contains(name, "..") {
		return nil, os.ErrNotExist
	}

	return sfs.fs.Open(name)
}

// loggingHandler wraps an http.Handler to log requests
type loggingHandler struct {
	handler http.Handler
	logger  *slog.Logger
}

func (lh *loggingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Create a response recorder to capture status code
	recorder := &responseRecorder{
		ResponseWriter: w,
		statusCode:     200, // Default to 200
	}

	start := time.Now()
	
	// Serve the request
	lh.handler.ServeHTTP(recorder, r)
	
	duration := time.Since(start)

	// Log the request
	lh.logger.Info("HTTP request",
		"method", r.Method,
		"path", r.URL.Path,
		"status", recorder.statusCode,
		"duration", duration.String(),
		"remote_addr", r.RemoteAddr,
	)
}

// responseRecorder captures the status code from the response
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}

func (rr *responseRecorder) Write(b []byte) (int, error) {
	return rr.ResponseWriter.Write(b)
}
