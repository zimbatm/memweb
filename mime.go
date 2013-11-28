package main

import (
	"fmt"
	"mime"
	"os/exec"
	"path/filepath"
	"strings"
)

func Mime(path string) (mimeType string, err error) {
	mimeType = mime.TypeByExtension(filepath.Ext(path))
	if mimeType != "" {
		return
	}

	out, err := exec.Command("file", "--mime", path).Output()
	if err != nil {
		return
	}

	out2 := strings.SplitN(
		strings.TrimSpace(string(out)),
		":",
		2,
	)

	if path != out2[0] {
		err = fmt.Errorf("%s != %s", path, out2[0])
		return
	}

	mimeType = strings.TrimSpace(out2[1])

	return
}
