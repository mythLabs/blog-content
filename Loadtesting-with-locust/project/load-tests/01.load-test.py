from locust import FastHttpUser, task

class RootTest(FastHttpUser):
    @task
    def home(self):
        self.client.get("/")