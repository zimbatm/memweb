package main

import (
	"github.com/zimbatm/httputil2"
	"net/http"
	"os"
	"time"
)

func NewServer(addr string, prefix string, mem bool, cors bool) *http.Server {
	var h http.Handler

	if mem {
		h = http.FileServer(httputil2.MemDir(prefix))
	} else {
		h = http.FileServer(http.Dir(prefix))
	}

	h = httputil2.GzipHandler(h)

	if cors {
		h = httputil2.CORSHandler(h, "*", 0)
	}

	h = httputil2.LogHandler(
		h,
		os.Stdout,
		httputil2.CommonLogFormatter(httputil2.CommonLogFormat),
	)

	return &http.Server{
		Addr:           addr,
		Handler:        h,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}
