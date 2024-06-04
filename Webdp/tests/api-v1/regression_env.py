"""

This test suite contains the tests present in the WebDP tumult server.

NOTE: The query tests assume that Tumult is set as the default engine.

Adaptations from the original tests:

- No sessions: Since the new architecture is built with persistance, 
  the environment is set up in a database.
- Most data structures are written in JSON -- data models are not imported.
- In test_update_dataset, the total_budget was increased from 2 to 4, because
  you cannot set a budget that is lower than what has been allocated already.

TODO: Teardowns/clearing of test data

"""


import pytest
import requests
import platform

URL = lambda endpoint: f"http://localhost:8000/v1/{endpoint}"

loginjson = lambda user, pwd: {"username": f"{user}", "password": f"{pwd}"}
PureDP    = lambda epsilon : {"epsilon": epsilon}
ApproxDP  = lambda epsilon, delta : {"epsilon": epsilon, "delta": delta}

@pytest.fixture
def roottoken():
    root = loginjson("root", "123")
    response = requests.post(URL("login"), json=root)
    assert response.status_code == 200
    return response.json()["jwt"]

@pytest.fixture
def headers_default(roottoken):
    return {
        "Accept": "application/json",
        "Content-Type": "application/json",
        "Authorization": "Bearer {key}".format(key=roottoken),
    }

@pytest.fixture
def headers_default_no_content(roottoken):
    return {
        "Accept": "application/json",
        "Authorization": "Bearer {key}".format(key=roottoken),
    }

@pytest.fixture
def headers_no_auth():
    return {
        "Accept": "application/json",
        "Content-Type": "application/json",
    }

@pytest.fixture
def file_salaries():
    return "./tests/testdata_salaries.csv"

@pytest.fixture
def file_fields():
    return "./tests/testdata_fields.csv"

@pytest.fixture
def testdata_salaries_schema():
    return [
        {"name":"name",   "type":{"name": "Text"}},
        {"name":"age",    "type":{"name": "Int", "low": 18, "high": 100 }},
        {"name":"salary", "type":{"name": "Int", "low": 0, "high": 100000}},
    ]

@pytest.fixture
def testdata_fields_schema():
    return [
        {"name":"fieldInt",    "type":{"name":"Int", "low":0, "high":10000}},
        {"name":"fieldDouble", "type":{"name":"Double", "low":-5, "high":5}},
        {"name":"fieldText",   "type":{"name":"Text"}},
        {"name":"fieldEnum",   "type":{"name":"Enum", "labels":["a", "b", "c"]}},
    ]

@pytest.fixture
def users():
    return [
        {
            "handle": "analystUser",
            "name": "Mr. analyst",
            "password": "foobar1",
            "roles": ["Analyst"],
        },
        {
            "handle": "analystUserNoBudget",
            "name": "Mr. analystNoBudget",
            "password": "foobar1",
            "roles": ["Analyst"],
        },
        {
            "handle": "curatorUser",
            "name": "Mr. curator",
            "password": "foobar2",
            "roles": ["Curator"],
        },
        {
            "handle": "adminUser",
            "name": "Mr. admin",
            "password": "foobar3",
            "roles": ["Admin"],
        },
    ]

@pytest.fixture
def datasets(testdata_salaries_schema):
    return [
        {
            "name": "dataset1",
            "owner": "root",
            "schema": testdata_salaries_schema,
            "privacy_notion": "PureDP",
            "total_budget": PureDP(5),
        },
        {
            "name": "dataset2",
            "owner": "root",
            "schema": testdata_salaries_schema,
            "privacy_notion": "ApproxDP",
            "total_budget": ApproxDP(4,0.001),
        },
    ]

@pytest.fixture
def budgets():
    return [
        {"user": "curatorUser", "dataset": "dataset1", "budget": PureDP(2)},
        {"user": "analystUser", "dataset": "dataset1", "budget": PureDP(1)},
        {"user": "analystUser", "dataset": "dataset2", "budget": ApproxDP(2, 0.0001)},
        {"user": "root", "dataset": "dataset1", "budget": PureDP(1)},
        {"user": "root", "dataset": "dataset2", "budget": ApproxDP(1, 0.0001)},
    ]

@pytest.fixture
def make_test_env(headers_default, file_salaries, users, datasets, budgets):
    """
    Initialize a test environment with some users and datasets preloaded.
    """

    clean_test_env(headers_default)
    env = {}

    # Initialize users
    for user in users:
        response = requests.post(
            URL("users"),
            headers=headers_default,
            json=user
        )
        assert response.status_code == 201

    # Initialize datasets
    for dataset in datasets:
        response = requests.post(
            URL("datasets"),
            headers=headers_default,
            json=dataset,
        )
        assert response.status_code == 201
        dataset["id"] = str(response.json()["id"])
        with open(file_salaries) as csv:
            data = csv.read()
            response = requests.post(
                URL("dataset/{did}/upload".format(did=dataset["id"])),
                headers=headers_default,
                data=data.encode(),
            )
            assert response.status_code == 204

    # Assign budgets
    for budget in budgets:
        did = next(filter(lambda dataset: dataset["name"] == budget["dataset"], datasets), None)[
            "id"
        ]
        response = requests.post(
            URL("budget/allocation/{uid}/{did}".format(uid=budget["user"], did=did)),
            headers=headers_default,
            json=budget["budget"],
        )
        assert response.status_code == 201 # was 204

    env["users"] = users
    env["datasets"] = datasets
    env["budgets"] = budgets
    return env

def clean_test_env(headers_default_no_content):
    response = requests.get(
        URL("datasets"),
        headers=headers_default_no_content
    )
    datasets = response.json()
    for dataset in datasets:
        response = requests.delete(
            URL("dataset/"+str(dataset["id"])),
            headers=headers_default_no_content
        )
    response = requests.get(
        URL("users"),
        headers=headers_default_no_content
    )
    users = response.json()
    for user in users:
        response = requests.delete(
            URL("user/"+user["handle"]),
            headers=headers_default_no_content
        )
