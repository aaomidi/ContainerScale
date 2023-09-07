package cli

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
)

func LoggingCommand(cmd *exec.Cmd) *exec.Cmd {
	cmd.Stdout = loggingWriter()
	cmd.Stderr = loggingWriter()

	return cmd
}

func Execute(cmd *exec.Cmd) error {
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("could not start command: %v", err)
	}

	return nil
}

func loggingWriter() io.Writer {
	reader, writer := io.Pipe()

	go func() {
		scanner := bufio.NewScanner(reader)

		for scanner.Scan() {
			log.Println(scanner.Text())
		}
	}()

	return writer
}
