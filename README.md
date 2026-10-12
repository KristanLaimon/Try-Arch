<p align="center"><img src="app/ArchIcon.svg" width="88" height="88" alt="Try Arch icon"></p>
<h1 align="center">Try Arch</h1>
<p align="center">Vanilla Arch Linux in a virtual machine on Windows.</p>

Try Arch runs an x86_64 Arch Linux guest with a terminal login. It installs
packages from the official Arch repositories. Choose your own account during
first boot, or select quick start to use `arch` / `arch`. Your Linux files
persist between sessions in a separate `%LOCALAPPDATA%\TryArch` installation.

There is no preinstalled desktop, theme collection, custom package repository,
or distribution updater. Install your preferred software using `pacman`.
Update Linux with `sudo pacman -Syu`.

Windows 10/11 on an Intel or AMD x86_64 PC and Windows Hypervisor Platform are
required. The launcher retains VM storage, backup, restore, snapshot, portable
copy, resource, USB and port-forward controls. The guest uses Btrfs and includes
text clipboard commands plus a Wayland clipboard bridge. Camera, Windows Hello
and application bridges require additional guest software.

## Build and run

To produce the downloadable executable and matching guest assets, use the
**Build Try Arch release** workflow or follow [RELEASING.md](docs/RELEASING.md).
The workflow creates a preview draft for review before publishing.

Build the Windows launcher from `app/`:

```powershell
go build -trimpath -ldflags "-H windowsgui -s -w" -o ../TryArch.exe .
```

Build the guest on Linux, or with Git Bash and Linux Docker on Windows:

```sh
bash scripts/release/build-guest.sh --output artifacts/arch-guest
bash scripts/release/prepare-assets.sh artifacts/arch-guest
```

The runtime assets are verified against `guest-build/runtime.lock.json`.
The guest image has no public release pin yet. Serve or publish the built
artifacts, compute the SHA256 of `SHA256SUMS`, and launch with that independent
trust pin:

For a local build, serve the artifacts in one terminal:

```powershell
python -m http.server 8765 --bind 127.0.0.1 --directory artifacts/arch-guest
```

Then launch from another terminal:

```powershell
$pin = (Get-FileHash artifacts/arch-guest/SHA256SUMS -Algorithm SHA256).Hash.ToLowerInvariant()
.\TryArch.exe -release http://127.0.0.1:8765 -sums-sha256 $pin
```

Default guest downloads and launcher updates are disabled until verified
Try Arch releases are published. The launcher will never use the former
distribution image as a fallback. Existing installations from that project
are left in their original data directory; back them up separately.

Read the [user guide](docs/USER-GUIDE.md), [guest build](guest-build/README.md),
[testing instructions](docs/TESTING.md), and [release instructions](docs/RELEASING.md).
The original launcher is MIT licensed; other components retain their own
licenses. See [third-party notices](THIRD_PARTY_NOTICES.md).
