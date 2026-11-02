"""Exercise the guest bridge with real sockets, without contacting Discord."""
from pathlib import Path
import runpy
import socket
import threading
import unittest

relay = runpy.run_path(str(Path(__file__).parent / 'overlay/usr/local/bin/try-arch-discord'))['relay']


class DiscordBridgeTests(unittest.TestCase):
    def test_bidirectional_rpc_and_app_exit(self):
        app, local = socket.socketpair()
        host, discord = socket.socketpair()
        with app, local, host, discord:
            app.settimeout(3)
            discord.settimeout(3)
            worker = threading.Thread(target=relay, args=(local, host), daemon=True)
            worker.start()
            app.sendall(b'handshake')
            self.assertEqual(discord.recv(9), b'handshake')
            discord.sendall(b'ready')
            self.assertEqual(app.recv(5), b'ready')
            app.shutdown(socket.SHUT_WR)
            worker.join(timeout=3)
            self.assertFalse(worker.is_alive())


if __name__ == '__main__':
    unittest.main()
