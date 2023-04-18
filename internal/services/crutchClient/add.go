package crutchClient

import (
	types "crutch/internal/services"
	"fmt"
	"github.com/google/uuid"
	"os"
	"path/filepath"
)

func addObserver(dirPath string, isEditRules bool) ([]byte, error) {
	homeDirPath, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	baseConfigPath := filepath.Join(homeDirPath, ".crutch", "rules.json")
	configPath := filepath.Join(homeDirPath, ".crutch", "conf."+filepath.Base(dirPath)+".json")
	if checkPath(configPath) {
		return nil, fmt.Errorf("\u001B[31m\u001B[1mDirectory '%v' is already being scanned\u001B[0m", dirPath)
	}

	socketMsg := types.SocketMsg{
		Command:   types.NewWorker,
		ChannelID: types.ChannelID(uuid.New().String()),
		Path:      dirPath,
	}
	configObj := newConfigFacade()

	if isEditRules {
		if err = configObj.copyFile(baseConfigPath, configPath); err != nil {
			return nil, err
		}
		if err = configObj.openNano(configPath); err != nil {
			return nil, err
		}

		config := types.DirData{
			ChannelID: socketMsg.ChannelID,
			Path:      dirPath,
		}

		rules, err := configObj.editConfig(configPath, config)
		if err != nil {
			return nil, err
		}
		socketMsg.Rules = rules

	} else {
		baseConfig, err := configObj.readConfig(baseConfigPath)
		if err != nil {
			return nil, err
		}

		config := types.DirData{
			ChannelID: socketMsg.ChannelID,
			Path:      dirPath,
			Rules:     baseConfig.Rules,
		}

		if err = configObj.saveConfig(configPath, config); err != nil {
			return nil, err
		}
		socketMsg.Rules = config.Rules
	}

	preparedMsg, err := configObj.prepareMsg(socketMsg)
	if err != nil {
		return nil, err
	}
	return preparedMsg, nil

}
