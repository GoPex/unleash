package unleash

import (
    "strings"
)

func EvaluateUrl(url string, branch string) string {
    return strings.Replace(strings.Replace(url, "{archive_format}", "tarball", 1), "{/ref}", "/" + branch, 1)
}
