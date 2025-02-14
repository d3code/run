package process

import (
	"errors"
	"fmt"
	"github.com/d3code/xlog"
	"os/exec"
	"strings"
	"time"
)

// KillPortProcess kills the process using the given port
func KillPortProcess(port int) {
	pid, err := getPortProcess(port)
	if err != nil {
		xlog.Error(err.Error())
		return
	}

	if pid == "" {
		xlog.Debugf("No process running on port %v", port)
		return
	}

	xlog.Debugf("Process with PID %s running on port %d", pid, port)
	killCommand := fmt.Sprintf("(kill -9 %s)", pid)
	kill := exec.Command("sh", "-c", killCommand)

	err = kill.Run()
	if err != nil {
		xlog.Errorf("Error killing process: %v", err)
		return
	}

	// Check if the process has quit every 200ms, timeout after 3 seconds
	timeout := time.After(3 * time.Second)
	tick := time.Tick(200 * time.Millisecond)

	for {
		select {
		case <-timeout:
			xlog.Errorf("Timeout waiting for process with PID %s to quit", pid)
			return
		case <-tick:
			// Check if the process is still running
			running, runningErr := isProcessRunning(pid)
			if runningErr != nil {
				xlog.Errorf("Error checking if process is running: %v", runningErr)
				return
			}

			if !running {
				xlog.Infof("Process with PID %s killed", pid)
				return
			}
		}
	}

}

func isProcessRunning(pid string) (bool, error) {
	run := exec.Command("sh", "-c", fmt.Sprintf("ps -p %s", pid))
	o, err := run.Output()
	if err != nil {
		// If the error is due to the process not existing, return false without an error
		var exitError *exec.ExitError
		if errors.As(err, &exitError) && exitError.ExitCode() == 1 {
			return false, nil
		}
		return false, err
	}

	// Check if the output contains the PID
	if strings.Contains(string(o), pid) {
		return true, nil
	}

	return false, nil
}

func getPortProcess(port int) (string, error) {
	lsofCommand := fmt.Sprintf("(lsof -i :%d | awk 'NR==2 {print $2}')", port)
	lsof := exec.Command("sh", "-c", lsofCommand)

	o, err := lsof.Output()
	if err != nil {
		return "", err
	}

	pid := strings.TrimSpace(string(o))
	return pid, err
}
