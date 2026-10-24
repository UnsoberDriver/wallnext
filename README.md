# wallnext

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

**Press Alt+Shift to go to the next wallpaper.**

Minimal wallpaper cycler for Wayland + KDE Plasma. Pick a folder of images, bind `Alt+Shift`, done.

No GUI, no daemon with 40 options. Just `next` / `random` / `prev` via hotkey or CLI.

## Why

Most wallpaper changers are either too heavy (Waypaper), X11-only (Variety), or require a full GUI to do something trivial. `wallnext` does one thing: cycle wallpapers from a directory with a global hotkey.

## Supported compositors/wallpaper setters

| Backend | Status | Notes
| --- | --- | --- |
| [swww](https://github.com/LGFae/swww) | **Recommended** | Works great on Wayland, smooth transitions, very popular on Hyprland/Sway/KDE Wayland |
| [hyprpaper](https://github.com/hyprwm/hyprpaper) | Supported | Good for Hyprland |
| [waypaper](https://github.com/anufrievroman/waypaper) | Supported (via IPC/existing setup) | If you already use waypaper |
| `plasma-apply-wallpaperimage` / KDE | Supported | Uses `plasma-apply-wallpaperimage` for KDE Plasma on Wayland/X11 |
| `feh` | Supported | X11 |
| `nitrogen` | Supported | X11 |

## Install

```bash
go install github.com/UnsoberDriver/wallnext/cmd/wallnext@latest
```

Or build from source:

```bash
git clone https://github.com/UnsoberDriver/wallnext
cd wallnext
go build ./cmd/wallnext
```

## Quick start

```bash
# 1) Create config
mkdir -p ~/.config/wallnext
cat > ~/.config/wallnext/config.toml << 'EOF'
[settings]
folder = "~/Images/Wallpapers"
backend = "swww"  # or "kde", "hyprpaper", "feh", ...
mode = "random"   # "random" or "sequential"
extensions = ["jpg", "jpeg", "png", "webp", "gif"]
EOF

# 2) Test it
wallnext next

# 3) Bind Alt+Shift
# KDE: System Settings > Shortcuts > Custom Shortcuts > Add Command > `wallnext next`
# Hyprland: bind = ALT SHIFT, N, exec, wallnext next
# Sway: bindsym Mod1+Shift+w exec wallnext next
```

## Commands

| Command | Description |
| --- | --- |
| `wallnext next` | Go to next wallpaper (respects `mode`) |
| `wallnext random` | Pick a random wallpaper |
| `wallnext prev` | Go to previous wallpaper |
| `wallnext list` | List wallpapers in the folder |
| `wallnext set <path>` | Set a specific wallpaper |
| `wallnext init` | Print a starter config |

## Config

```toml
[settings]
folder = "~/Images/Wallpapers"
backend = "auto"  # auto-detect swww/hyprpaper/kde/feh/nitrogen
mode = "random"   # random or sequential
recursive = false
extensions = ["jpg", "jpeg", "png", "webp", "avif"]

[swww]
transition_type = "fade"
transition_duration = 1.0
```

## License

MIT
