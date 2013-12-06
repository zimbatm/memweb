package main

import (
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

func Mime(path string) (mimeType string, err error) {
	mimeType = mime.TypeByExtension(filepath.Ext(path))
	if mimeType != "" {
		return
	}

	file, err := os.Open(path)
	if err != nil {
		return
	}

	// DetectContentType considers at most the first 512 bytes of data.
	data := make([]byte, 512)
	if _, err = file.Read(data); err != nil {
		return
	}

	return http.DetectContentType(data), nil
}
