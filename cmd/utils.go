package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"log"
	"math/rand"
	"os"
)

// Get the current directory where the download should happen.
func getCurrentDir() string {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return path
}

func printMessage(message string, kind string) {
	switch kind := kind; kind {
	case "error":
		_, _ = color.New(color.FgRed, color.Bold).Print("ERROR: ")
	case "warning":
		_, _ = color.New(color.FgYellow, color.Bold).Print("WARNING: ")
	case "success":
		_, _ = color.New(color.FgGreen, color.Bold).Print("SUCCESS: ")
	case "download":
		_, _ = color.New(color.FgGreen).Print("DOWNLOADED: ")
	case "info":
		_, _ = color.New(color.FgCyan, color.Bold).Print("INFO: ")
	default:
	}
	fmt.Println(message)
}

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
