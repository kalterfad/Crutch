package crutchClient

import (
	types "crutch/internal/services"
	"fmt"
)

type configFacade struct {
	handler *configHandler
}

func newConfigFacade() *configFacade {
	return &configFacade{
		handler: &configHandler{},
	}
}

func (c *configFacade) removeConfig(configPath string) error {
	err := c.handler.removeConfig(configPath)
	if err != nil {
		return fmt.Errorf("remove config: '%v' error: '%v'", configPath, err)
	}
	return nil
}

func (c *configFacade) readConfig(configPath string) (types.DirData, error) {
	data, err := c.handler.readConfig(configPath)
	if err != nil {
		return types.DirData{}, fmt.Errorf("read config: '%v' error: '%v'", configPath, err)
	}
	return data, nil
}

func (c *configFacade) saveConfig(configPath string, config types.DirData) error {
	err := c.handler.saveConfig(configPath, config)
	if err != nil {
		return fmt.Errorf("save config: '%v' error: '%v'", configPath, err)
	}
	return nil
}

func (c *configFacade) editConfig(configPath string, config types.DirData) (map[string]string, error) {
	rules, err := c.handler.editConfig(configPath, config)
	if err != nil {
		return nil, fmt.Errorf("edit config: '%v' error: '%v'", configPath, err)
	}
	return rules, nil
}

func (c *configFacade) openNano(configPath string) error {
	err := c.handler.openNano(configPath)
	if err != nil {
		return fmt.Errorf("open in Nano config: '%v' error: '%v'", configPath, err)
	}
	return nil
}

func (c *configFacade) copyFile(baseConfigPath, configPath string) error {
	if err := c.handler.copyFile(baseConfigPath, configPath); err != nil {
		return fmt.Errorf("copy file error: '%v'", err)
	}
	return nil
}

func (c *configFacade) prepareMsg(message types.SocketMsg) ([]byte, error) {
	data, err := c.handler.prepareMsg(message)
	if err != nil {
		return nil, fmt.Errorf("prepare messae: '%v' error: '%v'", message, err)
	}
	return data, nil
}
