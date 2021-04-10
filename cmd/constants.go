package cmd

import "fmt"

// Constants and variables related to main funtions
const version = 0.2

var longDescription = fmt.Sprintf(`GoTake provide fast, easy and reliable fast downloads.

Creator: Simeon Emanuilov
Link: https://github.com/simeonemanuilov/gotake

Version: %.1f`, version)

// List with chars for generate random names, prefixes and suffixes
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// Min and Max connections to restrict the range for using
const minConnections = 2
const maxConnections = 500

// Default name of the file
const defaultFilename = "gotake-file"
