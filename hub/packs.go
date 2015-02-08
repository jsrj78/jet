package main

import (
	"os"
	"os/exec"
	"path"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/surge/glog"
)

var configChan = make(chan configEvent)

// hubConfig represents the configuration file settings of the hub.
type hubConfig struct {
	Packs map[string]string
}

// launchAllPacks is called once on startup, and is driven by the config file.
func launchAllPacks() {
	config := hubConfig{make(map[string]string)}
	go config.processEvents()

	reload()
}

// reload can be called any time to force re-reading of the configuration file.
func reload() {
	newConf := &hubConfig{}
	if err := newConf.load(); err != nil {
		glog.Error(err)
		return
	}

	configChan <- configEvent{"", newConf}
}

// configEvent represents changes to the configuration and pack state changes
type configEvent struct {
	name string
	arg  interface{}
}

// processEvents is called as goroutine from launchAllPacks.
func (h *hubConfig) processEvents() {
	processes := make(map[string]*os.Process)

	for event := range configChan {
		name := event.name

		switch argVal := event.arg.(type) {

		case *hubConfig:
			// only act on what has changed
			for name, vOld := range h.Packs {
				vNew, ok := argVal.Packs[name]
				delete(argVal.Packs, name)

				if ok && vNew == vOld {
					continue
				}

				if ok {
					h.Packs[name] = vNew
				} else {
					delete(h.Packs, name)
				}

				glog.Infoln("quit:", name)
				if p, ok := processes[name]; ok {
					p.Signal(syscall.SIGQUIT)
					delete(processes, name)
					continue // exit will re-launch
				}

				if ok {
					h.launch(name, vNew)
				}
			}
			for kNew, vNew := range argVal.Packs {
				h.launch(kNew, vNew)
			}

		case *os.Process:
			glog.Infoln("started:", name, "pid:", argVal.Pid)
			processes[name] = argVal

		default:
			glog.Infoln("ended:", name, "status:", event.arg)
			delete(processes, name)

			time.Sleep(3 * time.Second)
			if cmd, ok := h.Packs[name]; ok {
				h.launch(name, cmd)
			}
		}
	}
}

// load the configuration settings from file.
func (h *hubConfig) load() error {
	configFile := path.Join(*runFlag, "hubtab.toml")
	if _, err := toml.DecodeFile(configFile, h); err != nil {
		return err
	}
	glog.Infoln("config loaded:", configFile)
	return nil
}

// launch the specified pack and report all process state changes as events.
func (h *hubConfig) launch(name, command string) {
	glog.Infoln("launch:", name, "cmd:", command)

	h.Packs[name] = command

	go func() {
		cmd := exec.Command("/bin/sh", "-c", command)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Start(); err != nil {
			glog.Error(err)
			return
		}

		configChan <- configEvent{name, cmd.Process} // lift-off!

		err := cmd.Wait()

		configChan <- configEvent{name, err} // touch-down!
	}()
}
