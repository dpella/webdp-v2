from typing import Dict
import requests 
import random as rd



BASE                     = "http://localhost:8000/v2"
URL_LOGIN                = BASE + "/login"
URL_DATASETS             = BASE + "/datasets"
URL_USERS                = f"{BASE}/users"
URL_DATASETS_ID          = lambda did: f"{BASE}/datasets/{did}"
URL_DATASETS_UPLOAD      = lambda did: f"{BASE}/datasets/{did}/upload"
URL_USER_DATASET_BUDGET  = lambda user, dataset: f"{BASE}/budgets/allocations/{user}/{dataset}"
URL_USER_BUDGETS         = lambda user: f"{BASE}/budgets/users/{user}"
URL_DATASET_BUDGETS      = lambda dataset: f"{BASE}/datasets/{dataset}"
URL_EVALUATE_F           = lambda engine: f"{BASE}/queries/evaluate?engine={engine}"
URL_EVALUATE             = f"{BASE}/queries/evaluate"
URL_ACCURACY_F           = lambda engine: f"{BASE}/queries/accuracy?engine={engine}"
URL_VALIDATE             = f"{BASE}/queries/validate"

def login_user(login_request: Dict[str, str]) -> Dict[str, str]:
    print(str(login_request))
    print(URL_LOGIN)
    resp = requests.post(url=URL_LOGIN, json=login_request)
    print(resp.status_code)
    print(resp.json())
    if resp.status_code > 299:
        return None
    token = resp.json()["token"]
    print(token)
    auth = {"Authorization" : "Bearer " + str(token)}
    return auth


job_sectors = ["Finance", "IT", "Education", "Unemployed", "Other"]

column_schema = [
    {
        "name" : "zip_code",
        "type" : {
            "name" : "Text"
        }
    },
    {
        "name" : "salary_SEK",
        "type" : {
            "name" : "Int",
            "low" : 10000,
            "high" : 100000
        }
    },
    {
        "name" : "job_sector",
        "type" : {
            "name" : "Enum",
            "labels" : job_sectors
        }
    },
    {
        "name" : "has_criminal_record",
        "type" : {
            "name" : "Bool"
        }
    },
    {
        "name" : "distance_to_closest_neighbor_km",
        "type" : {
            "name" : "Double",
            "low" : 0,
            "high" : 100
        }
    }
]


def choose_job(f: float) -> int:
    if f < 0.1:
        return 0
    elif f < 0.35:
        return 1
    elif f < 0.43:
        return 2
    elif f < 0.51:
        return 3
    else:
        return 4


def job_sal(job: str) -> int:
    if job == "Unemployed":
        return round(rd.randint(10000, 15000), 3)
    if job == "Education":
        return round(rd.randint(25000, 45000), 3)
    if job == "Finance":
        return round(rd.randint(30000, 100000), 3)
    if job == "IT":
        return round(rd.randint(25000, 95000), 3)
    else:
        return round(rd.randint(10000, 100000), 3)

def gen_csv_row():
    z = ["A", "B", "C", "D", "E", "F", "G"][rd.randint(0, 6)]
    c = rd.random() >= 0.95
    j = job_sectors[choose_job(rd.random())]
    d = round(rd.random() * 100.0, 3)
    s = job_sal(j)
    return f"{z},{s},{j},{c},{d}"

def gen_csv(rows: int):
    return "zip_code,salary_SEK,job_sector,has_criminal_record,distance_to_closest_neighbor_km\n" + \
            "\n".join(map(lambda _ : gen_csv_row(), range(0, rows)))



