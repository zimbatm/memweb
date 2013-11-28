Memweb - serve files from memory
================================

Memweb is a little web server that serves files from memory. Use it when you have a static website and that every second counts.


TODO
----

* Replace inotify-based update with direct disk access
http://golang.org/src/pkg/net/http/fs.go?s=3056:3160#L27

* Benchmark
* SSL support
* Apache Combined Logs
* CORS
* mapping support (with redirects). maybe
http://www.sitemaps.org/protocol.html
