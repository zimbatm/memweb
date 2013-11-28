package main

import (
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type VFSHandler struct {
	vfs   VFS
	index string
}

func (self *VFSHandler) findInode(path string) (inode *Inode, status int) {
	status = http.StatusOK

	if inode = self.vfs[path]; inode != nil {
		return
	}

	if inode = self.vfs[filepath.Join(path, self.index)]; inode != nil {
		return
	}

	status = http.StatusNotFound
	if inode = self.vfs["404.html"]; inode != nil {
		return
	}

	inode = &Inode{
		data: []byte("Page not found"),
		mime: "text/html",
	}

	return
}

func (self *VFSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")

	inode, status := self.findInode(path)

	w.Header().Set("Content-Type", inode.mime)
	w.Header().Add("Cache-Control", "must-revalidate")
	w.Header().Add("Vary", "Accept-Encoding")

	if inode.etag != "" {
		w.Header().Set("ETag", inode.etag)

		if r.Header.Get("If-None-Match") == inode.etag {
			w.WriteHeader(http.StatusNotModified)
			return
		}
	}

	var data []byte
	if len(inode.dataGz) > 0 && strings.Index(r.Header.Get("Accept-Encoding"), "gzip") >= 0 {
		w.Header().Set("Content-Encoding", "gzip")
		data = inode.dataGz
	} else {
		data = inode.data
	}

	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(status)
	w.Write(data)
}

func NewServer(vfs VFS, addr string) *http.Server {
	return &http.Server{
		Addr:           addr,
		Handler:        &VFSHandler{vfs, "index.html"},
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}
