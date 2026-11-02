# Upgrade an existing Try Arch VM

Shut down Linux normally and close Try Arch. Keep the existing data directory
(normally `%LOCALAPPDATA%\TryArch`). Replace its `TryArch.exe` with the new
launcher, saving the previous executable for rollback. Start it with the same
data directory. The launcher reuses the installed guest boot files and runtime;
the writable `vm/disk.raw` contains your Linux packages, desktop and files.
Updating the launcher does not reinstall Linux or update those packages.

Use Settings → Backup for an independent VM backup before changing Linux.
Do not use Reset, Uninstall or `-fresh` to upgrade. Keep the installed external
kernel and initramfs matched to the guest. Automatic factory boot-payload
replacement is disabled for existing guests; `-update-guest-payload` is an
explicit maintenance option, requiring a matching guest kernel and backup; see [the user guide](USER-GUIDE.md).

# Existing apps' Discord integration

The optional bridge relays Rich Presence from apps that already implement it
to the Windows Discord desktop client. It does not detect foreground apps and
needs no new application ID, bot or Discord token. The original app supplies
its application ID and activity. Account, authorization, subscription and voice
control RPC commands are blocked. Join/invite callbacks are not supported.

Fresh installations include the guest bridge. On an older guest, download
`tryarch-guest-integrations.tar.gz` from the release, extract it, and run as
your Linux desktop user from the extracted directory:

```sh
sh scripts/install-guest-integrations.sh
```

This installs per-user helpers and backs up any helper files it replaces.
It does not replace your desktop, packages or VM disk. Install any missing
recording dependencies with `sudo pacman -S wf-recorder slurp`.

Enable **Settings ? Apps ? Discord Rich Presence** on Windows and restart
Try Arch. It is off by default. `-discord-rpc` enables it for one launch;
`-discord-rpc=false` disables it. Existing opt-in markers are migrated when
no integration preferences have been saved. Preferences are stored separately
from older launchers' settings for rollback compatibility.

The bridge replaces every activity's status/location with `In arch-linux`,
so this field cannot expose the Linux account name. The app still controls
its other details, artwork and timestamps; review each app's own privacy settings.

Open Windows Discord, enable its activity-sharing setting, then restart the
Linux app. The guest exposes a private per-user `discord-ipc-0` Unix socket;
the launcher connects to a Windows `discord-ipc-*` named pipe. The host listener
binds only to `127.0.0.1:4457`, reachable through QEMU's host gateway. Other local
host processes and guest users can reach that listener, which is why it accepts
only bounded Rich Presence messages. Do not port-forward it to a LAN.

Quit Linux Discord before enabling the socket; the two cannot own the same
socket path. Sandboxed apps may need permission to access the guest IPC socket.
To disable the guest side:

```sh
systemctl --user disable --now try-arch-discord.socket
```

The guest socket is gated by the Windows opt-in kernel flag `tryarch.discord=1`.

# Recording a virtual GPU

`gsr error: unknown gpu vendor: Mesa` followed by `failed to load opengl` means
GPU Screen Recorder rejected the guest GPU vendor. Windows NVIDIA drivers are
not Linux NVIDIA drivers inside this VM: the guest sees a Virtio/Mesa GPU.
Changing the launcher or reinstalling Mesa does not give it NVIDIA's encoder.

Use CPU encoding with the installed Wayland recorder instead:

```sh
try-arch-record
# Optional: try-arch-record --region --audio
# Optional: try-arch-record --output ~/Videos/new-recording.mp4
```

Stop with Ctrl+C. This requires a running compatible Wayland compositor and
records video without audio by default. It uses more CPU than hardware encoding.
It is a separate recording path; Ryoku's GPU Screen Recorder buttons continue
to use that application's backend until Ryoku provides a compatible fallback.

On Windows OBS, choose **Window Capture**, select the QEMU/Try Arch display
window, and try **Windows 10 (1903 and up)** as the capture method. Display
Capture is another option. In Discord screen sharing, select the Try Arch
window or its display. Keep the VM visible while recording; minimizing a window
can stop updates. These capture programs run on Windows and capture the VM
display, rather than an individual Linux application's native window.

Capture-friendly Immersive is enabled by default for new and upgraded launchers.
It uses a borderless window that leaves one physical pixel visible at the
bottom of the selected monitor, avoiding exact fullscreen presentation.
GPU acceleration remains enabled. The Immersive checkbox and Ctrl+Alt+F both
switch between this mode and a regular window.

In OBS Window Capture, select **Windows 10 (1903 and up)** and disable
**Capture Cursor** if focused fullscreen capture goes black. The pointer stays
usable in the VM but may be absent from the recording. Other users must select
these settings in their own OBS; the launcher does not rewrite OBS profiles.

Disable **Settings ? Display ? Capture-friendly Immersive** or pass
`-capture-friendly=false` to use SDL's regular fullscreen. This may restore the
black capture issue on affected drivers. These options do not alter Linux files.

References: [Discord IPC and Rich Presence](https://github.com/discord/discord-api-docs/blob/main/developers/topics/rpc.mdx),
[OBS Window Capture](https://obsproject.com/kb/window-capture-sources).
