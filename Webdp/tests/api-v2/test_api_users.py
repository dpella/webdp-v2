"""
Run this file to test user APIs.


Tests covered:

    Logging in
    Logging out
LIN Logging in non-existing user

---------------------------------------------------------------
GET ALL (req: admin/curator)
---------------------------------------------------------------
GA1     admin/curator
GA2   ¬ admin/curator
---------------------------------------------------------------
GAA Consistent after adding
GAD Consistent after deleting
GAT Using old tokens
GAR Get, gain role, get again, lose role, get again
---------------------------------------------------------------

---------------------------------------------------------------
GET ONE (req: admin/curator or self-request)
---------------------------------------------------------------
GO1     admin/curator and   self
GO2   ¬ admin/curator and   self
GO3     admin/curator and ¬ self
GO4   ¬ admin/curator and ¬ self
GO5                                    ¬ exists
---------------------------------------------------------------

---------------------------------------------------------------
POST (req: admin/curator)
---------------------------------------------------------------
PO1     admin/curator
PO2   ¬ admin/curator
---------------------------------------------------------------
POB Malformed body, bad request
---------------------------------------------------------------

---------------------------------------------------------------
PATCH (req: admin)
---------------------------------------------------------------
PA1     admin
PA2   ¬ admin
---------------------------------------------------------------

---------------------------------------------------------------
DELETE (req: admin)
---------------------------------------------------------------
DE1     admin
DE2   ¬ admin
---------------------------------------------------------------
DER Delete root
---------------------------------------------------------------
"""

import requests
from server_env import *
from models import *
from fixtures import *

@pytest.fixture(autouse=True)
def setup(clean_users):
    clean_users

@pytest.fixture(autouse=True)
def teardown(clean_users):
    clean_users

