# nm-tui

![License](https://img.shields.io/github/license/alphameo/nm-tui)
![Go Version](https://img.shields.io/github/go-mod/go-version/alphameo/nm-tui)
![GitHub stars](https://img.shields.io/github/stars/alphameo/nm-tui)

Lightweight TUI wrapper for [NetworkManager](https://gitlab.freedesktop.org/NetworkManager/NetworkManager)

Why `nm-tui`: built-in `nmtui` doesn't look great and there aren't many TUI alternatives

## 📌 Table of Contents

- [💫 Features](#features)
- [📹 Demo](#demo)
- [🖼️ Screenshots](#screenshots)
- [🗃️ Requirements](#requirements)
- [📥 Installation](#installation)
- [⚙️ Configuration](#configuration)
- [👨‍💻 Tech Stack](#tech-stack)
- [🖲️ Contributing](#contributing)
- [⚖️ License](#license)
- [⭐ Inspirations](#inspirations)

## Features

- 😎 TUI style (looks cool)
- 📡 Scan and list available networks
- 🔑 Connect to networks with password
- 🔘 Activate connections to saved networks
- 📜 View detailed network information (signal strength, security, etc.)
- 🌐 Control device networking
- 📡 Create hotspot
- 🖥️ Clean, modern TUI built with Bubbletea
- ⚡ Fast and lightweight — single static binary
- 🐧 Linux only — designed specifically for NetworkManager
- 🌈 Customize look & feel

## Demo

![Demo](../assets/demo-conn.gif)

<details>
    <summary><h2>Screenshots</h2></summary>

### Main tabs

<div style="display: flex; gap: 10px;">
    <img src="../assets/wifi-tab.png" alt="wifi connector" width="400"/>
    <img src="../assets/networking-tab.png" alt="wifi info" width="400"/>
</div>

### WiFi connection and Network info

<div style="display: flex; gap: 10px;">
    <img src="../assets/connect-to-wifi.png" alt="wifi connector" width="400"/>
    <img src="../assets/network-info.png" alt="wifi info" width="400"/>
</div>

### Network and Access point creation

<div style="display: flex; gap: 10px;">
    <img src="../assets/create-wifi-profile.png" alt="wifi info" width="400"/>
    <img src="../assets/create-wifi-hotspot.png" alt="wifi connector" width="400"/>
</div>

</details>

## Requirements

- [`NetworkManager`](https://gitlab.freedesktop.org/NetworkManager/NetworkManager) as the main network manager
- [`Go`](https://github.com/golang/go) ![Go Version](https://img.shields.io/github/go-mod/go-version/alphameo/nm-tui?label=)
- (optional) `xdg-open` + `ip` on Linux -- opens captive portal for connecting to the public WiFi-networks
- [Nerd Font](https://www.nerdfonts.com/font-downloads)

## Installation

### Go install

Requires [Go](https://github.com/golang/go) installed.

```bash
go install github.com/alphameo/nm-tui@latest
```

The binary is placed in `$(go env GOPATH)/bin`.

```bash
# Run directly
$(go env GOPATH)/bin/nm-tui

# If you have the `$(go env GOPATH)/bin` directory in your `PATH` just run
nm-tui
```

To use just `nm-tui` add the `$(go env GOPATH)/bin` directory to your `PATH`

### Nix

With [Nix](https://nixos.org) and flakes enabled:

```bash
nix profile install github:alphameo/nm-tui
```

Or install into your profile:

```bash
nix run github:alphameo/nm-tui
```

### Binary

Move binary `nm-tui` into `/usr/bin` or add location to `PATH` if you want system-wide access.

### Prebuilt binary

Load the latest archive for your architecture from the
[releases page](https://github.com/alphameo/nm-tui/releases).

#### Manual

Clone repo

```bash
git clone https://github.com/alphameo/nm-tui.git
# or
git clone git@github.com:alphameo/nm-tui.git
```

Generate binary:

```bash
make deps
make build
# or
make clean-build
```

Binary generated at `./bin/nm-tui`

## Configuration

Config is placed at `$XDG_CONFIG_HOME/nm-tui/config.kdl` (e.g. `~/.config/nm-tui/config.kdl`).

All settings have default values, with which the user configuration is subsequently merged.

You don't need to copy the defaults — a fully-commented example covering every option is available in [`config.example.kdl`](./config.example.kdl). Only include the sections you want to override, e.g.:

```kdl
colors {
    accent "#865fff"
}
keys {
    main {
        quit "esc" "ctrl+c" "q" "ctrl+q"
    }
}
```

## Tech Stack

- Programming language [Go](https://github.com/golang/go) ![Go Version](https://img.shields.io/github/go-mod/go-version/alphameo/nm-tui?label=)
- TUI framework [Bubbletea](https://github.com/charmbracelet/bubbletea) with [Bubbles](https://github.com/charmbracelet/bubbles) and [Lipgloss](https://github.com/charmbracelet/lipgloss)

## Contributing

Pull requests are welcome! Please open an issue first to discuss what you would like to change.

## License

[MIT](LICENSE)

## Inspirations

- [`Lazygit`](https://github.com/jesseduffield/lazygit)
- [`Lazydocker`](https://github.com/jesseduffield/lazydocker)
- [`impala`](https://github.com/pythops/impala)
