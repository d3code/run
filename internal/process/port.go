package process

import (
	"fmt"
	"github.com/d3code/xlog"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// KillPortProcess kills the process using the given port
func KillPortProcess(port int) {
	pid, err := getPortProcess(port)
	if err != nil {
		xlog.Error(err.Error())
		return
	}

	if pid == 0 {
		xlog.Debugf("No process running on port %v", port)
		return
	}

	xlog.Debugf("Process with PID %d running on port %d", pid, port)
	err1 := syscall.Kill(pid, syscall.SIGTERM)
	if err1 != nil {
		xlog.Errorf("Error killing process: %v", err1)
		return
	}

	// Check if the process has quit every 200ms, timeout after 3 seconds
	timeout := time.After(3 * time.Second)
	tick := time.Tick(200 * time.Millisecond)

	for {
		select {
		case <-timeout:
			xlog.Errorf("Timeout waiting for process with PID %d to quit", pid)
			return
		case <-tick:
			// Check if the process is still running
			running, runningErr := isProcessRunning(pid)
			if runningErr != nil {
				xlog.Errorf("Error checking if process is running: %v", runningErr)
				return
			}

			if !running {
				xlog.Infof("Process with PID %d killed", pid)
				return
			}
		}
	}

}

func getPortProcess(port int) (int, error) {
	lsofCommand := fmt.Sprintf("(lsof -i :%d | awk 'NR==2 {print $2}')", port)
	lsof := exec.Command("sh", "-c", lsofCommand)

	output, err := lsof.Output()
	if err != nil {
		return 0, err
	}

	pid := strings.TrimSpace(string(output))
	if pid == "" {
		return 0, nil
	}

	num, err := strconv.Atoi(pid)

	return num, err
}
