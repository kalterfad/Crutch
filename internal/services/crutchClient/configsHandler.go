package crutchClient

import (
	types "crutch/internal/services"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func checkArgs(args []string) (string, error) {
	if len(args) < 3 {
		return "", fmt.Errorf("\u001B[31m\u001B[1mDirectory not specified\u001B[0m")
	}
	if args[2] == "." {
		args[2], _ = os.Getwd()
	}
	path, err := os.Stat(args[2])
	if err != nil {
		return "", fmt.Errorf("\u001B[31m\u001B[1mNo such directory\u001B[0m")
	}

	if path.IsDir() {
		return args[2], nil
	} else {
		return "", fmt.Errorf("\u001B[31m\u001B[1mThis is not directory\u001B[0m")
	}
}

func askUser() bool {
	for {
		fmt.Print("\x1b[1mDo you want to create your own rules for managing files? [\x1b[31myes/no\x1b[0m]: ")
		var isEdit string
		_, err := fmt.Scanln(&isEdit)
		if err != nil {
			continue
		}
		if isEdit == "yes" {
			return true
		} else if isEdit == "no" {
			return false
		}
	}
}

func checkPath(dirPath string) bool {
	file, err := os.Stat(dirPath)
	if err != nil || file.IsDir() {
		return false
	}
	return true
}

type configHandler struct{}

func (h *configHandler) removeConfig(configPath string) error {
	err := os.Remove(configPath)
	if err != nil {
		return err
	}
	return nil
}

func (h *configHandler) readConfig(configPath string) (types.DirData, error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return types.DirData{}, err
	}

	dirData := types.DirData{}
	err = json.Unmarshal(data, &dirData)
	if err != nil {
		return types.DirData{}, err
	}
	return dirData, nil
}

func (h *configHandler) saveConfig(configPath string, config types.DirData) error {
	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	rawDataOut, err := json.MarshalIndent(&config, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(configPath, rawDataOut, 0)
	if err != nil {
		return err
	}
	return nil
}

func (h *configHandler) copyFile(baseConfigPath, configPath string) error {
	log.Println("copy file")
	sourceFile, err := os.Open(baseConfigPath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	newFile, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, sourceFile)
	if err != nil {
		return err
	}
	return nil
}

func (h *configHandler) openNano(configPath string) error {
	cmd := exec.Command("nano", configPath)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func (h *configHandler) editConfig(configPath string, config types.DirData) (map[string]string, error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	dirData := types.DirData{}
	if err = json.Unmarshal(data, &dirData); err != nil {
		return nil, err
	}

	config.Rules = dirData.Rules

	rawDataOut, err := json.MarshalIndent(&config, "", "  ")
	if err != nil {
		return nil, err
	}

	if err = ioutil.WriteFile(configPath, rawDataOut, 0644); err != nil {
		return nil, err
	}

	return config.Rules, nil
}

func (h *configHandler) prepareMsg(message types.SocketMsg) ([]byte, error) {
	bytearrayMsg, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}
	return bytearrayMsg, nil
}
