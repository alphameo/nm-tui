{
    description = "Lightweight TUI wrapper for NetworkManager.";

    inputs = {
        nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
        systems.url = "github:nix-systems/triplet";
    };

    outputs =
        { self, nixpkgs, ... }@inputs:
        let
            forEachSystem = nixpkgs.lib.genAttrs (import inputs.systems);
            systemPackages = system: nixpkgs.legacyPackages.${system};
        in
        {
            packages = forEachSystem (
                system:
                let
                    pkgs = systemPackages system;
                in
                {
                    default = pkgs.callPackage ./nix/package.nix { inherit self; };
                }
            );

            devShells = forEachSystem (
                system:
                let
                    pkgs = systemPackages system;
                in
                {
                    default = pkgs.mkShell {
                        packages = with pkgs; [
                            go_1_27
                            goreleaser
                            golangci-lint
                            gnumake
                        ];
                    };
                }
            );
        };
}
