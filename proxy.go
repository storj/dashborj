// MIT Licence. Copyright 2018, Digineo GmbH (https://github.com/digineo/http-over-ssh/)

package main

import (
	"errors"
	"net"
	"net/http"
	"sync"

	"golang.org/x/crypto/ssh"
)

// Proxy holds the HTTP client and the SSH connection pool
type Proxy struct {
	clients map[clientKey]*client
	mtx     sync.Mutex
}

// NewProxy creates a new proxy
func NewProxy() *Proxy {
	return &Proxy{
		clients: make(map[clientKey]*client),
	}
}

// getClient returns a (un)connected SSH client
func (proxy *Proxy) getClient(key clientKey) *client {
	proxy.mtx.Lock()
	defer proxy.mtx.Unlock()

	// connection established?
	pClient := proxy.clients[key]
	if pClient != nil {
		return pClient
	}

	pClient = &client{
		key:       key,
		sshConfig: systems[key.SystemIndex].sshConfig, // make copy
	}
	sysIndex := key.SystemIndex
	pClient.sshConfig.HostKeyCallback = func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		if err := systems[sysIndex].sshConfig.HostKeyCallback(hostname, remote, key); err != nil {
			return err
		}
		if cert, ok := key.(*ssh.Certificate); ok && cert != nil {
			pClient.sshCert = cert
		}
		return nil
	}

	pClient.httpClient = &http.Client{
		Transport: &http.Transport{
			Dial: pClient.dial,
			DialTLS: func(network, addr string) (net.Conn, error) {
				return nil, errors.New("not implemented")
			},
		},
	}

	// set and return the new connection
	proxy.clients[key] = pClient
	return pClient
}
