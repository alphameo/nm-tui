{
    self,
    buildGo127Module,
}:
let
    rev = self.shortRev or self.dirtyShortRev;
in
buildGo127Module (finalAttrs: {
    pname = "nm-tui";
    version = "0.3.3+${rev}";
    src = ../.;

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
