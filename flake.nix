{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-25.05";
    systems.url = "github:nix-systems/default";
    devenv.url = "github:cachix/devenv";
  };

  nixConfig = {
    extra-trusted-public-keys = "devenv.cachix.org-1:w1cLUi8dv3hnoSPGAuibQv+f9TZLr6cv/Hm9XgU50cw=";
    extra-substituters = "https://devenv.cachix.org";
  };

  outputs =
    { self
    , nixpkgs
    , devenv
    , systems
    , ...
    }@inputs:
    let
      forEachSystem = nixpkgs.lib.genAttrs (import systems);
    in
    {

      packages = forEachSystem (system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
          nixPkg = pkgs.nixVersions.nix_2_30;
        in
        {
          default = pkgs.buildGoModule {
            name = "gonix";
            src = ./.;
            vendorHash = null;
            env.CGO_ENABLED = 1;
            nativeBuildInputs = [ pkgs.pkg-config ];
            buildInputs = [ nixPkg ];
            # this is not actually required for the build, but
            # for tests that require `import <nixpkgs>`
            NIX_PATH = "nixpkgs=${nixpkgs}";

            # gonix tests require instantiating a store, which cannot be done inside a nix build
            # (that would be recursive nix) so when building inside nix, we disable the tests
            #
            # you can still run tests with nix develop -c go test ./... -v
            doCheck = false;
          };
        });
    };
}
