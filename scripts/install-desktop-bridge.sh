#!/bin/sh
# Run from the guest's desktop terminal, without sudo.
set -eu
case "${1:-}" in ''|--offline) ;; *) echo 'Usage: install-desktop-bridge.sh [--offline]' >&2; exit 2 ;; esac
[ "$(id -u)" -ne 0 ] || { echo 'Run as your Linux desktop user, without sudo.' >&2; exit 1; }
repo=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
bin="$HOME/.local/bin"
units="${XDG_CONFIG_HOME:-$HOME/.config}/systemd/user"
stamp=$(date +%Y%m%d-%H%M%S)-$$
for name in clipboard-bridge.sh clipboard-files file-transfer try-arch-drop; do
    [ -f "$repo/scripts/guest/$name" ] || { echo "Missing helper: $name" >&2; exit 1; }
done

if [ "${1:-}" != --offline ]; then
    systemctl --user show-environment >/dev/null
    packages=''
    command -v wl-copy >/dev/null && command -v wl-paste >/dev/null || packages="$packages wl-clipboard"
    command -v socat >/dev/null || packages="$packages socat"
    command -v python3 >/dev/null || packages="$packages python"
    command -v xdg-user-dir >/dev/null || packages="$packages xdg-user-dirs"
    command -v notify-send >/dev/null || packages="$packages libnotify"
    if [ -n "$packages" ]; then
        echo "Installing missing dependencies:$packages"
        # Package names above are fixed; only package installation needs sudo.
        sudo pacman -S --needed $packages
    fi
fi

mkdir -p "$bin" "$units"
backup() {
    if [ -e "$1" ] || [ -L "$1" ]; then
        [ ! -L "$1" ] || { echo "Refusing to replace symlink: $1" >&2; exit 1; }
        cp -p -- "$1" "$1.before-tryarch-$stamp"
    fi
}
for name in clipboard-bridge clipboard-files file-transfer try-arch-drop; do
    backup "$bin/$name"
done
backup "$units/try-arch-clipboard.service"
install -m755 "$repo/scripts/guest/clipboard-bridge.sh" "$bin/clipboard-bridge"
for name in clipboard-files file-transfer try-arch-drop; do
    install -m755 "$repo/scripts/guest/$name" "$bin/$name"
done
cat > "$units/try-arch-clipboard.service" <<'UNIT'
[Unit]
Description=Clipboard and file-drop bridge to Windows (Try Arch)

[Service]
Environment="PATH=%h/.local/bin:/usr/local/bin:/usr/bin:/bin"
ExecStart=%h/.local/bin/clipboard-bridge
Restart=always
RestartSec=3

[Install]
WantedBy=default.target
UNIT
chmod 644 "$units/try-arch-clipboard.service"
if [ "${1:-}" = --offline ]; then
    mkdir -p "$units/default.target.wants"
    target="$units/default.target.wants/try-arch-clipboard.service"
    if [ ! -e "$target" ] && [ ! -L "$target" ]; then
        ln -s ../try-arch-clipboard.service "$target"
    fi
    echo 'Bridge installed for next login. Requires python, socat, wl-clipboard, xdg-user-dirs and libnotify.'
else
    # Retire the early launcher's differently named service to avoid two watchers.
    if systemctl --user cat clipboard-bridge.service >/dev/null 2>&1; then
        systemctl --user disable --now clipboard-bridge.service
    fi
    systemctl --user daemon-reload
    for variable in WAYLAND_DISPLAY HYPRLAND_INSTANCE_SIGNATURE XDG_CURRENT_DESKTOP; do
        if printenv "$variable" >/dev/null; then
            systemctl --user import-environment "$variable"
        fi
    done
    if [ ! -f "${XDG_CONFIG_HOME:-$HOME/.config}/user-dirs.dirs" ]; then
        xdg-user-dirs-update
    fi
    systemctl --user enable try-arch-clipboard.service
    systemctl --user restart try-arch-clipboard.service
    systemctl --user is-active --quiet try-arch-clipboard.service
    echo 'Bridge running. Copy text/images/files between Windows and Wayland, or drop Windows files onto the VM.'
    echo 'Drops go to Downloads unless a supported file-manager destination can be identified.'
fi
echo 'Replaced files were backed up alongside their originals as *.before-tryarch-*.'
