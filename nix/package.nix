{
    src, # flakelight
    inputs, # flakelight
    buildGo127Module,
}:
let
    inherit (inputs) self;
    rev = self.shortRev or self.dirtyShortRev;
in
buildGo127Module (finalAttrs: {
    inherit src;
    pname = "nm-tui";
    version = "0.3.2+${rev}";

    vendorHash = "sha256-qq7dbswbq+5h1iKBpohZBGEln+SsuZFnXs0RDYDhUy4=";

    env = {
        CGO_ENABLED = 0;
    };

    ldflags = [
        "-s"
        "-w"
        "-X main.version=v${finalAttrs.version}"
    ];
})
