package watch

import (
	"fmt"
	"github.com/d3code/run/internal/cfg"
	"github.com/d3code/xlog"
	"github.com/fsnotify/fsnotify"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func Watch(watcher *fsnotify.Watcher, build chan bool, errors chan error) {
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				errors <- fmt.Errorf("event error")
				continue
			}

			if event.Op&fsnotify.Chmod == fsnotify.Chmod {
				continue
			}

			AddFolderToWatcher(event, watcher)

			for _, x := range cfg.Config.Extension {
				if strings.HasSuffix(event.Name, x) {
					xlog.Debug(event.String())
					build <- true
				}
			}
		case errWatcher, ok := <-watcher.Errors:
			if !ok {
				errors <- fmt.Errorf("watcher error")
			}
			if errWatcher != nil {
				errors <- fmt.Errorf("watcher error: %s", errWatcher.Error())
			}
		}
	}
}

func AddFolderToWatcher(event fsnotify.Event, watcher *fsnotify.Watcher) {
	if event.Op&fsnotify.Create == fsnotify.Create {
		info, errCreate := os.Stat(event.Name)
		if info != nil && info.IsDir() {
			errCreate = watcher.Add(event.Name)
			if errCreate != nil {
				xlog.Error(errCreate.Error())
			} else {
				xlog.Infof("Watching directory: %s", event.Name)
			}
		}
	}
}

func AddDirectory(dir string, watcher *fsnotify.Watcher) error {
	fn := func(p string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			parts := strings.Split(p, string(os.PathSeparator))
			shouldIgnore := false
			for _, part := range parts {
				for _, ignoreDirectory := range cfg.Config.Ignore {
					if part == ignoreDirectory {
						shouldIgnore = true
						break
					}
				}
			}

			if !shouldIgnore {
				if cfg.Config.Verbose {
					xlog.Infof("{{ Watching directory | grey }} {{ %s | blue }}", p)
				}
				errWatch := watcher.Add(p)
				if errWatch != nil {
					return errWatch
				}
			}
		}

		return nil
	}

	return filepath.WalkDir(dir, fn)
}

func CloseWatcher(watcher *fsnotify.Watcher) {
	if watcher != nil {
		err := watcher.Close()
		if err != nil {
			xlog.Error(err.Error())
		}
	}
}
