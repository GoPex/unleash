package unleash

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// ReadUrl create a *tar.Reader from an url
func ExtractRepository(url string, destination string) error {
	// Get the tar file from url
	check := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			log.Error(url)
			log.Error("REDIRECT")
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	res, err := check.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Check the return code of our http request
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Not able to GET %s, status code is %s !", url, res.StatusCode)
	}

	// Copy the body to a bytes buffer that we'll use to read our tar from
	buf := &bytes.Buffer{}
	_, err = io.Copy(buf, res.Body)
	if err != nil {
		return err
	}

	// Create a gzip..Reader from our bytes.Buffer
	gr, err := gzip.NewReader(buf)
	if err != nil {
		return err
	}
	defer gr.Close()

	// Create a tar.Reader from our gzip.Reader
	tr := tar.NewReader(gr)

	// Extract files from the tar
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			// end of tar archive
			break
		}
		if err != nil {
			return err
		}
		splittedFileName := strings.Split(hdr.Name, string(os.PathSeparator))
		removedRootDirectory := strings.Join(splittedFileName[1:], string(os.PathSeparator))
		path := filepath.Join(destination, removedRootDirectory)
		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(path, hdr.FileInfo().Mode()); err != nil {
				return err
			}
		case tar.TypeReg:
			ow, err := os.Create(path)
			if err != nil {
				return err
			}
			defer ow.Close()

			if _, err := io.Copy(ow, tr); err != nil {
				return err
			}
		}
	}

	return nil
}
