package main

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"strings"
)

const PACKS_DIR = "./packs"

var packMap = map[string]*exec.Cmd{}

// listen to requests to launch or kill a JET "pack"
func packsListener(feed string) {
	if e := os.MkdirAll(PACKS_DIR, 0777); e != nil {
		log.Fatal(e)
	}

	for evt := range topicWatcher(feed) {
		packName := evt.Topic[6:]
		logTopic := evt.Topic + "/log"

		// kill the pack if it is currently running
		if cmd, ok := packMap[packName]; ok {
			if cmd.Process != nil {
				log.Println("stopping pack:", packName, "pid:", cmd.Process.Pid)
				if e := cmd.Process.Kill(); e != nil {
					log.Println("kill", packName, "error:", e)
				}
			}
			delete(packMap, packName)
		}

		// launch the new pack, with stdout & stderr output logged via pipes
		if len(evt.Payload) > 0 {
			var packReq []string
			if evt.Decode(&packReq) {
				var cmdName string
				if len(packReq) > 0 {
					cmdName = packReq[0]
				}
				if cmdName == "" || strings.Contains(cmdName, "/") {
					log.Println("pack name is not valid:", cmdName)
					continue
				}

				path, e := exec.LookPath(PACKS_DIR + "/" + cmdName)
				if e != nil {
					log.Println("can't find", cmdName, "in:", PACKS_DIR)
					continue
				}

				log.Println("starting pack:", packName, packReq)
				cmd := exec.Command(path, packReq[1:]...)

				// capture all stdout and log it to MQTT
				if pipe, e := cmd.StdoutPipe(); e != nil {
					log.Fatal(e)
				} else {
					go func() {
						scanner := bufio.NewScanner(pipe)
						for scanner.Scan() {
							msg := scanner.Text()
							log.Println("pack:", cmdName, "stdout:", msg)
							sendToHub(logTopic, msg, false)
						}
						log.Println("pack:", cmdName, "EOF on stdout")
					}()
				}

				// capture all stderr and log it to MQTT with "(stderr)" prefix
				if pipe, e := cmd.StderrPipe(); e != nil {
					log.Fatal(e)
				} else {
					go func() {
						scanner := bufio.NewScanner(pipe)
						for scanner.Scan() {
							msg := "(stderr) " + scanner.Text()
							log.Println("pack:", cmdName, "stderr:", msg)
							sendToHub(logTopic, msg, false)
						}
						log.Println("pack:", cmdName, "EOF on stderr")
					}()
				}

				if e := cmd.Start(); e != nil {
					log.Println("pack:", path, "can't start", e)
					continue
				}
				log.Println("started:", cmdName, "pid:", cmd.Process.Pid)
				packMap[packName] = cmd
			}
		}
	}
}
