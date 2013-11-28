package main

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"hash"
	"os"
	"path/filepath"
)

type Inode struct {
	data   []byte
	dataGz []byte
	mime   string
	etag   string
}

type VFS map[string]*Inode

func (f *Inode) Size() int {
	return len(f.data)
}

func gz(data []byte) ([]byte, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	_, err := gz.Write(data)
	if err != nil {
		return nil, err
	}
	gz.Close()
	return b.Bytes(), nil
}

func hexdigest(data []byte, h hash.Hash) string {
	h.Reset()
	h.Write(data)
	b := h.Sum([]byte{})

	return fmt.Sprintf("%x", b)
}

func (vfs VFS) Add(filename string, inode *Inode) (err error) {
	if len(inode.dataGz) == 0 {
		if inode.dataGz, err = gz(inode.data); err != nil {
			return
		}

		// Don't compress if it takes more bandwidth
		if len(inode.data) < len(inode.dataGz) {
			inode.dataGz = []byte{}
		}
	}

	if inode.etag == "" {
		inode.etag = hexdigest(inode.data, sha256.New())
	}

	vfs[filename] = inode
	log.Printf("File %s added", filename)

	return
}

func (vfs VFS) Has(filename string) bool {
	return (vfs[filename] != nil)
}

func (vfs VFS) Remove(filename string) {
	delete(vfs, filename)
}

func (vfs VFS) LoadFile(filename string, srcPath string) (err error) {
	inode := &Inode{}
	if inode.data, err = ioutil.ReadFile(srcPath); err != nil {
		return
	}

	if inode.mime, err = Mime(srcPath); err != nil {
		return
	}
	if inode.mime == "" {
		log.Printf("No mime type found for %s", filename)
		inode.mime = "application/octet-stream"
	}

	err = vfs.Add(filename, inode)

	return
}

func (vfs VFS) LoadDir(prefix string) (err error) {
	err = filepath.Walk(prefix, func(path string, info os.FileInfo, errIn error) (err error) {
		var filename string

		if errIn != nil {
			log.Println("errIn", errIn)
			return errIn
		}
		if info.IsDir() {
			return
		}

		if filename, err = filepath.Rel(prefix, path); err != nil {
			log.Println("filename")
			return
		}

		err = vfs.LoadFile(filename, path)
		return
	})
	return
}
