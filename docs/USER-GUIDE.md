# Using Try Arch

Choose your installation location and account mode in setup. The default data
folder is `%LOCALAPPDATA%\TryArch`. First boot opens a console account form.
Quick start uses username `arch` and password `arch`; change the password with
`passwd`. Sign in at the terminal login prompt.

- Update packages: `sudo pacman -Syu`.
- Install packages: `sudo pacman -S package-name`.
- Shut down: `sudo poweroff`.
- Reboot: `sudo reboot`.
- Toggle fullscreen: Ctrl+Alt+F. Release input: Ctrl+Alt+G.

There are no distribution menu or desktop shortcuts. Install and configure a
desktop yourself if you want one. The Arch installation guide and wiki describe
normal account, locale, time-zone, networking and desktop configuration.

QEMU currently boots the kernel and initramfs from the launcher's guest folder.
After a kernel package upgrade, those external boot files must be updated through
a matching verified guest release before rebooting. Keep a backup first.

Backup, restore, snapshots, portable copy, reset, move and uninstall remain in
Settings. Close the VM before changing storage. Reset creates a fresh Arch
disk; save a backup first. Existing data from the predecessor app is not imported.

For SSH, add `tcp:2222:22` and select your SSH public key in the launcher, then
restart and connect using `ssh -p 2222 your-user@localhost`. SSH password login
and root login are disabled. Do not expose the quick-start account on untrusted
networks.

Choose a Windows folder in Settings and restart Arch to access it at `/mnt/host`
and `~/<Windows folder name>`. Existing files at the home-folder link are kept;
use `/mnt/host` if that name is already occupied. Files can be read and written
from Windows and Arch. Older installations need the `try-arch-share` script and
service from `guest-build/overlay` installed once, then enabled with
`sudo systemctl enable --now try-arch-share`.

Guest desktop integrations are not installed. Create diagnostics
from the tray when reporting problems and review the ZIP before sharing it.
