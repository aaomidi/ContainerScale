package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
)

func SetupLogging(path string) (err error) {
	if path == "" {
		path, err = defaultPath()
		if err != nil {
			return err
		}
	}

	out, err := writer(path)
	if err != nil {
		return err
	}

	log.SetOutput(out)
	return nil
}

func writer(p string) (io.Writer, error) {
	if err := os.MkdirAll(path.Dir(p), 0777); err != nil {
		return nil, fmt.Errorf("could not create containerscale directory: %v", err)
	}
	file, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("could not open log file: %v", err)
	}

	return file, nil
}

func defaultPath() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("could not get cache directory: %v", err)
	}
	return path.Join(cacheDir, "containerscale", "containerscale.log"), nil
}
