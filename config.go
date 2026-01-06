package main

import (
	"log"
	"path/filepath"

	"gopkg.in/ini.v1"
)

type config struct {
	libdir    string // NOTE: the local directory to pull books from (server) OR the directory to store books into.
	remoteURL string // remote target: ideally set once and forget
	multiURL  string // multicast target: dynamically change based on mdns queries
}

func parseCfg(conf *config) {
	cfg, err := ini.Load(filepath.Join(".", "kindlesync.ini")) // we expect the configuration file to be in the same directory as binary
	if err != nil {
		log.Fatal(err)
	}
	conf.libdir = cfg.Section("storage").Key("library").String() //Always the absolute path (/mnt/us/path/to/dir)
	conf.remoteURL = "http://" + cfg.Section("").Key("address").String() + ":" + cfg.Section("").Key("port").String()
}
