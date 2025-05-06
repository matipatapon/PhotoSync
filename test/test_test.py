import requests
def test_user_should_register_and_login() -> None:
    response = requests.put("http://localhost/register", json={"username":"user", "password": "strongpassword123"})
    print(response.json())
