package main

import (
	"net/http"
	"strings"
	"fmt"
	"time"
)

const CommonLogFormat = `%h %l %u %t "%r" %>s %b`
const CommonDateFormat = "02/Jan/2006:15:04:05 -0700"

type logMiddleware struct {
	parent http.Handler
	format string
}

type logResponseWriter struct {
	http.ResponseWriter
	status int
	bytes int
}

func LogMiddleware(parent http.Handler, f string) http.Handler {
	return &logMiddleware{parent, f}
}

func (self *logMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	lw := &logResponseWriter{w, -1, 0}
	self.parent.ServeHTTP(lw, r)


	line := self.format
	// %h is the IP address of the remote host
	line = strings.Replace(line, "%h", r.RemoteAddr, -1)
	// %l is the ident and will not be implemented
	line = strings.Replace(line, "%l", "-", -1)
	// TODO: %u is the REMOTE_USER
	line = strings.Replace(line, "%u", "-", -1)
	// %t 10/Oct/2000:13:55:36 -0700
	line = strings.Replace(line, "%t", start.Format(CommonDateFormat), -1)
  // %r GET /apache_pb.gif HTTP/1.0
	line = strings.Replace(line, "%r",
		fmt.Sprintf("%s %s %s", r.Method, r.RequestURI, r.Proto), -1)
	// %>s HTTP status
	line = strings.Replace(line, "%>s", fmt.Sprintf("%d", lw.status), -1)

	bytes := "-"
	Bytes := fmt.Sprintf("%d", lw.bytes)
	if lw.bytes > 0 {
		bytes = Bytes
	}
	// %b size of the body. - for no content
	line = strings.Replace(line, "%b", bytes, -1)
	// %B like %b. 0 for no content
	line = strings.Replace(line, "%B", Bytes, -1)

	fmt.Println(line)
}

func (self *logResponseWriter) WriteHeader(status int) {
	self.status = status
	self.ResponseWriter.WriteHeader(status)
}

func (self *logResponseWriter) Write(data []byte) (n int, err error) {
	n, err = self.ResponseWriter.Write(data)
	self.bytes += n
	return
}
