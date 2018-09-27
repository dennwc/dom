# Go DOM binding (and more) for WebAssembly

This library provides a Go API for different Web APIs for WebAssembly target.

It's in an active development, but an API will be carefully versioned to
avoid breaking users.
Use Go dependency management tools to lock a specific version.

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