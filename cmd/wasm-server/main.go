package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/dennwc/dom/internal/goenv"
)

var (
	host      = flag.String("host", ":8080", "host to serve on")
	directory = flag.String("d", "./", "the directory of static file to host")
	cmds      = flag.String("apps", "cmd", "the root directory for apps")
	def       = flag.String("main", "app", "default app name")
)

var (
	indexTmpl = template.Must(template.New("index").Parse(indexHTML))
	GOROOT    string
)

func main() {
	flag.Parse()

	GOROOT = goenv.GOROOT()

	h := http.FileServer(http.Dir(*directory))

	http.Handle("/", buildHandler{h: h, dir: *directory})

	fmt.Println("goroot: ", GOROOT)
	fmt.Println("apps folder: ", *cmds)
	fmt.Println("default app: ", *def)
	fmt.Println("static files:", *directory)
	fmt.Println("serving on:  ", *host)
	err := http.ListenAndServe(*host, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

type buildHandler struct {
	h   http.Handler
	dir string
}

func (h buildHandler) appPath(name string) string {
	return filepath.Join(h.dir, *cmds, name, "main.go")
}

func (h buildHandler) buildWASM(name string) {
	dst := filepath.Join(h.dir, name+".wasm")
	bin := "go"
	if GOROOT != "" {
		bin = filepath.Join(GOROOT, "bin", "go")
	}
	cmd := exec.Command(bin, "build",
		"-o", dst,
		h.appPath(name),
	)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env,
		"GOOS=js",
		"GOARCH=wasm",
	)
	if GOROOT != "" {
		cmd.Env = append(cmd.Env,
			"GOROOT="+GOROOT,
		)
	}
	start := time.Now()
	out, err := cmd.CombinedOutput()
	dt := time.Since(start)
	if err != nil {
		os.Remove(dst)
		log.Println(err)
		log.Println("logs:\n" + string(out))
	} else {
		log.Printf("built %q in %v", name, dt)
	}
}

func (h buildHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.Path)
	rpath := r.URL.Path
	if rpath == "/" {
		rpath = "/index.html"
	}
	switch rpath {
	case "/" + execName:
		http.ServeFile(w, r, filepath.Join(GOROOT, "misc", "wasm", execName))
		return
	}
	switch ext := path.Ext(rpath); ext {
	case "", ".html":
		appName := *def
		if rpath != "/index.html" {
			appName = strings.Trim(strings.TrimSuffix(rpath, ext), "/")
		}
		if _, err := os.Stat(h.appPath(appName)); err == nil {
			err = indexTmpl.Execute(w, struct {
				AppName string
			}{
				AppName: appName,
			})
			if err != nil {
				log.Println(err)
			}
			return
		}
	case ".wasm":
		h.buildWASM(strings.TrimSuffix(r.URL.Path, ext))
	}
	h.h.ServeHTTP(w, r)
}

const (
	execName = "wasm_exec.js"
)

const indexHTML = `<!doctype html>
<!--
Copyright 2018 The Go Authors. All rights reserved.
Use of this source code is governed by a BSD-style
license that can be found in the LICENSE file.
-->
<html>

<head>
	<meta charset="utf-8">
	<title>{{ .AppName }}</title>
</head>

<body>
	<script src="/` + execName + `"></script>
	<script>
		if (!WebAssembly.instantiateStreaming) { // polyfill
			WebAssembly.instantiateStreaming = async (resp, importObject) => {
				const source = await (await resp).arrayBuffer();
				return await WebAssembly.instantiate(source, importObject);
			};
		}

		const go = new Go();
		WebAssembly.instantiateStreaming(fetch("{{ .AppName }}.wasm"), go.importObject).then(async (result) => {
			let mod = result.module;
			let inst = result.instance;

			let spin = document.getElementById('spin');
			spin.parentNode.removeChild(spin);

			// run
			await go.run(inst);
			// reset instance
			// inst = await WebAssembly.instantiate(mod, go.importObject);
		});
	</script>
	<style>` + spinCSS + `</style>
	<div id="spin">
		<div class="sk-cube-grid">
			<div class="sk-cube sk-cube1"></div>
			<div class="sk-cube sk-cube2"></div>
			<div class="sk-cube sk-cube3"></div>
			<div class="sk-cube sk-cube4"></div>
			<div class="sk-cube sk-cube5"></div>
			<div class="sk-cube sk-cube6"></div>
			<div class="sk-cube sk-cube7"></div>
			<div class="sk-cube sk-cube8"></div>
			<div class="sk-cube sk-cube9"></div>
		</div>
		<div class="sk-text">Loading...</div>
	</div>
</body>

</html>
`

const spinCSS = `.sk-cube-grid {
  width: 40px;
  height: 40px;
  margin: 100px auto 20px auto;
}

.sk-text {
  font-family: Arial, Helvetica, sans-serif;
  text-align: center;
  margin: 0 auto 30px auto;
}

.sk-cube-grid .sk-cube {
  width: 33%;
  height: 33%;
  background-color: #333;
  float: left;
  -webkit-animation: sk-cubeGridScaleDelay 1.3s infinite ease-in-out;
          animation: sk-cubeGridScaleDelay 1.3s infinite ease-in-out; 
}
.sk-cube-grid .sk-cube1 {
  -webkit-animation-delay: 0.2s;
          animation-delay: 0.2s; }
.sk-cube-grid .sk-cube2 {
  -webkit-animation-delay: 0.3s;
          animation-delay: 0.3s; }
.sk-cube-grid .sk-cube3 {
  -webkit-animation-delay: 0.4s;
          animation-delay: 0.4s; }
.sk-cube-grid .sk-cube4 {
  -webkit-animation-delay: 0.1s;
          animation-delay: 0.1s; }
.sk-cube-grid .sk-cube5 {
  -webkit-animation-delay: 0.2s;
          animation-delay: 0.2s; }
.sk-cube-grid .sk-cube6 {
  -webkit-animation-delay: 0.3s;
          animation-delay: 0.3s; }
.sk-cube-grid .sk-cube7 {
  -webkit-animation-delay: 0s;
          animation-delay: 0s; }
.sk-cube-grid .sk-cube8 {
  -webkit-animation-delay: 0.1s;
          animation-delay: 0.1s; }
.sk-cube-grid .sk-cube9 {
  -webkit-animation-delay: 0.2s;
          animation-delay: 0.2s; }

@-webkit-keyframes sk-cubeGridScaleDelay {
  0%, 70%, 100% {
    -webkit-transform: scale3D(1, 1, 1);
            transform: scale3D(1, 1, 1);
  } 35% {
    -webkit-transform: scale3D(0, 0, 1);
            transform: scale3D(0, 0, 1); 
  }
}

@keyframes sk-cubeGridScaleDelay {
  0%, 70%, 100% {
    -webkit-transform: scale3D(1, 1, 1);
            transform: scale3D(1, 1, 1);
  } 35% {
    -webkit-transform: scale3D(0, 0, 1);
            transform: scale3D(0, 0, 1);
  } 
}`
