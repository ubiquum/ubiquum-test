package main

import (
	"flag"

	log "github.com/sirupsen/logrus"
	"github.com/ubiquum/ubiquum/rfb"
)

var (
	listen   = flag.String("listen", ":5900", "listen on [ip]:port")
	logLevel = flag.String("log", "", "log level")
)

const (
	width  = 1920
	height = 1080
)

func main() {
	flag.Parse()

	if *logLevel != "" {
		level, err := log.ParseLevel(*logLevel)
		if err != nil {
			log.Fatalf("Failed to parse %s", *logLevel)
			return
		}
		log.SetLevel(level)
		log.Infof("Set log level to %s", *logLevel)
	}

	rfb.Serve(*listen)
}
