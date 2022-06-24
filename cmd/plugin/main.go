/*
* Copyright (c) 2022 Intel Corporation.
*
* SPDX-License-Identifier: Apache-2.0
*
 */
package main

import (
	_ "ep/pkg/epplugins"
	plugin "ep/pkg/plugin"
	"errors"
	"flag"
	"os"

	log "github.com/sirupsen/logrus"
)

var (
	errMaxResource = errors.New("Exceed max resource ")
)

//nolint:unused,deadcode
const (
	ADDRESS   = ":50088"
	MAX_COUNT = 100
)

func main() {
	flag.Parse()

	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: false,
	})
	lens := len(flag.Args())
	if lens > MAX_COUNT {
		log.Errorln(errMaxResource)
		os.Exit(1)
	}
	for _, p := range flag.Args() {
		log.Infof("Enable plugin remote Log: %v\n", p)
		if err := plugin.EnablePluginRemoteLog(p); err != nil {
			log.Fatal(err)
		}
		log.Infof("Start Plugin: %v\n", p)
		if err := plugin.StartPlugin(p, nil); err != nil {
			log.Errorln(err)
			os.Exit(1)
		}
	}
	for _, p := range flag.Args() {
		if err := plugin.WaitPluginFinished(p); err != nil {
			log.Errorln(err)
			os.Exit(1)
		}
		log.Infof("Plugin Finished: %v\n", p)
	}
}
