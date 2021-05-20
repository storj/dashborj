DashBorj
========

DashBorj provides a dashboard for internal Storj dev-ops.

Setup
-----

Things you need:
 - A Hetzner token set in your environment as HCLOUD_TOKEN
 - ~/.ssh private key certificates corresponding to these servers
 - ~/.ssh known_hosts entries corresponding to these servers

Usage
-----

`go run .` then open a browser to http://localhost:8090

Design
------

Mostly it just resolves host names and SSH'es into boxes for you.  ¯\\_(ツ)_/¯

Pretty much everything of value here is copped from https://github.com/digineo/http-over-ssh/.

## License

MIT Licence. Copyright 2021, Storj

MIT Licence. Copyright 2018, Digineo GmbH
