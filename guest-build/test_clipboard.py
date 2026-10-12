import base64
import importlib.machinery
import importlib.util
import io
from pathlib import Path
import unittest
from unittest.mock import MagicMock, patch

loader = importlib.machinery.SourceFileLoader(
    'console_clipboard', str(Path(__file__).parent / 'overlay/usr/local/bin/try-arch-clipboard'))
spec = importlib.util.spec_from_loader(loader.name, loader)
clipboard = importlib.util.module_from_spec(spec)
loader.exec_module(clipboard)


class ClipboardTests(unittest.TestCase):
    def test_round_trip_preserves_unicode_and_newlines(self):
        text = 'hello é\n\n'.encode()
        connection = MagicMock()
        connection.__enter__.return_value = connection
        with patch.object(clipboard.socket, 'create_connection', return_value=connection) as connect:
            clipboard.exchange('copy', io.BytesIO(text), io.BytesIO())
            connect.assert_called_once_with(('10.0.2.2', 4448), timeout=5)
            connection.sendall.assert_called_with(base64.b64encode(text) + b'\n')
            connection.makefile.return_value = io.BytesIO(base64.b64encode(text) + b'\r\n')
            output = io.BytesIO()
            clipboard.exchange('paste', io.BytesIO(), output)
            connection.sendall.assert_called_with(b'text-only\n')
            self.assertEqual(output.getvalue(), text)

    def test_copy_rejects_invalid_input_before_connecting(self):
        with patch.object(clipboard, 'LIMIT', 4), patch.object(clipboard.socket, 'create_connection') as connect:
            for text in (b'', b'12345', b'\xff'):
                with self.subTest(text=text), self.assertRaises(ValueError):
                    clipboard.exchange('copy', io.BytesIO(text), io.BytesIO())
            connect.assert_not_called()

    def test_paste_rejects_non_text_and_oversized_frames(self):
        connection = MagicMock()
        connection.__enter__.return_value = connection
        with patch.object(clipboard, 'LIMIT', 4), patch.object(clipboard.socket, 'create_connection', return_value=connection):
            for frame in (b'', b'png:YWJj\n', b'!!!\n', b'/w==\n', b'MTIzNDU=\n', b'A' * 20):
                with self.subTest(frame=frame), self.assertRaises(ValueError):
                    connection.makefile.return_value = io.BytesIO(frame)
                    output = io.BytesIO()
                    clipboard.exchange('paste', io.BytesIO(), output)
                self.assertEqual(output.getvalue(), b'')


if __name__ == '__main__':
    unittest.main()
