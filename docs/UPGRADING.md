# Upgrade an existing Try Arch VM

Shut down Linux normally and close Try Arch. Keep the existing data directory
(normally `%LOCALAPPDATA%\TryArch`). Replace its `TryArch.exe` with the new
launcher, saving the previous executable for rollback. Start it with the same
data directory. The launcher reuses the installed guest boot files and runtime;
the writable `vm/disk.raw` contains your Linux packages, desktop and files.
Updating the launcher does not reinstall Linux or update those packages.

Use Settings → Backup for an independent VM backup before changing Linux.
Do not use Reset, Uninstall or `-fresh` to upgrade. Keep the installed external
kernel and initramfs matched to the guest; see [the user guide](USER-GUIDE.md).

# Existing apps' Discord integration

The optional bridge relays Rich Presence from apps that already implement it
to the Windows Discord desktop client. It does not detect foreground apps and
needs no new application ID, bot or Discord token. The original app supplies
its application ID and activity. Account, authorization, subscription and voice
control RPC commands are blocked. Join/invite callbacks are not supported.

On an existing guest, run this from a copy of this repository as your Linux
desktop user:

```sh
sh scripts/install-discord-bridge.sh
```

Start the Windows launcher with `-discord-rpc`. To enable it for ordinary
shortcut launches, create an empty `discord-rpc-enabled` file inside the data
directory. Remove that marker to disable it, or use `-discord-rpc=false` for
one launch. Older launchers ignore the marker; their settings stay compatible.

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

Future guest builds contain the helper and units but leave this feature off.

# Recording a virtual GPU

`gsr error: unknown gpu vendor: Mesa` followed by `failed to load opengl` means
GPU Screen Recorder rejected the guest GPU vendor. Windows NVIDIA drivers are
not Linux NVIDIA drivers inside this VM: the guest sees a Virtio/Mesa GPU.
Changing the launcher or reinstalling Mesa does not give it NVIDIA's encoder.

Use CPU encoding with the installed Wayland recorder instead:

```sh
wf-recorder --no-dmabuf -c libx264 -p preset=ultrafast -x yuv420p -r 30 -f recording.mp4
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

If capture goes black only while the fullscreen VM is focused, start with
`-fullscreen -capture-friendly` instead. This uses a borderless window that
leaves one physical pixel of the Windows desktop visible at the bottom of the
selected monitor, preventing exact fullscreen presentation. GPU acceleration
remains enabled. In OBS Window Capture, select **Windows 10 (1903 and up)** and
disable **Capture Cursor**; cursor capture also produced black frames on the
affected host. The pointer remains usable in the VM but may be absent from the
recording.

To enable this for normal launches, create an empty
`capture-friendly-immersive` file in the installation data directory. The
Immersive checkbox still controls whether it opens across the monitor. Use
Settings to switch between windowed and Immersive in this mode; Ctrl+Alt+F
enters SDL's regular fullscreen path and can bring the capture problem back.
Remove the marker or pass `-capture-friendly=false` to return to regular
fullscreen. The flag does not reinstall Linux or alter its files.

References: [Discord IPC and Rich Presence](https://github.com/discord/discord-api-docs/blob/main/developers/topics/rpc.mdx),
[OBS Window Capture](https://obsproject.com/kb/window-capture-sources).
