// Package cli wraps around command execution peculiarities.
package cli

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os/exec"
)

func LoggingCommand(ctx context.Context, cmd *exec.Cmd) *exec.Cmd {
	cmd.Stdout = loggingWriter(ctx)
	cmd.Stderr = loggingWriter(ctx)

	return cmd
}

func Execute(cmd *exec.Cmd, wait bool) error {
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("could not start command: %v", err)
	}

	if wait {
		if err := cmd.Wait(); err != nil {
			return fmt.Errorf("error when waiting for command: %v", err)
		}
	}

	return nil
}

func loggingWriter(ctx context.Context) io.Writer {
	reader, writer := io.Pipe()
	scanner := bufio.NewScanner(reader)
	go func() {
		defer writer.Close()

		for {
			select {
			// Handle cancels, timeouts, etc.
			case <-ctx.Done():
				return
			default:
				// Technically this can cause the outer for loop from executing, but eh, probably not.
				// Either way the caller is going to exit(0) before context gets called anyway.
				for scanner.Scan() {
					log.Println(scanner.Text())
				}
			}
		}
	}()

	return writer
}
