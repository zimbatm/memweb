Memweb - serve files from memory
================================

Memweb is a little web server that serves files from memory. Use it when you have a static website and that every second counts.

Custom builds
-------------

The [gobuild service](https://github.com/ddollar/gobuild) allows to compile
and download custom version of memweb. It's useful if you don't have the go
toolchain installed.

Example:

```bash
curl https://gobuild.herokuapp.com/zimbatm/memweb/master/darwin/amd64 -o memweb
chmod +x memweb
./memweb
```

TODO
----

* ETag handler
* Benchmark
* SSL support

* mapping support (with redirects). maybe
http://www.sitemaps.org/protocol.html
