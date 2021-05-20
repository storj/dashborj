// MIT Licence. Copyright 2021, Storj (https://github.com/storj/dashborj/)

package main

import (
	"context"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

const (
	Satellite   = "satellite"
	Auth        = "auth"
	Linksharing = "linksharing"
	Gateway     = "gateway"
	DNS         = "dns"
	Hetzner     = "hetzner"

	AuthDebugPort = 5998
	GateDebugPort = 5999
	LinkDebugPort = 5997
	SatDebugPort  = 5996
)

type System struct {
	Kind       string
	Host       string
	IPResolver string
	IPs        []string

	SSHUser   string
	SSHPort   int
	sshConfig ssh.ClientConfig
}

func applySystemsDefaults() {
	// load SSH configs from ~/.ssh
	authMethods, hostKeyCallback := loadUserSSHFiles()

	for x, _ := range systems {
		//load SSH defaults
		if systems[x].SSHPort == 0 {
			systems[x].SSHPort = 22
		}
		if systems[x].SSHUser == "" {
			systems[x].SSHUser = "root"
		}
		if systems[x].IPResolver == "" {
			systems[x].IPResolver = DNS
		}
		// todo: move this?
		systems[x].sshConfig = ssh.ClientConfig{
			Timeout:         10 * time.Second,
			User:            systems[x].SSHUser,
			Auth:            authMethods,
			HostKeyCallback: hostKeyCallback,
		}

		// load IPs
		if systems[x].IPResolver == Hetzner {
			hCloudToken := os.Getenv("HCLOUD_TOKEN")
			if hCloudToken == "" {
				log.Fatal("environment variable HCLOUD_TOKEN must be defined")
			}
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
		} else if systems[x].IPResolver == DNS {
			iprecords, err := net.LookupIP(systems[x].Host)
			if err != nil {
				log.Fatalf("error resolving host: %s\n", err)
			}
			for _, ip := range iprecords {
				systems[x].IPs = append(systems[x].IPs, ip.String())
			}
		} else {
			log.Fatal("Uknown system resolver", systems[x].IPResolver)
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
		if err == nil {
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
