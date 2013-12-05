export GOPATH=$(PWD)

all: memweb

dist:
	./make-dist

memweb: *.go
	go build -o $@

.PHONY: all dist
