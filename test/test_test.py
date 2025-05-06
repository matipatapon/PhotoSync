import requests
import subprocess
import time
def test_user_should_register_and_login() -> None:
    app = subprocess.Popen("go run main.go")
    # while True:
    #     time.sleep(1)
    response = requests.post("http://localhost:8080/register", json={"username":"user", "password": "strongpassword123"}, timeout=30)
    assert response.status_code == 200
    app.kill()