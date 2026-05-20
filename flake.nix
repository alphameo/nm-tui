{
    description = "Lightweight TUI wrapper for NetworkManager.";

    outputs =
        { flakelight, systems, ... }@inputs:
        flakelight ./. {
            inherit inputs;
            systems = import systems;
            devShell.inputsFrom = pkgs: [ pkgs.nm-tui ];
        };

    inputs = {
        nixpkgs.url = "github:NixOS/nixpkgs/master"; # go 1.26.3
        systems.url = "github:nix-systems/default-linux";

        flakelight = {
            url = "github:nix-community/flakelight";
            inputs.nixpkgs.follows = "nixpkgs";
        };
    };
}
