package command

import (
	"fmt"
	"github.com/d3code/run/internal/cfg"
	"github.com/d3code/run/internal/process"
	"github.com/d3code/xlog"
	"os"
	"os/exec"
	"syscall"
)

func Command(commandCh chan bool, errorCh chan error) {
	for {
		select {
		case <-commandCh:
			if len(cfg.Config.Run) == 0 {
				xlog.Fatalf("No command(s) specified")
				return
			}

			process.KillAllProcessGroups()
			for _, x := range cfg.Config.Run {
				go ExecuteCommand(x, errorCh)
			}
		}
	}
}

func ExecuteCommand(build string, errors chan error) {
	build = fmt.Sprintf("(%s)", build)
	c := exec.Command("sh", "-c", build)
	c.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	process.AddProcess(c)

	err := c.Start()
	if err != nil {
		xlog.Errorf("Error starting command [%s]: %v", build, err)
	} else {
		xlog.Infof("Command [%s] with PID %v started", build, c.Process.Pid)
	}
}
