"""
Run this file to test dataset APIs.


Tests covered:

---------------------------------------------------------------
GET ALL (req: admin/curator or granted access)
   - GA1 & GA2 covered in budget tests
---------------------------------------------------------------
...     admin/curator and   granted access
...   ¬ admin/curator and   granted access
GA3     admin/curator and ¬ granted access
GA4   ¬ admin/curator and ¬ granted access
---------------------------------------------------------------

---------------------------------------------------------------
GET ONE (req: admin/curator or granted access)
   - GO1 & GO2 covered in budget tests
---------------------------------------------------------------
...  (  admin/curator and   granted access) and   exists
...  (¬ admin/curator and   granted access) and   exists
GO3  (  admin/curator and ¬ granted access) and   exists
GO4  (¬ admin/curator and ¬ granted access) and   exists
GO5                                             ¬ exists
---------------------------------------------------------------

---------------------------------------------------------------
POST (req: curator and curator owner)
---------------------------------------------------------------
PO1     curator       and   new owner curator
PO2   ¬ curator       and   new owner curator
PO3     curator       and ¬ new owner curator
PO4   ¬ curator       and ¬ new owner curator

---------------------------------------------------------------


---------------------------------------------------------------
POST (req: curator)
---------------------------------------------------------------
PNB1    curator     and     PureDP      and     delta > 0
PNB2    curator     and     ApproxDP    and     delta is Null
PNB3    curator     and     PureDP      and     epsilon < 0
PNB4    curator     and     ApproxDP    and     delta < 0
---------------------------------------------------------------


---------------------------------------------------------------
PATCH (req: curator and is owner)
---------------------------------------------------------------
PA1  (  curator       and   owner)          and   exists
PA2  (¬ curator       and   owner)          and   exists
PA3  (  curator       and ¬ owner)          and   exists
PA4  (¬ curator       and ¬ owner)          and   exists
PA5                                             ¬ exists
---------------------------------------------------------------

---------------------------------------------------------------
DELETE (req: is owner)
---------------------------------------------------------------
DE1                         owner           and   exists
DE2                       ¬ owner           and   exists
DE3                                             ¬ exists
---------------------------------------------------------------

---------------------------------------------------------------
POST UPLOAD (req: is owner)
---------------------------------------------------------------
PU1                         owner           and   exists
PU2                       ¬ owner           and   exists
PU3                                             ¬ exists
PU4 badly formatted schema
PU5 badly formatted data
---------------------------------------------------------------
"""

import requests
from server_env import *
from models import *
from fixtures import *

@pytest.fixture(autouse=True)
def setup(clean_datasets, clean_users, setup_users):
    clean_datasets
    clean_users
    setup_users

@pytest.fixture(autouse=True)
def teardown(clean_datasets, clean_users):
    clean_datasets
    clean_users

