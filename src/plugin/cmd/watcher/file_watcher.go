package watcher

import (
	"time"

	"github.com/radovskyb/watcher"
)

const (
	Create = watcher.Create
	Write  = watcher.Write
	Remove = watcher.Remove

	usbDeviceDir = "/dev/bus/usb"
)

type FileWatcher struct {
	watcher *watcher.Watcher
}

func NewFileWatcher() *FileWatcher {
	return &FileWatcher{
		watcher: watcher.New(),
	}
}

func (fw *FileWatcher) Watch(pluginPath, kubeletPath string, timeout time.Duration) error {
	fw.watcher.FilterOps(watcher.Create, watcher.Remove, watcher.Write)

	if err := fw.watcher.Add(pluginPath); err != nil {
		return err
	}

	if err := fw.watcher.Add(kubeletPath); err != nil {
		return err
	}

	if err := fw.watcher.AddRecursive(usbDeviceDir); err != nil {
		return err
	}

	go func() {
		_ = fw.watcher.Start(timeout)
	}()

	fw.watcher.Wait()

	return nil
}

func (fw *FileWatcher) Events() <-chan watcher.Event {
	return fw.watcher.Event
}

func (fw *FileWatcher) Errors() <-chan error {
	return fw.watcher.Error
}

func (fw *FileWatcher) Close() {
	fw.watcher.Close()
}
