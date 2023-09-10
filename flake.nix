{
  inputs = {
    nix.url = "github:tweag/nix/nix-c-bindings";
    nixpkgs.follows = "nix/nixpkgs";
    # nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    systems.url = "github:nix-systems/default";
    devenv.url = "github:cachix/devenv";
  };

  nixConfig = {
    extra-trusted-public-keys = "devenv.cachix.org-1:w1cLUi8dv3hnoSPGAuibQv+f9TZLr6cv/Hm9XgU50cw=";
    extra-substituters = "https://devenv.cachix.org";
  };

  outputs = { self, nixpkgs, nix, devenv, systems, ... } @ inputs:
    let
      forEachSystem = nixpkgs.lib.genAttrs (import systems);
    in
    {
      devShells = forEachSystem
        (system:
          let
            pkgs = nixpkgs.legacyPackages.${system};
            nixPkg = nix.packages.${system}.nix;
          in
          {
            default = devenv.lib.mkShell {
              inherit inputs pkgs;

              modules = [
                {
                  # env.LDFLAGS = "-F${pkgs.darwin.CF}/Library/Frameworks";
                  # env.LDFLAGS = "-F/Library/Developer/CommandLineTools/SDKs/MacOSX.sdk/System/Library/Frameworks";

                  languages.nix.enable = true;
                  languages.go.enable = true;
                  
                  packages = [
                    nixPkg
                    nixPkg.dev
                    pkgs.pkg-config
                  ];
                }
              ];
            };
          });
    };
}
