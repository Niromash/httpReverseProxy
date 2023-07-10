package main

import (
	"context"
	"git.niromash.me/odyssey/reverse-proxy/api"
	"git.niromash.me/odyssey/reverse-proxy/config"
	"git.niromash.me/odyssey/reverse-proxy/utils"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

func main() {

	sigCh := make(chan os.Signal, 1)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	signal.Notify(sigCh, syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sigCh
		cancel()
		log.Println("The application is stopping, please wait few seconds")
		os.Exit(0)
	}()

	//atoi, err := strconv.Atoi(os.Getenv("PORT"))
	//if err != nil {
	//	log.Fatalln(err)
	//}

	//cfg := api.Config{
	//	HttpServerOptions: api.HttpServerOptions{
	//		Port: atoi,
	//	},
	//	Preferences: api.Preferences{
	//		HostnameNotFoundMessage: "Not found",
	//	},
	//	ProxyHost: []api.ProxyHost{
	//		{
	//			Hostname:   "localhost",
	//			RedirectTo: "https://www.twitch.tv",
	//			PathType:   api.PathTypePrefix,
	//			Path:       "/directory",
	//		},
	//		{
	//			Hostname:   "localhost",
	//			RedirectTo: "https://blog.filador.fr",
	//			//RedirectTo: "https://emeralds-hub.com",
	//			PathType: api.PathTypePrefix,
	//			Path:     "/",
	//		},
	//	},
	//}

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

	getProxyHostFromHostnameAndPath := func(hostname string, path string) *api.ProxyHost {
		for _, host := range cfg.ProxyHost {
			if strings.EqualFold(host.Hostname, hostname) &&
				(host.PathType == api.PathTypePrefix && strings.HasPrefix(path, host.Path)) ||
				(host.PathType == api.PathTypeExact && path == host.Path) {
				return &host
			}
		}
		return nil
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		hostname := strings.Split(r.Host, ":")[0]
		proxyHost := getProxyHostFromHostnameAndPath(hostname, r.URL.Path)
		if proxyHost != nil {
			remote, err := url.Parse(proxyHost.RedirectTo)
			if err != nil {
				panic(err)
			}
			proxy := httputil.NewSingleHostReverseProxy(remote)

			proxy.Director = func(req *http.Request) {
				// Remove trailing slash
				//wildcardIndex := strings.IndexAny("/*catchall", "*")
				//proxyPath := singleJoiningSlash(remote.Path, r.URL.Path[wildcardIndex:])
				//if strings.HasSuffix(proxyPath, "/") && len(proxyPath) > 1 {
				//	proxyPath = proxyPath[:len(proxyPath)-1]
				//}
				//fmt.Println(proxyPath)

				req.URL.Scheme = remote.Scheme
				req.URL.Host = remote.Host
				req.Host = remote.Host
			}

			log.Printf("Redirecting %s to %s", hostname, remote.Host)
			r.Header.Set("Cache-control", "no-cache")
			r.Header.Set("Host", utils.IfThenElse(r.TLS == nil, "http://", "https://")+remote.Host)
			r.Header.Set("Referer", r.Referer())
			r.Header.Set("X-Forwarded-Host", proxyHost.Hostname)
			r.Header.Set("X-Forwarded-Proto", r.Proto)
			r.Header.Set("X-Forwarded-Port", strconv.Itoa(cfg.HttpServerOptions.Port))
			r.Header.Add("X-Origin-Host", remote.Host)

			if _, err = net.LookupIP(remote.Hostname()); err != nil {
				log.Println("Target cannot be reached:", err.Error())
				w.WriteHeader(500)
				_, _ = w.Write([]byte("Cannot reach the target"))
				return
			}
			if strings.HasPrefix(hostname, "stats.") {
				r.Header.Set("Origin", "") // For Grafana
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

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
