# from locust import HttpUser, task, between
#
# class KeyValueCacheUser(HttpUser):
#     wait_time = between(0.005, 0.01)  # Reduce wait time to send faster requests
#
#     @task(3)  # 3x more GET requests than PUT requests
#     def get_request(self):
#         self.client.get("/get?key=testKey")
#
#     @task(1)
#     def put_request(self):
#         self.client.post("/put", json={"key": "testKey", "value": "testValue"})
#
# # Run Locust with:
# # locust -f locustfile.py --host=http://localhost:7171
#



import time
import socket
from locust import User, task, between

class CustomRedisUser(User):
    wait_time = between(0.05, 0.2)  # Avoid CPU saturation

    def on_start(self):
        """Initialize a persistent TCP connection to the custom Redis server"""
        self.client = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self.client.connect(("localhost", 7171))  # Connect to your custom Redis server

    def send_command(self, command):
        """Sends a raw RESP command and receives a response"""
        try:
            self.client.sendall(command.encode())
            return self.client.recv(1024).decode()
        except Exception as e:
            return f"ERROR: {e}"

    @task(3)
    def get_request(self):
        """Sends a GET command"""
        start_time = time.time()
        key = "testKey"

        # Use RESP protocol format
        command = f"*2\r\n$3\r\nGET\r\n${len(key)}\r\n{key}\r\n"
        response = self.send_command(command)

        self.environment.events.request.fire(
            request_type="GET",
            name="redis_get",
            response_time=(time.time() - start_time) * 1000,
            response_length=len(response),
            exception=None if response else Exception("Key not found"),
        )

    @task(1)
    def put_request(self):
        """Sends a PUT (SET) command"""
        start_time = time.time()
        key, value = "testKey", "testValue"

        # Use RESP protocol format
        command = f"*3\r\n$3\r\nSET\r\n${len(key)}\r\n{key}\r\n${len(value)}\r\n{value}\r\n"
        response = self.send_command(command)

        self.environment.events.request.fire(
            request_type="SET",
            name="redis_set",
            response_time=(time.time() - start_time) * 1000,
            response_length=len(response),
            exception=None if "OK" in response else Exception("SET command failed"),
        )

    def on_stop(self):
        """Close the TCP connection when the user stops"""
        self.client.close()

