package internal

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

type NotifyFile struct {
	watch *fsnotify.Watcher
	Path  chan ActionPath
}

// ActionPath 文件操作
type ActionPath struct {
	Path       string
	ActionType fsnotify.Op
	Desc       string
	SourcePath string
	TargetPath string
}

func NewNotifyFile() *NotifyFile {
	n := new(NotifyFile)
	n.watch, _ = fsnotify.NewWatcher()
	n.Path = make(chan ActionPath, 10)
	return n
}

// WatchDir 监控目录
func (n *NotifyFile) WatchDir(dir, target string) {
	// 通过walk便利目录下的所有子目录
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {

		if info.IsDir() {
			abs, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			err = n.watch.Add(abs)
			if err != nil {
				return err
			}
			fmt.Printf("watching: %s\n", path)
		}
		return nil
	})
	go n.WatchEvent(dir, target)
}

func (n *NotifyFile) WatchEvent(dir, target string) {
	for {
		select {
		case event := <-n.watch.Events:
			{
				if event.Op&fsnotify.Create == fsnotify.Create {
					fmt.Printf("New file: %s\n", event.Name)
					// If the new file is directory, add watching
					info, err := os.Stat(event.Name)
					if err == nil && info.IsDir() {
						n.watch.Add(event.Name)
						fmt.Printf("Add watch: %s\n", event.Name)
					}
					go n.PushEventChannel(event.Name, fsnotify.Create, "Add watch", dir, target)
				}
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				fmt.Printf("Writing file: %s\n", event.Name)
				go n.PushEventChannel(event.Name, fsnotify.Write, "Writing file", dir, target)
			}
			if event.Op&fsnotify.Remove == fsnotify.Remove {
				fmt.Printf("Removing file: %s\n", event.Name)
				// If the file being deleted is a directory, remove watching
				info, err := os.Stat(event.Name)
				if err == nil && info.IsDir() {
					n.watch.Remove(event.Name)
					fmt.Printf("Remove watch: %s\n", event.Name)
				}
			}
			if event.Op&fsnotify.Rename == fsnotify.Rename {
				// Because we cannot use `os.Stat` to see if the file being renamed is a directory,
				// so we simply remove watching
				fmt.Printf("Renaming file: %s\n", event.Name)
				n.watch.Remove(event.Name)
			}
			if event.Op&fsnotify.Chmod == fsnotify.Chmod {
				fmt.Printf("Chmod: %s\n", event.Name)
			}

		case err := <-n.watch.Errors:
			{
				fmt.Printf("error: %v\n", err)
				return
			}

		}
	}

}

func (n *NotifyFile) PushEventChannel(path string, actionType fsnotify.Op, desc, source, target string) {
	n.Path <- ActionPath{
		Path:       path,
		ActionType: actionType,
		Desc:       desc,
		SourcePath: source,
		TargetPath: target,
	}
}
