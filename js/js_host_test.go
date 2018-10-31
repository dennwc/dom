//+build !js

package js

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/dennwc/testproxy"

	"github.com/dennwc/dom/internal/goenv"
	"github.com/stretchr/testify/require"
)

func TestJS(t *testing.T) {
	// compile js tests
	GOROOT := goenv.GOROOT()
	f, err := ioutil.TempFile("", "js_test_")
	require.NoError(t, err)
	f.Close()
	defer os.Remove(f.Name())

	cmd := exec.Command(goenv.Go(), "test", "-c", "-o", f.Name(), ".")
	cmd.Env = append(cmd.Env, os.Environ()...)
	cmd.Env = append(cmd.Env, []string{
		"GOROOT=" + GOROOT,
		"GOARCH=wasm",
		"GOOS=js",
	}...)
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	require.NoError(t, err)
	err = os.Chmod(f.Name(), 0755)
	require.NoError(t, err)

	// run tests with Node.js
	wd, err := os.Getwd()
	require.NoError(t, err)
	cmd = exec.Command(filepath.Join(GOROOT, "misc/wasm/go_js_wasm_exec"), f.Name())
	cmd.Dir = wd
	testproxy.RunTestBinary(t, cmd)
}
