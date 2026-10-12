import json
import hashlib
from pathlib import Path
import tempfile
import unittest

from spec import ARTIFACT_NAMES, SPEC, main

class VanillaGuestTests(unittest.TestCase):
    def test_boot_contract(self):
        with tempfile.TemporaryDirectory() as directory:
            output = Path(directory)
            (output / 'packages.txt').write_text('base 3-2\nlinux 6.18-1\n')
            for name in ARTIFACT_NAMES:
                (output / name).write_bytes(('fixture ' + name).encode())
            main(output)
            self.assertEqual(json.loads((output / 'build-spec.json').read_text()), SPEC)
            manifest = json.loads((output / 'guest-manifest.json').read_text())
            self.assertEqual(manifest['schemaVersion'], 1)
            for artifact in manifest['artifacts']:
                data = (output / artifact['path']).read_bytes()
                self.assertEqual(artifact['bytes'], len(data))
                self.assertEqual(artifact['sha256'], hashlib.sha256(data).hexdigest())
            self.assertEqual({entry['path'] for entry in manifest['artifacts']}, set(ARTIFACT_NAMES))
            self.assertIsNone(manifest['desktop'])
            self.assertEqual(manifest['packages'], ['base 3-2', 'linux 6.18-1'])
            self.assertIn('root=/dev/vda', SPEC['runtime']['kernelCommandLine'])
            self.assertIn('console=tty0', SPEC['runtime']['kernelCommandLine'])
            self.assertEqual(SPEC['runtime']['optionalDevices'], [])

    def test_builder_uses_only_official_packages(self):
        source = Path(__file__).parent
        builder = (source / 'build.sh').read_text()
        packages = next(line for line in builder.splitlines() if line.startswith('pacstrap '))
        self.assertEqual(packages.split('"$work/root" ')[1].split(), ['base', 'linux', 'linux-firmware', 'btrfs-progs', 'openssh', 'sudo', 'nano', 'python', 'socat', 'wl-clipboard'])
        self.assertNotIn('git clone', builder)
        self.assertNotIn('yay', builder)
        self.assertIn('-K "$work/root"', builder)
        self.assertIn("sed '/^NoExtract/d'", builder)
        self.assertIn('cp -r --no-preserve=mode,ownership', builder)
        self.assertIn('mkfs.btrfs -f -L tryarch --rootdir "$work/root" "$image"', builder)
        self.assertIn('btrfs check --readonly "$image"', builder)
        self.assertIn('/dev/vda / btrfs defaults 0 0', builder)
        self.assertIn('virtio_gpu btrfs)', builder)
        self.assertIn('systemctl --global enable try-arch-clipboard', builder)
        self.assertIn('/src/scripts/guest/clipboard-bridge.sh', builder)
        self.assertIn('/src/scripts/guest/clipboard-files /src/scripts/guest/file-transfer', builder)
        self.assertIn('PasswordAuthentication no', builder)
        setup = (source / 'overlay/usr/local/bin/try-arch-setup').read_text()
        self.assertIn('btrfs filesystem resize max /', setup)
        self.assertIn('tryarch.instant=1', setup)
        self.assertIn('ssh-keygen -l', setup)
        self.assertIn('/etc/ssh/authorized_keys/$username', setup)
        self.assertNotIn('$home/.ssh', setup)
        self.assertIn("printf 'ready\\n'", setup)

if __name__ == '__main__':
    unittest.main()
