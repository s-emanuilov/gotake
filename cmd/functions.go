package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	wg       sync.WaitGroup
	filename string
)

func getCurrentDir() string {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return path
}

func downloadRanges(url string) {
	if verbose {
		fmt.Println("Starting downloading: " + url)
	}
	// Init the counter for the current time
	start := time.Now()

	// Get the current dir, where the file should be stored
	currentDir := getCurrentDir()

	// Perform a HEAD request to fetch the headers
	res, err := http.Head(url)
	if err != nil {
		fmt.Printf("Error fetching %s! Please try with a different method and not Range Requests.\n", url)
		return
	}

	// Extract the filename
	filename := filepath.Base(url)

	// Make the temp dir to hold the chunks
	tempDir := path.Join(currentDir, filename+"_"+randStringBytes(4)+"_tmp")

	// Create the temp directory to hold intermediate data
	err = os.Mkdir(tempDir, 0755)
	if err != nil {
		printMessage("You don't have permissions to download in this folder.", "error")
		return
	}

	// Check if there is an error with getting name for the file
	if len(filename) == 0 {
		filename = "gotake-download"
	}

	// Make a maps from the header
	maps := res.Header

	acceptRanges := maps["Accept-Ranges"][0]

	if strings.Contains(acceptRanges, "bytes") == false {
		log.Fatalf("%s doesn't support Range headers. Please try another method.", url)
	}

	// Get the content length from the header request
	contentLength, _ := strconv.Atoi(maps["Content-Length"][0])

	// Get the Content-Type of the resource
	contentType := maps["Content-Type"][0]

	// Check if the Content-Length is 0, so most likely we have a problem
	if contentLength == 0 {
		fmt.Printf("Error fetching %s! Content-Length is 0. Please try with a different method and not Range Requests.\n", url)
		return
	}

	// Number of  Go-routines for the process so each downloads
	// Bytes for each Go-routine
	lenSub := contentLength / connections

	// Get the remaining for the last request
	diff := contentLength % connections

	// Make up a temporary array to hold the data to be written to the file
	body := make([]string, connections+1)

	for i := 0; i < connections; i++ {
		// Add to a Waiting Group
		wg.Add(1)

		// Calculate min range
		min := lenSub * i

		// Calculate max range
		max := lenSub * (i + 1)

		// Add the remaining bytes in the last request
		if i == connections-1 {
			max += diff
		}

		// Fire the Goroutine
		go func(min int, max int, i int) {
			client := &http.Client{}
			req, _ := http.NewRequest("GET", url, nil)

			// Add the data for the Range header of the form "bytes=0-100"
			rangeHeader := "bytes=" + strconv.Itoa(min) + "-" + strconv.Itoa(max-1)
			// Add the Range request header
			req.Header.Add("Range", rangeHeader)

			// Get the body response
			resp, _ := client.Do(req)
			defer resp.Body.Close()
			reader, _ := ioutil.ReadAll(resp.Body)
			body[i] = string(reader)

			// Make the chunks as small files in a temp dir
			path := path.Join(tempDir, strconv.Itoa(i)+".temp")
			// Write to the file i as a byte array
			ioutil.WriteFile(path, []byte(string(body[i])), 0x777)
			wg.Done()
			if verbose {
				printMessage("chunk of "+filename+" with range "+rangeHeader, "download")
			}
		}(min, max, i)
	}

	wg.Wait()

	// Create the target file
	targetFile := path.Join(currentDir, filename)
	out, err := os.OpenFile(targetFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln("failed to open output file:", err)
	}

	// Close the file
	defer out.Close()

	for i := 0; i < connections; i++ {
		// Temporary file chunk open
		tempFile := path.Join(tempDir, strconv.Itoa(i)+".temp")

		// Open the current file and get the chunk
		chunk, err := os.Open(tempFile)
		if err != nil {
			log.Fatalln("failed to open zip for reading:", err)
		}

		// Copy the bytes from chunk into the result file
		if _, err := io.Copy(out, chunk); err != nil {
			log.Fatal(err)
		}

		// Close the chunk and remove from memory
		if err := chunk.Close(); err != nil {
			log.Fatal(err)
		}

		// Remove the temp file
		if err := os.Remove(tempFile); err != nil {
			log.Fatal(err)
		}
	}

	// Delete the temp dir
	defer os.Remove(tempDir)

	// Measure the time elapsed
	elapsed := time.Since(start)

	if summary {
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"URL", "SAVED IN", "MAX CONNECTIONS", "BYTES", "TYPE"})
		t.AppendRows([]table.Row{
			{url, targetFile, connections, contentLength, contentType},
		})
		t.AppendFooter(table.Row{"", "", "", "DOWNLOADED IN", elapsed})
		t.Render()
	} else {
		printMessage("Saved to "+path.Join(currentDir, filename), "success")
	}
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
		_, _ = color.New(color.FgGreen, color.Bold).Print("DOWNLOADED: ")
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
