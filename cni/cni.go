package cni

import (
	"encoding/json"
	"fmt"
	"github.com/aaomidi/containerscale/secret"
	"github.com/aaomidi/containerscale/ts"
	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/types"
	current "github.com/containernetworking/cni/pkg/types/100"
	"github.com/containernetworking/cni/pkg/version"
	"log"
	"time"
)

type Config struct {
	types.NetConf

	RuntimeConfig RuntimeConfig `json:"runtimeConfig"`
}

type RuntimeConfig struct {
	AuthKey secret.PrivateString `json:"authKey"`
}

type cniFunc func(_ *skel.CmdArgs) error

func Enable() {
	// Wraps the `cniFunc`s with logs.
	withErr := func(action string, cb cniFunc) cniFunc {
		return func(args *skel.CmdArgs) error {
			log.Printf("cmdAdd for %s", args.ContainerID)
			if err := cb(args); err != nil {
				log.Printf("%s-%s errored: %v", action, args.ContainerID, err)
				return err
			}
			time.Sleep(1 * time.Second)
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

	if err := ts.StartSession(input.ContainerID, input.Netns, config.RuntimeConfig.AuthKey); err != nil {
		return fmt.Errorf("unable to start session: %v", err)
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
