// MIT Licence. Copyright 2018, Digineo GmbH (https://github.com/digineo/http-over-ssh/)

package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func (proxy *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	key, uri, err := parseRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, err)
		return
	}

	r.Close = false
	r.Host = ""
	r.URL, _ = url.Parse(uri)
	r.RequestURI = ""
	removeHopHeaders(r.Header)

	// do the request
	client := proxy.getClient(*key)
	res, err := client.httpClient.Do(r)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintln(w, err.Error())
		return
	}

	// copy response header and body
	copyHeader(w.Header(), res.Header)
	w.WriteHeader(res.StatusCode)
	_, _ = io.Copy(w, res.Body)
	res.Body.Close()
}

func parseRequest(r *http.Request) (*clientKey, string, error) {
	parts := strings.SplitN(r.RequestURI, "/", 6)
	if len(parts) != 6 { // the initial slash makes parts[0] empty
		return nil, "", errors.New("bad request URI " + r.RequestURI)
	}
	sysIndex, err := strconv.Atoi(parts[2])
	if err != nil || sysIndex < 0 || sysIndex >= len(systems) {
		return nil, "", errors.New("bad request URI sysIndex " + r.RequestURI)
	}
	ipIndex, err := strconv.Atoi(parts[3])
	if err != nil || ipIndex < 0 || ipIndex >= len(systems[sysIndex].IPs) {
		return nil, "", errors.New("bad request URI ipIndex " + r.RequestURI)
	}
	port, err := strconv.Atoi(parts[4])
	if err != nil {
		return nil, "", errors.New("bad request URI port " + r.RequestURI)
	}
	path := parts[5]

	key := clientKey{
		// host:     systems[hostIndex].IPs[ipIndex],
		// port:     uint16(systems[hostIndex].SSHPort),
		// username: systems[hostIndex].SSHUser,
		SystemIndex: sysIndex,
		IPIndex:     ipIndex,
	}

	// todo: extract scheme from URL
	return &key, fmt.Sprintf("%s://%s:%d/%s", "http", "127.0.0.1", port, path), nil
}

// Hop-by-hop headers. These are removed when sent to the backend.
// http://www.w3.org/Protocols/rfc2616/rfc2616-sec13.html
var hopHeaders = []string{
	"Connection",
	"Proxy-Connection", // non-standard but still sent by libcurl and rejected by e.g. google
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te",      // canonicalized version of "TE"
	"Trailer", // not Trailers per URL above; http://www.rfc-editor.org/errata_search.php?eid=4522
	"Transfer-Encoding",
	"Upgrade",
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

// removeHopHeaders removes hop-by-hop headers to the backend. Especially
// important is "Connection" because we want a persistent
// connection, regardless of what the client sent to us.
func removeHopHeaders(header http.Header) {
	for _, h := range hopHeaders {
		header.Del(h)
	}
}
