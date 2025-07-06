{
  description = "Juicer dev environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";

    gomod2nix = {
      url = "github:nix-community/gomod2nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = { self, nixpkgs, flake-utils, gomod2nix, ... }:
    flake-utils.lib.eachDefaultSystem
      (system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
        in
        {
          devShells = {
            default = pkgs.mkShell {
              packages = with pkgs; [
                # go
                # go-tools
                # gotools
                # golangci-lint
                # golangci-lint-langserver
                # gopls
                # delve
                # gomodifytags
                # errcheck
                # goreleaser
                # golines
                # gotests
                # gotestsum
                # impl
                # docker
                # docker-compose
                # protobuf
                protoc-gen-go
                air
                # atlas
                oapi-codegen
                go-jet
                # just
                # nodejs_24
                # caddy
                # jdk
              ];

              shellHook = ''
                echo "*** YOU ARE IN JUICER DEV ENVIRONMENT ***"
              '';
            };
          };
        }
      );
}
