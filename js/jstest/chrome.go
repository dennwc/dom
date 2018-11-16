//+build !js

package jstest

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	clog "github.com/chromedp/cdproto/log"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/client"
	"github.com/dennwc/dom/internal/goenv"
	"github.com/dennwc/testproxy"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/stretchr/testify/require"
)

const chromeImage = "chromedp/headless-shell:70.0.3526.1"

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

// RunTestChrome compiles all tests in the current working directory for WASM+JS, and runs them in a headless Chrome
// browser using Docker. It will stream test results back to t.
//
// Optionally, the default HTTP handler can be set to serve additional content to the the script.
//
// The caller should specify the "!js" build tag, while all JS tests in the package should include "js" build tag.
func RunTestChrome(t *testing.T, def http.Handler) {
	testFile := buildTestJS(t)
	defer os.Remove(testFile)

	addr, closer := runChromeInDocker(t)
	defer closer()

	rn := newChromeRunner(addr, testFile, def)
	testproxy.RunAndReplay(t, rn)
}

func newChromeRunner(addr, testfile string, def http.Handler) testproxy.Runner {
	c := client.New(client.URL("http://" + addr + "/json"))
	return &chromeRunner{
		c: c, bin: testfile, def: def,
	}
}

type chromeRunner struct {
	c   *client.Client
	bin string
	def http.Handler
}

func (rn *chromeRunner) RunAndWait(stdout, stderr io.Writer) error {
	goroot := goenv.GOROOT()
	port := strconv.Itoa(rnd.Intn(10000) + 10000)

	srv := &http.Server{
		Addr: ":" + port,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := strings.Trim(r.URL.Path, "/")
			switch path {
			case "", "index.html":
				w.Write([]byte(indexHTML))
			case "wasm_exec.js":
				http.ServeFile(w, r, filepath.Join(goroot, wasmDir, wasmExecJS))
			case "test.wasm":
				http.ServeFile(w, r, rn.bin)
			default:
				if rn.def != nil {
					rn.def.ServeHTTP(w, r)
				} else {
					w.WriteHeader(http.StatusNotFound)
				}
			}
		}),
	}
	defer srv.Close()

	errc := make(chan error, 1)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			errc <- err
		}
	}()

	// TODO: discover Docker host IP
	surl := "http://172.17.0.1:" + port

	ctx := context.TODO()
	c := rn.c
	tg, err := c.NewPageTargetWithURL(ctx, "about:blank")
	if err != nil {
		return err
	}
	defer c.CloseTarget(ctx, tg)

	logf := func(format string, args ...interface{}) {
		log.Printf(format, args...)
	}
	noop := func(format string, args ...interface{}) {}
	_ = noop

	h, err := chromedp.NewTargetHandler(tg, logf, noop, logf)
	if err != nil {
		return err
	}

	h.OnEvent(func(ev interface{}) {
		switch ev := ev.(type) {
		case *runtime.EventConsoleAPICalled:
			if ev.Type != "log" || len(ev.Args) != 1 {
				return
			}
			var line string
			if err := json.Unmarshal(ev.Args[0].Value, &line); err != nil {
				select {
				case errc <- err:
				default:
				}
				return
			}
			line += "\n"
			stdout.Write([]byte(line))
		case *clog.EventEntryAdded:
			stderr.Write([]byte(ev.Entry.Text + "\n"))
		}
	})

	err = h.Run(ctx)
	if err != nil {
		return err
	}

	err = chromedp.Navigate(surl).Do(ctx, h)
	if err != nil {
		return err
	}

	err = chromedp.WaitEnabled("#runButton").Do(ctx, h)
	if err != nil {
		return err
	}

	// TODO: better wait logic
	var res interface{}
	err = chromedp.Evaluate("run().then(() => {done = true}).catch(() => {done = true})", &res).Do(ctx, h)
	if err != nil {
		return err
	}

	var done bool
	err = chromedp.Evaluate("done", &done).Do(ctx, h)
	if err != nil {
		return err
	}
	ticker := time.NewTicker(time.Second / 3)
	defer ticker.Stop()
	for !done {
		err = chromedp.Evaluate("done", &done).Do(ctx, h)
		if err != nil {
			return err
		}
		select {
		case <-ticker.C:
		case err = <-errc:
			return err
		}
	}
	return nil
}

func runChromeInDocker(t testing.TB) (string, func()) {
	p, err := dockertest.NewPool("")
	require.NoError(t, err)
	cli := p.Client

	now := time.Now()
	if !pullIfNotExists(cli, chromeImage) {
		t.SkipNow()
		return "", func() {}
	}
	t.Logf("pulled image %q in %v", chromeImage, time.Since(now))

	c, err := cli.CreateContainer(docker.CreateContainerOptions{
		Config: &docker.Config{
			Image: chromeImage,
		},
	})
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	cw, err := cli.AttachToContainerNonBlocking(docker.AttachToContainerOptions{
		Container:    c.ID,
		OutputStream: buf,
		ErrorStream:  buf,
		Stdout:       true, Stderr: true,
		Logs: true, Stream: true,
	})
	if err != nil {
		cli.RemoveContainer(docker.RemoveContainerOptions{
			ID: c.ID, RemoveVolumes: true, Force: true,
		})
		require.NoError(t, err)
	}

	remove := func() {
		cw.Close()
		cli.RemoveContainer(docker.RemoveContainerOptions{
			ID: c.ID, RemoveVolumes: true, Force: true,
		})
	}
	t.Log("running in Chrome in Docker")

	err = cli.StartContainer(c.ID, nil)
	if err != nil {
		remove()
		require.NoError(t, err)
	}

	info, err := cli.InspectContainer(c.ID)
	if err != nil {
		remove()
		require.NoError(t, err)
	}
	addr := info.NetworkSettings.IPAddress + ":9222"
	if !waitPort(addr) {
		remove()
		require.Fail(t, "timeout", "logs:\n%v", buf.String())
	}
	return addr, remove
}

func waitPort(addr string) bool {
	for i := 0; i < 10; i++ {
		c, err := net.DialTimeout("tcp", addr, time.Second)
		if err == nil {
			c.Close()
			return true
		}
		time.Sleep(time.Second * 3)
	}
	return false
}

const indexHTML = `<!doctype html>
<!--
Copyright 2018 The Go Authors. All rights reserved.
Use of this source code is governed by a BSD-style
license that can be found in the LICENSE file.
-->
<html>

<head>
	<meta charset="utf-8">
	<title>Go wasm</title>
</head>

<body>
	<!--
	Add the following polyfill for Microsoft Edge 17/18 support:
	<script src="https://cdn.jsdelivr.net/npm/text-encoding@0.7.0/lib/encoding.min.js"></script>
	(see https://caniuse.com/#feat=textencoder)
	-->
	<script src="wasm_exec.js"></script>
	<script>
		if (!WebAssembly.instantiateStreaming) { // polyfill
			WebAssembly.instantiateStreaming = async (resp, importObject) => {
				const source = await (await resp).arrayBuffer();
				return await WebAssembly.instantiate(source, importObject);
			};
		}

		const go = new Go();
		go.argv.push('-test.v');
		let mod, inst;
		WebAssembly.instantiateStreaming(fetch("test.wasm"), go.importObject).then((result) => {
			mod = result.module;
			inst = result.instance;
			document.getElementById("runButton").disabled = false;
		}).catch((err) => {
			console.error(err);
		});

		async function run() {
			console.clear();
			await go.run(inst);
		}
		let done = false;
	</script>

	<button onClick="run();" id="runButton" disabled>Run</button>
</body>

</html>`
