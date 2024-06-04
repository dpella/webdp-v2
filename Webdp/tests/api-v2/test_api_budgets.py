"""
Run this file to test budget APIs.


Tests covered:

---------------------------------------------------------------
GET ALL DATASETS (req: admin/curator or granted access)
   - GA3 & GA4 covered in dataset tests
---------------------------------------------------------------
GA1     admin/curator and   granted access
GA2   ¬ admin/curator and   granted access
...     admin/curator and ¬ granted access
...   ¬ admin/curator and ¬ granted access
---------------------------------------------------------------

---------------------------------------------------------------
GET ONE DATASET (req: admin/curator or granted access)
   - GA3 & GA4 covered in dataset tests
---------------------------------------------------------------
GO1  (  admin/curator and   granted access)
GO2  (¬ admin/curator and   granted access)
...  (  admin/curator and ¬ granted access)
...  (¬ admin/curator and ¬ granted access)
---------------------------------------------------------------

---------------------------------------------------------------
GET USER BUDGETS (req: admin/curator or self-request)
---------------------------------------------------------------
GUB1    admin/curator   and   self
GUB2  ¬ admin/curator   and   self
GUB3    admin/curator   and ¬ self
GUB4  ¬ admin/curator   and ¬ self
---------------------------------------------------------------

---------------------------------------------------------------
GET DATASET BUDGET (req: curator or granted access)
---------------------------------------------------------------
GDB1    curator         and   granted access
GDB2  ¬ curator         and   granted access
GDB3    curator         and ¬ granted access
GDB4  ¬ curator         and ¬ granted access
---------------------------------------------------------------

---------------------------------------------------------------
GET USER DATASET BUDGET (req: curator or self-request)
---------------------------------------------------------------
GUT1    curator         and   self
GUD2  ¬ curator         and   self
GUD3    curator         and ¬ self
GUD4  ¬ curator         and ¬ self
GUD5                                      ¬ exists
---------------------------------------------------------------

---------------------------------------------------------------
POST USER DATASET BUDGET (req: curator/analyst and is owner)
---------------------------------------------------------------
POB1    curator/analyst and   owner
POB2  ¬ curator/analyst and   owner
POB3    curator/analyst and ¬ owner
POB4  ¬ curator/analyst and ¬ owner
---------------------------------------------------------------


---------------------------------------------------------------
PATCH USER DATASET BUDGET (req: is owner)
---------------------------------------------------------------
PAB1                          owner
PAB2                        ¬ owner
---------------------------------------------------------------

---------------------------------------------------------------
DELETE USER DATASET BUDGET (req: is owner)
---------------------------------------------------------------
 DUD1                         owner
 DUD2                       ¬ owner
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

class Test_BudgetAdmin():

    # Get own budget (empty)
    def test_GUB1(self):
        head = do_login(admin_login)
        response = requests.get(URL_USER_BUDGET(admin["handle"]), headers=head)
        assert response.status_code in SUCCESS
        assert len(response.json()) == 0
        do_logout(head)

    # Get curator's budget (fail)
    def test_GUB3(self):
        head = do_login(admin_login)
        response = requests.get(URL_USER_BUDGET(curator["handle"]), headers=head)
        assert response.status_code in FAIL
        do_logout(head)

    # Post dataset budget (fail)
    def test_POB2(self, admin_dataset):
        head = do_login(admin_login)
        response = requests.post(URL_USER_DATASET_BUDGET(admin["handle"], admin_dataset), json=admin_budget, headers=head)
        assert response.status_code in FAIL
        do_logout(head)

    # Post dataset budget (fail)
    def test_POB4(self, curator_dataset):
        head = do_login(admin_login)
        response = requests.post(URL_USER_DATASET_BUDGET(curator["handle"], curator_dataset), json=curator_budget, headers=head)
        assert response.status_code in FAIL
        do_logout(head)

class Test_BudgetAnalyst():

    # Get own budget
    def test_GUB2(self):
        head = do_login(analyst_login)
        response = requests.get(URL_USER_BUDGET(analyst["handle"]), headers=head)
        assert response.status_code in SUCCESS
        assert len(response.json()) == 0
        do_logout(head)

    # Get curator's budget (fail)
    def test_GUB4(self):
        head = do_login(analyst_login)
        response = requests.get(URL_USER_BUDGET(curator["handle"]), headers=head)
        assert response.status_code in FAIL
        do_logout(head)

    # Get dataset budget (fail)
    def test_GDB4(self, curator_dataset):
        head = do_login(analyst_login)
        response = requests.get(URL_DATASET_BUDGET(curator_dataset), headers=head)
        assert response.status_code in FAIL
        do_logout(head)

    # Get user dataset budget (does not exist)
    def test_GUD5(self, analyst_dataset):
        head = do_login(analyst_login)
        response = requests.get(URL_USER_DATASET_BUDGET(analyst["handle"], analyst_dataset), headers=head)
        assert response.status_code in FAIL # TODO gives 500
        do_logout(head)

    # Get user dataset budget
    def test_GUD4(self,curator_dataset):
        head = do_login(analyst_login)
        response = requests.get(URL_USER_DATASET_BUDGET(curator["handle"], curator_dataset), headers=head)
        assert response.status_code in FAIL
        do_logout(head)

    # Post dataset budget
    def test_POB1(self, analyst_dataset):
        head = do_login(analyst_login)
        response = requests.post(URL_USER_DATASET_BUDGET(analyst["handle"], analyst_dataset), json=analyst_budget, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    # Post dataset budget (fail)
    def test_POB3(self, curator_dataset):
        head = do_login(analyst_login)
        response = requests.post(URL_USER_DATASET_BUDGET(curator["handle"], curator_dataset), json=curator_budget, headers=head)
        assert response.status_code in FAIL
        do_logout(head)

    def test_GUD2(self, analyst_dataset_w_budget):
        head = do_login(analyst_login)
        response = requests.get(URL_USER_DATASET_BUDGET(analyst["handle"], analyst_dataset_w_budget), headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    # Get user dataset budget
    def test_GUD4(self, curator_dataset_w_budget):
        head = do_login(analyst_login)
        response = requests.get(URL_USER_DATASET_BUDGET(curator["handle"], curator_dataset_w_budget), headers=head)
        assert response.status_code in FAIL
        do_logout(head)

class Test_BudgetCurator():

    # Post dataset budget
    def test_POB1(self, curator_dataset):
        head = do_login(curator_login)
        response = requests.post(URL_USER_DATASET_BUDGET(curator["handle"], curator_dataset), json=curator_budget, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    # Get dataset budget
    def test_GDB3(self, curator_dataset):
        head = do_login(curator_login)
        response = requests.get(URL_DATASET_BUDGET(curator_dataset), headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    # Get user dataset budget
    def test_GUD1(self, curator_dataset_w_budget):
        head = do_login(curator_login)
        response = requests.get(URL_USER_DATASET_BUDGET(curator["handle"], curator_dataset_w_budget), headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    # Get user dataset budget
    def test_GUD3(self, analyst_dataset_w_budget):
        head = do_login(curator_login)
        response = requests.get(URL_USER_DATASET_BUDGET(analyst["handle"], analyst_dataset_w_budget), headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    # Patch dataset budget
    def test_PAB1(self, curator_dataset_w_budget):
        head = do_login(curator_login)
        response = requests.patch(URL_USER_DATASET_BUDGET(curator["handle"], curator_dataset_w_budget), json=curator_budget, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    # Patch dataset budget (fail)
    def test_PAB2(self, analyst_dataset_w_budget):
        head = do_login(curator_login)
        response = requests.patch(URL_USER_DATASET_BUDGET(analyst["handle"], analyst_dataset_w_budget), json=analyst_budget, headers=head)
        assert response.status_code in FAIL
        do_logout(head)

class Test_BudgetGrantedAccess():

    def test_POB1(self, curator_dataset, analyst_dataset):
        curhead = do_login(curator_login) 
        anahead = do_login(analyst_login)

        # Post dataset budget (analyst budget on curator's dataset)
        response = requests.post(URL_USER_DATASET_BUDGET(analyst["handle"], curator_dataset), json=analyst_budget, headers=curhead)
        assert response.status_code in SUCCESS

        # Post dataset budget (curator budget on analyst's dataset)
        response = requests.post(URL_USER_DATASET_BUDGET(curator["handle"], analyst_dataset), json=curator_budget, headers=anahead)
        assert response.status_code in SUCCESS

        do_logout(curhead)
        do_logout(anahead)

    def test_GA1_GA2(self, granted_access):
        curhead = do_login(curator_login) 
        anahead = do_login(analyst_login)

        # Get all datasets (curator)
        response = requests.get(URL_DATASETS, headers=curhead)
        assert response.status_code in SUCCESS
        assert len(response.json()) == 2

        # Get all datasets (analyst)
        response = requests.get(URL_DATASETS, headers=anahead)
        assert response.status_code in SUCCESS
        assert len(response.json()) == 1 
    
        do_logout(curhead)
        do_logout(anahead)

    def test_GO1_GO2(self, granted_access):
        curator_dataset = granted_access[0]
        analyst_dataset = granted_access[1]

        curhead = do_login(curator_login) 
        anahead = do_login(analyst_login)

        # Curator gets analyst's dataset (granted access)
        response = requests.get(URL_DATASET(analyst_dataset), headers=curhead)
        assert response.status_code in SUCCESS

        # Analyst gets curators's dataset (granted access)
        response = requests.get(URL_DATASET(curator_dataset), headers=anahead)
        assert response.status_code in SUCCESS

        do_logout(curhead)
        do_logout(anahead)

    def test_GDB1_GDB2(self, granted_access):
        curator_dataset = granted_access[0]
        analyst_dataset = granted_access[1]
        
        curhead = do_login(curator_login) 
        anahead = do_login(analyst_login)

        # Curator gets budget for analyst's dataset (granted access)
        response = requests.get(URL_DATASET_BUDGET(analyst_dataset), headers=curhead)
        assert response.status_code in SUCCESS

        # Analyst gets budget for curators's dataset (granted access)
        response = requests.get(URL_DATASET_BUDGET(curator_dataset), headers=anahead)
        assert response.status_code in SUCCESS
    
        do_logout(curhead)
        do_logout(anahead)

    def test_DUD1_DUD2(self, granted_access):
        curator_dataset = granted_access[0]
        analyst_dataset = granted_access[1]
        
        curhead = do_login(curator_login) 
        anahead = do_login(analyst_login)

        # Delete dataset budget (analyst budget on curator's dataset)
        response = requests.delete(URL_USER_DATASET_BUDGET(analyst["handle"], curator_dataset), headers=curhead)
        assert response.status_code in SUCCESS

        # Delete dataset budget (analyst buget on analyst's dataset)
        response = requests.delete(URL_USER_DATASET_BUDGET(analyst["handle"], analyst_dataset), headers=curhead)
        assert response.status_code in FAIL

        # Logout
        do_logout(curhead)
        do_logout(anahead)
