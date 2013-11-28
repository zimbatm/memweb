package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
)

var addr string
var prefix string
var watch bool

func init() {
	flag.StringVar(&addr, "addr", ":8484", "[host]:port to listen to")
	flag.StringVar(&prefix, "prefix", "", "source folder")
	flag.BoolVar(&watch, "watch", false, "refresh the cache when a file changes")
	flag.Parse()
}

func main() {
	var err error
	var vfs VFS
	var cwd string

	if cwd, err = os.Getwd(); err != nil {
		log.Fatal(err)
	}

	if prefix == "" {
		prefix = cwd
	}

	if prefix, err = filepath.Abs(prefix); err != nil {
		log.Fatal(err)
	}

	var rel string
	if rel, err = filepath.Rel(cwd, prefix); err != nil {
		log.Fatal(err)
	}

	log.Printf("Loading %s", rel)

	vfs = make(VFS)
	err = vfs.LoadDir(prefix)
	if err != nil {
		log.Fatal(err)
	}

	if watch {
		go WatchFS(vfs, prefix)
	}

	s := NewServer(vfs, addr)

	log.Printf("Starting on %s", addr)
	log.Fatal(s.ListenAndServe())
}
