{
    src, # flakelight

    lib,
    buildGo127Module,
    makeWrapper,

    networkmanager,
}:
let
    version = "0.2.3";
in
buildGo127Module {
    inherit src version;
    pname = "nm-tui";

    vendorHash = "sha256-qq7dbswbq+5h1iKBpohZBGEln+SsuZFnXs0RDYDhUy4=";

    nativeBuildInputs = [ makeWrapper ];

    env = {
        CGO_ENABLED = 0;
    };

    ldflags = [
        "-X main.version=${version}-nix"
    ];

    postInstall = ''
        wrapProgram $out/bin/nm-tui \
            --prefix PATH : ${lib.makeBinPath [ networkmanager ]}
    '';
}
