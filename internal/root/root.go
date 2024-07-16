package root

import (
    "github.com/d3code/clog"
    "github.com/d3code/clog/color"
    "github.com/d3code/go-reload/internal/cfg"
    "github.com/d3code/go-reload/internal/command"
    "github.com/d3code/go-reload/internal/embed_text"
    "github.com/d3code/go-reload/internal/process"
    "github.com/d3code/go-reload/internal/watch"
    "github.com/fsnotify/fsnotify"
    "github.com/spf13/cobra"
    "os"
)

func init() {
    Root.Flags().BoolP("verbose", "v", false, "show additional information about command execution")

    Root.Flags().IntSliceP("port", "p", []int{}, "kill processes running on port")
    Root.Flags().StringSliceP("directory", "d", []string{"."}, "directory to watch")
    Root.Flags().StringSliceP("extension", "e", []string{"."}, "extension to watch")
    Root.Flags().StringSliceP("ignore", "i", []string{".git"}, "files or sub-directories to ignore")
    Root.Flags().StringSliceP("run", "r", []string{}, "command to run and restart on file change")
}

var Root = &cobra.Command{
    Use:  "grl",
    Long: color.Template(embed_text.Root),
    Run:  Run,
}

func Run(cmd *cobra.Command, args []string) {
    clog.Info(color.Template(embed_text.Root))

    // Get configuration
    cfg.GetConfiguration(cmd)

    // Kill processes running on port
    for _, x := range cfg.Config.Port {
        clog.Infof("Killing process running on port %d", x)
        process.KillPortProcess(x)
    }

    // Create watcher
    watcher, err := fsnotify.NewWatcher()
    defer watch.CloseWatcher(watcher)
    if err != nil {
        clog.Error(err.Error())
        return
    }

    // Add directories to watcher
    for _, x := range cfg.Config.Directory {
        err = watch.AddDirectory(x, watcher)
        if err != nil {
            clog.Warn(err.Error())
        }
    }

    errorCh := make(chan error)
    commandCh := make(chan bool)

    // Watch for changes
    go watch.Watch(watcher, commandCh, errorCh)

    // Run commands
    go command.Command(commandCh, errorCh)

    // Run
    commandCh <- true

    // Wait for errors
    for {
        select {
        case x := <-errorCh:
            clog.Error(x.Error())
            os.Exit(1)
        }
    }
}
