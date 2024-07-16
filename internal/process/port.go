package process

import (
    "fmt"
    "github.com/d3code/clog"
    "os/exec"
    "strings"
)

// KillPortProcess kills the process using the given port
func KillPortProcess(port int) {
    // Run the lsof and awk commands to get the PID
    lsofCommand := fmt.Sprintf("(lsof -i :%v | awk 'NR==2 {print $2}')", port)
    lsof := exec.Command("sh", "-c", lsofCommand)

    o, err := lsof.Output()
    if err != nil {
        fmt.Printf("Error running commands: %v\n", err)
        return
    }

    // Trim any whitespace from the output and extract the PID
    pid := strings.TrimSpace(string(o))

    if pid == "" {
        clog.Warnf("No process running on port %v", port)
        return
    }

    clog.Infof("PID: %s", pid)

    // Kill the process
    killCommand := fmt.Sprintf("(kill -9 %s)", o)
    kill := exec.Command("sh", "-c", killCommand)
    err = kill.Run()
    if err != nil {
        fmt.Printf("Error killing process: %v\n", err)
        return
    }
}
