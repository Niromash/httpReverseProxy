package api

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

type PathType string

const (
	PathTypePrefix PathType = "Prefix"
	PathTypeExact  PathType = "Exact"
)

type ProxyHost struct {
	Hostname   string   `toml:"hostname"`
	RedirectTo string   `toml:"redirectTo"`
	PathType   PathType `toml:"pathType"` // Prefix or Exact
	Path       string   `toml:"path"`     // Path to match for the redirect
}
