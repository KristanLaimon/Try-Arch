#!/usr/bin/env bash
# Exercise the shipped CPU command against a real headless Wayland compositor.
set -euo pipefail
repo=$(cd "$(dirname "$0")/../.." && pwd)
output=${1:?Usage: test-cpu-recording.sh OUTPUT_DIR}
mkdir -p "$output"
output=$(cd "$output" && pwd)
work=$(mktemp -d /tmp/tryarch-cpu.XXXXXX)
compositor=
cleanup() {
  if [[ -n $compositor ]]; then kill "$compositor" 2>/dev/null || true; wait "$compositor" 2>/dev/null || true; fi
  [[ $work == /tmp/tryarch-cpu.* ]] && rm -rf -- "$work"
}
trap cleanup EXIT
mkdir -m700 "$work/runtime"
export XDG_RUNTIME_DIR="$work/runtime"
export WLR_BACKENDS=headless WLR_RENDERER=pixman WLR_LIBINPUT_NO_DEVICES=1
printf 'output HEADLESS-1 mode 640x480\noutput * bg #338866 solid_color\n' > "$work/sway.conf"
dbus-run-session -- sway --unsupported-gpu -c "$work/sway.conf" > "$output/cpu-compositor.log" 2>&1 &
compositor=$!
for ((attempt=0;attempt<100;attempt++)); do
  for candidate in "$XDG_RUNTIME_DIR"/wayland-*; do
    if [[ -S $candidate ]]; then export WAYLAND_DISPLAY=${candidate##*/}; break 2; fi
  done
  kill -0 "$compositor" || { cat "$output/cpu-compositor.log"; exit 1; }
  sleep .1
done
[[ -n ${WAYLAND_DISPLAY:-} ]]
set +e
timeout --signal=INT 5s python3 "$repo/guest-build/overlay/usr/local/bin/try-arch-record" --output "$output/cpu-recording.mkv" > "$output/cpu-recorder.log" 2>&1
status=$?
set -e
[[ $status == 0 || $status == 124 ]] || { cat "$output/cpu-recorder.log"; exit 1; }
python3 - "$output/cpu-recording.mkv" <<'PY'
import json,subprocess,sys
name=sys.argv[1]
info=json.loads(subprocess.check_output(['ffprobe','-v','error','-show_streams','-of','json',name]))
video=next(s for s in info['streams'] if s['codec_type']=='video')
assert video['codec_name']=='h264' and video['width']==640 and video['height']==480,video
pixels=subprocess.check_output(['ffmpeg','-v','error','-i',name,'-frames:v','1','-vf','scale=32:24,format=rgb24','-f','rawvideo','-'])
assert len(pixels)==32*24*3 and sum(pixels)/len(pixels)>30,'recording is black'
print('PASS: shipped CPU recorder produced a nonblack H.264 Wayland video')
PY
