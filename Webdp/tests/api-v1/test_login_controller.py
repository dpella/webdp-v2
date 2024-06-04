import pytest
import requests
from regression_env import *

def test_login(make_test_env, headers_no_auth):
    """Test case for login

    Login using user/password credentials
    """
    test_user = make_test_env["users"][0]

    login_request = loginjson(test_user["handle"], test_user["password"])
    response = requests.post(
        URL("login"),
        headers=headers_no_auth,
        json=login_request,
    )
    assert response.status_code == 200

def test_logout(make_test_env, headers_no_auth):
    """Test case for logout

    Logout from the current session
    """
    # Loggin in
    test_user = make_test_env["users"][0]

    login_request = loginjson(test_user["handle"], test_user["password"])
    response_login = requests.post(
        URL("login"),
        headers=headers_no_auth,
        json=login_request,
    )
    assert response_login.status_code == 200
    test_user_JWT = response_login.json()["jwt"]

    # Loggin out
    headers_user_auth = {
        "Accept": "application/json",
        "Authorization": "Bearer {key}".format(key=test_user_JWT),
    }
    response_logout = requests.post(
        URL("logout"),
        headers=headers_user_auth
    )
    assert response_logout.status_code == 204

