# Go bindings to nix API

[![Go Report
Card](https://goreportcard.com/badge/github.com/farcaller/gonix)](https://goreportcard.com/report/github.com/farcaller/gonix)
[![GoDoc](https://godoc.org/github.com/farcaller/gonix?status.svg)](https://godoc.org/github.com/farcaller/gonix)

## Building

These bindings depend on the Nix C API, which is currently only available in nix
master. You must pull nix from `github:nixos/nix`, and make
both `nix.dev` and `pkg-config` available in your `buildInputs` for the CGO
bindings to work. Consult the `flake.nix` for an example.

## API Docs

See [godoc](https://godoc.org/github.com/farcaller/gonix) for API and examples.
