import subprocess
from time import sleep

with subprocess.Popen(
    ["./bscript.sh", "3"], stdout=subprocess.PIPE
) as process:

    def poll_and_read():
        print(f"polling: {process.poll()}")
        print(f"stdout: {process.stdout.read1().decode('utf-8')}")

    poll_and_read()
    sleep(2)
    poll_and_read()
    sleep(2)
    poll_and_read()