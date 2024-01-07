import subprocess

completed_process = subprocess.run(
    ["./bscript.sh", "4"], timeout=10, check=True, capture_output=True
)
print("In main python script")
print(completed_process.stdout)
