package config

import (
	"github.com/BurntSushi/toml"
	"httpReverseProxy/api"
	"log"
	"os"
)

const configSample = `[httpServer]
port = 80

[preferences]
useIndexFileWhenHostnameNotFound = true
# Only if useIndexFileWhenHostnameNotFound is set to true
indexFilePath = "index.html" # Only html files accepted
# Only if useIndexFileWhenHostnameNotFound is set to false false
hostnameNotFoundMessage = "Web app linked to the hostname you provided does not exist!"

[[proxyHost]]
hostname = "yourproject.com"
redirectTo = "http://web_app:3000"

[[proxyHost]]
hostname = "api.yourproject.com"
redirectTo = "http://api_app:3000"
`

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

func Decode(data []byte) (*api.Config, error) {
	var conf *api.Config
	if _, err := toml.Decode(string(data), &conf); err != nil {
		return nil, err
	}
	return conf, nil
}
