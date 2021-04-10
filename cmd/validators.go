package cmd

import "strings"

// Validate if string looks like a normal URL
func validateUrl(url string) bool {
	if strings.HasPrefix(url, "http") && strings.Contains(url, "://") {
		return true
	}
	return false
}

func validateConnections() {
	// Sanitize the flag inputs
	if connections <= minConnections {
		connections = minConnections
	} else if connections > maxConnections {
		connections = maxConnections
	}
}
