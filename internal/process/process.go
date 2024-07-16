package process

import (
    "github.com/d3code/clog"
    "os/exec"
    "sync"
    "syscall"
)

var (
    RunningProcesses      []*exec.Cmd
    RunningProcessesMutex = sync.Mutex{}
)

func AddProcess(process *exec.Cmd) {
    RunningProcessesMutex.Lock()
    defer RunningProcessesMutex.Unlock()

    RunningProcesses = append(RunningProcesses, process)
}

func KillAllProcessGroups() {
    RunningProcessesMutex.Lock()
    defer RunningProcessesMutex.Unlock()

    clog.Infof("Killing %v processes", len(RunningProcesses))

    for _, process := range RunningProcesses {
        if process == nil || process.Process == nil {
            continue
        }

        processGroupId, err := syscall.Getpgid(process.Process.Pid)
        if err != nil {
            clog.Debugf("no process group for pid [ %v ]: %v", process.Process.Pid, err)
            continue
        }

        clog.Debugf("killing process group [ %v ]", processGroupId)
        err = syscall.Kill(-processGroupId, syscall.SIGTERM)
        if err != nil {
            clog.Error(err.Error())
        }

        // Wait for process to exit
        err = process.Wait()
        if err != nil {
            clog.Debug(err.Error())
        }
    }

    RunningProcesses = nil
}

func RemoveProcess(process *exec.Cmd) {
    RunningProcessesMutex.Lock()
    defer RunningProcessesMutex.Unlock()
    for i, p := range RunningProcesses {
        if p == process {
            RunningProcesses = append(RunningProcesses[:i], RunningProcesses[i+1:]...)
        }
    }
}
