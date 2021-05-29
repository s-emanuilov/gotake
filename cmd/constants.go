package cmd

import "fmt"

// Constants and variables related to main functions
const version = 0.5

var longDescription = fmt.Sprintf(`GoTake provide fast, easy and reliable fast downloads.

Creator: Simeon Emanuilov
Link: https://github.com/simeonemanuilov/gotake

Version: %.1f`, version)

// List with chars for generate random names, prefixes and suffixes
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// Min and Max connections to restrict the range for using
const minConnections = 2
const maxConnections = 250

// Default name of the file
const defaultFilename = "gotake-file"

// Default Content Type of the file
const defaultContentType = "text/plain"

// The amount in megabytes for each connections in auto mode
const chunkSize = 1 * 1024 * 1024
