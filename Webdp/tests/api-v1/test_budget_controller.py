import pytest
import requests
from regression_env import *

def test_allocate_user_dataset_budget(make_test_env, headers_default):
    """Test case for allocate_user_dataset_budget

    Allocate some budget to a user for a given dataset
    """

    test_dataset = make_test_env["datasets"][0]

    budget = PureDP(1)
    response = requests.post(
        URL("budget/allocation/{user_handle}/{dataset_id}".format(
            user_handle="analystUserNoBudget", 
            dataset_id=test_dataset["id"]
        )),
        headers=headers_default,
        json=budget,
    )
    assert response.status_code == 201 # was 204

def test_delete_user_dataset_budget(make_test_env, headers_default_no_content):
    """Test case for delete_user_dataset_budget

    Dealocate the budget assigned to a user for a given dataset
    """
    test_env = make_test_env
    test_alloc = test_env["budgets"][0]
    test_user_handle = test_alloc["user"]
    test_did = next(
        filter(lambda dataset: dataset["name"] == test_alloc["dataset"], test_env["datasets"]),
        None,
    )["id"]

    response = requests.delete(
        URL("budget/allocation/{user_handle}/{dataset_id}".format(
            user_handle=test_user_handle, 
            dataset_id=test_did
        )),
        headers=headers_default_no_content,
    )
    assert response.status_code == 204

def test_get_dataset_budget_allocation(make_test_env, headers_default_no_content):
    """Test case for get_dataset_budget_allocation

    Get the dataset budget allocation accross users
    """
    test_dataset = make_test_env["datasets"][0]

    response = requests.get(
        URL("budget/dataset/{dataset_id}".format(dataset_id=test_dataset["id"])),
        headers=headers_default_no_content,
    )
    assert response.status_code == 200

def test_get_user_budget_allocation(make_test_env, headers_default_no_content):
    """Test case for get_user_budget_allocation

    Get the user budget allocation across datasets
    """
    _test_env = make_test_env

    response = requests.get(
        URL("budget/user/{user_handle}".format(user_handle="analystUser")),
        headers=headers_default_no_content,
    )
    assert response.status_code == 200

def test_get_user_dataset_budget(make_test_env, headers_default_no_content):
    """Test case for get_user_dataset_budget

    Get the budget allocated to a user on a given dataset
    """
    test_env = make_test_env
    test_user = test_env["users"][0]
    test_dataset = test_env["datasets"][0]

    response = requests.get(
        URL("budget/allocation/{user_handle}/{dataset_id}".format(
            user_handle=test_user["handle"], dataset_id=test_dataset["id"]
        )),
        headers=headers_default_no_content,
    )
    assert response.status_code == 200

def test_update_user_dataset_budget(make_test_env, headers_default):
    """Test case for update_user_dataset_budget

    Update the budget alocated to a user for a given dataset
    """
    test_env = make_test_env
    test_alloc = test_env["budgets"][0]
    test_user_handle = test_alloc["user"]
    test_did = next(
        filter(lambda dataset: dataset["name"] == test_alloc["dataset"], test_env["datasets"]),
        None,
    )["id"]

    budget = PureDP(0.5)
    response = requests.patch(
        URL("budget/allocation/{user_handle}/{dataset_id}".format(
            user_handle=test_user_handle, dataset_id=test_did
        )),
        headers=headers_default,
        json=budget,
    )
    assert response.status_code == 204
