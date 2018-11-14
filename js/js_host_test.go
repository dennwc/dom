//+build !js

package js

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dennwc/testproxy"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"

	"github.com/dennwc/dom/internal/goenv"
	"github.com/stretchr/testify/require"
)

const (
	image = "node:8-slim"

	wasmDir      = "misc/wasm"
	wasmExec     = "go_js_wasm_exec"
	wasmExecJS   = "wasm_exec.js"
	wasmExecHTML = "wasm_exec.html"
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

	if runInDocker(t, f.Name(), GOROOT) {
		return
	}
	t.Log("running locally")

	// run tests with Node.js
	if _, err := exec.LookPath("node"); err != nil {
		t.Skipf("cannot find NodeJS: %v", err)
	}

	wd, err := os.Getwd()
	require.NoError(t, err)
	cmd = exec.Command(filepath.Join(GOROOT, wasmDir, wasmExec), f.Name())
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

	if !pullIfNotExists(cli, image) {
		return false
	}

	c, err := cli.CreateContainer(docker.CreateContainerOptions{
		Config: &docker.Config{
			Image:      image,
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

func pullIfNotExists(cli *docker.Client, image string) bool {
	_, err := cli.InspectImage(image)
	if err == nil {
		return true
	}
	i := strings.Index(image, ":")
	err = cli.PullImage(docker.PullImageOptions{
		Repository: image[:i], Tag: image[i+1:],
	}, docker.AuthConfiguration{})
	return err == nil
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
