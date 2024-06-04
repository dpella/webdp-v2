import requests
import pytest
from regression_env import *

def test_create_user(make_test_env, headers_default):
    """Test case for create_user

    Create a new user
    """
    new_user = {
        "handle": "foo",
        "name": "Mr. foo",
        "password": "foo123",
        "roles": ["Curator"],
    }

    response = requests.post(
        URL("users"),
        headers=headers_default,
        json=new_user,
    )
    assert response.status_code == 201

def test_delete_user(make_test_env, headers_default_no_content):
    """Test case for delete_user

    Delete an existing user
    """
    test_user = make_test_env["users"][0]

    response = requests.delete(
        URL("user/{user_handle}".format(user_handle=test_user["handle"])),
        headers=headers_default_no_content,
    )
    assert response.status_code == 204

def test_get_user(make_test_env, headers_default_no_content):
    """Test case for get_user

    Get user information
    """
    test_user = make_test_env["users"][0]

    response = requests.get(
        URL("user/{user_handle}".format(user_handle=test_user["handle"])),
        headers=headers_default_no_content,
    )
    assert response.status_code == 200

def test_get_users(make_test_env, headers_default_no_content):
    """Test case for get_users

    List all users
    """
    _test_env = make_test_env

    response = requests.get(URL("users"), headers=headers_default_no_content)
    assert response.status_code == 200

def test_update_user(make_test_env, headers_default):
    """Test case for update_user

    Update an existing user
    """
    test_user = make_test_env["users"][0]
    update_user_request = test_user
    update_user_request["password"]="new pass"

    response = requests.patch(
        URL("user/{user_handle}".format(user_handle=test_user["handle"])),
        headers=headers_default,
        json=update_user_request,
    )
    assert response.status_code == 204

