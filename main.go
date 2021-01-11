package main

import (
	"errors"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

var wg sync.WaitGroup

type source struct {
	url         string
	destination string
	connections int64
}

func getSourceData() (source, error) {
	// Validate that we're getting the correct number of arguments
	if len(os.Args) < 2 {
		return source{}, errors.New("A source URL argument is required")
	}

	flag.Parse()             // This will parse all the arguments from the terminal
	sourceUrl := flag.Arg(0) // The only argument (that is not a flag option) is the file location (CSV file)

	// Defining option flags. For this, we're using the Flag package from the standard library
	// We need to define three arguments: the flag's name, the default value, and a short description (displayed whith the option --help)
	destination := flag.String("destination", "change.txt", "Column separator")
	connections := flag.Int64("connections", 50, "Generate pretty JSON")

	// If we get to this endpoint, our program arguments are validated
	// We return the corresponding struct instance with all the required data
	return source{sourceUrl, *destination, *connections}, nil
}

func main() {
	start := time.Now()
	targetUrl := "http://sample.li/FromTheAir.mp4"
	res, _ := http.Head(targetUrl) // 187 MB file of random numbers per line
	maps := res.Header
	length, _ := strconv.Atoi(maps["Content-Length"][0]) // Get the content length from the header request
	limit := 65                                          // 10 Go-routines for the process so each downloads 18.7MB
	len_sub := length / limit                            // Bytes for each Go-routine
	diff := length % limit                               // Get the remaining for the last request
	body := make([]string, limit+1)                      // Make up a temporary array to hold the data to be written to the file
	for i := 0; i < limit; i++ {
		wg.Add(1)

		min := len_sub * i       // Min range
		max := len_sub * (i + 1) // Max range

		if i == limit-1 {
			max += diff // Add the remaining bytes in the last request
		}

		go func(min int, max int, i int) {
			client := &http.Client{}
			req, _ := http.NewRequest("GET", targetUrl, nil)
			range_header := "bytes=" + strconv.Itoa(min) + "-" + strconv.Itoa(max-1) // Add the data for the Range header of the form "bytes=0-100"
			req.Header.Add("Range", range_header)
			resp, _ := client.Do(req)
			defer resp.Body.Close()
			reader, _ := ioutil.ReadAll(resp.Body)
			body[i] = string(reader)
			ioutil.WriteFile(strconv.Itoa(i)+".temp", []byte(string(body[i])), 0x777) // Write to the file i as a byte array
			wg.Done()
		}(min, max, i)
	}

	wg.Wait()

	// Combine files
	out, err := os.OpenFile("video.mp4", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln("failed to open output file:", err)
	}
	defer out.Close()

	for i := 0; i < limit; i++ {
		// Temporary file chunk open
		tempFile := strconv.Itoa(i) + ".temp"

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

	elapsed := time.Since(start)
	log.Printf("Download took %s", elapsed)
}
