package main

import (
	"httpReverseProxy/config"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
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

			//log.Printf("Redirecting %s to %s", hostname, redirectTo)
			r.Host = remote.Host
			w.Header().Set("Cache-control", "no-cache")
			proxy.ServeHTTP(w, r)
		} else {
			w.WriteHeader(404)
			_, err := w.Write([]byte("You entered an invalid hostname!"))
			if err != nil {
				w.WriteHeader(500)
				_, _ = w.Write([]byte("Internal error"))
				log.Println(err)
				return
			}
		}
	})

	err = http.ListenAndServe(":" + strconv.Itoa(cfg.HttpServerOptions.Port), nil)
	if err != nil {
		panic(err)
	}
}
