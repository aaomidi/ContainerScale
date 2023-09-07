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

func StartSession(containerID, nsPath string, authKey secret.PrivateString) error {
	var err error
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = func() error {
			log.Println("Starting tailscale session")
			runtime.LockOSThread()

			nsHandle, err := netns.Path(nsPath)
			if err != nil {
				return fmt.Errorf("failed to get ns: %v", err)
			}
			if err := nsHandle.Set(); err != nil {
				return fmt.Errorf("failed to set ns: %v", err)
			}

			if err := tailscaled(containerID); err != nil {
				return fmt.Errorf("failed to start tailscaled: %v", err)
			}

			if err := tailscale(containerID, authKey); err != nil {
				return fmt.Errorf("failed to start tailscale: %v", err)
			}

			return nil
		}()
	}()
	wg.Wait()
	return err
}

func tailscale(containerID string, authKey secret.PrivateString) error {
	socket := socket(containerID)
	cmd := exec.Command("tailscale", "--socket", socket, "up", "--authkey", authKey.AccessPrivateString())
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	cli.LoggingCommand(cmd)

	if err := cli.Execute(cmd); err != nil {
		return err
	}
	return nil
}

func tailscaled(containerID string) error {
	socket := socket(containerID)
	state := path.Join(dir(containerID), "tailscaled.state")

	cmd := exec.Command("tailscaled", "--socket", socket, "--statedir", state)
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
