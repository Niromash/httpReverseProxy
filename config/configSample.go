package config

const configSample = `[httpServer]
port = 80

[[proxyHost]]
hostname = "yourproject.com"
redirectTo = "http://web_app:3000"

[[proxyHost]]
hostname = "api.yourproject.com"
redirectTo = "http://api_app:3000"
`

type Config struct {
	HttpServerOptions HttpServerOptions `toml:"httpServer"`
	ProxyHost         []ProxyHost       `toml:"proxyHost"`
}

type HttpServerOptions struct {
	Port int `toml:"port"`
}

type ProxyHost struct {
	Hostname   string `toml:"hostname"`
	RedirectTo string `toml:"redirectTo"`
}
