# New release

You need to bump version in `package.nix`

## If dependencies have changed

1. Build and get `vendorHash` on fail

    ```bash
    nix build '.#' --no-link 2>&1 | grep "got:"
    ```

1. Paste result inside [`package.nix`](./package.nix) into `vendorHash` field

1. Optionally update `flake.lock`

    ```bash
    nix flake update
    ```
