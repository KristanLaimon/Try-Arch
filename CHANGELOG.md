# Changelog

## v0.1.1

- Mount selected Windows shared folders automatically at `/mnt/host` and link them into the Linux account's home, including names with spaces.
- Preserve existing files at the home-folder link and verify normal-user file exchange during guest smoke tests.
- Include a shared-folder repair archive for existing Arch installations.

## Unreleased — Try Arch

- Vanilla Arch terminal guest built from official repositories.
- Try Arch launcher, Windows resources, icon, shortcuts, storage and translations.
- Removed the distribution desktop patch series, importer, promotional images,
  shell canary and old desktop acceptance records.
- Separate data directory; no automatic reuse of predecessor disks or update feed.
- Guest and automatic launcher release pins await verified Try Arch artifacts.
