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
	"time"
)

func Watch(watcher *fsnotify.Watcher, build chan bool, errors chan error) {
	var changes bool
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			if changes {
				xlog.Debug("Changes detected")
				build <- true
				changes = false
			}
		}
	}()

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

			AddCreatedDirectory(event, watcher)

			for _, x := range cfg.Config.Extension {
				if strings.HasSuffix(event.Name, x) {
					xlog.Debug(event.String())
					changes = true
				}
			}
		case errWatcher, ok := <-watcher.Errors:
			if !ok {
				errors <- fmt.Errorf("watcher error")
				continue
			}
			if errWatcher != nil {
				errors <- fmt.Errorf("watcher error: %s", errWatcher.Error())
			}
		}
	}
}

// AddCreatedDirectory adds a directory to the watcher if it was created
func AddCreatedDirectory(event fsnotify.Event, watcher *fsnotify.Watcher) {
	if event.Op&fsnotify.Create == fsnotify.Create {
		info, err := os.Stat(event.Name)
		if err != nil {
			xlog.Error(err.Error())
			return
		}
		if info.IsDir() {
			err = watcher.Add(event.Name)
			if err != nil {
				xlog.Error(err.Error())
			} else {
				xlog.Infof("Watching directory: %s", event.Name)
			}
		}
	}
}

// SetWatchDirectory walks through a directory and adds all subdirectories to the watcher
func SetWatchDirectory(dir string, watcher *fsnotify.Watcher) error {
	fn := func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if shouldIgnore(p) {
				return nil
			}

			err = watcher.Add(p)
			if err != nil {
				return err
			}

			xlog.Tracef("Watching directory [%s]", p)
		}
		return nil
	}

	return filepath.WalkDir(dir, fn)
}

func shouldIgnore(path string) bool {
	parts := strings.Split(path, string(os.PathSeparator))
	for _, part := range parts {
		for _, ignoreDirectory := range cfg.Config.Ignore {
			if part == ignoreDirectory {
				return true
			}
		}
	}
	return false
}

func CloseWatcher(watcher *fsnotify.Watcher) {
	if watcher != nil {
		err := watcher.Close()
		if err != nil {
			xlog.Error(err.Error())
		}
	}
}
