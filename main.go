package main

import (
	"flag"

	"github.com/ubiquum/ubiquum/rfb"
)

var (
	listen = flag.String("listen", ":5900", "listen on [ip]:port")
)

const (
	width  = 1920
	height = 1080
)

func main() {
	flag.Parse()
	rfb.Serve(*listen)
}
