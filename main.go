// MIT Licence. Copyright 2021, Storj (https://github.com/storj/dashborj/)

package main

import (
	"log"
	"net/http"
	"sync"
)

var systems = []System{
	{Kind: Satellite, Host: "us1.storj.io"},
	{Kind: Satellite, Host: "eu1.storj.io"},
	{Kind: Satellite, Host: "ap1.storj.io"},
	{Kind: Satellite, Host: "us2.storj.io"},
	{Kind: Satellite, Host: "saltlake.tardigrade.io"},
	{Kind: Satellite, Host: "europe-north-1.tardigrade.io"},

	{Kind: Auth, Host: "auth.us1.storjshare.io"},
	{Kind: Auth, Host: "auth.eu1.storjshare.io"},
	{Kind: Auth, Host: "auth.ap1.storjshare.io"},

	{Kind: Linksharing, Host: "link.us1.storjshare.io"},
	{Kind: Linksharing, Host: "link.eu1.storjshare.io", SSHUser: "ubuntu", SSHPort: 2222},
	{Kind: Linksharing, Host: "link.ap1.storjshare.io"},

	{Kind: Gateway, Host: "gateway.us1.storjshare.io"},
	{Kind: Gateway, Host: "gateway.eu1.storjshare.io", SSHUser: "ubuntu", SSHPort: 2222},
	{Kind: Gateway, Host: "gateway.ap1.storjshare.io"},

	{Kind: Gateway, IPResolver: Hetzner, Host: "gateway.tardigradeshare.io"},
}
var systemsMtx sync.Mutex

func main() {
	applySystemsDefaults()

	http.Handle("/p/", NewProxy())      // ssh tunnels
	http.HandleFunc("/", handleRequest) // service list
	log.Println("listening on", ":8090")
	log.Fatal(http.ListenAndServe(":8090", nil))
}
