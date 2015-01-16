package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Proxy struct {
		SSL struct {
			Enabled  bool   `json:"enabled"`
			Listener string `json:"listener"`
			Key      string `json:"key"`
			Cert     string `json:"cert"`
		} `json:"ssl"`
		HTTP struct {
			Enabled  bool   `json:"enabled"`
			Listener string `json:"listener"`
		} `json:"http"`
	} `json:"proxy"`
	Admin struct {
		Listener string `json:"listener"`
		Key      string `json:"key"`
		Cert     string `json:"cert"`
	} `json:"admin"`
	JWTKey     string `json:"jwt_key"`
	BoltDBFile string `json:"bolt_db_file"`
}

// ParseConfig parses json from the file provided by filename into a Config struct.
func ParseConfig(filename string) (*Config, error) {
	config := &Config{}
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return config, err
	}

	if err := json.Unmarshal(file, &config); err != nil {
		return config, err
	}

	return config, nil
}
