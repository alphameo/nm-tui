# nm-tui Fedora/Copr packaging

`nm-tui.spec` builds a static binary on the Fedora/EPEL targets of a COPR project.

## Publishing to COPR

Prereqs: `copr-cli` installed and logged in (`copr login`).

1. Create the project (one time):

   ```bash
   copr create nm-tui --chroot fedora-latest-x86_64 \
     --chroot epel-9-x86_64 \
     --repo https://download.copr.fedorainfracloud.org/results/alphameo/nm-tui/fedora-latest-x86_64/ \
     --description "Lightweight TUI wrapper for NetworkManager" \
     --instructions "dnf install copr-cli && dnf copr enable alphameo/nm-tui"
   ```

2. Build from this spec (needs the release tarball reachable at `Source0`):

   ```bash
   VER=0.2.0
   sed -i "s/^Version:.*/Version:        ${VER}/" nm-tui.spec
   curl -sL "https://github.com/alphameo/nm-tui/archive/refs/tags/v${VER}.tar.gz" -o "nm-tui-${VER}.tar.gz"
   copr-cli build alphameo/nm-tui nm-tui.spec
   ```

3. Tell users to enable the repo:

   ```bash
   sudo dnf copr enable alphameo/nm-tui
   sudo dnf install nm-tui
   ```

## Notes

- `%autosetup -n nm-tui-%{version}` matches the GitHub tarball layout
  (`nm-tui-<version>/` directory).
- `BuildRequires: golang >= 1.26` must track `go.mod`'s `go` directive.
- Static build with `CGO_ENABLED=0`; version baked in via `-X main.version`.
- Fedora guidelines require a `networkmanager` runtime dependency when upstream
  builds against it; for COPR this is optional and can be added via
  `Requires: networkmanager` if needed.