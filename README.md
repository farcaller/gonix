# Go bindings to nix API

## Building

These bindings depend on the unstable nix C API, that's not currently merged
upstream. You must pull nix from `github:tweag/nix/nix-c-bindings`, and make
both `nix.dev` and `pkg-config` available in your `buildInputs` for the cgo
buindings to work. Consult the `flake.nix` for an expample.
