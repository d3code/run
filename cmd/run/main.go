package main

import (
	"fmt"
	"github.com/d3code/run/internal/process"
	"github.com/d3code/run/internal/root"
	"github.com/d3code/xlog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		err := root.Root.Execute()
		if err != nil {
			xlog.Error(err.Error())
			os.Exit(1)
		}
	}()

	sig := <-cancelChan
	fmt.Println()
	xlog.Warnf("Shutting down due to signal [%v]", sig)
	process.KillAllProcessGroups()
}
