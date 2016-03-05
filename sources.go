package unleash

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// ExtractRepository from an url hosting a tar.gz of this repository with an
// added root directory that will be flatten
func ExtractRepository(url string, destination string) error {
	// Get the tar reader from the URL
	tr, err := ReadURL(url)
	if err != nil {
		return err
	}

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

// ReadURL reads a tar.gz archive from an URL and return the associated reader
func ReadURL(url string) (*tar.Reader, error) {
	// Get the tar file from url
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.SetBasicAuth(Config.GitUsername, Config.GitPassword)

	res, err := (&http.Client{}).Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Check the return code of our http request
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Not able to GET %s, status code is %d !", url, res.StatusCode)
	}

	// Copy the body to a bytes buffer that we'll use to read our tar from
	buf := &bytes.Buffer{}
	_, err = io.Copy(buf, res.Body)
	if err != nil {
		return nil, err
	}

	// Create a gzip..Reader from our bytes.Buffer
	gr, err := gzip.NewReader(buf)
	if err != nil {
		return nil, err
	}
	defer gr.Close()

	// Return a tar.Reader from our gzip.Reader
	return tar.NewReader(gr), nil
}
