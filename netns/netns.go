package netns

import (
	"fmt"
	"golang.org/x/sys/unix"
)

const InvalidNsHandle NsHandle = -999

type NsHandle int

func (ns NsHandle) validate() error {
	if ns == InvalidNsHandle {
		return fmt.Errorf("operation on invalid ns handle")
	}
	return nil
}

func (ns NsHandle) Set() error {
	if err := ns.validate(); err != nil {
		return err
	}

	if err := unix.Setns(int(ns), unix.CLONE_NEWNET); err != nil {
		return fmt.Errorf("could not set network namespace: %v", err)
	}

	return nil
}

func Path(path string) (NsHandle, error) {
	fd, err := unix.Open(path, unix.O_CLOEXEC|unix.O_RDONLY, 0)
	if err != nil {
		return InvalidNsHandle, fmt.Errorf("could not open network namespace %s: %v", path, err)
	}

	return NsHandle(fd), nil
}
