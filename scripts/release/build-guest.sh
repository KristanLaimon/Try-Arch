#!/usr/bin/env bash
set -euo pipefail
repo_root=$(cd "$(dirname "$0")/../.." && pwd)
if [[ ${1:-} == --contract-only ]]; then
  python3 -m unittest discover -s "$repo_root/guest-build" -p 'test_*.py'
  for script in "$repo_root/guest-build/build.sh" "$repo_root/guest-build/overlay/usr/local/bin/try-arch-setup" "$repo_root/guest-build/overlay/usr/local/bin/try-arch-share" "$repo_root/scripts/guest/clipboard-bridge.sh"; do
    bash -n "$script"
  done
  exit
fi
[[ ${1:-} == --output && -n ${2:-} && $# == 2 ]] || {
  echo 'Usage: build-guest.sh --output DIR | --contract-only' >&2; exit 2;
}
mkdir -p "$2"
output=$(cd "$2" && pwd)
docker run --rm --privileged -v "$repo_root:/src:ro" -v "$output:/output" \
  --mount type=volume,destination=/var/cache/pacman/pkg archlinux:base \
  bash -c 'pacman -Syu --noconfirm arch-install-scripts btrfs-progs python zstd && bash /src/guest-build/build.sh /output && chown --reference=/output /output/*'
