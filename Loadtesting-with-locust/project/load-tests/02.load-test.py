from locust import HttpUser, task, between, events
import time
from json import JSONDecodeError

@events.test_start.add_listener
def on_test_start(environment, **kwargs):
    print("Logging test start. Sending notification")

@events.test_stop.add_listener
def on_test_stop(environment, **kwargs):
    print("Logging test stop. Sending test report")

@task
def home(self):
    with self.client.get("/", catch_response=True, name="home-page") as response:
        if response.elapsed.total_seconds() > 1:
            response.failure("Request took too long")

@task
def view_products(self):
    for item_id in range(10):
        with self.client.get(f"/products?id={item_id}", catch_response=True, name="/products") as response:
            try:
                if response.json()["productId"] != str(item_id):
                    response.failure("Did not get expected value")
            except JSONDecodeError:
                response.failure("Response could not be decoded as JSON")
            except KeyError:
                response.failure("Response did not contain expected key 'greeting'")
            time.sleep(1)

class NormalUser(HttpUser):
    wait_time = between(1, 5)

    tasks = {home: 2, view_products: 1}
    
    def on_start(self):
        self.client.post("/signin", json={"username":"foo", "password":"bar"})

class WindowShopperUser(HttpUser):
    wait_time = between(1, 2)

    tasks = {home: 1, view_products: 3}
    
    def on_start(self):
        self.client.post("/signin", json={"username":"foo", "password":"bar"})