# Vanilla Arch guest

`scripts/release/build-guest.sh --output DIR` uses a privileged Linux Docker
container and `pacstrap -K` to install `base`, `linux`, `linux-firmware`,
`btrfs-progs`, `openssh`, `sudo`, `nano`, `python`, `socat`, and `wl-clipboard`
from the official signed repositories.
No custom repository, AUR helper or desktop is installed.

The factory image is a 4 GiB Btrfs filesystem. The launcher expands its working
disk to 16 GiB by default; the guest grows the filesystem on boot. QEMU boots
the matching external kernel and initramfs with `root=/dev/vda`. The graphical
VM window shows a text console, including the first-boot account form.

The launcher boots external kernel files. If pacman installs a newer kernel,
publish matching kernel/initramfs assets and their authenticated manifest before
the next launch. Rebuilding only the launcher does not update those boot files.

Quick start creates `arch` / `arch`, with password-required sudo. Personal
setup asks for a valid unused username and password at tty1. Root login stays
locked. SSH is off unless requested by a forwarded port; password authentication
is disabled, and the launcher's public key is validated before installation.

Networking uses systemd-networkd DHCP and systemd-resolved. The default locale
is en_US.UTF-8 and time zone is UTC; change these using normal Arch tools.

`try-arch-share` mounts the launcher's selected Windows folder at `/mnt/host`
after account setup and links it into the account's home using the folder name
passed in `tryarch.sharename`. It preserves existing files at that name and
skips boots without a selected share.

Copy console text to Windows with `printf 'hello' | try-arch-clipboard copy`,
or `cat file.txt | try-arch-clipboard copy`. Read Windows clipboard text with
`try-arch-clipboard paste` (redirect it to a file if needed). Text is limited
to 8 MiB. The text console has no desktop selection clipboard or Ctrl+V paste.
After installing a Wayland desktop that starts `graphical-session.target` and
supports wl-clipboard, the included user service shares selections automatically.
Other Wayland sessions can start it with `systemctl --user start try-arch-clipboard`.

Artifacts: `rootfs.ext4`, `rootfs.ext4.zst`, `vmlinuz-linux`,
`initramfs-linux.img`, `build-spec.json`, `guest-manifest.json`, `packages.txt`,
and `SHA256SUMS`. The package inventory records the build transaction.
The legacy `rootfs.ext4` artifact names are retained for launcher, backup, and
release compatibility; their contents are Btrfs. Existing disks are not converted.

This is a rolling build: the Docker base and official packages resolve at build
time. Save the emitted artifacts and checksum pin for reproducible deployment;
rebuilding later may select newer packages.

`--contract-only` runs the local boot-contract checks and shell syntax validation
without downloading packages or building an image.
