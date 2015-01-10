// TODO messy code, all the clean stuff is in main.go

package main

import (
	"flag"
	"github.com/sevlyar/go-daemon"
	"log"
	"os"
	"syscall"
)

var ctx = &daemon.Context{
	PidFileName: "hub.pid",
	LogFileName: "hub.log",
}

func daemonSetup() {
	var cmd string
	if flag.NArg() > 0 {
		cmd = flag.Args()[0]
	}

	if cmd == "" && !daemon.WasReborn() {
		printVersion()
		return
	}

	daemon.AddCommand(daemon.StringFlag(&cmd, "stop"), syscall.SIGTERM, onQuit)
	daemon.AddCommand(daemon.StringFlag(&cmd, "quit"), syscall.SIGQUIT, onQuit)
	daemon.AddCommand(daemon.StringFlag(&cmd, "reload"), syscall.SIGHUP, onHup)

	if len(daemon.ActiveFlags()) > 0 {
		d, err := ctx.Search()
		if err != nil {
			log.Fatalln(err)
		}
		daemon.SendCommands(d)
		return
	}

	dispatch(cmd)
}

func daemonStart() {
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
	log.Println("signal:", sig)
	termHandler(sig == syscall.SIGTERM)
	return daemon.ErrStop
}

func onHup(sig os.Signal) error {
	reload()
	return nil
}
