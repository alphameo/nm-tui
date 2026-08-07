# Formula/nm-tui.rb — Homebrew tap

**Caveat:** nm-tui is Linux-only (wraps NetworkManager). Homebrew-core builds macOS
bottles, which can never run this, so ship it via a self-hosted tap
(e.g. `alphameo/homebrew-tap`) rather than core.

## Publishing

1. Create the tap repo:

   ```bash
   brew tap-new alphameo/homebrew-tap
   ```

2. Copy the formula in and set the release version + URL + sha256:

   ```bash
   VER=0.2.0
   cp Formula/nm-tui.rb ~/homebrew-tap/Formula/nm-tui.rb
   curl -sL "https://github.com/alphameo/nm-tui/archive/refs/tags/v${VER}.tar.gz" | sha256sum
   ```

   Update the `url` to `.../tags/v${VER}.tar.gz` and paste the hash into `sha256`.

3. Verify the formula builds from source:

   ```bash
   brew install --build-from-source ~/homebrew-tap/Formula/nm-tui.rb
   ```

4. Commit and push:

   ```bash
   git -C ~/homebrew-tap add Formula/nm-tui.rb
   git -C ~/homebrew-tap commit -m "nm-tui ${VER}"
   git -C ~/homebrew-tap push
   ```

## Notes

- `depends_on "networkmanager"` — present in homebrew-core, brew resolves it.
- `std_go_args` injects the standard trimpath/output flags; we add
  `-X main.version` so `nm-tui --version` reports the release version.
- The `man1.install` line is a no-op until a man page is added under `docs/`.