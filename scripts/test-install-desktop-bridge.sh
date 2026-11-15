#!/usr/bin/env bash
# Exercise installation against a disposable home and mocked desktop services.
set -euo pipefail
repo=$(cd "$(dirname "$0")/.." && pwd)
work=$(mktemp -d "${TMPDIR:-/tmp}/tryarch-bridge-test.XXXXXX")
trap 'case "$work" in "${TMPDIR:-/tmp}"/tryarch-bridge-test.*) rm -rf -- "$work" ;; esac' EXIT
mkdir -p "$work/commands" "$work/home/.local/bin" "$work/home/.config/systemd/user"
cat > "$work/commands/id" <<'MOCK'
#!/bin/sh
echo 1000
MOCK
cat > "$work/commands/systemctl" <<'MOCK'
#!/bin/sh
printf '%s\n' "$*" >> "$BRIDGE_TEST_LOG"
[ "$*" != '--user cat clipboard-bridge.service' ]
MOCK
for name in wl-copy wl-paste socat python3 xdg-user-dir xdg-user-dirs-update notify-send; do
    printf '#!/bin/sh\nexit 0\n' > "$work/commands/$name"
done
chmod +x "$work/commands/"*
export BRIDGE_TEST_LOG="$work/systemctl.log"
export HOME="$work/home" XDG_CONFIG_HOME="$work/home/.config"
export PATH="$work/commands:$PATH"
printf 'previous helper\n' > "$HOME/.local/bin/clipboard-bridge"
printf 'previous unit\n' > "$XDG_CONFIG_HOME/systemd/user/try-arch-clipboard.service"
printf 'important personal file\n' > "$HOME/keep-me"
sh "$repo/scripts/install-desktop-bridge.sh"
cmp "$repo/scripts/guest/clipboard-bridge.sh" "$HOME/.local/bin/clipboard-bridge"
for name in clipboard-files file-transfer try-arch-drop; do
    cmp "$repo/scripts/guest/$name" "$HOME/.local/bin/$name"
    test -x "$HOME/.local/bin/$name"
done
grep -qx 'previous helper' "$HOME/.local/bin/clipboard-bridge.before-tryarch-"*
grep -qx 'previous unit' "$XDG_CONFIG_HOME/systemd/user/try-arch-clipboard.service.before-tryarch-"*
grep -qx 'important personal file' "$HOME/keep-me"
grep -qx 'WantedBy=default.target' "$XDG_CONFIG_HOME/systemd/user/try-arch-clipboard.service"
grep -qx -- '--user restart try-arch-clipboard.service' "$BRIDGE_TEST_LOG"
grep -qx -- '--user is-active --quiet try-arch-clipboard.service' "$BRIDGE_TEST_LOG"
# Reinstallation must keep earlier backups and the user's own files.
sh "$repo/scripts/install-desktop-bridge.sh"
grep -qx 'previous helper' "$HOME/.local/bin/clipboard-bridge.before-tryarch-"*
grep -qx 'important personal file' "$HOME/keep-me"
echo 'PASS: desktop bridge installs, starts, backs up replacements and preserves personal files.'
