// MIT Licence. Copyright 2021, Storj (https://github.com/storj/dashborj/)

package main

import (
	"fmt"
	"net/http"
)

const (
	Satellite   = "satellite"
	Auth        = "auth"
	Linksharing = "linksharing"
	Gateway     = "gateway"
	DNS         = "dns"
	Hetzner     = "hetzner"
)

type System struct {
	Kind     string
	Host     string
	Resolver string
	IPs      []string
}

func handleRequest(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte(`<html><body>`))

	for _, kind := range []string{Auth, Gateway, Linksharing, Satellite} {
		w.Write([]byte(fmt.Sprintf(`<h1>%s</h1>`, kind)))
		for _, sys := range systems {
			if sys.Kind != kind {
				continue
			}
			w.Write([]byte(fmt.Sprintf(`<h2><a href="/%s/">%s</a></h2>`, sys.Host, sys.Host)))
			switch sys.Kind {
			case Auth:
				writeAuthResponse(w, sys)
			case Gateway:
				writeGatewayResponse(w, sys)
			case Linksharing:
				writeLinkResponse(w, sys)
			case Satellite:
				writeSatResponse(w, sys)
			}
		}
	}

	w.Write([]byte(`</body></html>`))
}

func writeAuthResponse(w http.ResponseWriter, sys System) {
	for _, ip := range sys.IPs {
		w.Write([]byte(fmt.Sprintf(`<p><a href="/p/%s/">%s</a></p>`, ip, ip)))
	}
}

func writeGatewayResponse(w http.ResponseWriter, sys System) {
	for _, ip := range sys.IPs {
		w.Write([]byte(fmt.Sprintf(`<p><a href="/p/%s/">%s</a></p>`, ip, ip)))
	}
}

func writeLinkResponse(w http.ResponseWriter, sys System) {
	for _, ip := range sys.IPs {
		w.Write([]byte(fmt.Sprintf(`<p><a href="/p/%s/">%s</a></p>`, ip, ip)))
	}
}
func writeSatResponse(w http.ResponseWriter, sys System) {
	for _, ip := range sys.IPs {
		w.Write([]byte(fmt.Sprintf(`<p><a href="/p/%s/">%s</a></p>`, ip, ip)))
	}
}
