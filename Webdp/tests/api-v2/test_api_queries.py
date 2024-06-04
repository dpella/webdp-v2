"""
Run this file to test dataset APIs.


Tests covered:

- Queries on default engine
- Supported queries on Tumult
- Supported queries on OpenDP
- Supported queries on GoogleDP

Tests not covered: 

- Failing queries

"""

import pytest
import requests
import math

from server_env import *
from models import *
from fixtures import *
from queries import *

@pytest.fixture(autouse=True)
def setup(clean_datasets):
    clean_datasets

@pytest.fixture(autouse=True)
def teardown(clean_datasets):
    clean_datasets

@pytest.fixture
def did():
  head = do_login(root_login)

  # create dataset
  response = requests.post(URL_DATASETS, json=data_root, headers=head)
  assert response.status_code in SUCCESS
  did = str(response.json()["id"])

  # allocate budget
  response = requests.post(URL_USER_DATASET_BUDGET("root", did), json=PureDP(5), headers=head)
  assert response.status_code in SUCCESS

  # upload data
  head['Content-Type'] = 'text/csv'
  response = requests.post(URL_DATASET(did)+"/upload", data=FILE, headers=head)
  assert response.status_code in SUCCESS

  do_logout(head)

  return int(did)

@pytest.fixture
def did_approx():
  head = do_login(root_login)

  # create dataset
  response = requests.post(URL_DATASETS, json=data_root_approx, headers=head)
  assert response.status_code in SUCCESS
  did = str(response.json()["id"])

  # allocate budget
  response = requests.post(URL_USER_DATASET_BUDGET("root", did), json=ApproxDP(5,0.1), headers=head)
  assert response.status_code in SUCCESS

  # upload data
  head['Content-Type'] = 'text/csv'
  response = requests.post(URL_DATASET(did)+"/upload", data=FILE, headers=head)
  assert response.status_code in SUCCESS

  do_logout(head)

  return int(did)

