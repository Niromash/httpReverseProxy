package config

import (
	"github.com/BurntSushi/toml"
	"log"
	"os"
)

func HandleConfigFile() ([]byte, error) {
	file, err := os.ReadFile("./config.toml")
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("./config.toml file not found, trying to create it...")
			file = []byte(configSample)
			if err := os.WriteFile("./config.toml", file, 777); err != nil {
				log.Fatal("Unable to create default config file: " + err.Error())
				return nil, err
			}
			log.Println("./config.toml file created with sample config.")
			return file, nil
		}

		if os.IsPermission(err) {
			log.Fatal("No permission to retrieve config file: " + err.Error())
			return nil, err
		}
		log.Fatal("Unable to retrieve config file: " + err.Error())
		return nil, err
	}

	return file, nil
}

func Decode(data []byte) (*Config, error) {
	var conf *Config
	if _, err := toml.Decode(string(data), &conf); err != nil {
		return nil, err
	}
	return conf, nil
}