class Test_UserRoot():

    # Get all users
    def test_GA1_root(self):
        head=do_login(root_login)
        response = requests.get(URL_USERS, headers=head)
        do_logout(head)
        assert response.status_code in SUCCESS
        assert len(response.json()) >= 1

    # Get root
    def test_GO1(self):
        head=do_login(root_login)
        response = requests.get(URL_USER("root"), headers=head)
        do_logout(head)
        assert response.status_code in SUCCESS

    # Delete self (fail)
    def test_DER(self):
        head=do_login(root_login)
        response = requests.delete(URL_USER("root"), headers=head)
        do_logout(head)
        assert response.status_code in FAIL

    # Create users
    def test_PO1_POB(self):
        head=do_login(root_login)
        # curator
        response = requests.post(URL_USERS, json=curator, headers=head)
        assert response.status_code in SUCCESS
        # analyst
        response = requests.post(URL_USERS, json=analyst, headers=head)
        assert response.status_code in SUCCESS
        # tester (fail)
        response = requests.post(URL_USERS, json=tester, headers=head)
        assert response.status_code in FAIL
        do_logout(head)

    # Get user x2
    def test_GO3_root(self, setup_users):
        head=do_login(root_login)
        # curator
        response = requests.get(URL_USER(curator["handle"]), headers=head)
        assert response.status_code in SUCCESS
        # analyst
        response = requests.get(URL_USER(analyst["handle"]), headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    # Get all, check len += 2
    def test_GAA(self, setup_users):
        head=do_login(root_login)
        response = requests.get(URL_USERS, headers=head)
        do_logout(head)
        assert response.status_code in SUCCESS
        assert len(response.json()) >= 3

    # Update users
    def test_PA1(self, setup_users):
        head=do_login(root_login)
        # self
        response = requests.patch(URL_USER("root"), json=root_patch, headers=head)
        assert response.status_code in SUCCESS
        head=do_login(root_login)
        # curator
        response = requests.patch(URL_USER(curator["handle"]), json=curator_patch, headers=head)
        assert response.status_code in SUCCESS
        # analyst
        response = requests.patch(URL_USER(analyst["handle"]), json=analyst_patch, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)


class Test_UserCurator():

    # Get all users
    def test_GA1_curator(self, setup_users):
        head=do_login(curator_login)
        response = requests.get(URL_USERS, headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    # Get Analyst
    def test_GO3_curator(self, setup_users):
        head=do_login(curator_login)
        response = requests.get(URL_USER(analyst["handle"]), headers=head)
        assert response.status_code in SUCCESS
        do_logout(head)

    def test_PA2_curator(self, setup_users):
        head=do_login(curator_login)
        # Update self (fail)
        curator_patch["password"] = "123123"
        response = requests.patch(URL_USER(curator["handle"]), json=curator_patch, headers=head)
        assert response.status_code in FAIL
        # Update Analyst (fail)
        analyst_patch["password"] = "123123"
        response = requests.patch(URL_USER(analyst["handle"]), json=analyst_patch, headers=head)
        assert response.status_code in FAIL
        do_logout(head)

    def test_DE2_curator(self, setup_users):
        head=do_login(curator_login)
        # Delete Analyst (fail)
        response = requests.delete(URL_USER(analyst["handle"]), headers=head)
        assert response.status_code in FAIL
        # Delete self (fail)
        response = requests.delete(URL_USER(curator["handle"]), headers=head)
        assert response.status_code in FAIL
        do_logout(head)

    # Create user (fail)
    def test_PO2_curator(self, setup_users):
        head=do_login(curator_login)
        response = requests.post(URL_USERS, json=curana, headers=head)
        assert response.status_code in FAIL
        do_logout(head)

class Test_UserAnalyst():

    # Get all users (fail)
    def test_GA2(self, setup_users):
        head = do_login(analyst_login)
        response = requests.get(URL_USERS, headers=head)
        assert response.status_code in FAIL
        do_logout(head)

    def test_GO2_GO4(self, setup_users):
        head = do_login(analyst_login)
        # Get Analyst (self)
        response = requests.get(URL_USER(analyst["handle"]), headers=head)
        assert response.status_code in SUCCESS
        # Get Curator (fail) 
        response = requests.get(URL_USER(curator["handle"]), headers=head)
        assert response.status_code in FAIL
        do_logout(head)

    def test_PA2_analyst(self, setup_users):
        head = do_login(analyst_login)
        # Update self (fail)
        analyst_patch["password"] = "123123"
        response = requests.patch(URL_USER(analyst["handle"]), json=analyst_patch, headers=head)
        assert response.status_code in FAIL
        # Update Curator (fail)
        curator_patch["password"] = "123123"
        response = requests.patch(URL_USER(curator["handle"]), json=curator_patch, headers=head)
        assert response.status_code in FAIL
        do_logout(head)

    def test_DE2_analyst(self, setup_users):
        head = do_login(analyst_login)
        # Delete Curator (fail)
        response = requests.delete(URL_USER(curator["handle"]), headers=head)
        assert response.status_code in FAIL
        # Delete self (fail)
        response = requests.delete(URL_USER(analyst["handle"]), headers=head)
        assert response.status_code in FAIL
        do_logout(head)

    # Create user (fail)
    def test_PO2_analyst(self, setup_users):
        head = do_login(analyst_login)
        response = requests.post(URL_USERS, json=curana, headers=head)
        assert response.status_code in FAIL
        do_logout(head)

class Test_UserRoleUpdate():

    def test_GAR(self):
        roothead = do_login(root_login)

        # Root: Create curator/analyst
        response = requests.post(URL_USERS, json=curana, headers=roothead)
        assert response.status_code in SUCCESS

        # Curator/analyst: get users (fail)
        curanahead = do_login(curana_login)
        response = requests.get(URL_USERS, headers=curanahead)
        assert response.status_code in FAIL

        # Root: Update curator/analyst
        response = requests.patch(URL_USER(curana["handle"]), json=curana_patch, headers=roothead)
        assert response.status_code in SUCCESS
        # login again since update
        curanahead = do_login(curana_login)
        # Curator/analyst: get users (as curator)
        response = requests.get(URL_USERS, headers=curanahead)
        assert response.status_code in SUCCESS

        # Root: Update curator/analyst
        curana_patch["roles"] = ["Analyst"]
        response = requests.patch(URL_USER(curana["handle"]), json=curana_patch, headers=roothead)
        assert response.status_code in SUCCESS
        # login again since update
        curanahead = do_login(curana_login)
        # Curator/analyst: get users (as analyst)
        response = requests.get(URL_USERS, headers=curanahead)
        assert response.status_code in FAIL

        # Logout
        do_logout(curanahead)
        do_logout(roothead)
    

class Test_UserOldTokens():

    def test_GAT(self, setup_users):
        head=do_login(root_login)
        do_logout(head)
        # Get users as root, even though logged out
        response = requests.get(URL_USERS, headers=head)
        assert response.status_code in FAIL

        head=do_login(curator_login)
        do_logout(head)
        # Get users as curator, even though logged out
        response = requests.get(URL_USERS, headers=head)
        assert response.status_code in FAIL
    

class Test_UserPostClean():

    def test_GO5_GAD(self):
        head = do_login(root_login)

        # Get Curator (empty)
        response = requests.get(URL_USER(curator["handle"]), headers=head)
        assert response.status_code in FAIL

        # Get Analyst (empty)
        response = requests.get(URL_USER(analyst["handle"]), headers=head)
        assert response.status_code in FAIL

        # Get curator/analyst (empty)
        response = requests.get(URL_USER(curana["handle"]), headers=head)
        assert response.status_code in FAIL

        # Get all, count users
        response = requests.get(URL_USERS, headers=head)
        assert response.status_code in SUCCESS
        # assert len(response.json()) == 1 # TODO other test suits leave side effects

        do_logout(head)

    def test_LIN(self):
        # Login Curator (fail)
        response = requests.post(URL_LOGIN, json=curator_login)
        assert response.status_code in FAIL

        # Login Analyst (fail)
        response = requests.post(URL_LOGIN, json=analyst_login)
        assert response.status_code in FAIL

        # Login curator/analyst (fail)
        response = requests.post(URL_LOGIN, json=curana_login)
        assert response.status_code in FAIL
