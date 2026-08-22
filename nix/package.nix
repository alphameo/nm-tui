{
    src, # flakelight

    lib,
    buildGoModule,
    makeWrapper,

    networkmanager,
}:
buildGoModule {
    inherit src;
    pname = "nm-tui";
    version = "0.2.1";

    vendorHash = "sha256-AyYJFuuURz+QDe69iAgWlN1Xd7+Ofh4hVL0Xya706N8=";

    nativeBuildInputs = [ makeWrapper ];
    subPackages = [ "cmd/nm-tui" ];

    env = {
        CGO_ENABLED = 0;
    };

    postInstall = ''
        wrapProgram $out/bin/nm-tui \
            --prefix PATH : ${lib.makeBinPath [ networkmanager ]}
    '';
}
