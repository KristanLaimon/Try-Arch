# Vanilla Arch guest

`scripts/release/build-guest.sh --output DIR` uses a privileged Linux Docker
container and `pacstrap -K` to install `base`, `linux`, `linux-firmware`,
`e2fsprogs`, `openssh`, `sudo`, and `nano` from the official signed repositories.
No custom repository, AUR helper or desktop is installed.

The factory image is a 4 GiB ext4 filesystem. The launcher expands its working
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

Artifacts: `rootfs.ext4`, `rootfs.ext4.zst`, `vmlinuz-linux`,
`initramfs-linux.img`, `build-spec.json`, `guest-manifest.json`, `packages.txt`,
and `SHA256SUMS`. The package inventory records the build transaction.

This is a rolling build: the Docker base and official packages resolve at build
time. Save the emitted artifacts and checksum pin for reproducible deployment;
rebuilding later may select newer packages.

`--contract-only` runs the local boot-contract checks and shell syntax validation
without downloading packages or building an image.
