import requests
import subprocess
import os
import time

def log_subprocess_output(pipe):
    for line in iter(pipe.readline, b''): # b'\n'-separated lines
        print('got line from subprocess: %r', line)

def test_user_should_register_and_login() -> None:
    app = subprocess.Popen("go run main.go", stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
    time.sleep(10)
    response = requests.post("http://localhost:8080/register", json={"username":"user", "password": "strongpassword123"}, timeout=30)
    assert response.status_code == 200
    response = requests.post("http://localhost:8080/exit", json={"username":"user", "password": "strongpassword123"}, timeout=30)

    with app.stdout:
        log_subprocess_output(app.stdout)