class Test_DatasetRoot():

    # Create dataset
    def test_PO1(self):
        head = do_login(root_login)
        response = requests.post(URL_DATASETS, json=data_root, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)
    
    # Create a bad dataset (PureDP & delta > 0)
    def test_PNB1(self):
        head = do_login(root_login)
        response = requests.post(URL_DATASETS, json=data_root_pure_notion_delta_positive, headers=head)
        assert response.status_code in FAIL
        do_logout(head)

    # Create a bad dataset (ApproxDP & delta is Null)
    def test_PNB2(self):
        head = do_login(root_login)
        response = requests.post(URL_DATASETS, json=data_root_approx_notion_delta_null, headers=head)
        assert response.status_code in FAIL
        do_logout(head)

    # Create a bad dataset (epsilon is negative)
    def test_PNB3(self):
        head = do_login(root_login)
        response = requests.post(URL_DATASETS, json=data_root_bad_epsilon, headers=head)
        assert response.status_code in FAIL 
        do_logout(head)
    
    # Create a bad dataset (delta is negative)
    def test_PNB4(self):
        head = do_login(root_login)
        response = requests.post(URL_DATASETS, json=data_root_bad_delta, headers=head)
        assert response.status_code in FAIL 
        do_logout(head)
        
    # Get all datasets
    def test_GA3(self):
        head = do_login(root_login)
        response = requests.get(URL_DATASETS, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    # get dataset
    def test_GO3_owner(self, root_dataset):
        head = do_login(root_login)
        response = requests.get(URL_DATASET(root_dataset), headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    # get dataset (not owner)
    def test_GO3_not_owner(self, curator_dataset):
        head = do_login(root_login)
        response = requests.get(URL_DATASET(curator_dataset), headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    # Update dataset  (root owner -> curator owner)
    def test_PA1(self, root_dataset):
        head = do_login(root_login)
        response = requests.patch(URL_DATASET(root_dataset), json=data_patch_curator, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    # Update dataset (not owner)
    def test_PA3(self, curator_dataset):
        head = do_login(root_login)
        response = requests.patch(URL_DATASET(curator_dataset), json=data_patch_curator, headers=head)
        assert response.status_code in FAIL
        do_logout(head)

    # upload data to dataset (not owner)
    def test_PU2(self, curator_dataset):
        head = do_login(root_login)
        response = requests.post(URL_DATASET(curator_dataset)+"/upload", data=FILE, headers=head)
        assert response.status_code in FAIL
        do_logout(head)

    # Delete dataset (not owner)
    def test_DE2(self, curator_dataset):
        head = do_login(root_login)
        response = requests.delete(URL_DATASET(curator_dataset), headers=head)
        assert response.status_code in FAIL
        do_logout(head)

    # Upload data to dataset (bad schema)
    def test_PU4(self, root_bad_dataset):
        head = do_login(root_login)
        response = requests.post(URL_DATASET(root_bad_dataset)+"/upload", data=FILE, headers=head)
        assert response.status_code in FAIL
        do_logout(head)

    # Upload data to dataset (bad data)
    def test_PU4(self, root_dataset):
        head = do_login(root_login)
        response = requests.post(URL_DATASET(root_dataset)+"/upload", data=FILE_BAD, headers=head)
        assert response.status_code in FAIL
        do_logout(head)

class Test_DatasetCurator():

    # Get all datasets
    def test_GA3(self):
        head = do_login(curator_login)
        response = requests.get(URL_DATASETS, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    # get dataset
    def test_GO3(self, curator_dataset):
        head = do_login(curator_login)
        response = requests.get(URL_DATASET(curator_dataset), headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    # Update dataset (curator owner)
    def test_PA1(self, curator_dataset):
        head = do_login(curator_login)
        response = requests.patch(URL_DATASET(curator_dataset), json=data_patch_curator, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    # Upload data
    def test_PU1(self, curator_dataset):
        head = do_login(curator_login)
        response = requests.post(URL_DATASET(curator_dataset)+"/upload", data=FILE, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    # Delete dataset
    def test_DE1_DE3(self, curator_dataset):
        head = do_login(curator_login)
        response = requests.delete(URL_DATASET(curator_dataset), headers=head)
        assert response.status_code in SUCCESS
        # Delete again (fail)
        response = requests.delete(URL_DATASET(curator_dataset), headers=head)
        assert response.status_code in FAIL
        do_logout(head)

    # get dataset that does not exist (fail)
    def test_GO5(self):
        head = do_login(curator_login)
        response = requests.get(URL_DATASET("0"), headers=head)
        assert response.status_code in FAIL
        do_logout(head)

    # Update dataset that does not exist (fail)
    def test_PA5(self):
        head = do_login(curator_login)
        response = requests.patch(URL_DATASET("0"), json=data_patch_curator, headers=head)
        assert response.status_code in FAIL
        do_logout(head)

    # upload data to dataset (fail)
    def test_PU3(self):
        head = do_login(curator_login)
        response = requests.post(URL_DATASET("0")+"/upload", data=FILE, headers=head)
        assert response.status_code in FAIL
        do_logout(head)

    # Create dataset 
    def test_PO1_curator(self):
        head = do_login(curator_login)
        response = requests.post(URL_DATASETS, json=data_curator, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    # Create dataset with wrong owner (fail)
    def test_PO3(self):
        head = do_login(curator_login)
        response = requests.post(URL_DATASETS, json=data_analyst, headers=head)
        assert response.status_code in FAIL
        do_logout(head)

    # Update dataset (curator owner -> analyst owner)
    def test_PA1(self, curator_dataset):
        head = do_login(curator_login)
        response = requests.patch(URL_DATASET(curator_dataset), json=data_patch_analyst, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)


class Test_DatasetAnalyst():

    # Get all datasets (fail)
    def test_GA4(self):
        head = do_login(analyst_login)
        response = requests.get(URL_DATASETS, headers=head)
        assert response.status_code in FAIL
        do_logout(head)

    # get dataset (fail)
    def test_GO4(self, curator_dataset):
        head = do_login(analyst_login)
        response = requests.get(URL_DATASET(curator_dataset), headers=head)
        assert response.status_code in FAIL
        do_logout(head)

    # Create dataset (fail)
    def test_PO2(self):
        head = do_login(analyst_login)
        response = requests.post(URL_DATASETS, json=data_curator, headers=head)
        assert response.status_code in FAIL
        do_logout(head)

    # Update dataset (fail)
    def test_PA2(self, analyst_dataset):
        head = do_login(analyst_login)
        response = requests.patch(URL_DATASET(analyst_dataset), json=data_patch_curator, headers=head)
        assert response.status_code in FAIL
        do_logout(head)

    # Delete dataset
    def test_DE1(self, analyst_dataset):
        head = do_login(analyst_login)
        response = requests.delete(URL_DATASET(analyst_dataset), headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    # Update dataset (fail)
    def test_PA4(self, curator_dataset):
        head = do_login(analyst_login)
        response = requests.patch(URL_DATASET(curator_dataset), json=data_patch_curator, headers=head)
        assert response.status_code in FAIL
        do_logout(head)

    # Create dataset with wrong owner (fail)
    def test_PO4(self):
        head = do_login(analyst_login)
        response = requests.post(URL_DATASETS, json=data_analyst, headers=head)
        assert response.status_code in FAIL
        do_logout(head)
