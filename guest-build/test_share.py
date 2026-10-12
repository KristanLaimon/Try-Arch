"""Run the Linux shared-folder regression check without mounting a disk."""
from pathlib import Path
import shutil
import subprocess
import unittest


class ShareTests(unittest.TestCase):
    def test_mount_link_and_existing_files(self):
        bash = shutil.which('bash')
        if not bash:
            self.skipTest('bash is required')
        result = subprocess.run([bash, str(Path(__file__).with_name('test-share.sh'))],
                                capture_output=True, text=True)
        self.assertEqual(result.returncode, 0, result.stderr)


if __name__ == '__main__':
    unittest.main()
