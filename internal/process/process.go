package process

import (
	"github.com/d3code/xlog"
	"os/exec"
	"sync"
	"syscall"
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

		processGroupId, err := syscall.Getpgid(process.Process.Pid)
		if err != nil {
			xlog.Debugf("no process group for pid [ %v ]: %v", process.Process.Pid, err)
			continue
		}

		xlog.Debugf("killing process group [ %v ]", processGroupId)
		err = syscall.Kill(-processGroupId, syscall.SIGTERM)
		if err != nil {
			xlog.Error(err.Error())
		}

		// Wait for process to exit
		err = process.Wait()
		if err != nil {
			xlog.Debug(err.Error())
		}
	}

	runningProcesses = nil
}

func RemoveProcess(process *exec.Cmd) {
	runningProcessesMutex.Lock()
	defer runningProcessesMutex.Unlock()
	for i, p := range runningProcesses {
		if p == process {
			runningProcesses = append(runningProcesses[:i], runningProcesses[i+1:]...)
		}
	}
}
