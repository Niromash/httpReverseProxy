package main

import (
	"httpReverseProxy/config"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func main() {

	cfgRaw, err := config.HandleConfigFile()
	if err != nil {
		return
	}

	cfg, err := config.Decode(cfgRaw)
	if err != nil {
		log.Fatalln("Unable to parse config file: ", err)
		return
	}

	for _, host := range cfg.ProxyHost {
		log.Printf("Redirect %s to %s", host.Hostname, host.RedirectTo)
	}

	getProxyHostFromHostname := func(hostname string) *config.ProxyHost {
		for _, host := range cfg.ProxyHost {
			if strings.EqualFold(host.Hostname, hostname) {
				return &host
			}
		}
		return nil
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		hostname := strings.Split(r.Host, ":")[0]
		proxyHost := getProxyHostFromHostname(hostname)
		if proxyHost != nil {
			remote, err := url.Parse(proxyHost.RedirectTo)
			if err != nil {
				panic(err)
			}
			proxy := httputil.NewSingleHostReverseProxy(remote)

			// log.Printf("Redirecting %s to %s", hostname, redirectTo)
			r.Host = remote.Host
			w.Header().Set("Cache-control", "no-cache")
			if _, err = net.LookupIP(remote.Hostname()); err != nil {
				log.Println("Target cannot be reached:", err.Error())
				w.WriteHeader(500)
				_, _ = w.Write([]byte("Cannot reach the target"))
				return
			}
			proxy.ServeHTTP(w, r)
		} else {
			w.WriteHeader(404)
			if cfg.Preferences.UseIndexFileWhenHostnameNotFound {
				indexFile, err := os.OpenFile("./"+cfg.Preferences.IndexFilePath, os.O_RDONLY, 777)
				if err == nil {
					indexFileContent, err := ioutil.ReadAll(indexFile)
					if err == nil {
						_, err = w.Write(indexFileContent)
						return
					}
				}
				log.Println("Unable to send file when hostname not found: ", err)
				return
			}
			_, err := w.Write([]byte(cfg.Preferences.HostnameNotFoundMessage))
			if err != nil {
				log.Println("Unable to send message when hostname not found")
				return
			}
		}
	})

	err = http.ListenAndServe(":"+strconv.Itoa(cfg.HttpServerOptions.Port), nil)
	if err != nil {
		panic(err)
	}
}
