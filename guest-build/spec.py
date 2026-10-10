"""Write the vanilla Arch boot contract consumed by the Windows launcher."""
import json
from pathlib import Path
import sys

SPEC = {"distribution": "archlinux", "runtime": {
    "console": True,
    "kernelCommandLine": "root=/dev/vda rw console=tty0 console=ttyS0",
    "optionalDevices": [], "storage": {"expandedSizeMiB": 16384}}}

def main(output):
    output = Path(output)
    (output / "build-spec.json").write_text(json.dumps(SPEC, indent=2) + "\n")
    (output / "guest-manifest.json").write_text(json.dumps({
        "distribution": "archlinux", "desktop": None,
        "packages": (output / "packages.txt").read_text().splitlines()}, indent=2) + "\n")

if __name__ == "__main__":
    main(sys.argv[1])
