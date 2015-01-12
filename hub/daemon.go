// TODO messy code, all the clean stuff is in main.go

package main

import (
	"flag"
	"os"
	"syscall"

	"github.com/dataence/glog"
	"github.com/sevlyar/go-daemon"
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
			glog.Fatal(err)
		}
		daemon.SendCommands(d)
		return
	}

	dispatch(cmd)
}

func daemonStart() {
	d, err := ctx.Reborn()
	if err != nil {
		glog.Fatal(err)
	}
	if d == nil {
		// FIXME ctx.Release fails on MacOSX, no /proc/ for lockfile GetFdName
		// use explicit remove instead
		defer os.Remove(ctx.PidFileName)
		defer ctx.Release()

		glog.Info("- - - - - - - - - - - - - - -")
		glog.Infoln("starting", os.Args)

		go worker()

		err = daemon.ServeSignals()
		if err != nil {
			glog.Fatal(err)
		}
		glog.Info("terminated")
	}
}

func onQuit(sig os.Signal) error {
	glog.Infoln("signal:", sig)
	termHandler(sig == syscall.SIGTERM)
	return daemon.ErrStop
}

func onHup(sig os.Signal) error {
	reload()
	return nil
}
