package command

import (
    "fmt"
    "github.com/d3code/clog"
    "github.com/d3code/go-reload/internal/cfg"
    "github.com/d3code/go-reload/internal/process"
    "os"
    "os/exec"
    "syscall"
)

func Command(commandCh chan bool, errorCh chan error) {
    for {
        select {
        case <-commandCh:

            clog.Infof("\n{{ Restarting... | green }}")

            // Checking if any process is running
            process.KillAllProcessGroups()

            // Execute commands
            for _, x := range cfg.Config.Run {
                go ExecuteCommand(x, errorCh)
            }
        }
    }
}

func ExecuteCommand(build string, errors chan error) {
    if len(cfg.Config.Run) == 0 {
        clog.Debug("no build command specified")
        return
    }

    build = fmt.Sprintf("(%s)", build)

    c := exec.Command("sh", "-c", build)
    c.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
    c.Stdout = os.Stdout
    c.Stderr = os.Stderr

    process.AddProcess(c)

    err := c.Start()
    if err != nil {
        clog.Warn(err.Error())
    }
}
