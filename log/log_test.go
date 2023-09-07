package log

import (
	"strings"
	"testing"
)

func TestSetupLogging(t *testing.T) {
	if err := SetupLogging(""); err != nil {
		t.Errorf("SetupLogging() returned an error: %v", err)
	}

}

func TestDefaultPath(t *testing.T) {
	path, err := defaultPath()
	if err != nil {
		t.Errorf("defaultPath() returned an error: %v", err)
	}
	if !strings.HasSuffix(path, "containerscale/containerscale.log") {
		t.Errorf("defaultPath() returned an unexpected path: %v", path)
	}
}

func TestWriters(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		expectErr bool
	}{
		{
			name:      "success",
			path:      "/tmp/containerscale.log",
			expectErr: false,
		},
		{
			name:      "fail",
			path:      "/",
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := writers(tc.path)
			if (err != nil) != tc.expectErr {
				t.Errorf("writers() expected error: %v, got %v", tc.expectErr, err)
			}
		})
	}
}
