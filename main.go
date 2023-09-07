package main

import (
	clog "github.com/aaomidi/containerscale/log"
)

func main() {
	if err := clog.SetupLogging(""); err != nil {
		panic("unable to setup logging")
	}
}
