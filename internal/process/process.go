package process

import (
	"errors"
	"fmt"
	"github.com/d3code/xlog"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

var (
	runningProcesses      []*exec.Cmd
	runningProcessesMutex = sync.Mutex{}
)

func AddProcess(process *exec.Cmd) {
	runningProcessesMutex.Lock()
	defer runningProcessesMutex.Unlock()

	runningProcesses = append(runningProcesses, process)
}

func KillAllProcessGroups() {
	runningProcessesMutex.Lock()
	defer runningProcessesMutex.Unlock()

	if len(runningProcesses) == 0 {
		return
	}

	xlog.Infof("Killing %v processes", len(runningProcesses))
	for _, process := range runningProcesses {
		if process == nil || process.Process == nil {
			continue
		}

		groupProcesses := getGroupProcesses(process.Process.Pid)
		xlog.Tracef("Group processes: %v", groupProcesses)

		for _, pid := range groupProcesses {
			xlog.Tracef("Killing group process %d", pid)
			err := syscall.Kill(pid, syscall.SIGTERM)
			if err != nil && err.Error() != "no such process" {
				xlog.Error(err.Error())
			}

			waitForProcessToExit(pid)
		}
	}

	runningProcesses = nil
}

func waitForProcessToExit(pid int) {
	// Check if the process has quit every 200ms
	tick := 200
	for {
		running, runningErr := isProcessRunning(pid)
		if runningErr != nil {
			xlog.Errorf("Error checking if process is running: %v", runningErr)
			return
		}

		if !running {
			xlog.Tracef("Process with PID %d killed", pid)
			return
		}

		// Sleep for tick milliseconds
		<-time.After(time.Duration(tick) * time.Millisecond)
	}

}

func getGroupProcesses(pid int) []int {
	var childProcesses []int
	out, err := exec.Command("pgrep", "-g", fmt.Sprintf("%d", pid)).Output()
	if err != nil {
		xlog.Warnf("Failed to get group processes: %v", err)
		childProcesses = append(childProcesses, pid)
		return childProcesses
	}

	for _, pidStr := range strings.Fields(string(out)) {
		childProcess, errConvert := strconv.Atoi(pidStr)
		if errConvert == nil {
			childProcesses = append(childProcesses, childProcess)
		}
	}

	return childProcesses
}

func isProcessRunning(pid int) (bool, error) {
	run := exec.Command("sh", "-c", fmt.Sprintf("ps -p %d -o stat=", pid))
	o, err := run.Output()
	if err != nil {
		// If the error is due to the process not existing, return false without an error
		var exitError *exec.ExitError
		if errors.As(err, &exitError) && exitError.ExitCode() == 1 {
			return false, nil
		}
		return false, err
	}

	// Check if the output contains the status
	status := strings.TrimSpace(string(o))
	if strings.Contains(status, "Z") {
		xlog.Tracef("Process with PID %d is a zombie process", pid)
		return false, nil
	}

	return true, nil
}
