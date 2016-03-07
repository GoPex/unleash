package helpers

import (
	"strings"
)

// EvaluateURL evaluates github archive download url
func EvaluateURL(url string, branch string) string {
	return strings.Replace(strings.Replace(url, "{archive_format}", "tarball", 1), "{/ref}", "/"+branch, 1)
}
