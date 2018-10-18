# Go DOM binding (and more) for WebAssembly

This library provides a Go API for different Web APIs for WebAssembly target.

It's in an active development, but an API will be carefully versioned to
avoid breaking users.
Use Go dependency management tools to lock a specific version.

More information about Go's WebAssembly support can be found on [Go's WebAssembly wiki page](https://github.com/golang/go/wiki/WebAssembly).

**Features:**

- Better JS API (wrappers for `syscall/js`)
- Basic DOM manipulation, styles, events
- Input elements
- SVG elements and transforms
- `LocalStorage` and `SessionStorage`
- Extension APIs (tested on Chrome):
    - Native Messaging
    - Bookmarks
    - Tabs
- `net`-like library for WebSockets
    - Tested with gRPC
- `wasm-server` for fast prototyping

## Quickstart

Pull the library and install `wasm-server` (optional):

```
go get -u github.com/dennwc/dom
go install github.com/dennwc/dom/cmd/wasm-server
```

Run an example app:

```
cd $GOPATH/src/github.com/dennwc/dom
wasm-server
```

Check result: http://localhost:8080/

The source code is recompiled on each page refresh, so feel free to experiment!

# Similar Projects

- [go-js-dom](https://github.com/dominikh/go-js-dom)

# Editor Configuration

If you are using Visual Studio Code, you can use [workspace settings](https://code.visualstudio.com/docs/getstarted/settings#_creating-user-and-workspace-settings) to configure the environment variables for the go tools.

Your settings.json file should look something like this:

```
{
    "go.toolsEnvVars": { "GOARCH": "wasm", "GOOS": "js" }
}
```
