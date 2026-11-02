#!/bin/sh
set -eu
repo=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
sh "$repo/scripts/install-discord-bridge.sh" "$@"
destination="$HOME/.local/bin/try-arch-record"
if [ -e "$destination" ] || [ -L "$destination" ]; then
    [ ! -L "$destination" ] || { echo "Refusing to replace symlink: $destination" >&2; exit 1; }
    cp -p -- "$destination" "$destination.before-tryarch-$(date +%Y%m%d-%H%M%S)"
fi
install -m755 "$repo/guest-build/overlay/usr/local/bin/try-arch-record" "$destination"
echo 'CPU recording installed: try-arch-record (inside a compatible Wayland session).'
echo 'If needed, install recording dependencies with: sudo pacman -S wf-recorder slurp'
