package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"com.kvs/deviceorbit/controller/cmd/docker"
	"com.kvs/deviceorbit/controller/cmd/k8s"
	"com.kvs/deviceorbit/controller/cmd/server"
	"github.com/golang/glog"
)

func main() {
	flag.Parse()
	glog.Info("starting mobile-device-controller")

	httpPort := os.Getenv("API_PORT")
	grpcPort := os.Getenv("RPC_PORT")

	glog.Info("exposed http port: ", httpPort)
	glog.Info("exposed grpc port: ", grpcPort)

	ctx := context.Background()
	kubeClient, err := k8s.KubernetesClient(ctx)
	if err != nil {
		glog.Errorf(err.Error())
		os.Exit(1)
	}

	// to track 2 server goroutines
	var wg sync.WaitGroup
	wg.Add(2)

	containerName := os.Getenv("CONTAINER_NAME")
	containerTag := os.Getenv("CONTAINER_TAG")
	containerCommand := os.Getenv("CONTAINER_COMMAND")
	containerArgs := os.Getenv("CONTAINER_ARGS")

	containerConfig := docker.NewContainerConfig(containerName, containerTag, strings.Split(containerCommand, ","), strings.Split(containerArgs, ","))
	podProvider := k8s.DefaultPodProvider(containerConfig)
	podController := k8s.NewPodController(ctx, kubeClient, "default", podProvider)

	httpServer := server.NewHttpServer(podController)

	// start Http Server
	go func(httpPort string) {
		defer wg.Done()

		err := httpServer.StartHttpServer(httpPort)
		if err != nil {
			glog.Errorf(err.Error())
			os.Exit(1)
		}
		glog.Infof("started http server on port %s", httpPort)
	}(httpPort)

	grpcServer := server.NewGrpcServer(podController)

	// start gRPC server
	go func(grpcPort string) {
		defer wg.Done()

		err := grpcServer.StartGrpcServer(grpcPort)
		if err != nil {
			glog.Errorf(err.Error())
			os.Exit(1)
		}
		glog.Infof("started grpc server on port %s", grpcPort)
	}(grpcPort)

	// to gracefully shutdown servers handle OS signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	glog.Info("shutting down services")

	err = grpcServer.ShutdownGrpcServer()
	if err != nil {
		glog.Errorf(err.Error())
	}

	err = httpServer.ShutdownHttpServer(ctx)
	if err != nil {
		glog.Errorf(err.Error())
	}

	wg.Wait()
	glog.Info("shutdown complete")
	os.Exit(0)
}
