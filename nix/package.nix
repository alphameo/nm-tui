{
    src, # flakelight

    lib,
    buildGo127Module,
    makeWrapper,

    networkmanager,
}:
buildGo127Module {
    inherit src;
    pname = "nm-tui";
    version = "0.2.3";

    vendorHash = "sha256-qq7dbswbq+5h1iKBpohZBGEln+SsuZFnXs0RDYDhUy4=";

    nativeBuildInputs = [ makeWrapper ];

    env = {
        CGO_ENABLED = 0;
    };

    postInstall = ''
        wrapProgram $out/bin/nm-tui \
            --prefix PATH : ${lib.makeBinPath [ networkmanager ]}
    '';
}
