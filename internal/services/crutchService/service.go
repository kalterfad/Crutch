package crutchService

import (
	"crutch/internal/services"
	"log"
	"os"
	"path/filepath"
)

func StartService() {
	configPath, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(err)
	}

	configPath = filepath.Join(configPath, ".crutch", "rules.json")
	newWorkerCh := make(chan services.DirData)
	cancelWorkerCh := make(chan services.ChannelID)
	facade := newFacade(newWorkerCh, cancelWorkerCh)

	baseRules, err := readConfig(configPath)
	if err != nil {
		log.Fatalf("read config error: '%v'", err)
	}

	dirDataList, err := newScanner(baseRules.Rules).run()
	if err != nil {
		log.Fatalf("scanner error: '%v'", err)
	}
	if len(dirDataList) > 0 {
		go facade.restart(dirDataList)
	}

	socket := newSocket(newWorkerCh, cancelWorkerCh)
	if err = socket.run(); err != nil {
		log.Fatalf("socket error: '%v'", err)
	}
}
