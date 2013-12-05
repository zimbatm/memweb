package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
)

var addr string
var cors bool
var ignores []string
var index string
var prefix string
var watch bool

func init() {
	//flag.BoolVar(&cors, "cors", false, "allow CORS requests")
	flag.BoolVar(&watch, "watch", false, "refresh the cache when a file changes")
	flag.StringVar(&addr, "addr", ":8484", "[host]:port to listen to")
	flag.StringVar(&index, "index", "index.html", "Name of the index file")
	flag.StringVar(&prefix, "prefix", "", "source folder")
	flag.Parse()
	ignores = []string{".git", ".hg", ".svn"}
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
	err = vfs.LoadDir(prefix, ignores)
	if err != nil {
		log.Fatal(err)
	}

	if watch {
		log.Printf("Starting watcher...")
		go WatchSync(vfs, prefix, ignores)
	}

	s := NewServer(vfs, addr, index, cors)

	log.Printf("Listening on %s", addr)
	log.Fatal(s.ListenAndServe())
}
