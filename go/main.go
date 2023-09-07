package main

import (
	"github.com/aaomidi/containerscale/cni"
	clog "github.com/aaomidi/containerscale/log"
)

func main() {
	// TODO: Allow the logging path to be configurable.
	if err := clog.SetupLogging(""); err != nil {
		panic("unable to setup logging")
	}
	cni.Enable()
}
