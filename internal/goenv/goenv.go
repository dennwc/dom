//+build !js

package goenv

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
)

func GOROOT() string {
	if s := os.Getenv("GOROOT"); s != "" {
		return s
	}
	out, err := exec.Command("go", "env", "GOROOT").Output()
	if err != nil {
		panic(err)
	}
	return string(bytes.TrimSpace(out))
}

func Go() string {
	return filepath.Join(GOROOT(), "bin/go")
}
