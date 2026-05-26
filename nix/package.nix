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
    version = "0.0.2";

    vendorHash = "sha256-JYO6UHZwOmudADOiTAPM+mSm3KMKOERFTiM1beMN+MI=";

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
