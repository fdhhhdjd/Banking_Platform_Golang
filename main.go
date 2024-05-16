package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
)

func main() {
	// Define the project root directory
	dir := "/home/fdhhhdjd/Documents/Code/Banking-Platform-Golang"

	// Create a new watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Walk the project directory and add each directory to the watcher
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return watcher.Add(path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// Channel to handle signals
	done := make(chan bool)

	// Start listening for events
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
					restartServer()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// Initial server start
	restartServer()

	<-done
}

var serverCmd *exec.Cmd

func restartServer() {
	if serverCmd != nil && serverCmd.Process != nil {
		log.Println("Stopping server...")
		err := serverCmd.Process.Signal(syscall.SIGTERM)
		if err != nil {
			log.Printf("Error stopping server: %v\n", err)
		}
		_, err = serverCmd.Process.Wait()
		if err != nil {
			log.Printf("Error waiting for server to stop: %v\n", err)
		}
	}

	log.Println("Starting server...")
	serverCmd = exec.Command("go", "run", "./cmd/main.go")
	serverCmd.Stdout = os.Stdout
	serverCmd.Stderr = os.Stderr
	err := serverCmd.Start()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// Give some time for the server to start
	time.Sleep(2 * time.Second)
}

// Start runs the HTTP server on a specific address.
