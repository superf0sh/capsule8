// Copyright 2017 Capsule8 Inc. All rights reserved.

package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/capsule8/reactive8/pkg/sensor"
	"github.com/capsule8/reactive8/pkg/version"
	"github.com/golang/glog"
)

func main() {
	// Set "alsologtostderr" flag so that glog messages go stderr as well as /tmp.
	flag.Set("alsologtostderr", "true")
	flag.Parse()

	// Log version and build at "Starting ..." for debugging
	version.InitialBuildLog("sensor")

	s, err := sensor.GetSensor()
	if err != nil {
		glog.Fatalf("Couldn't create Sensor: %s", err)
	}

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		glog.Infof("Caught signal %v, stopping Sensor", sig)
		s.Stop()
	}()

	err = s.Serve()
	if err != nil {
		glog.Error(err)
	}

	glog.Info("Exiting...")
}
