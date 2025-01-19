package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"com.kvs/deviceorbit/plugin/cmd/client"
	"com.kvs/deviceorbit/plugin/cmd/resource"
	"com.kvs/deviceorbit/plugin/cmd/watcher"
	"github.com/golang/glog"

	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: mobile-device-plugin\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func init() {
	flag.Usage = usage
	// NOTE: This next line is key you have to call flag.Parse() for the command line
	// options or "flags" that are defined in the glog module to be picked up.
	// flag.StringVar(&confFileName, "config", "config/conf.yaml", "set the configuration file to use")
	flag.Parse()
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func newOSWatcher(sigs ...os.Signal) chan os.Signal {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, sigs...)

	return sigChan
}

func main() {
	defer glog.Flush()

	glog.V(0).Info("starting FS watcher.")
	fileWatcher := watcher.NewFileWatcher()

	err := fileWatcher.Watch(pluginapi.DevicePluginPath, pluginapi.KubeletSocket, 100*time.Millisecond)
	if err != nil {
		glog.Errorf(err.Error())
		os.Exit(1)
	}

	defer fileWatcher.Close()

	glog.V(0).Info("starting OS watcher.")
	sigs := newOSWatcher(syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	ctx := context.Background()
	client, err := client.NewGrpcCLient(ctx)
	if err != nil {
		glog.Errorf(err.Error())
		os.Exit(1)
	}

	defer client.Close()

	glog.V(0).Info("starting device regsitry.")
	deviceCount, err := strconv.Atoi(os.Getenv("DEVICES_REQUESTED"))
	if err != nil {
		glog.Warningf("fallback to default requse amount: 1, cause", err.Error())
		deviceCount = 1
	}

	deviceRegsitry := resource.NewDeviceRegistry(client, uint(deviceCount))
	err = deviceRegsitry.RunDevices(ctx)
	if err != nil {
		glog.Warningf("found 0 devices", err.Error())
	}

	deviceStrings := listRunningDevices(deviceRegsitry)
	glog.V(0).Infof("starting serving device resources:\n%s", strings.Join(deviceStrings, "\n"))

	watcherEvents := fileWatcher.Events()
	watcherError := fileWatcher.Errors()

	glog.V(0).Info("starting listeninng for resource events")
	for {
		select {
		case event := <-watcherEvents:
			device, err := deviceRegsitry.HandleDeviceEvent(ctx, event)
			if err != nil {
				glog.Warning(err.Error())
			}
			if !deviceRegsitry.IsEmptyResource(device) {
				glog.V(0).Infof("device: %s event: %s", device, event.Op)

				deviceStrings := listRunningDevices(deviceRegsitry)
				glog.V(0).Infof("updated device resources:\n%s", strings.Join(deviceStrings, "\n"))
			}

		case err := <-watcherError:
			glog.V(0).Infof("inotify: %s", err)

		case s := <-sigs:
			switch s {
			case syscall.SIGHUP:
				glog.V(0).Info("received SIGHUP, restarting.")
				err := deviceRegsitry.RunDevices(ctx)
				if err != nil {
					glog.Error(err.Error())
				}
			default:
				glog.V(0).Infof("received signal \"%v\", shutting down.", s)
				err := deviceRegsitry.StopDevices()
				if err != nil {
					glog.Error(err.Error())
				}
				break
			}
		}
	}
}

func listRunningDevices(deviceRegsitry *resource.DeviceRegistry) []string {
	var deviceStrings []string
	for _, d := range deviceRegsitry.ListRunningDevices() {
		deviceStrings = append(deviceStrings, d.String())
	}
	return deviceStrings
}
