package crutchClient

import (
	types "crutch/internal/services"
	"fmt"
	"os"
	"path/filepath"
)

func rmObserve(dirPath string) ([]byte, error) {
	homeDirPath, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(homeDirPath, ".crutch", "conf."+filepath.Base(dirPath)+".json")
	if !checkPath(configPath) {
		return nil, fmt.Errorf("\u001B[31m\u001B[1mDirectory '%v' is not being scanned\u001B[0m", dirPath)
	}

	configObj := newConfigFacade()
	dirData, err := configObj.readConfig(configPath)
	if err != nil {
		return nil, err
	}

	socketMsg := types.SocketMsg{
		Command:   types.RmWorker,
		ChannelID: dirData.ChannelID,
	}

	err = configObj.removeConfig(configPath)
	if err != nil {
		return nil, err
	}

	preparedMsg, err := configObj.prepareMsg(socketMsg)
	if err != nil {
		return nil, err
	}
	return preparedMsg, nil
}
