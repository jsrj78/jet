// TODO messy code, all the clean stuff is in main.go

package main

import (
	"github.com/sevlyar/go-daemon"
	"log"
	"os"
	"syscall"
)

var ctx *daemon.Context

func daemonAndSignalSetup() {
	var cmd string
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	if cmd == "" && !daemon.WasReborn() {
		usage(false)
		return
	}

	ctx = &daemon.Context{
		PidFileName: "jethub.pid",
		LogFileName: "jethub.log",
	}

	daemon.AddCommand(daemon.StringFlag(&cmd, "quit"), syscall.SIGQUIT, onQuit)
	daemon.AddCommand(daemon.StringFlag(&cmd, "stop"), syscall.SIGTERM, onQuit)
	daemon.AddCommand(daemon.StringFlag(&cmd, "reload"), syscall.SIGHUP, reload)

	if len(daemon.ActiveFlags()) > 0 {
		d, err := ctx.Search()
		if err != nil {
			log.Fatalln(err)
		}
		daemon.SendCommands(d)
		return
	}

	performCmd(cmd)
}

func startDaemon() {
	d, err := ctx.Reborn()
	if err != nil {
		log.Fatalln(err)
	}
	if d == nil {
		defer ctx.Release()
		defer os.Remove(ctx.PidFileName)

		log.Println("- - - - - - - - - - - - - - -")
		log.Println("starting", os.Args)

		go worker()

		err = daemon.ServeSignals()
		if err != nil {
			log.Println(err)
		}
		log.Println("terminated")
	}
}

func onQuit(sig os.Signal) error {
	termHandler(sig == syscall.SIGQUIT)
	return daemon.ErrStop
}

func reload(sig os.Signal) error {
	reloadHandler()
	return nil
}
