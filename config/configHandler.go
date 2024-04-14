package config

import (
	"git.niromash.me/odyssey/reverse-proxy/api"
	"github.com/BurntSushi/toml"
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
pathType = "Prefix"
path = "/"

[[proxyHost]]
hostname = "api.yourproject.com"
redirectTo = "http://api_app:3000"
pathType = "Prefix"
path = "/"
`

func HandleConfigFile() ([]byte, error) {
	file, err := os.ReadFile("/etc/httpReverseProxy/config.toml")
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("/etc/httpReverseProxy/config.toml file not found, trying to create it...")
			file = []byte(configSample)
			if err := os.WriteFile("/etc/httpReverseProxy/config.toml", file, 777); err != nil {
				log.Fatal("Unable to create default config file: " + err.Error())
				return nil, err
			}
			log.Println("/etc/httpReverseProxy/config.toml file created with sample config.")
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
