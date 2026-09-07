{
    description = "Lightweight TUI wrapper for NetworkManager.";

    outputs =
        { flakelight, systems, ... }@inputs:
        flakelight ./. {
            inherit inputs;
            systems = import systems;
            devShell.packages =
                pkgs: with pkgs; [
                    go_1_27
                    goreleaser
                    golangci-lint
                    gnumake
                ];
        };

    inputs = {
        nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
        systems.url = "github:nix-systems/triplet";

        flakelight = {
            url = "github:nix-community/flakelight";
            inputs.nixpkgs.follows = "nixpkgs";
        };
    };
}
