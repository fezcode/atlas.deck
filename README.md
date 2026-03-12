# atlas.deck 🎚️⚡

![Banner Image](./banner-image.png)

**atlas.deck** is a high-performance, interactive TUI command deck for your terminal. Part of the **Atlas Suite**, it allows you to organize, trigger, and monitor your daily project workflows through a customizable grid of "pads."

![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)
![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20Linux%20%7C%20macOS-lightgrey)

## ✨ Features

- 🏗️ **Project-Aware:** Automatically loads `deck.piml` from your current directory or falls back to a global configuration.
- 🎚️ **Interactive Grid:** A responsive TUI layout that maps single keypresses to complex terminal commands.
- 📜 **Integrated Logs:** Capture and view command output directly within the interface in a dedicated viewport.
- 🎨 **Aesthetic Minimalism:** A clean, high-fidelity interface built with Bubble Tea and Lip Gloss.
- 📦 **Cross-Platform:** Pure Go implementation that works seamlessly on Windows (PowerShell) and Linux/macOS (Bash).

## 🚀 Installation

### Recommended: Via Atlas Hub
```bash
atlas.hub
```
Select `atlas.deck` from the list and confirm.

### From Source
```bash
git clone https://github.com/fezcode/atlas.deck
cd atlas.deck
gobake build
```

## ⌨️ Usage

Simply run the binary in any project directory containing a `deck.piml`:
```bash
./atlas.deck
```

### 🚩 Options
| Flag | Description |
|------|-------------|
| `-c`, `--create` | Create a default `deck.piml` file in the current directory. |
| `-h`, `--help` | Show help information. |
| `-v`, `--version` | Show version information. |

### 🕹️ Controls
| Key | Action |
|-----|--------|
| `[Mapped Key]` | **Trigger:** Execute the command associated with that pad. |
| `Ctrl+L` | **Clear:** Wipe the current log viewport. |
| `Ctrl+C` | **Exit:** Close the Mission Control. |

## 🛠️ Blueprint (`deck.piml`)

Define your workflow in a simple PIML file. The deck supports custom labels and colors (`gold`, `cyan`, `green`, `red`).

```piml
(name) "My Project Deck"
(version) "1.0.0"
(pads)
  > (pad)
    (key) "d"
    (label) "Dev Server"
    (cmd) "npm run dev"
    (color) "cyan"

  > (pad)
    (key) "b"
    (label) "Build Prod"
    (cmd) "docker-compose build"
    (color) "gold"
```

## 📄 License
MIT License - see [LICENSE](LICENSE) for details.
