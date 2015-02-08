package main

import (
	"path"

	"github.com/BurntSushi/toml"
	"github.com/surge/glog"
)

var config hubConfig

type hubConfig struct {
	Packs map[string]string
}

func (h *hubConfig) load() error {
	configFile := path.Join(*runFlag, "hubtab.toml")
	if _, err := toml.DecodeFile(configFile, h); err != nil {
		return err
	}
	glog.Infoln("config loaded:", configFile)
	return nil
}

func launchAllPacks() {
	if err := config.load(); err != nil {
		glog.Fatal(err)
	}
	for k, v := range config.Packs {
		glog.Infoln(k, "=>", v)
	}
}

func reload() {
	var newConf hubConfig
	if err := newConf.load(); err != nil {
		glog.Error(err)
		return
	}

	// only act on what has changed
	for kOld, vOld := range config.Packs {
		vNew, ok := newConf.Packs[kOld]
		if !ok {
			glog.Infoln("delete", kOld)
		} else {
			if vOld != vNew {
				glog.Infoln("modify", kOld)
				config.Packs[kOld] = vNew
			}
			delete(newConf.Packs, kOld)
		}
	}
	for kNew, vNew := range newConf.Packs {
		glog.Infoln("insert", kNew)
		config.Packs[kNew] = vNew
	}
}
