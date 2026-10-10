# Contributing

Describe the problem, the user-visible change and your validation. Use the Go
version in app/go.mod. Run `go test ./...` from app/ and build TryArch.exe.
Run `bash scripts/release/build-guest.sh --contract-only` for guest changes.
Build and boot the full guest before claiming image compatibility.

Repository: app/ is the Windows launcher; guest-build/ contains the vanilla
Arch build and overlay; runtime-build/ contains the QEMU recipe; scripts/
contains development tools. See docs/TESTING.md and docs/RELEASING.md.

Report launcher issues at https://github.com/KristanLaimon/Try-Arch/issues .
Review diagnostics before sharing them. Never upload guest disks or backups.
