{
  inputs = {
    nix.url = "github:nixos/nix";
    nixpkgs.follows = "nix/nixpkgs";
    # nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    systems.url = "github:nix-systems/default";
    devenv.url = "github:cachix/devenv";
  };

  nixConfig = {
    extra-trusted-public-keys = "devenv.cachix.org-1:w1cLUi8dv3hnoSPGAuibQv+f9TZLr6cv/Hm9XgU50cw=";
    extra-substituters = "https://devenv.cachix.org";
  };

  outputs =
    {
      self,
      nixpkgs,
      nix,
      devenv,
      systems,
      ...
    }@inputs:
    let
      forEachSystem = nixpkgs.lib.genAttrs (import systems);
    in
    {
      devShells = forEachSystem (
        system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
          nixPkg = nix.packages.${system}.nix;
        in
        {
          default = devenv.lib.mkShell {
            inherit inputs pkgs;

            modules = [
              {
                env.hardeningDisable = [ "fortify" ];
                languages.nix.enable = true;
                languages.go.enable = true;
                languages.c.enable = true;

                packages = [
                  nixPkg.dev
                  pkgs.pkg-config
                  pkgs.delve
                ];
              }
            ];
          };
        }
      );
    };
}
