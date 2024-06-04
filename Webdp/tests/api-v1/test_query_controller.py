import pytest
import requests
from regression_env import *

@pytest.mark.skip("Tumult Analytics does not support computing query accuracy")
def test_query_accuracy(headers_default):
    """Test case for query_accuracy

    Compute the accuracy a differential privacy query over a dataset
    """
    query_accuracy_request = {
        "dataset":42, 
        "query":[], 
        "budget":PureDP(1), 
        "confidence":0.05
    }
    response = requests.post(
        URL("query/accuracy"),
        headers=headers_default,
        json=query_accuracy_request,
    )
    assert response.status_code == 200

@pytest.mark.skip("Tumult Analytics does not support computing custom queries")
def test_query_custom(headers_default):
    """Test case for query_custom

    Evaluate a custom differential privacy query over a dataset
    """
    query_custom_request = {
        "dataset":42, 
        "query":"some query", 
        "budget":PureDP(1)
    }
    response = requests.post(
        URL("query/custom"),
        headers=headers_default,
        json=query_custom_request,
    )
    assert response.status_code == 200

def test_query_evaluate(make_test_env, headers_default):
    """Test case for query_evaluate

    Evaluate a differential privacy query over a dataset
    """
    test_dataset = make_test_env["datasets"][0]

    # query_steps = [
    #     QueryStep(filter=["age >= 20", "age <= 80", "salary >= 50000"]),
    #     QueryStep(select=["age", "salary"]),
    #     QueryStep(bin={"age":[20, 30, 40, 50, 60, 70]}),
    #     QueryStep(groupby={"age":[30, 40, 50, 60, 70]}),
    #     QueryStep(mean=MeasurementParams(column="salary")),
    # ]
    query_steps = [
        {"filter": ["age >= 20", "age <= 80", "salary >= 50000"]},
        {"select":["age", "salary"]},
        {"bin":{"age":[20, 30, 40, 50, 60, 70]}},
        {"groupby":{"age":[30, 40, 50, 60, 70]}},
        {"mean":{"column":"salary"}}
    ]

    query_evaluate_request = {
        "dataset":int(test_dataset["id"]), 
        "query":query_steps, 
        "budget":PureDP(0.1)
    }
    response = requests.post(
        URL("query/evaluate"), # NOTE this should go directly to Tumult in v1 testing
        headers=headers_default,
        json=query_evaluate_request,
    )
    print(response.json())
    assert response.status_code == 200
