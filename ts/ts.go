package ts

import (
	"fmt"
	"github.com/aaomidi/containerscale/cli"
	"github.com/aaomidi/containerscale/netns"
	"github.com/aaomidi/containerscale/secret"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"sync"
	"syscall"
)

type CreateSessionRequest struct {
	ContainerID          string
	NetworkNamespacePath string
	AuthKey              secret.PrivateString
	TailscaledFlags      []string
	TailscaleFlags       []string
}

func CreateSession(req CreateSessionRequest) error {
	var err error
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = func() error {
			log.Println("Starting tailscale session")
			runtime.LockOSThread()

			nsHandle, err := netns.Path(req.NetworkNamespacePath)
			if err != nil {
				return fmt.Errorf("failed to get ns: %v", err)
			}
			if err := nsHandle.Set(); err != nil {
				return fmt.Errorf("failed to set ns: %v", err)
			}

			if err := tailscaled(req.ContainerID, req.TailscaledFlags); err != nil {
				return fmt.Errorf("failed to start tailscaled: %v", err)
			}

			if err := tailscale(req.ContainerID, req.AuthKey, req.TailscaleFlags); err != nil {
				return fmt.Errorf("failed to start tailscale: %v", err)
			}

			return nil
		}()
	}()
	wg.Wait()
	return err
}

func tailscale(containerID string, authKey secret.PrivateString, flags []string) error {
	socket := socket(containerID)
	defaultFlags := []string{
		"--socket", socket,
		"up",
		"--authkey", authKey.AccessPrivateString(),
	}
	flags = append(defaultFlags, flags...)

	cmd := exec.Command("tailscale", flags...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	cli.LoggingCommand(cmd)

	if err := cli.Execute(cmd); err != nil {
		return err
	}
	return nil
}

func tailscaled(containerID string, flags []string) error {
	socket := socket(containerID)
	state := path.Join(dir(containerID), "tailscaled.state")
	defaultFlags := []string{
		"--socket", socket,
		"--statedir", state,
	}
	flags = append(defaultFlags, flags...)

	cmd := exec.Command("tailscaled", flags...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	cli.LoggingCommand(cmd)

	if err := cli.Execute(cmd); err != nil {
		return err
	}
	return nil
}

func socket(containerID string) string {
	return path.Join(dir(containerID), "tailscaled.sock")
}

func dir(containerID string) string {
	return path.Join(os.TempDir(), "containerscale", containerID)
}
