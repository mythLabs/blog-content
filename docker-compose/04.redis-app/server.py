import time
import os
import redis
from pathlib import Path
from flask import Flask

app = Flask(__name__)
cache = redis.Redis(host='redis', port=6379)

def get_hit_count():
    retries = 5
    while True:
        try:
            return cache.incr('hits')
        except redis.exceptions.ConnectionError as exc:
            if retries == 0:
                raise exc
            retries -= 1
            time.sleep(0.5)

@app.route('/')
def hello():
    count = get_hit_count()
    contents = Path("/run/secrets/app_secret").read_text()
    return f'Hello World! I have been seen {count} times.\n Environment key is {os.environ.get("ENVKEY")}\n Secret is {contents}'