package main

import (
	"github.com/howeyc/fsnotify"
	"log"
	"path/filepath"
)

func WatchFS(vfs VFS, prefix string) {
	var err error
	
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return
	}
	defer watcher.Close()

	done := make(chan bool)

	// Process events
	go func() {
		defer func() { done <- true }()

		for {
			select {
			case ev := <-watcher.Event:
				log.Println("event:", ev)
				if ev == nil {
					continue
				}
				if ev.IsCreate() || ev.IsModify() {
					var filename string
					if filename, err = filepath.Rel(prefix, ev.Name); err != nil {
						return
					}

					if err = vfs.LoadFile(filename, ev.Name); err != nil {
						return
					}
				} else if ev.IsDelete() || ev.IsRename() {
					vfs.Remove(ev.Name)
				} else {
					panic("BUG")
				}
			case err := <-watcher.Error:
				log.Println("error:", err)
				return
			}
		}
	}()

	if err = watcher.Watch(prefix); err != nil {
		log.Fatal(err)
	}

	<-done

	if err != nil {
		log.Fatal(err)
	}

	return
}
