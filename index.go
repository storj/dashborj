// MIT Licence. Copyright 2021, Storj (https://github.com/storj/dashborj/)

package main

import (
	"fmt"
	"net/http"
)

func handleRequest(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte(`<html><body>`))

	for _, kind := range []string{Auth, Gateway, Linksharing, Satellite} {
		w.Write([]byte(fmt.Sprintf(`<h1>%s</h1>`, kind)))
		for sysNum, sys := range systems {
			if sys.Kind != kind {
				continue
			}
			w.Write([]byte(fmt.Sprintf(`<h2><a href="/%s/">%s</a></h2>`, sys.Host, sys.Host)))
			switch sys.Kind {
			case Auth:
				writeAuthResponse(w, sys, sysNum)
			case Gateway:
				writeGatewayResponse(w, sys, sysNum)
			case Linksharing:
				writeLinkResponse(w, sys, sysNum)
			case Satellite:
				writeSatResponse(w, sys, sysNum)
			}
		}
	}

	w.Write([]byte(`</body></html>`))
}

func writeAuthResponse(w http.ResponseWriter, sys System, sysNum int) {
	for i, ip := range sys.IPs {
		w.Write([]byte(fmt.Sprintf(`<p><a href="/p/%d/%d/%d/debug/pprof/">pprof %s</a></p>`, sysNum, i, AuthDebugPort, ip)))
	}
}

func writeGatewayResponse(w http.ResponseWriter, sys System, sysNum int) {
	for i, ip := range sys.IPs {
		w.Write([]byte(fmt.Sprintf(`<p><a href="/p/%d/%d/%d/debug/pprof/">pprof %s</a></p>`, sysNum, i, GateDebugPort, ip)))
	}
}

func writeLinkResponse(w http.ResponseWriter, sys System, sysNum int) {
	for i, ip := range sys.IPs {
		w.Write([]byte(fmt.Sprintf(`<p><a href="/p/%d/%d/%d/debug/pprof/">pprof %s</a></p>`, sysNum, i, LinkDebugPort, ip)))
	}
}
func writeSatResponse(w http.ResponseWriter, sys System, sysNum int) {
	for i, ip := range sys.IPs {
		w.Write([]byte(fmt.Sprintf(`<p><a href="/p/%d/%d/%d/debug/pprof/">pprof %s</a></p>`, sysNum, i, SatDebugPort, ip)))
	}
}
