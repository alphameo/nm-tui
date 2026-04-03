# nm-tui

![License](https://img.shields.io/github/license/alphameo/nm-tui)
![Go Version](https://img.shields.io/github/go-mod/go-version/alphameo/nm-tui)
![GitHub stars](https://img.shields.io/github/stars/alphameo/nm-tui)

Lightweight TUI wrapper for [NetworkManager](https://gitlab.freedesktop.org/NetworkManager/NetworkManager)

Why `nm-tui`: builtin `nmtui` doesn't look great and there aren't many TUI alternatives

## 📌 Table of Contents

- [💫 Features](#features)
- [📹 Demo](#demo)
- [🖼️ Screenshots](#screenshots)
- [🗃️ Requirements](#requirements)
- [📥 Installation](#installation)
- [⚙️ Tech Stack](#tech-stack)
- [🖲️ Contributing](#contributing)
- [⚖️ License](#license)
- [⭐ Inspirations](#inspirations)

## 💫 Features

- 😎 TUI style looks cool
- 📡 Scan and list available WiFi networks
- 🔑 Connect to WiFi networks with password input
- 📋 View detailed network information (signal strength, security, etc.)
- 🖥️ Clean, modern TUI built with Bubbletea
- ⚡ Fast and lightweight — single static binary
- 🐧 Linux only — designed specifically for NetworkManager

## 📹 Demo

![Demo](../assets/demo1.gif)

## 🖼️ Screenshots

### Main window

<img src="../assets/main.png" alt="main window" />

### Wifi connector and info

<div style="display: flex; gap: 10px;">
    <img src="../assets/wifi-connector.png" alt="wifi connector" width="400"/>
    <img src="../assets/wifi-info.png" alt="wifi info" width="400"/>
</div>

## 🗃️ Requirements

- [NetworkManager](https://gitlab.freedesktop.org/NetworkManager/NetworkManager) as the main network manager
- [Go](https://github.com/golang/go) v1.24.4

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

## ⚙️ Tech Stack

- Programming language [Go](https://github.com/golang/go) v1.24.4
- TUI framework [Bubbletea](https://github.com/charmbracelet/bubbletea) with [Bubbles](https://github.com/charmbracelet/bubbles) and [Lipgloss](https://github.com/charmbracelet/lipgloss)

## 🖲️ Contributing

Pull requests are welcome! Please open an issue first to discuss what you would like to change.

## ⚖️ License

This project is licensed under the [MIT License](LICENSE).

## ⭐ Inspirations

- [`Lazygit`](https://github.com/jesseduffield/lazygit)
- [`Lazydocker`](https://github.com/jesseduffield/lazydocker)
- [`impala`](https://github.com/pythops/impala)
