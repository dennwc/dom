//+build !js

package jstest

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/dennwc/testproxy"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"

	"github.com/dennwc/dom/internal/goenv"
	"github.com/stretchr/testify/require"
)

const nodejsImage = "node:8-slim"

// RunTestNodeJS compiles all tests in the current working directory for WASM+JS, and runs them either in Docker or
// locally with NodeJS. It will stream test results back to t.
//
// The caller should specify the "!js" build tag, while all JS tests in the package should include "js" build tag.
func RunTestNodeJS(t *testing.T) {
	GOROOT := goenv.GOROOT()
	testFile := buildTestJS(t)
	defer os.Remove(testFile)

	if runInDocker(t, testFile, GOROOT) {
		return
	}
	t.Log("running locally")

	// run tests with Node.js
	if _, err := exec.LookPath("node"); err != nil {
		t.Skipf("cannot find NodeJS: %v", err)
	}

	wd, err := os.Getwd()
	require.NoError(t, err)
	cmd := exec.Command(filepath.Join(GOROOT, wasmDir, wasmExec), testFile)
	cmd.Dir = wd
	testproxy.RunTestBinary(t, cmd)
}

type dockerRunner struct {
	cli *docker.Client
	c   *docker.Container
}

func (r dockerRunner) RunAndWait(stdout, stderr io.Writer) error {
	cw, err := r.cli.AttachToContainerNonBlocking(docker.AttachToContainerOptions{
		Container:    r.c.ID,
		OutputStream: stdout, ErrorStream: stderr,
		Stdout: true, Stderr: true,
		Logs: true, Stream: true,
	})
	if err != nil {
		return err
	}
	defer cw.Close()

	err = r.cli.StartContainer(r.c.ID, nil)
	if err != nil {
		return err
	}
	return cw.Wait()
}

func runInDocker(t *testing.T, fname, goroot string) bool {
	dir, err := ioutil.TempDir("", "js_test_")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	copyFile(t, fname, filepath.Join(dir, "test"))
	for _, name := range []string{
		wasmExec, wasmExecJS, wasmExecHTML,
	} {
		copyFile(t, filepath.Join(goroot, wasmDir, name), filepath.Join(dir, name))
	}

	p, err := dockertest.NewPool("")
	if err != nil {
		return false
	}
	cli := p.Client

	now := time.Now()
	if !pullIfNotExists(cli, nodejsImage) {
		return false
	}
	t.Logf("pulled image %q in %v", nodejsImage, time.Since(now))

	c, err := cli.CreateContainer(docker.CreateContainerOptions{
		Config: &docker.Config{
			Image:      nodejsImage,
			WorkingDir: "/test",
			Entrypoint: []string{"./" + wasmExec, "./test", "-test.v"},
		},
		HostConfig: &docker.HostConfig{
			Binds: []string{
				dir + ":" + "/test",
			},
		},
	})
	require.NoError(t, err)
	defer cli.RemoveContainer(docker.RemoveContainerOptions{
		ID: c.ID, RemoveVolumes: true, Force: true,
	})
	t.Log("running in Docker")

	testproxy.RunAndReplay(t, dockerRunner{
		cli: cli, c: c,
	})
	return true
}
