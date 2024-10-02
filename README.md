# Go bindings to nix API

[![Go Report
Card](https://goreportcard.com/badge/github.com/farcaller/gonix)](https://goreportcard.com/report/github.com/farcaller/gonix)
[![GoDoc](https://godoc.org/github.com/farcaller/gonix?status.svg)](https://godoc.org/github.com/farcaller/gonix)

## Building

These bindings depend on the Nix C API, which is currently only available in nix
master. It's a moving target, and while we try to catch up, sometimes the API is
broken so you must pull nix from the **same revision** that's tracked in this
flake's `nix`, and make both `nix.dev` and `pkg-config` available in your
`buildInputs` for the CGO bindings to work. Consult the `flake.nix` for an
example.

## API Docs

See [godoc](https://godoc.org/github.com/farcaller/gonix) for API and examples.