class Test_QueryDefault():

    def test_engines(self):
        # Without parameters: will return list of all engines
        head = do_login(root_login)
        response = requests.get(URL_Q_ENGINES, headers=head)
        assert response.status_code in SUCCESS
        assert len(response.json()) > 0
        do_logout(head)

    def test_docs(self):
        # Without parameters: will return all docs
        head = do_login(root_login)
        response = requests.get(URL_Q_DOCS, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)
    
    def test_functions(self):
        # Without parameters: will return list of all engines' features
        head = do_login(root_login)
        response = requests.get(URL_Q_FUNC, headers=head)
        assert response.status_code in SUCCESS
        assert len(response.json()) > 0
        do_logout(head)

    def test_validate(self, did): 
        # Without parameters: will return list of all compatible engines
        head = do_login(root_login)
        query = QUERY(did, COUNT) # default is tumult
        response = requests.post(URL_Q_VAL, json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_evaluate(self, did): 
        # Without parameters: will run on default engine
        head = do_login(root_login)
        query = QUERY(did, COUNT) # default is tumult
        response = requests.post(URL_Q_EVAL, json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_evaluate_unsupported(self, did): 
        # Testing response for unsupported queries
        head = do_login(root_login)
        query = QUERY(did, TUM_MIN) # only supported on tumult; not openDP, googleDP
        response = requests.post(URL_Q_EVAL_E("opendp"), json=query, headers=head)
        assert response.status_code in FAIL
        do_logout(head)

    def test_evaluate_bad(self, did): 
        # Testing response for unsupported queries
        head = do_login(root_login)
        query = [{"sum": {}}] # invalid query
        response = requests.post(URL_Q_EVAL_E("opendp"), json=query, headers=head)
        assert response.status_code in FAIL
        do_logout(head)

class Test_QueryTumult():

    def test_docs(self):
        head = do_login(root_login)
        response = requests.get(URL_Q_DOCS_E("tumult"), headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_functions(self):
        head = do_login(root_login)
        response = requests.get(URL_Q_FUNC_E("tumult"), headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_validate_count(self, did):
        head = do_login(root_login)
        query = QUERY(did, COUNT)
        response = requests.post(URL_Q_VAL_E("tumult"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_evaluate_count(self, did):
        head = do_login(root_login)
        query = QUERY(did, COUNT)
        response = requests.post(URL_Q_EVAL_E("tumult"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_validate_sum(self, did):
        head = do_login(root_login)
        query = QUERY(did, SUM)
        response = requests.post(URL_Q_VAL_E("tumult"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_evaluate_sum(self, did):
        head = do_login(root_login)
        query = QUERY(did, SUM)
        response = requests.post(URL_Q_EVAL_E("tumult"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_validate_mean(self, did):
        head = do_login(root_login)
        query = QUERY(did, MEAN)
        response = requests.post(URL_Q_VAL_E("tumult"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_evaluate_mean(self, did):
        head = do_login(root_login)
        query = QUERY(did, MEAN)
        response = requests.post(URL_Q_EVAL_E("tumult"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_validate_min(self, did):
        head = do_login(root_login)
        query = QUERY(did, TUM_MIN)
        response = requests.post(URL_Q_VAL_E("tumult"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_evaluate_min(self, did):
        head = do_login(root_login)
        query = QUERY(did, TUM_MIN)
        response = requests.post(URL_Q_EVAL_E("tumult"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_validate_max(self, did):
        head = do_login(root_login)
        query = QUERY(did, TUM_MAX)
        response = requests.post(URL_Q_VAL_E("tumult"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_evaluate_max(self, did):
        head = do_login(root_login)
        query = QUERY(did, TUM_MAX)
        response = requests.post(URL_Q_EVAL_E("tumult"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_validate_filter_count(self, did):
        head = do_login(root_login)
        query = QUERY(did, FILTER_COUNT)
        response = requests.post(URL_Q_VAL_E("tumult"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_evaluate_filter_count(self, did):
        head = do_login(root_login)
        query = QUERY(did, FILTER_COUNT)
        response = requests.post(URL_Q_EVAL_E("tumult"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_validate_filter_sum(self, did):
        head = do_login(root_login)
        query = QUERY(did, FILTER_SUM)
        response = requests.post(URL_Q_VAL_E("tumult"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_evaluate_filter_sum(self, did):
        head = do_login(root_login)
        query = QUERY(did, FILTER_SUM)
        response = requests.post(URL_Q_EVAL_E("tumult"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_validate_filter_mean(self, did):
        head = do_login(root_login)
        query = QUERY(did, FILTER_MEAN)
        response = requests.post(URL_Q_VAL_E("tumult"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_evaluate_filter_mean(self, did):
        head = do_login(root_login)
        query = QUERY(did, FILTER_MEAN)
        response = requests.post(URL_Q_EVAL_E("tumult"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_validate_bin_count(self, did):
        head = do_login(root_login)
        query = QUERY(did, TUM_BIN_COUNT)
        response = requests.post(URL_Q_VAL_E("tumult"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_evaluate_bin_count(self, did):
        head = do_login(root_login)
        query = QUERY(did, TUM_BIN_COUNT)
        response = requests.post(URL_Q_EVAL_E("tumult"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_validate_bin_sum(self, did):
        head = do_login(root_login)
        query = QUERY(did, TUM_BIN_SUM)
        response = requests.post(URL_Q_VAL_E("tumult"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_evaluate_bin_sum(self, did):
        head = do_login(root_login)
        query = QUERY(did, TUM_BIN_SUM)
        response = requests.post(URL_Q_EVAL_E("tumult"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_validate_bin_mean(self, did):
        head = do_login(root_login)
        query = QUERY(did, TUM_BIN_MEAN)
        response = requests.post(URL_Q_VAL_E("tumult"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_evaluate_bin_mean(self, did):
        head = do_login(root_login)
        query = QUERY(did, TUM_BIN_MEAN)
        response = requests.post(URL_Q_EVAL_E("tumult"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_validate_filter_bin_count(self, did):
        head = do_login(root_login)
        query = QUERY(did, TUM_FILTER_BIN_COUNT)
        response = requests.post(URL_Q_VAL_E("tumult"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_evaluate_filter_bin_count(self, did):
        head = do_login(root_login)
        query = QUERY(did, TUM_FILTER_BIN_COUNT)
        response = requests.post(URL_Q_EVAL_E("tumult"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_validate_filter_bin_sum(self, did):
        head = do_login(root_login)
        query = QUERY(did, TUM_FILTER_BIN_SUM)
        response = requests.post(URL_Q_VAL_E("tumult"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_evaluate_filter_bin_sum(self, did):
        head = do_login(root_login)
        query = QUERY(did, TUM_FILTER_BIN_SUM)
        response = requests.post(URL_Q_EVAL_E("tumult"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_validate_filter_bin_mean(self, did):
        head = do_login(root_login)
        query = QUERY(did, TUM_FILTER_BIN_MEAN)
        response = requests.post(URL_Q_VAL_E("tumult"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_evaluate_filter_bin_mean(self, did):
        head = do_login(root_login)
        query = QUERY(did, TUM_FILTER_BIN_MEAN)
        response = requests.post(URL_Q_EVAL_E("tumult"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_validate_groupby_mean(self, did):
        head = do_login(root_login)
        query = QUERY(did, TUM_GBY_MEAN)
        response = requests.post(URL_Q_VAL_E("tumult"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_evaluate_groupby_mean(self, did):
        head = do_login(root_login)
        query = QUERY(did, TUM_GBY_MEAN)
        response = requests.post(URL_Q_EVAL_E("tumult"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_validate_bin_groupby_mean(self, did):
        head = do_login(root_login)
        query = QUERY(did, TUM_BIN_GBY_MEAN)
        response = requests.post(URL_Q_VAL_E("tumult"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_evaluate_bin_groupby_mean(self, did):
        head = do_login(root_login)
        query = QUERY(did, TUM_BIN_GBY_MEAN)
        response = requests.post(URL_Q_EVAL_E("tumult"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

class Test_QueryOpenDP():

    def test_docs(self):
        head = do_login(root_login)
        response = requests.get(URL_Q_DOCS_E("opendp"), headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)
    
    def test_functions(self):
        head = do_login(root_login)
        response = requests.get(URL_Q_FUNC_E("opendp"), headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_validate_count(self, did):
        head = do_login(root_login)
        query = QUERY(did, COUNT)
        response = requests.post(URL_Q_VAL_E("opendp"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_evaluate_count(self, did):
        head = do_login(root_login)
        query = QUERY(did, COUNT)
        response = requests.post(URL_Q_EVAL_E("opendp"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_validate_sum(self, did):
        head = do_login(root_login)
        query = QUERY(did, SUM)
        response = requests.post(URL_Q_VAL_E("opendp"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_evaluate_sum(self, did):
        head = do_login(root_login)
        query = QUERY(did, SUM)
        response = requests.post(URL_Q_EVAL_E("opendp"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_sum_accuracy_pure(self, did):
        head = do_login(root_login)
        query = QUERY(did, SUM)
        query["confidence"] = math.sqrt(0.95)
        response = requests.post(URL_Q_ACC_E("opendp"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_count_accuracy_pure(self, did):
        head = do_login(root_login)
        query = QUERY(did, COUNT)
        query["confidence"] = math.sqrt(0.95)
        response = requests.post(URL_Q_ACC_E("opendp"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_sum_accuracy_approx(self, did_approx):
        head = do_login(root_login)
        query = QUERY_approx(did_approx, SUM_approx)
        query["confidence"] = math.sqrt(0.95)
        response = requests.post(URL_Q_ACC_E("opendp"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_count_accuracy_approx(self, did_approx):
        head = do_login(root_login)
        query = QUERY_approx(did_approx, COUNT_approx)
        query["confidence"] = math.sqrt(0.95)
        response = requests.post(URL_Q_ACC_E("opendp"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

class Test_QueryGoogleDP():

    def test_docs(self):
        head = do_login(root_login)
        response = requests.get(URL_Q_DOCS_E("googledp"), headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_functions(self):
        head = do_login(root_login)
        response = requests.get(URL_Q_FUNC_E("googledp"), headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_validate_count(self, did):
        head = do_login(root_login)
        query = QUERY(did, COUNT)
        response = requests.post(URL_Q_VAL_E("googledp"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_evaluate_count(self, did):
        head = do_login(root_login)
        query = QUERY(did, COUNT)
        response = requests.post(URL_Q_EVAL_E("googledp"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_validate_sum(self, did):
        head = do_login(root_login)
        query = QUERY(did, SUM)
        response = requests.post(URL_Q_VAL_E("googledp"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_evaluate_sum(self, did):
        head = do_login(root_login)
        query = QUERY(did, SUM)
        response = requests.post(URL_Q_EVAL_E("googledp"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)
    
    def test_validate_mean(self, did):
        head = do_login(root_login)
        query = QUERY(did, MEAN)
        response = requests.post(URL_Q_VAL_E("googledp"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_evaluate_mean(self, did):
        head = do_login(root_login)
        query = QUERY(did, MEAN)
        response = requests.post(URL_Q_EVAL_E("googledp"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_validate_filter_count(self, did):
        head = do_login(root_login)
        query = QUERY(did, FILTER_COUNT)
        response = requests.post(URL_Q_VAL_E("googledp"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_evaluate_filter_count(self, did):
        head = do_login(root_login)
        query = QUERY(did, FILTER_COUNT)
        response = requests.post(URL_Q_EVAL_E("googledp"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_validate_filter_sum(self, did):
        head = do_login(root_login)
        query = QUERY(did, FILTER_SUM)
        response = requests.post(URL_Q_VAL_E("googledp"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_evaluate_filter_sum(self, did):
        head = do_login(root_login)
        query = QUERY(did, FILTER_SUM)
        response = requests.post(URL_Q_EVAL_E("googledp"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_validate_filter_mean(self, did):
        head = do_login(root_login)
        query = QUERY(did, FILTER_MEAN)
        response = requests.post(URL_Q_VAL_E("googledp"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_evaluate_filter_mean(self, did):
        head = do_login(root_login)
        query = QUERY(did, FILTER_MEAN)
        response = requests.post(URL_Q_EVAL_E("googledp"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_validate_bin_count(self, did):
        head = do_login(root_login)
        query = QUERY(did, GDP_BIN_COUNT)
        response = requests.post(URL_Q_VAL_E("googledp"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_evaluate_bin_count(self, did):
        head = do_login(root_login)
        query = QUERY(did, GDP_BIN_COUNT)
        response = requests.post(URL_Q_EVAL_E("googledp"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_validate_bin_sum(self, did):
        head = do_login(root_login)
        query = QUERY(did, GDP_BIN_SUM)
        response = requests.post(URL_Q_VAL_E("googledp"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_evaluate_bin_sum(self, did):
        head = do_login(root_login)
        query = QUERY(did, GDP_BIN_SUM)
        response = requests.post(URL_Q_EVAL_E("googledp"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_validate_bin_mean(self, did):
        head = do_login(root_login)
        query = QUERY(did, GDP_BIN_MEAN)
        response = requests.post(URL_Q_VAL_E("googledp"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_evaluate_bin_mean(self, did):
        head = do_login(root_login)
        query = QUERY(did, GDP_BIN_MEAN)
        response = requests.post(URL_Q_EVAL_E("googledp"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_validate_filter_bin_count(self, did):
        head = do_login(root_login)
        query = QUERY(did, GDP_FILTER_BIN_COUNT)
        response = requests.post(URL_Q_VAL_E("googledp"), json=query, headers=head)

        assert response.status_code in SUCCESS
        do_logout(head)

    def test_evaluate_filter_bin_count(self, did):
        head = do_login(root_login)
        query = QUERY(did, GDP_FILTER_BIN_COUNT)
        response = requests.post(URL_Q_EVAL_E("googledp"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_validate_filter_bin_sum(self, did):
        head = do_login(root_login)
        query = QUERY(did, GDP_FILTER_BIN_SUM)
        response = requests.post(URL_Q_VAL_E("googledp"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_evaluate_filter_bin_sum(self, did):
        head = do_login(root_login)
        query = QUERY(did, GDP_FILTER_BIN_SUM)
        response = requests.post(URL_Q_EVAL_E("googledp"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_validate_filter_bin_mean(self, did):
        head = do_login(root_login)
        query = QUERY(did, GDP_FILTER_BIN_MEAN)
        response = requests.post(URL_Q_VAL_E("googledp"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_evaluate_filter_bin_mean(self, did):
        head = do_login(root_login)
        query = QUERY(did, GDP_FILTER_BIN_MEAN)
        response = requests.post(URL_Q_EVAL_E("googledp"), json=query, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)
