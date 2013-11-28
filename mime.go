package main

import (
	"os/exec"
	"strings"
	"fmt"
	"mime"
	"path/filepath"
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
