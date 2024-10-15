dpc_tox in
=====

[![GoDoc](https://godoc.org/github.com/calvindc/dpc-tox?status.png)](http://github.com/calvindc/dpc-tox)


gotox is a Go wrapper for the [c-toxcore](https://github.com/irungentoo/toxcore) library.

Pull requests, bug reportings and feature requests (via github issues) are always welcome. :)

For a list of supported toxcore features see [PROGRESS.md](PROGRESS.md).

## Installation
First, install the [c-toxcore](https://github.com/calvindc/dpc-tox) library.

Next, download `go-tox` using go:
```
go get github.com/calvindc/dpc-tox
```

## License
gotox is licensed under the [GPLv3](COPYING).

## How to use
See [bindings.go](bindings.go) for details about supported API functions and [callbacks.go](callbacks.go) for the supported callbacks.

The best place to get started are the examples in [examples/](examples/).

```
go run ...
```

Feel free to ask for help in the issue tracker. ;)
