"""Write the vanilla Arch boot contract consumed by the Windows launcher."""
import json
import hashlib
from pathlib import Path
import sys

SPEC = {"distribution": "archlinux", "runtime": {
    "console": True,
    "kernelCommandLine": "root=/dev/vda rw console=tty0 console=ttyS0",
    "optionalDevices": [], "storage": {"expandedSizeMiB": 16384}}}

ARTIFACT_NAMES = ('rootfs.ext4', 'rootfs.ext4.zst', 'vmlinuz-linux', 'initramfs-linux.img')

def artifact_metadata(path):
    with path.open('rb') as source:
        digest = hashlib.file_digest(source, 'sha256').hexdigest()
    return {'path': path.name, 'bytes': path.stat().st_size, 'sha256': digest}

def main(output):
    output = Path(output)
    (output / "build-spec.json").write_text(json.dumps(SPEC, indent=2) + "\n")
    (output / "guest-manifest.json").write_text(json.dumps({
        "schemaVersion": 1, "distribution": "archlinux", "desktop": None,
        "artifacts": [artifact_metadata(output / name) for name in ARTIFACT_NAMES],
        "packages": (output / "packages.txt").read_text().splitlines()}, indent=2) + "\n")

if __name__ == "__main__":
    main(sys.argv[1])
