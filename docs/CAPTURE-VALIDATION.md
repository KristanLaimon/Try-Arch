# Focused Windows capture check — 2026-11-02

Tested on Windows 11 build 26200.9711 with an NVIDIA GeForce RTX 3050 installed,
driver 32.0.16.1088, the installed WINQ-EMU runtime and OBS 32.2.2.

The affected host recorded windowed mode but produced black captures in regular
fullscreen. A screenshot after switching focus away from the VM does not
validate this behavior. Acceptance requires inspecting a video recorded with
the VM foreground, including the signed-in desktop.

The installed launcher was rebuilt with `-capture-friendly`. The existing
installation opted in through `capture-friendly-immersive`, retaining its
Immersive setting and virgl rendering path. The actual native window was
1920×1079, borderless, on a 1920×1080 monitor. The guest disk was retained.

Only the existing Try Arch Window Capture source was changed in the user's OBS
collection: Windows Graphics Capture (`method=2`), cursor capture off. Its
previous scene file and launcher are backed up under
`local-backups/capture-mode-20261102`. An independent, verified VM disk backup
from before the upgrade remains under `local-backups/before-upgrade-20261102`.

An isolated portable copy of the installed OBS recorded the signed-in Ryoku
desktop without changing the user's audio or other scene settings. The VM was
foreground for eight of ten one-second focus observations, including the first
seven consecutive observations. The 11.4-second H.264 recording was decoded
with FFmpeg. Ten one-second samples after initialization were nonblack, with
at least 64.8% of pixels above luma 20. All eleven sampled frames differed.
Decoded frames at 3, 6 and 9 seconds were inspected visually and showed the
Ryoku desktop.

Local artifacts, excluded from source control:

- `.tmp/2026-11-02 07-29-40.mkv`
- `.tmp/installed-capture-validation.json`
- `.tmp/recorded-desktop-active-3.png`, `-6.png`, `-9.png`

The full Go test suite passed before building the installed launcher, including
the capture flag, bounds and opt-in backup checks. The VM remained running after
the final recording. NVIDIA driver preferences were restored after an
unsuccessful presentation-override experiment.

This validates the current host and short recording. It does not establish
Discord screen sharing, extended recordings, other GPUs or mixed-DPI displays.
This recording used the earlier opt-in build. The release now enables this
mode by default and routes Ctrl+Alt+F through the same borderless toggle;
shortcut routing is covered by regression tests.
