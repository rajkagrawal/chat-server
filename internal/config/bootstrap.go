package config

import (
	"encoding/json"
	"io/ioutil"
)

// AppConfig defines the application configuration
type appConfig struct {
	ServerIP     string `json:"ip"`
	HTTPServerIP string `json:"http_ip"`
	PortNum      string `json:"port"`
	HTTPPortNum  string `json:"http_port"`
	LogFile      string `json:"logfile"`
}

// BS struct is just to contain the method
type BS struct{}

// Boot is bootstrap function
func (b BS) Boot(configFileLoc string) (*appConfig, error) {
	configData, err := ioutil.ReadFile(configFileLoc)
	if err != nil {
		return nil, err
	}
	var appConf appConfig
	err = json.Unmarshal(configData, &appConf)

	return &appConf, err
}
