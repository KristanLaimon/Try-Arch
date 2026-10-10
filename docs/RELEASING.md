# Releasing Try Arch

The manual **Build Try Arch release** workflow builds and boots the vanilla
guest, verifies the pinned runtime archives, builds the Windows launcher with
that exact manifest, and creates a preview draft containing all required assets.
It does not publish the draft or sign the executable. Automatic updates are off.

Build the guest using `scripts/release/build-guest.sh --output DIR`. Add the
verified runtime with `scripts/release/prepare-assets.sh DIR`. Boot and validate
the finished guest; see TESTING.md. Copy LICENSE and THIRD_PARTY_NOTICES.md into
DIR and add their checksums to SHA256SUMS. Build the pinned release launcher:

```powershell
./scripts/release/build-launcher.ps1 -Artifacts release/v0.1.0 -ReleaseUrl https://github.com/KristanLaimon/Try-Arch/releases/download/v0.1.0
```

The release tag must match currentVersion and the compiled Windows resources.
Bump app/update.go, app/versioninfo.rc and regenerate the resource object together
for a new version. Write .github/release-notes/TAG.md before creating its draft.

Publish the compressed guest, boot files, metadata, runtime and SHA256SUMS only
after validation. Keep the uncompressed rootfs.ext4 hash in the manifest for
decompression verification; the 4 GiB raw image exceeds GitHub's asset limit.
The build helper embeds SHA256SUMS and pins both guest and runtime downloads to
the release URL using Go linker flags. Development builds keep guest downloads
disabled. It updates the checked-in guest manifest and its digest together;
commit those files with the release source. SHA256SUMS excludes the executable
to avoid a circular checksum;
TryArch.exe.sha256 contains its separate digest. Keep checksum and signature
validation enabled. Automatic updates require this fork's own signing key and
signed manifest feed; do not enable the predecessor's update feed.

After reviewing all assets in the draft, publish it on GitHub. The release URL
embedded in TryArch.exe becomes downloadable only once the draft is published.
Remove raw rootfs.ext4, smoke logs and temporary test/build resources after
verification. Keep the compressed release assets that users need to install.

Preserve the runtime licenses and corresponding-source archive. The runtime
lock continues to reference authenticated existing QEMU assets; this does not
install any predecessor guest software.
