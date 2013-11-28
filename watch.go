package main

import (
	"github.com/romanoff/fsmonitor"
	"log"
	"path/filepath"
)

func WatchSync(vfs VFS, prefix string, ignores []string) {
	if err := watchSync(vfs, prefix, ignores); err != nil {
		log.Fatal(err)
	}
}

func watchSync(vfs VFS, prefix string, ignores []string) (err error) {
	var watcher *fsmonitor.Watcher

	watcher, err = fsmonitor.NewWatcherWithSkipFolders(ignores)
	if err != nil {
		return
	}

	if err = watcher.Watch(prefix); err != nil {
		return
	}

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
		case err = <-watcher.Error:
			log.Println("error:", err)
			return
		}
	}
	return
}
