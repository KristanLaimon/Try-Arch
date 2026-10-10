# Third-party notices

Try Arch is a fork of the MIT-licensed Try Omarchy Windows launcher. Its original
copyright and license in LICENSE are preserved. Original launcher sources:
https://github.com/omacom/try-omarchy-windows . No Omarchy software or artwork is
installed in the vanilla guest.

- Arch Linux packages retain their individual licenses; the build records the
  installed versions in packages.txt. The guest uses the official repositories.
- The Arch Linux icon comes from https://archlinux.org/art/ and remains subject
  to https://terms.archlinux.org/docs/trademark-policy/ . Try Arch is an independent
  project and is not an official Arch Linux product.
- QEMU and the WINQ-EMU QEMU fork retain GPL-2.0 and other component licenses.
  Distribute the corresponding-source archive with the runtime.
- virglrenderer is MIT licensed. MSYS2 and the runtime's other libraries retain
  their original licenses; see the runtime archive's inventory and notices.
- klauspost/compress is BSD-3-Clause and is used by the launcher for decompression.
- dockur/windows is MIT licensed and is used only for development.

See guest-build/runtime.lock.json and runtime-build/sources.lock.json for
runtime provenance. Historical checksum manifests in app/testdata retain the
original asset names for authentication tests; they are not an Arch guest pin.
