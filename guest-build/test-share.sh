#!/usr/bin/env bash
set -euo pipefail
case $(uname -s) in MINGW*|MSYS*) export MSYS="${MSYS:-} winsymlinks:nativestrict" ;; esac
script=$(cd "$(dirname "$0")" && pwd)/overlay/usr/local/bin/try-arch-share
test_home=$(mktemp -d)
trap 'rm -rf -- "$test_home"' EXIT
cat() { case $1 in /proc/cmdline) printf 'tryarch.sharename=%s' "$encoded_name" ;; *) printf 'arch\n' ;; esac; }
getent() { printf 'arch:x:1000:1000::%s:/bin/bash\n' "$test_home"; }
mkdir() { :; }
mountpoint() { return 1; }
mount() { printf '%s\n' "$*" > "$test_home/mount-call"; }
chown() { :; }
export -f cat getent mkdir mountpoint mount chown
export test_home
run() { encoded_name=$(printf '%s' "$1" | base64); export encoded_name; bash "$script"; }
run ''
[[ ! -e $test_home/mount-call ]]
run 'VMS SHARED'
[[ $(readlink "$test_home/VMS SHARED") == /mnt/host ]]
grep -q 'trans=virtio,version=9p2000.L,cache=none,nodev,nosuid hostshare /mnt/host' "$test_home/mount-call"
run 'VMS SHARED'
rm "$test_home/VMS SHARED"
printf 'keep me' > "$test_home/VMS SHARED"
run 'VMS SHARED'
[[ $(command cat "$test_home/VMS SHARED") == 'keep me' ]]
rm "$test_home/VMS SHARED"
ln -s missing "$test_home/VMS SHARED"
run 'VMS SHARED'
[[ $(readlink "$test_home/VMS SHARED") == missing ]]
if run '../outside'; then exit 1; fi
