package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
)

var addr string
var prefix string
var mem bool
var cors bool

func init() {
	
	flag.StringVar(&addr, "addr", ":8484", "[host]:port to listen to")
	flag.StringVar(&prefix, "prefix", "", "source folder")
	flag.BoolVar(&mem, "mem", true, "load files into memory")
	flag.BoolVar(&cors, "cors", false, "allow CORS requests")
	
	flag.Parse()
}

func main() {
	var err error
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

	s := NewServer(addr, prefix, mem, cors)

	log.Printf("Listening on %s", addr)
	log.Fatal(s.ListenAndServe())
}
