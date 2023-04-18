package crutchService

import (
	types "crutch/internal/services"
	"log"
	"os"
)

type fileHandler struct {
	pathData types.DirPathData
}

func newFileHandler(pathData types.DirPathData) fileHandler {
	return fileHandler{
		pathData: pathData,
	}
}

func (h fileHandler) move() {
	newPath := h.pathData.NewPath + h.pathData.FileName
	_ = os.Rename(h.pathData.OldPath, newPath)
}

func (h fileHandler) makeDir() {
	_, err := os.Stat(h.pathData.NewPath)

	if os.IsNotExist(err) {
		errDir := os.MkdirAll(h.pathData.NewPath, 0777)
		if errDir != nil {
			log.Fatal(err)
		}
	}
}

func (h fileHandler) handle() {
	h.makeDir()
	h.move()
}
