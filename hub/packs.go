package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

var packMap = map[string]*exec.Cmd{}

// listen to requests to launch or kill a JET "pack"
func packsListener(feed, dir string) {
	if e := os.MkdirAll(dir, 0777); e != nil {
		log.Fatal(e)
	}

	for evt := range topicWatcher(feed) {
		packName := evt.Topic[6:] // TODO wrong if feed isn't "packs/+"
		logTopic := evt.Topic + "/log"

		killRunningPack(packName)
		delete(packMap, packName)

		if len(evt.Payload) == 0 {
			continue
		}

		// launch the new pack, with stdout & stderr output logged via pipes
		var packReq []string
		if evt.Decode(&packReq) {
			if len(packReq) == 0 {
				packReq = append(packReq, "") // avoid indexing error
			}
			path := validatePack(packReq[0], dir)
			if path != "" {
				log.Println("starting", packName, "pack:", packReq)
				cmd := exec.Command(path, packReq[1:]...)

				// capture stdout and log it to MQTT
				if pipe, e := cmd.StdoutPipe(); e != nil {
					log.Fatal(e)
				} else {
					go reportPackOutput(pipe, packName, logTopic, "")
				}

				// capture stderr and log it to MQTT with "(stderr)" prefix
				if pipe, e := cmd.StderrPipe(); e != nil {
					log.Fatal(e)
				} else {
					go reportPackOutput(pipe, packName, logTopic, "(stderr) ")
				}

				if e := cmd.Start(); e != nil {
					log.Println("pack:", path, "can't start", e)
					continue
				}
				log.Println("started:", path, "pid:", cmd.Process.Pid)
				packMap[packName] = cmd
			}
		}
	}
}

func killRunningPack(name string) {
	if cmd, ok := packMap[name]; ok && cmd.Process != nil {
		log.Println("stopping pack:", name, "pid:", cmd.Process.Pid)
		if e := cmd.Process.Kill(); e != nil {
			log.Println("kill", name, "error:", e)
		}
	}
}

func validatePack(name, dir string) string {
	if name == "" || strings.Contains(name, "/") {
		log.Println("pack name is not valid:", name)
		return ""
	}
	path, e := exec.LookPath(dir + "/" + name)
	if e != nil {
		log.Println(e)
		return ""
	}
	return path
}

func reportPackOutput(pipe io.ReadCloser, name, topic, prefix string) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		msg := prefix + scanner.Text()
		log.Println(name, "pack:", msg)
		sendToHub(topic, msg, false)
	}
	log.Println(name, "pack:", prefix + "EOF")
}
