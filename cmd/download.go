package cmd

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"io"
	"io/ioutil"
	"log"
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

// Download from URL with range requests
func downloadRanges(url string) {
	// Validate connections
	validateConnections()

	if verbose {
		message := fmt.Sprintf("Starting downloading (ranges method): %s", url)
		printMessage(message, "info")
	}

	// Init the counter for the current time
	start := time.Now()

	// Get the current dir, where the file should be stored
	currentDir := getCurrentDir()

	// Perform a HEAD request to fetch the headers
	res, err := http.Head(url)
	if err != nil {
		message := fmt.Sprintf("Error fetching %s! Please try with a different method and not Range Requests.\n", url)
		printMessage(message, "error")
		return
	}

	// Check if the status code is different than 200 OK
	if res.StatusCode != 200 {
		message := fmt.Sprintf("Status code of %s is %s. Please check input data.", url, res.Status)
		printMessage(message, "error")
		return
	}

	// Extract the filename
	if filename == "" {
		filename = filepath.Base(url)
	}

	// Make the temp dir to hold the chunks
	tempFilename := fmt.Sprintf("%s_%s_tmp", filename, randStringBytes(4))
	tempDir := path.Join(currentDir, tempFilename)
	// Create the temp directory to hold intermediate data
	err = os.Mkdir(tempDir, 0755)
	if err != nil {
		printMessage("You don't have permissions to download in this folder.", "error")
		return
	}

	if verbose {
		message := fmt.Sprintf("Downloading in a temp dir %s", tempDir)
		printMessage(message, "info")
	}

	// Check if there is an error with getting name for the file
	if len(filename) == 0 {
		filename = defaultFilename
	}

	// Make a maps from the header
	maps := res.Header

	acceptRangesHeader := maps["Accept-Ranges"]
	// If we don't have a range request explicitly defined, we assume that server
	// doesnt support those
	if acceptRangesHeader == nil {
		// Remove the temp dir
		if err := os.Remove(tempDir); err != nil {
			message := fmt.Sprintf("%s cannot be removed.", tempDir)
			printMessage(message, "error")
		}

		// Return an error mesage
		message := fmt.Sprintf("%s doesn't support Range Headers. Please try another method. (-s flag)", url)
		printMessage(message, "error")
		return
	}

	// Extract the value of Range Requests header
	acceptRanges := acceptRangesHeader[0]
	if strings.Contains(acceptRanges, "bytes") == false {
		// Remove the temp dir
		if err := os.Remove(tempDir); err != nil {
			message := fmt.Sprintf("%s cannot be removed.", tempDir)
			printMessage(message, "error")
		}

		// Return an error mesage
		message := fmt.Sprintf("%s doesn't support Range Headers. Please try another method. (-s flag)", url)
		printMessage(message, "error")
		return
	}

	// Get the content length from the header request
	contentLength := 0
	contentLengthHeader := maps["Content-Length"]

	if contentLengthHeader != nil {
		contentLength, _ = strconv.Atoi(contentLengthHeader[0])
	} else {
		// Remove the temp dir
		if err := os.Remove(tempDir); err != nil {
			message := fmt.Sprintf("%s cannot be removed.", tempDir)
			printMessage(message, "error")
		}

		// Return an error mesage
		message := fmt.Sprintf("%s Content-Length headers is 0. Most likely origin server doesn't support Range requests.", url)
		printMessage(message, "error")
		return
	}

	// Get the Content-Type of the resource
	contentType := defaultContentType
	contentTypeHeader := maps["Content-Type"]

	if contentTypeHeader != nil {
		contentType = maps["Content-Type"][0]
	}

	// Calculate optimal connections number
	if auto {
		autoConnections := contentLength / chunkSize
		// Make sure we will add good amount of connections in small files
		if autoConnections < 1 {
			autoConnections = 1
		}
		connections = autoConnections

		if verbose {
			message := fmt.Sprintf("Auto connections set to: %d", connections)
			printMessage(message, "info")
		}
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
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {

				}
			}(resp.Body)
			reader, _ := ioutil.ReadAll(resp.Body)
			body[i] = string(reader)

			// Make the chunks as small files in a temp dir
			tempPath := path.Join(tempDir, strconv.Itoa(i)+".temp")
			// Write to the file i as a byte array
			err := ioutil.WriteFile(tempPath, []byte(string(body[i])), 0x777)
			if err != nil {
				return
			}
			wg.Done()
			if verbose {
				message := fmt.Sprintf("chunk of %s with range %s", filename, rangeHeader)
				printMessage(message, "download")
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
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {

		}
	}(out)

	if verbose {
		printMessage("Started to combine chunks into the result file", "info")
	}

	// Combining chunks into the end file
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
			message := fmt.Sprintf("%s file cannot be removed.", tempFile)
			printMessage(message, "error")
		}
	}

	// Delete the temp dir
	if err := os.Remove(tempDir); err != nil {
		message := fmt.Sprintf("%s cannot be removed.", tempDir)
		printMessage(message, "error")
	}

	if verbose {
		printMessage("Temp dir removed", "info")
	}

	// Measure the time elapsed
	elapsed := time.Since(start)

	// Check if the Content-Length from header matches the result file
	fi, err := os.Stat(targetFile)
	if err != nil || fi == nil {
		message := fmt.Sprintf("Cannot check the size of downloaded file")
		printMessage(message, "error")
	}

	// get the size
	size := fi.Size()

	if size != int64(contentLength) {
		message := fmt.Sprintf("Missmatch in Content-Length and downloaded file size.")
		printMessage(message, "error")
		return
	}

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
		message := fmt.Sprintf("Saved to %s in %s", targetFile, elapsed)
		printMessage(message, "success")
	}
}

