// TODO messy code, all the clean stuff is in main.go

package main

import (
	"github.com/sevlyar/go-daemon"
	"log"
	"os"
	"syscall"
)

var ctx = &daemon.Context{
	PidFileName: "jethub.pid",
	LogFileName: "jethub.log",
}

func daemonAndSignalSetup() {
	var cmd string
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	if cmd == "" && !daemon.WasReborn() {
		usage(false)
		return
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
		// FIXME ctx.Release fails on MacOSX, no /proc/ for lockfile GetFdName
		// use explicit remove instead
		defer os.Remove(ctx.PidFileName)
		defer ctx.Release()

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
