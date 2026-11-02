#!/bin/sh
# Run as the Linux desktop user. --offline enables it for the next login.
set -eu
case "${1:-}" in ''|--offline) ;; *) echo 'Usage: install-discord-bridge.sh [--offline]' >&2; exit 2 ;; esac
[ "$(id -u)" -ne 0 ] || { echo 'Run this as your Linux desktop user, without sudo.' >&2; exit 1; }
repo=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
overlay="$repo/guest-build/overlay"
units="${XDG_CONFIG_HOME:-$HOME/.config}/systemd/user"
bin="$HOME/.local/bin"
stamp=$(date +%Y%m%d-%H%M%S)
mkdir -p "$bin" "$units/sockets.target.wants"
backup() {
    if [ -e "$1" ] || [ -L "$1" ]; then
        [ ! -L "$1" ] || { echo "Refusing to replace symlink: $1" >&2; exit 1; }
        cp -p -- "$1" "$1.before-tryarch-$stamp"
    fi
}
backup "$bin/try-arch-discord"
backup "$units/try-arch-discord.socket"
backup "$units/try-arch-discord@.service"
install -m755 "$overlay/usr/local/bin/try-arch-discord" "$bin/try-arch-discord"
install -m644 "$overlay/usr/lib/systemd/user/try-arch-discord.socket" "$units/try-arch-discord.socket"
sed 's|ExecStart=/usr/local/bin/try-arch-discord|ExecStart=%h/.local/bin/try-arch-discord|' \
    "$overlay/usr/lib/systemd/user/try-arch-discord@.service" > "$units/try-arch-discord@.service"
chmod 644 "$units/try-arch-discord@.service"
if [ "${1:-}" = --offline ]; then
    target="$units/sockets.target.wants/try-arch-discord.socket"
    if [ ! -e "$target" ] && [ ! -L "$target" ]; then
        ln -s ../try-arch-discord.socket "$target"
    fi
else
    systemctl --user daemon-reload
    systemctl --user enable --now try-arch-discord.socket
fi
echo 'Discord bridge installed. Enable the Windows bridge with TryArch.exe -discord-rpc.'
