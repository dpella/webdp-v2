
import platform


URL = "http://localhost:8000/v2/"

URL_LOGIN               =                  URL + "login"
URL_LOGOUT              =                  URL + "logout"

URL_USERS               =                  URL + "users"
URL_USER                = lambda user:     URL_USERS + f"/{user}"

URL_DATASETS            =                  URL + "datasets"
URL_DATASET             = lambda id:       URL_DATASETS + f"/{id}"

URL_BUDGETS             =                  URL + "budgets"
URL_USER_BUDGET         = lambda user:     URL_BUDGETS + f"/users/{user}"
URL_DATASET_BUDGET      = lambda id:       URL_BUDGETS + f"/datasets/{id}"
URL_USER_DATASET_BUDGET = lambda user, id: URL_BUDGETS + f"/allocations/{user}/{id}"

URL_Q                   =                  URL + "queries"
URL_Q_ENGINES           =                  URL_Q + "/engines"
URL_Q_DOCS              =                  URL_Q + "/docs"
URL_Q_DOCS_E            = lambda engine:   URL_Q_DOCS  + f"?engine={engine}" 
URL_Q_FUNC              =                  URL_Q + "/functions"
URL_Q_FUNC_E            = lambda engine:   URL_Q_FUNC  + f"?engine={engine}" 
URL_Q_EVAL              =                  URL_Q + "/evaluate"
URL_Q_EVAL_E            = lambda engine:   URL_Q_EVAL  + f"?engine={engine}" 
URL_Q_VAL               =                  URL_Q + "/validate"
URL_Q_VAL_E             = lambda engine:   URL_Q_VAL   + f"?engine={engine}" 
URL_Q_ACC               =                  URL_Q + "/accuracy"
URL_Q_ACC_E             = lambda engine:   URL_Q_ACC   + f"?engine={engine}" 

PATH = "./tests/testdata_jobs.csv"
PATH_BAD = "./tests/testdata_jobs_bad.csv"

with open(PATH) as csv:
    data = csv.read()
FILE = data.encode()

with open(PATH_BAD) as csv:
    data_bad = csv.read()
FILE_BAD = data_bad.encode()

SUCCESS = range(200, 210)
FAIL = range(400, 510)



