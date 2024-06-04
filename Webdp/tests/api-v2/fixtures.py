import pytest
import requests

from server_env import *
from models import *

#######################################
# FUNCTIONS
#######################################

def do_login(user):
  response = requests.post(URL_LOGIN, json=user)
  assert response.status_code == 200
  return {"Authorization": "Bearer " + response.json()["jwt"]}

def do_logout(head):
  response = requests.post(URL_LOGOUT, headers=head)
  assert response.status_code == 204

#######################################
# USERS
#######################################

@pytest.fixture
def setup_users():
  head = do_login(root_login)
  # Create Admin user
  response = requests.post(URL_USERS, json=admin, headers=head)
  assert response.status_code in SUCCESS
  # Create Curator user
  response = requests.post(URL_USERS, json=curator, headers=head)
  assert response.status_code in SUCCESS
  # Create Analyst user
  response = requests.post(URL_USERS, json=analyst, headers=head)
  assert response.status_code in SUCCESS
  do_logout(head)

@pytest.fixture
def clean_users():
  head = do_login(root_login)
  response = requests.delete(URL_USER(admin["handle"]), headers=head)
  response = requests.delete(URL_USER(curator["handle"]), headers=head)
  response = requests.delete(URL_USER(analyst["handle"]), headers=head)
  response = requests.delete(URL_USER(curana["handle"]), headers=head)
  response = requests.delete(URL_USER("curatorUser"), headers=head)
  do_logout(head)

#######################################
# DATASETS
#######################################

@pytest.fixture
def root_dataset():
  head = do_login(root_login)
  response = requests.post(URL_DATASETS, json=data_root, headers=head)
  assert response.status_code in SUCCESS
  do_logout(head)
  return str(response.json()["id"])

@pytest.fixture
def curator_dataset():
  head = do_login(curator_login)
  response = requests.post(URL_DATASETS, json=data_curator, headers=head)
  assert response.status_code in SUCCESS
  do_logout(head)
  return str(response.json()["id"])

@pytest.fixture
def curator_dataset_w_budget(curator_dataset):
  head = do_login(curator_login)
  response = requests.post(URL_USER_DATASET_BUDGET(curator["handle"], curator_dataset), json=PureDP(3), headers=head)
  assert response.status_code in SUCCESS
  do_logout(head)
  return curator_dataset

@pytest.fixture
def analyst_dataset():
  head = do_login(root_login)
  # create with root owner
  response = requests.post(URL_DATASETS, json=data_root, headers=head)
  assert response.status_code in SUCCESS
  did = str(response.json()["id"])
  # patch to analyst owner
  response = requests.patch(URL_DATASET(did), json=data_patch_analyst, headers=head)
  assert response.status_code in SUCCESS
  do_logout(head)
  return did

@pytest.fixture
def analyst_dataset_w_budget(analyst_dataset):
  head = do_login(analyst_login)
  response = requests.post(URL_USER_DATASET_BUDGET(analyst["handle"], analyst_dataset), json=PureDP(3), headers=head)
  assert response.status_code in SUCCESS
  do_logout(head)
  return analyst_dataset

@pytest.fixture
def admin_dataset():
  head = do_login(root_login)
  # create with root owner
  response = requests.post(URL_DATASETS, json=data_root, headers=head)
  assert response.status_code in SUCCESS
  did = str(response.json()["id"])
  # patch to analyst owner
  response = requests.patch(URL_DATASET(did), json=data_patch_admin, headers=head)
  assert response.status_code in SUCCESS
  do_logout(head)
  return did

@pytest.fixture
def granted_access(curator_dataset, analyst_dataset):
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
  return [curator_dataset, analyst_dataset]

@pytest.fixture
def clean_datasets():
  h = do_login(root_login)
  print("deleting root datasets (if present)")
  response = requests.get(URL_DATASETS, headers=h)
  if response.status_code == 200:
    data = response.json()
    print("- number of datasets:", len(data))
    for i in range(0, len(data)):
      response = requests.delete(URL_DATASET(str(data[i]["id"])), headers=h)
      if response.status_code == 204:
        print("-- removed id", data[i]["id"])
      else:
        print("-- did not remove id", data[i]["id"], "- has owner", data[i]["owner"])
        print("--- will be deleted upon user delete (cascade)")
  do_logout(h)

