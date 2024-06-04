
#######################################
# USERS
#######################################

root_login = { "username" : "root", "password" : "123" }

root_patch = {
  "name": "Ruth the Root",
  # "password": "123",
  "roles": ["Admin", "Curator", "Analyst"]
}

admin_login = { "username" : "adde", "password" : "add123" }

admin = {
  "handle": "adde",
  "name": "Adrian the Admin",
  "password": "add123",
  "roles": ["Admin"]
}

curator_login = { "username" : "curt", "password" : "cur123" }

curator = {
  "handle": "curt",
  "name": "Curt the Curator",
  "password": "cur123",
  "roles": ["Curator"]
}

curator_patch = {
  "name": "Curt the Creative Curator",
  "password": "cur123",
  "roles": ["Curator"]
}

analyst_login = { "username" : "anna", "password" : "ana123" }

analyst = {
  "handle": "anna",
  "name": "Anna the Analyst",
  "password": "ana123",
  "roles": ["Analyst"]
}

analyst_patch = {
  "name": "Anna the Amazing Analyst",
  "password": "ana123",
  "roles": ["Analyst"]
}

curana_login = { "username" : "curtarne", "password" : "curry" }

curana = {
  "handle": "curtarne",
  "name": "Curt-Arne the Analyst",
  "password": "curry",
  "roles": ["Analyst"]
}

curana_patch = {
  "name": "Curt-Arne the Curator and Analyst",
  "password": "curry",
  "roles": ["Curator", "Analyst"]
}

tester_login = { "username" : "timmy", "password" : "pass" }

tester = {
  "handle": "timmy",
  "name": "Timmy the Tester",
  "password": "pass",
  "roles": ["root"] # invalid; bad request
}

#######################################
# BUDGETS
#######################################

admin_budget = {
  "epsilon": 0.2,
  # "delta": 0.3
}

curator_budget = {
  "epsilon": 0.5,
  # "delta": 0.1
}

analyst_budget = {
  "epsilon": 3,
  # "delta": 0.1
}

PureDP    = lambda epsilon        : {"epsilon": epsilon}
ApproxDP  = lambda epsilon, delta : {"epsilon": epsilon, "delta": delta}

#######################################
# DATASETS
#######################################

schema_jobs =  [
        { "name": "name",   "type": { "name": "Text" } },
        { "name": "age",    "type": { "name": "Int", "low": 18, "high": 100 } },
        { "name": "job",    "type": { "name": "Enum", "labels": ["Accountant", "Dentist", "High School Teacher", "Software Engineer"] } },
        { "name": "salary", "type": { "name": "Int", "low": 0, "high": 100000 } }
    ]

schema_jobs_not_good = [
        { "name": "name",   "type": { "name": "Text" } },
        { "name": "age",    "type": { "name": "Int", "low": 18, "high": 100 } },
        { "name": "job",    "type": { "name": "Enum", "labels": ["Accountant", "Dentist", "High School Teacher", "Software Engineer"] } },
    ]

schema_jobs_approx = [
        { "name": "name",   "type": { "name": "Text" } },
        { "name": "age",    "type": { "name": "Double", "low": 23, "high": 65 } },
        { "name": "job",    "type": { "name": "Enum", "labels": ["Accountant", "Dentist", "High School Teacher", "Software Engineer"] } },
        { "name": "salary", "type": { "name": "Double", "low": 4000, "high": 24100 } }
    ]

data_root = {
    "name": "salaries",
    "owner": "root",
    "schema": schema_jobs,
    "privacy_notion": "PureDP",
    "total_budget": PureDP(5)
}

data_root_not_good = {
    "name": "salaries",
    "owner": "root",
    "schema": schema_jobs_not_good,
    "privacy_notion": "PureDP",
    "total_budget": PureDP(5)
}

data_root_pure_notion_delta_positive = {
    "name": "salaries",
    "owner" : "root",
    "schema" : schema_jobs,
    "privacy_notion" : "PureDP",
    "total_budget" : ApproxDP(1, 0.1)
}

data_root_approx_notion_delta_null = {
    "name": "salaries",
    "owner" : "root",
    "schema" : schema_jobs,
    "privacy_notion" : "ApproxDP",
    "total_budget" : PureDP(1)
}

data_root_bad_epsilon = {
    "name": "salaries",
    "owner" : "root",
    "schema" : schema_jobs,
    "privacy_notion" : "PureDP",
    "total_budget" : {"epsilon" : -1.0}
}

data_root_bad_delta = {
    "name": "salaries",
    "owner" : "root",
    "schema" : schema_jobs,
    "privacy_notion" : "ApproxDP",
    "total_budget" : {"epsilon" : 1.0, "delta" : -0.1}
}

data_root_approx = {
    "name": "anothersalaries",
    "owner": "root",
    "schema": schema_jobs_approx,
    "privacy_notion": "ApproxDP",
    "total_budget": ApproxDP(5,0.1)
}

data_curator = {
    "name": "salaries",
    "owner": "curt",
    "schema": schema_jobs,
    "privacy_notion": "PureDP",
    "total_budget": PureDP(5)
}

data_analyst = {
    "name": "salaries",
    "owner": "anna",
    "schema": schema_jobs,
    "privacy_notion": "PureDP",
    "total_budget": PureDP(5)
}

data_patch_admin = {
    "name": "salaries",
    "owner": "adde",
    "total_budget": PureDP(6)
}

data_patch_curator = {
    "name": "salaries",
    "owner": "curt",
    "total_budget": PureDP(6)
}

data_patch_analyst = {
    "name": "salaries",
    "owner": "anna",
    "total_budget": PureDP(6)
}


