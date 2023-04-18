package crutchService

import (
	types "crutch/internal/services"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type observer struct {
	watcher  *fsnotify.Watcher
	dirEvent chan types.DirPathData

	dirData  types.DirData
	cancelCh chan struct{}
}

func newObserver(dirEvent chan types.DirPathData, dirData types.DirData, cancelCh chan struct{}) *observer {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("new watcher error: '%v'", err)
	}
	return &observer{
		watcher:  watcher,
		dirEvent: dirEvent,
		dirData:  dirData,
		cancelCh: cancelCh,
	}
}

func (o *observer) run() {
	go o.findFiles()
	go o.watch()
}

func (o *observer) findFiles() {
	files, err := os.ReadDir(o.dirData.Path)
	if err != nil {
		log.Fatalf("directory read error: '%v'", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		extension := filepath.Ext(file.Name())

		if path, ok := o.dirData.Rules[extension]; ok {
			o.dirEvent <- types.DirPathData{
				OldPath:  filepath.Join(o.dirData.Path, file.Name()),
				NewPath:  filepath.Join(o.dirData.Path, path, "/"),
				FileName: file.Name(),
			}
		}
	}
}

func (o *observer) watch() {
	if err := o.watcher.Add(o.dirData.Path); err != nil {
		log.Fatalf("error adding directory: '%v' to scan: '%v'", o.dirData.Path, err)
	}
	defer func(watcher *fsnotify.Watcher) {
		err := watcher.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(o.watcher)

	log.Printf("watcher '%v' working", o.dirData.ChannelID)

	for {
		select {
		case event := <-o.watcher.Events:

			extension := filepath.Ext(event.Name)

			if path, ok := o.dirData.Rules[extension]; ok {
				stringSplit := strings.Split(event.Name, "/")
				pathData := types.DirPathData{
					OldPath:  event.Name,
					NewPath:  o.dirData.Path + "/" + path + "/",
					FileName: stringSplit[len(stringSplit)-1],
				}
				o.dirEvent <- pathData
			}

		case <-o.cancelCh:
			return
		case err := <-o.watcher.Errors:
			log.Fatalf("watcher error: '%v'", err)
		}
	}
}
