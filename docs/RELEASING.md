# Releasing Try Arch

The manual Build workflow uploads the Windows launcher and the vanilla guest
artifacts for review. It does not publish a release or sign the launcher.
No public guest release or automatic-update feed is configured yet.

Build the guest using `scripts/release/build-guest.sh --output DIR`. Add the
verified runtime with `scripts/release/prepare-assets.sh DIR`. Build and test
the launcher. Boot and validate the finished guest; see TESTING.md.

Publish the compressed guest, boot files, metadata, runtime and SHA256SUMS only
after validation. Keep the uncompressed rootfs.ext4 hash in the manifest for
decompression verification; the 4 GiB raw image exceeds GitHub's asset limit.
Set `defaultReleaseURL` to the actual release URL and `defaultSumsSHA256` to
the independently computed SHA256 of its manifest, then rebuild the launcher.
Replace the embedded manifest when changing the guest pin. Keep checksum and
signature validation enabled. Automatic updates require this fork's own signing
key and signed manifest feed; do not enable the predecessor's update feed.

Preserve the runtime licenses and corresponding-source archive. The runtime
lock continues to reference authenticated existing QEMU assets; this does not
install any predecessor guest software.
