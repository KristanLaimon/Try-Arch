import argparse
from pathlib import Path
import runpy
import subprocess
import tempfile
import unittest

recording_command = runpy.run_path(str(Path(__file__).parent / 'overlay/usr/local/bin/try-arch-record'))['recording_command']


class RecordingTests(unittest.TestCase):
    def test_virtual_gpu_uses_cpu_and_preserves_existing_recordings(self):
        with tempfile.TemporaryDirectory() as directory:
            output=Path(directory)/'video.mp4'
            options=argparse.Namespace(output=str(output),region=False,audio=None)
            command=recording_command(options,which=lambda name:'/usr/bin/'+name)
            self.assertIn('--no-dmabuf',command)
            self.assertEqual(command[command.index('-c')+1],'libx264')
            self.assertNotIn('-y',command)
            output.write_bytes(b'previous recording')
            with self.assertRaisesRegex(ValueError,'overwrite'):
                recording_command(options,which=lambda name:'/usr/bin/'+name)
            self.assertEqual(output.read_bytes(),b'previous recording')

    def test_region_cancel_does_not_create_output(self):
        with tempfile.TemporaryDirectory() as directory:
            output=Path(directory)/'missing/video.mp4'
            options=argparse.Namespace(output=str(output),region=True,audio='microphone')
            result=recording_command(options,which=lambda name:'/usr/bin/'+name,
                run=lambda *a,**k:subprocess.CompletedProcess(a,1,stdout=''))
            self.assertIsNone(result)
            self.assertFalse(output.parent.exists())


if __name__ == '__main__':
    unittest.main()
