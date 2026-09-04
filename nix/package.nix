{
    src, # flakelight
    buildGo127Module,
}:
buildGo127Module (finalAttrs: {
    inherit src;
    pname = "nm-tui";
    version = "0.3.1";

    vendorHash = "sha256-qq7dbswbq+5h1iKBpohZBGEln+SsuZFnXs0RDYDhUy4=";

    env = {
        CGO_ENABLED = 0;
    };

    ldflags = [
        "-X main.version=${finalAttrs.version}-nix"
    ];
})
