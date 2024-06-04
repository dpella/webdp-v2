import pytest
import requests
from regression_env import *

def test_create_dataset(testdata_fields_schema, headers_default):
    """Test case for create_dataset

    Create a new dataset
    """
    dataset_request = {
        "name": "test_dataset",
        "owner": "root",
        "privacy_notion":"PureDP",
        "total_budget":PureDP(3),
        "schema":testdata_fields_schema,
    }
    response = requests.post(
        URL("datasets"),
        headers=headers_default,
        json=dataset_request,
    )
    assert response.status_code == 201

def test_delete_dataset(make_test_env, headers_default_no_content):
    """Test case for delete_dataset

    Delete an existing dataset
    """
    test_dataset = make_test_env["datasets"][0]

    response = requests.delete(
        URL("dataset/{dataset_id}".format(dataset_id=test_dataset["id"])),
        headers=headers_default_no_content,
    )
    assert response.status_code == 204

def test_get_dataset(make_test_env, headers_default_no_content):
    """Test case for get_dataset

    Get dataset information
    """
    test_dataset = make_test_env["datasets"][0]

    response = requests.get(
        URL("dataset/{dataset_id}".format(dataset_id=test_dataset["id"])),
        headers=headers_default_no_content,
    )
    assert response.status_code == 200

def test_get_datasets(make_test_env, headers_default_no_content):
    """Test case for get_datasets

    List all the datasets visible to the requester
    """
    _test_env = make_test_env

    response = requests.get(
        URL("datasets"), 
        headers=headers_default_no_content
    )
    assert response.status_code == 200

def test_update_dataset(make_test_env, headers_default):
    """Test case for update_dataset

    Update an existing dataset
    """
    test_env = make_test_env
    test_dataset = test_env["datasets"][0]
    test_curator_handle = next(
        filter(lambda user: "Curator" in user["roles"], test_env["users"]),
        None,
    )["handle"]

    dataset_request = test_dataset
    dataset_request["name"] = "updated_dataset"
    dataset_request["owner"] = test_curator_handle
    dataset_request["total_budget"] = PureDP(4)

    response = requests.patch(
        URL("dataset/{dataset_id}".format(dataset_id=test_dataset["id"])),
        headers=headers_default,
        json=dataset_request,
    )
    assert response.status_code == 204

def test_upload_dataset(headers_default, file_fields, testdata_fields_schema):
    """Test case for upload_dataset

    Upload dataset data
    """
    # Create a new dataset
    dataset_request = {
        "name":"test_dataset",
        "owner":"root",
        "privacy_notion":"PureDP",
        "total_budget":PureDP(3),
        "schema":testdata_fields_schema,
    }

    response_create = requests.post(
        URL("datasets"),
        headers=headers_default,
        json=dataset_request,
    )
    assert response_create.status_code == 201
    test_did = str(response_create.json()["id"])

    with open(file_fields) as csv:
        data = csv.read()
        response_upload = requests.post(
            URL("dataset/{dataset_id}/upload".format(dataset_id=test_did)),
            headers=headers_default,
            data=data.encode(),
        )
    assert response_upload.status_code == 204
