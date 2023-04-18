package crutchService

import (
	types "crutch/internal/services"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type dirScanner struct {
	dirInfoList []types.DirData
	pathList    []string
	baseRules   map[string]string
}

func newScanner(baseRules map[string]string) *dirScanner {
	return &dirScanner{baseRules: baseRules}
}

func (d *dirScanner) run() ([]types.DirData, error) {
	if err := d.scanRootDir(); err != nil {
		return nil, fmt.Errorf("scan root dir error: '%v'", err)
	}

	if err := d.collectConfigs(); err != nil {
		return nil, fmt.Errorf("collect config error: '%v'", err)
	}

	return d.dirInfoList, nil
}

func (d *dirScanner) scanRootDir() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Println(err)
		return err
	}

	path := filepath.Join(homeDir, ".crutch")
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Println(err)
		return err
	}

	for _, file := range files {
		if file.Name() != "rules.json" {
			d.pathList = append(d.pathList, filepath.Join(path, file.Name()))
		}
	}
	return nil
}

func (d *dirScanner) collectConfigs() error {
	d.dirInfoList = make([]types.DirData, 0, len(d.pathList))
	for _, configPath := range d.pathList {
		dirData, err := readConfig(configPath)
		if err != nil {
			log.Println(err)
			return err
		}
		d.dirInfoList = append(d.dirInfoList, dirData)
	}
	return nil
}

func readConfig(configPath string) (types.DirData, error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return types.DirData{}, fmt.Errorf("config read error: '%v' in file: '%v'", err, configPath)
	}

	dirData := types.DirData{}
	err = json.Unmarshal(data, &dirData)
	if err != nil {
		return types.DirData{}, fmt.Errorf("config unmarshalling error: '%v' in file: '%v'", err, configPath)
	}

	return dirData, nil
}
