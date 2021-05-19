// MIT Licence. Copyright 2021, Storj (https://github.com/storj/dashborj/)

package main

import (
	"context"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

var systems = []System{
	{Kind: Satellite, Resolver: DNS, Host: "us1.storj.io"},
	{Kind: Satellite, Resolver: DNS, Host: "eu1.storj.io"},
	{Kind: Satellite, Resolver: DNS, Host: "ap1.storj.io"},
	{Kind: Satellite, Resolver: DNS, Host: "us2.storj.io"},
	{Kind: Satellite, Resolver: DNS, Host: "saltlake.tardigrade.io"},
	{Kind: Satellite, Resolver: DNS, Host: "europe-north-1.tardigrade.io"},

	{Kind: Auth, Resolver: DNS, Host: "auth.us1.storjshare.io"},
	{Kind: Auth, Resolver: DNS, Host: "auth.eu1.storjshare.io"},
	{Kind: Auth, Resolver: DNS, Host: "auth.ap1.storjshare.io"},

	{Kind: Linksharing, Resolver: DNS, Host: "link.us1.storjshare.io"},
	{Kind: Linksharing, Resolver: DNS, Host: "link.eu1.storjshare.io"},
	{Kind: Linksharing, Resolver: DNS, Host: "link.ap1.storjshare.io"},

	{Kind: Gateway, Resolver: DNS, Host: "gateway.us1.storjshare.io"},
	{Kind: Gateway, Resolver: DNS, Host: "gateway.eu1.storjshare.io"},
	{Kind: Gateway, Resolver: DNS, Host: "gateway.ap1.storjshare.io"},

	{Kind: Gateway, Resolver: Hetzner, Host: "gateway.tardigradeshare.io"},
}

func main() {
	resolveIPs()
	authMethods, hostKeyCallback := loadUserSSHFiles()

	proxy := NewProxy()
	proxy.sshConfig = ssh.ClientConfig{
		Timeout:         10 * time.Second,
		User:            "root",
		Auth:            authMethods,
		HostKeyCallback: hostKeyCallback,
	}

	http.Handle("/p/", proxy)
	http.HandleFunc("/", handleRequest)
	log.Println("listening on", ":8090")
	log.Fatal(http.ListenAndServe(":8090", nil))
}

func resolveIPs() {
	for x, sys := range systems {
		if sys.Resolver == Hetzner {
			hCloudToken := os.Getenv("HCLOUD_TOKEN")
			client := hcloud.NewClient(hcloud.WithToken(hCloudToken))
			servers, _, err := client.Server.List(context.Background(), hcloud.ServerListOpts{})
			if err != nil {
				log.Fatalf("error retrieving server: %s\n", err)
			}
			for _, server := range servers {
				if server != nil && strings.HasPrefix(server.Name, "stargate") {
					systems[x].IPs = append(systems[x].IPs, server.PublicNet.IPv4.IP.String())
				}
			}
		} else if sys.Resolver == DNS {
			iprecords, err := net.LookupIP(sys.Host)
			if err != nil {
				log.Fatalf("error resolving host: %s\n", err)
			}
			for _, ip := range iprecords {
				systems[x].IPs = append(systems[x].IPs, ip.String())
			}
		} else {
			log.Fatal("Uknown system resolver", sys.Resolver)
		}
	}
}

func loadUserSSHFiles() (methods []ssh.AuthMethod, hostKeyCallback ssh.HostKeyCallback) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("home dir not found", err)
	}
	sshDir := filepath.Join(home, ".ssh")
	files, err := ioutil.ReadDir(sshDir)
	if err != nil {
		log.Fatal("error enumerating .ssh directory", err)
	}
	for _, file := range files {
		buf, err := ioutil.ReadFile(filepath.Join(sshDir, file.Name()))
		if err != nil {
			continue
		}
		signer, err := ssh.ParsePrivateKey(buf)
		if err != nil {
			log.Println("Found public key " + file.Name())
			methods = append(methods, ssh.PublicKeys(signer))
		}
	}
	if len(methods) == 0 {
		log.Fatal("no SSH keys found")
	}
	knownHosts := filepath.Join(sshDir, "known_hosts")
	hostKeyCallback, err = knownhosts.New(knownHosts)
	if err != nil {
		log.Fatal("error parsing known_hosts", err)
	}
	return methods, hostKeyCallback
}
