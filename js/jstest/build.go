//+build !js

package jstest

import (
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"github.com/dennwc/dom/internal/goenv"
)

const (
	wasmDir      = "misc/wasm"
	wasmExec     = "go_js_wasm_exec"
	wasmExecJS   = "wasm_exec.js"
	wasmExecHTML = "wasm_exec.html"
)

func buildTestJS(t testing.TB) string {
	// compile js tests
	GOROOT := goenv.GOROOT()
	f, err := ioutil.TempFile("", "js_test_")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	cmd := exec.Command(goenv.Go(), "test", "-c", "-o", f.Name(), ".")
	cmd.Env = append(cmd.Env, os.Environ()...)
	cmd.Env = append(cmd.Env, []string{
		"GOROOT=" + GOROOT,
		"GOARCH=wasm",
		"GOOS=js",
	}...)
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err == nil {
		err = os.Chmod(f.Name(), 0755)
	}
	if err != nil {
		os.Remove(f.Name())
		t.Fatal(err)
	}
	return f.Name()
}

func copyFile(t testing.TB, src, dst string) {
	r, err := os.Open(src)
	require.NoError(t, err)
	defer r.Close()

	w, err := os.Create(dst)
	require.NoError(t, err)
	defer w.Close()

	_, err = io.Copy(w, r)
	require.NoError(t, err)

	fi, err := r.Stat()
	require.NoError(t, err)

	err = w.Chmod(fi.Mode())
	require.NoError(t, err)

	err = w.Close()
	require.NoError(t, err)
}
