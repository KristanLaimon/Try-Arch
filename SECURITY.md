# Security

Report security issues privately through the repository's Security tab:
https://github.com/KristanLaimon/Try-Arch/security . Include reproduction steps
and the launcher, guest and runtime versions. Do not share guest disks, backups,
private keys or unreviewed diagnostics.

The launcher authenticates SHA256 manifests before using guest or runtime
downloads. Custom releases require an independent checksum pin. Automatic
updates are disabled until this fork has its own signing key and signed feed.
Locally built TryArch.exe is unsigned; do not expect the predecessor's signature.

The guest uses official signed Arch packages. Root login is locked, SSH is
off unless requested, and SSH password authentication is disabled. Quick start
uses arch / arch locally; change its password with passwd.

Arch and QEMU retain their upstream security policies. This project owns the
launcher, guest configuration and integration changes it ships.
