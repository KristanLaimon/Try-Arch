#!/usr/bin/env bash
# Run inside a privileged Arch Linux container; see scripts/release/build-guest.sh.
set -euo pipefail
umask 022
output=${1:?Usage: build.sh OUTPUT}
[[ $EUID == 0 ]] || { echo 'The image build requires root.' >&2; exit 1; }
mkdir -p "$output"
work=$(mktemp -d)
image="$work/rootfs.ext4"
cleanup() {
  pkill gpg-agent || true
  pkill dirmngr || true
  rm -rf -- "$work"
}
trap cleanup EXIT
mkdir "$work/root"
truncate -s 4G "$image"
# The Docker base omits docs and pacman.conf; a normal VM must install them.
sed '/^NoExtract/d' /etc/pacman.conf > "$work/pacman.conf"
pacstrap -C "$work/pacman.conf" -c -K "$work/root" base linux linux-firmware e2fsprogs openssh sudo nano
cp -r --no-preserve=mode,ownership /src/guest-build/overlay/. "$work/root/"
chmod 755 "$work/root/usr/local/bin/try-arch-setup" "$work/root/usr/local/bin/try-arch-share"
printf 'tryarch\n' > "$work/root/etc/hostname"
printf 'en_US.UTF-8 UTF-8\n' > "$work/root/etc/locale.gen"
printf 'LANG=en_US.UTF-8\n' > "$work/root/etc/locale.conf"
printf '/dev/vda / ext4 defaults 0 1\n' > "$work/root/etc/fstab"
printf '%%wheel ALL=(ALL:ALL) ALL\n' > "$work/root/etc/sudoers.d/wheel"
chmod 440 "$work/root/etc/sudoers.d/wheel"
mkdir -p "$work/root/etc/ssh/sshd_config.d"
printf 'PasswordAuthentication no\nPermitRootLogin no\nAuthorizedKeysFile .ssh/authorized_keys /etc/ssh/authorized_keys/%%u\n' > "$work/root/etc/ssh/sshd_config.d/10-try-arch.conf"
printf 'MODULES=(virtio_pci virtio_blk virtio_gpu ext4)\nHOOKS=(base systemd autodetect modconf kms keyboard sd-vconsole block filesystems fsck)\n' > "$work/root/etc/mkinitcpio.conf"
arch-chroot "$work/root" locale-gen
arch-chroot "$work/root" passwd -l root
arch-chroot "$work/root" ln -sf /usr/share/zoneinfo/UTC /etc/localtime
arch-chroot "$work/root" systemctl enable systemd-networkd systemd-resolved try-arch-setup try-arch-share
ln -sf /run/systemd/resolve/stub-resolv.conf "$work/root/etc/resolv.conf"
arch-chroot "$work/root" mkinitcpio -P
# Recent Arch kernels keep the boot image in their module directory.
kernel=$(find "$work/root/usr/lib/modules" -name vmlinuz -print -quit)
# shortcut: QEMU boots external kernel files, update these together after guest kernel upgrades.
[[ -n $kernel ]] || { echo 'Installed kernel is missing.' >&2; exit 1; }
install -m644 "$kernel" "$output/vmlinuz-linux"
install -m644 "$work/root/boot/initramfs-linux.img" "$output/initramfs-linux.img"
arch-chroot "$work/root" pacman -Q > "$output/packages.txt"
rm -f "$work/root/etc/machine-id" "$work/root/var/lib/systemd/random-seed"
rm -f "$work/root/boot/"*.img
gpgconf --homedir "$work/root/etc/pacman.d/gnupg" --kill all
pkill gpg-agent || true
pkill dirmngr || true
mkfs.ext4 -F -L tryarch -d "$work/root" "$image"
e2fsck -fn "$image"
zstd -T0 -10 "$image" -o "$image.zst"
cp --sparse=always "$image" "$output/rootfs.ext4"
cp "$image.zst" "$output/rootfs.ext4.zst"
python /src/guest-build/spec.py "$output"
(cd "$output" && sha256sum build-spec.json guest-manifest.json packages.txt vmlinuz-linux initramfs-linux.img rootfs.ext4 rootfs.ext4.zst > SHA256SUMS)
