package cni

import (
	"encoding/json"
	"fmt"
	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/types"
	current "github.com/containernetworking/cni/pkg/types/100"
	"github.com/containernetworking/cni/pkg/version"
	"log"
)

type Config struct {
	types.NetConf

	RuntimeConfig RuntimeConfig `json:"runtimeConfig"`
}

type PrivateString string

func (PrivateString) String() string {
	return "private"
}

type RuntimeConfig struct {
	AuthKey PrivateString `json:"authKey"`
}

type cniFunc func(_ *skel.CmdArgs) error

func Enable() {
	// Wraps the cniFuncs with logs.
	withErr := func(action string, cb cniFunc) cniFunc {
		return func(args *skel.CmdArgs) error {
			log.Printf("cmdAdd for %s", args.ContainerID)
			if err := cb(args); err != nil {
				log.Printf("%s-%s errored: %v", action, args.ContainerID, err)
				return err
			}
			return nil
		}
	}

	skel.PluginMain(
		withErr("add", cmdAdd),
		withErr("check", cmdCheck),
		withErr("del", cmdDel),
		version.All, "containerscale",
	)
}

func cmdAdd(input *skel.CmdArgs) error {
	log.Printf("cmdAdd for %s", input.ContainerID)

	config, err := parseConfig(input.StdinData)
	if err != nil {
		return err
	}

	return types.PrintResult(&current.Result{}, config.CNIVersion)
}

func cmdCheck(_ *skel.CmdArgs) error {
	return types.PrintResult(&current.Result{}, "")
}

func cmdDel(_ *skel.CmdArgs) error {
	return types.PrintResult(&current.Result{}, "")
}

func parseConfig(stdin []byte) (*Config, error) {
	conf := Config{}

	if err := json.Unmarshal(stdin, &conf); err != nil {
		return nil, fmt.Errorf("failed to load plugin configuration: %v", err)
	}

	if conf.RawPrevResult != nil {
		if err := version.ParsePrevResult(&conf.NetConf); err != nil {
			return nil, fmt.Errorf("failed to parse prevResult: %v", err)
		}

		if _, err := current.NewResultFromResult(conf.PrevResult); err != nil {
			return nil, fmt.Errorf("could not convert previous result to current version: %v", err)
		}
	}

	return &conf, nil
}
