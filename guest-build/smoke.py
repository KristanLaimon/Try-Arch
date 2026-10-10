"""Linux/QEMU smoke check. Boots a temporary overlay; never changes the factory image."""
import os
from pathlib import Path
import selectors
import socket
import subprocess
import sys
import tempfile
import time


def main(directory):
    directory = Path(directory).resolve()
    with tempfile.TemporaryDirectory(prefix='try-arch-smoke-') as temporary, socket.create_server(('127.0.0.1', 4450)) as lifecycle:
        disk = Path(temporary) / 'disk.qcow2'
        subprocess.run(['qemu-img', 'create', '-f', 'qcow2', '-F', 'raw', '-b',
                        str(directory / 'rootfs.ext4'), str(disk)], check=True)
        subprocess.run(['qemu-img', 'resize', str(disk), '16G'], check=True)
        process = subprocess.Popen([
            'qemu-system-x86_64', '-machine', 'q35,accel=tcg',
            '-cpu', 'qemu64,+ssse3,+sse4.1,+sse4.2,+popcnt,+aes',
            '-m', '2048', '-smp', '2', '-display', 'none', '-serial', 'stdio',
            '-monitor', 'none', '-no-reboot', '-drive', f'file={disk},format=qcow2,if=virtio',
            '-kernel', str(directory / 'vmlinuz-linux'),
            '-initrd', str(directory / 'initramfs-linux.img'),
            '-append', 'root=/dev/vda rw console=ttyS0 tryarch.instant=1',
            '-device', 'virtio-net-pci,netdev=n0', '-netdev', 'user,id=n0',
            '-device', 'virtio-rng-pci'], stdin=subprocess.PIPE,
            stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
        selector = selectors.DefaultSelector()
        selector.register(process.stdout, selectors.EVENT_READ)
        transcript = bytearray()
        stage = 0
        deadline = time.monotonic() + 300
        checks = (
            'test "$SUDO_USER" = arch; '
            '. /etc/os-release; test "$ID" = arch; '
            'test -z "$(pacman -Qm)"; '
            'test "$(stat -c %U:%a /etc)" = root:755; '
            'test "$(stat -c %U:%a /usr)" = root:755; '
            'test "$(df --output=size / | tail -1)" -gt 8000000; '
            'passwd -S root | grep -q " L "; '
            'ssh-keygen -A; '
            'sshd -T | grep -qix "passwordauthentication no"; '
            'sshd -T | grep -qix "permitrootlogin no"; '
            'until systemctl is-active --quiet try-arch-setup; do sleep 1; done; '
            'systemctl is-active --quiet systemd-networkd systemd-resolved; '
            'test -f /usr/lib/modules/$(uname -r)/modules.dep.bin; '
            'pacman -Syu --noconfirm; '
            'printf "SMOKE:%s:%s\\n" "$ID" "verified"; systemctl poweroff'
        )
        try:
            while time.monotonic() < deadline:
                for key, _ in selector.select(1):
                    block = os.read(key.fileobj.fileno(), 65536)
                    if not block:
                        break
                    transcript.extend(block)
                if process.poll() is not None:
                    break
                tail = bytes(transcript[-4096:])
                message = None
                if stage == 0 and b'login: ' in tail:
                    message, stage = 'arch\n', 1
                elif stage == 1 and b'Password:' in tail:
                    message, stage = 'arch\n', 2
                elif stage == 2 and b'$ ' in tail:
                    message, stage = "printf 'arch\\n' | sudo -S bash -exc '" + checks + "'\n", 3
                elif stage == 3 and b'exit=failure' in tail:
                    break
                if message:
                    process.stdin.write(message.encode())
                    process.stdin.flush()
                    transcript.extend(b'\n')
            log = bytes(transcript).replace(b'\r', b'')
            (directory / 'smoke.log').write_bytes(log)
            assert b'\nSMOKE:arch:verified\n' in log, log[-10000:].decode(errors='replace')
            lifecycle.settimeout(3)
            connection, _ = lifecycle.accept()
            with connection:
                assert connection.recv(64) == b'ready\n', 'Guest did not report readiness'
            assert process.wait(timeout=20) == 0, 'QEMU did not shut down cleanly'
            print('PASS: vanilla login, sudo, networking, package update, root/SSH policy, kernel modules, disk growth and shutdown')
        finally:
            selector.close()
            if process.poll() is None:
                process.kill()
            process.wait()


if __name__ == '__main__':
    main(sys.argv[1])
