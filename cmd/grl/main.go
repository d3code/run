package main

import (
    "github.com/d3code/clog"
    "github.com/d3code/go-reload/internal/process"
    "github.com/d3code/go-reload/internal/root"
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
            clog.Error(err.Error())
            os.Exit(1)
        }
    }()

    sig := <-cancelChan
    clog.Debugf("received signal %s", sig.String())
    process.KillAllProcessGroups()
}
