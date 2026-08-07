# nm-tui

![License](https://img.shields.io/github/license/alphameo/nm-tui)
![Go Version](https://img.shields.io/github/go-mod/go-version/alphameo/nm-tui)
![GitHub stars](https://img.shields.io/github/stars/alphameo/nm-tui)

Lightweight TUI wrapper for [NetworkManager](https://gitlab.freedesktop.org/NetworkManager/NetworkManager)

Why `nm-tui`: built-in `nmtui` doesn't look great and there aren't many TUI alternatives

## 📌 Table of Contents

- [💫 Features](#-features)
- [📹 Demo](#-demo)
- [🖼️ Screenshots](#%EF%B8%8F-screenshots)
- [🗃️ Requirements](#%EF%B8%8F-requirements)
- [📥 Installation](#-installation)
- [⚙️ Configuration](#-configuration)
- [⚙️ Tech Stack](#%EF%B8%8F-tech-stack)
- [🖲️ Contributing](#%EF%B8%8F-contributing)
- [⚖️ License](#%EF%B8%8F-license)
- [⭐ Inspirations](#-inspirations)

## 💫 Features

- 😎 TUI style looks cool
- 📡 Scan and list available WiFi-networks
- 🔑 Connect to WiFi-networks with password input
- 📜 View detailed network information (signal strength, security, etc.)
- 🌐 Control device networking
- 📡 Create hotspot
- 🖥️ Clean, modern TUI built with Bubbletea
- ⚡ Fast and lightweight — single static binary
- 🐧 Linux only — designed specifically for NetworkManager

## 📹 Demo

![Demo](../assets/demo-conn.gif)

## 🖼️ Screenshots

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

## 🗃️ Requirements

- [`NetworkManager`](https://gitlab.freedesktop.org/NetworkManager/NetworkManager) as the main network manager
- [`Go`](https://github.com/golang/go) v1.24.3
- (optional) `xdg-open` + `ip` on Linux -- opens captive portal for connecting to the public WiFi-networks
- [Nerd Font](https://www.nerdfonts.com/font-downloads)

## 📥 Installation

### Manual

#### Clone repo

```bash
git clone https://github.com/alphameo/nm-tui.git
```

or

```bash
git clone git@github.com:alphameo/nm-tui.git
```

#### Generate binary

```bash
make deps
make build
```

or

```bash
make clean-build
```

#### Use binary

```bash
./bin/nm-tui
```

## ⚙️ Configuration

Config is placed at `$XDG_CONFIG_HOME/nm-tui/config.kdl` (e.g. `~/.config/nm-tui/config.kdl`).

All settings have default values, with which the user configuration is subsequently merged.

You don't need to copy default config.

<details>
    <summary>Default <code>config.kdl</code> for reference (no need to copy)</summary>

```kdl
notification_close_time 50

// Colors support:
// 1. rgb-format: e.g. "#000000"
// 2. default value: "default"
// 3. colors, defined in your terminal:
//    - "none" - bg color for background, text color for foreground
//    - "black", "red", "green", "yellow",
//      "blue", "magenta", "cyan", "white",
//      "bright_black", "bright_red", "bright_green", "bright_yellow",
//      "bright_blue", "bright_magenta", "bright_cyan", "bright_white",
colors {
    text "none"
    accent "blue"
    muted "bright_black"
    error "red"
    notification "yellow"
}

icons {
    nerd_preset false // set icons to nerd variant

    border_style "ascii" // default for nerd: "rounded"
                         // non-nerd variants: "ascii", "markdown",
                         // nerd variants: "rounded", "square", "thick_square",
                         //                "double_square", "block",
                         //                "outer_half_block", "inner_half_block"

    spinner_style "line" // default for nerd: "meter"
                         // non-nerd variants: "line", "ellipsis"
                         // nerd variants: "dot", "mini_dot", "jump", "pulse",
                         //                "points", "meter", "hamburger"

    input_cursor_shape "bar" // variants: "bar", "underline", "block"

    toggle_off "[ ]"              // default for nerd: " "
    toggle_on "[x]"               // default for nerd: " "
    password_hidden_character "*" // default for nerd: "•", limited to 1 character
    error "!"                     // default for nerd: "✗"
    check "v"                     // default for nerd: " "
    connection "sig"              // default for nerd: " "
    signal "con"                  // default for nerd: "󱘖 "
    access_point "ap"             // default for nerd: "󰀃 "
    infra "infr"                  // default for nerd: "🖳 "
    mesh "#"                      // default for nerd: " "
    ad_hoc "ah"                   // default for nerd: ""
}

logging {
    level "error" // variants: "debug", "info", "warn", "error"
    file_path "~/.local/state/nm-tui/nm-tui.log"
}

// Every mapping can be present in several variants
// Overlapping: dialog -> main -> no-section -> other...
keys {
    toggle "space"
    rescan "r"
    rescan_focused "ctrl+r"
    focus_1 "1"
    focus_2 "2"
    focus_3 "3"
    focus_4 "4"
    focus_5 "5"
    focus_6 "6"
    focus_7 "7"
    focus_8 "8"
    focus_9 "9"
    focus_10 "10"
    main {
        next_tab "]"
        prev_tab "["
        focus_next "tab"
        focus_prev "shift+tab"
        quit "esc" "ctrl+c" "q" "ctrl+q"
    }
    dialog {
        focus_down "ctrl+j"
        focus_up "ctrl+k"
        toggle_pw_visibility "ctrl+p"
        accept "ctrl+enter"
        close "ctrl+q"
    }
    wifi {
        create_profile "a"
        open_network_login "l"
        enable_hotspot "ctrl+h"
        create_hotspot "h"
    }
    wifi_available {
        connect "enter"
    }
    wifi_saved {
        edit "enter"
        connect "space"
        disconnect "ctrl+space"
        delete "d" "delete"
    }
}
```

</details>

## ⚙️ Tech Stack

- Programming language [Go](https://github.com/golang/go) v1.26.3
- TUI framework [Bubbletea](https://github.com/charmbracelet/bubbletea) with [Bubbles](https://github.com/charmbracelet/bubbles) and [Lipgloss](https://github.com/charmbracelet/lipgloss)

## 🖲️ Contributing

Pull requests are welcome! Please open an issue first to discuss what you would like to change.

## ⚖️ License

This project is licensed under the [MIT License](LICENSE).

## ⭐ Inspirations

- [`Lazygit`](https://github.com/jesseduffield/lazygit)
- [`Lazydocker`](https://github.com/jesseduffield/lazydocker)
- [`impala`](https://github.com/pythops/impala)
