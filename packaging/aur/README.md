# AUR package: nm-tui

This directory holds the AUR package source for `nm-tui`.

## Publishing

1. Create the AUR repo:

   ```bash
   git clone ssh://aur@aur.archlinux.org/nm-tui.git /tmp/aur-nm-tui
   cp PKGBUILD .SRCINFO /tmp/aur-nm-tui/
   cd /tmp/aur-nm-tui
   ```

2. Update the release version in `PKGBUILD` and `.SRCINFO`:

   ```bash
   pkgver=0.2.0   # next release tag without the leading 'v'
   ```

3. Get the real checksum for the release tarball and put it in both files:

   ```bash
   curl -sL https://github.com/alphameo/nm-tui/archive/refs/tags/v0.2.0.tar.gz | sha256sum
   ```

4. Validate and regenerate `.SRCINFO` (must stay in sync with `PKGBUILD`):

   ```bash
   makepkg -f   # also builds; run a real build locally first
   makepkg --printsrcinfo > .SRCINFO
   ```

5. Commit and push:

   ```bash
   git add PKGBUILD .SRCINFO
   git commit -m "nm-tui 0.2.0"
   git push
   ```

## Notes

- The leading `v` is stripped from the git tag for `pkgver` (Arch packaging convention).
- `depends=('networkmanager')` — the binary shells out to `nmcli`, so it must be present at runtime.
- The binary is statically built with `CGO_ENABLED=0`; `-trimpath` keeps builds reproducible.
- Version is baked in via `-X main.version=$pkgver` so `nm-tui --version` reports it.
