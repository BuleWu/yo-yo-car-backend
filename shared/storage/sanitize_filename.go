package storage

import (
	"regexp"
	"strings"
)

func sanitizeFilename(filename string) string {
	// Remove any path components (just keep the base filename)
	parts := strings.Split(filename, "/")
	filename = parts[len(parts)-1]

	// Replace spaces with underscores
	filename = strings.ReplaceAll(filename, " ", "_")

	// Allow only alphanumerics, dash, underscore, and dot
	re := regexp.MustCompile(`[^a-zA-Z0-9._-]`)
	filename = re.ReplaceAllString(filename, "")

	return filename
}
