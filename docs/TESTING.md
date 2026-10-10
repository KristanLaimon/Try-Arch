# Testing

With Python available on PATH, run `go test ./...` in `app/` and build
`TryArch.exe` for Windows amd64. The tests run the real Python guest builder
against tiny fixtures through the download/install path, and validate the
actual pinned release metadata against the launcher's schema and checksums.
Run `python -m unittest discover -s guest-build -p 'test_*.py'`.
Run `bash scripts/release/build-guest.sh --contract-only` for shell syntax and
guest metadata checks. These checks do not prove that a Linux image boots.

On Linux with QEMU and Python installed, run `python guest-build/smoke.py DIR`
against the generated artifacts. It uses a temporary copy-on-write disk and
checks quick-start login, sudo, official packages, SSH policy, networking,
package updates, disk growth, matching kernel modules and clean shutdown.
It writes `smoke.log` to DIR and removes the temporary VM disk when finished.

Build and boot the complete guest before publishing. Check personal account
setup, quick start, normal login, sudo, root lock, package updates, networking,
SSH with a public key, disk growth, shutdown and reboot. Check both the CPU
QEMU runtime and the accelerated Windows runtime. Inspect `/etc/pacman.conf`
and `pacman -Q` for official packages only.

Old desktop-specific acceptance records were removed because they do not
validate this vanilla image. Test Windows backup, restore, snapshots and reset
against the new image before making a release claim.
