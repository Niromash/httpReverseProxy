package config

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

type Config struct {
	HttpServerOptions HttpServerOptions `toml:"httpServer"`
	Preferences       Preferences       `toml:"preferences"`
	ProxyHost         []ProxyHost       `toml:"proxyHost"`
}

type HttpServerOptions struct {
	Port int `toml:"port"`
}

type Preferences struct {
	UseIndexFileWhenHostnameNotFound bool   `toml:"useIndexFileWhenHostnameNotFound"`
	IndexFilePath                    string `toml:"indexFilePath"`
	HostnameNotFoundMessage          string `toml:"hostnameNotFoundMessage"`
}

type ProxyHost struct {
	Hostname   string `toml:"hostname"`
	RedirectTo string `toml:"redirectTo"`
}