func downloadStandard(url string) {
	if verbose {
		message := fmt.Sprintf("Starting downloading (standart method): %s", url)
		printMessage(message, "info")
	}

	// Init the counter for the current time
	start := time.Now()

	// Get the current dir, where the file should be stored
	currentDir := getCurrentDir()

	// Extract the filename
	if filename == "" {
		filename = filepath.Base(url)
	}

	targetFile := path.Join(currentDir, filename)
	out, err := os.Create(targetFile)

	if err != nil {
		message := fmt.Sprintf("Error creating file %s", targetFile)
		printMessage(message, "error")
		return
	}

	defer func(out *os.File) {
		err := out.Close()
		if err != nil {

		}
	}(out)

	res, err := http.Get(url)
	if err != nil {
		message := fmt.Sprintf("Error fetching %s!", url)
		printMessage(message, "error")
		return
	}

	// Check if the status code is different than 200 OK
	if res.StatusCode != 200 {
		message := fmt.Sprintf("Status code of %s is %s. Please check input data.", url, res.Status)
		printMessage(message, "error")
		return
	}

	// Make a maps from the header
	maps := res.Header

	// Get the content length from the header request
	contentLength := 0
	contentLengthHeader := maps["Content-Length"]

	if contentLengthHeader != nil {
		contentLength, _ = strconv.Atoi(contentLengthHeader[0])
	}

	// Get the Content-Type of the resource
	contentType := defaultContentType
	contentTypeHeader := maps["Content-Type"]

	if contentTypeHeader != nil {
		contentType = maps["Content-Type"][0]
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)

	if _, err := io.Copy(out, res.Body); err != nil {
		message := fmt.Sprintf("Error downloading body in file %s", targetFile)
		printMessage(message, "error")
		return
	}

	if verbose {
		message := fmt.Sprintf("chunk of %s with no ranges (standard method)", filename)
		printMessage(message, "download")
	}

	// Measure the time elapsed
	elapsed := time.Since(start)

	if summary {
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"URL", "SAVED IN", "MAX CONNECTIONS", "BYTES", "TYPE"})
		t.AppendRows([]table.Row{
			{url, targetFile, "1 (NO RANGES)", contentLength, contentType},
		})
		t.AppendFooter(table.Row{"", "", "", "DOWNLOADED IN", elapsed})
		t.Render()
	} else {
		message := fmt.Sprintf("Saved to %s in %s", targetFile, elapsed)
		printMessage(message, "success")
	}
}
