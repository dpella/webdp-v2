
#######################################
# QUERIES
#######################################

def QUERY(ds, q):
    return {
        "budget": {
            "epsilon": 0.01
        },
        "dataset": ds,
        "query": q
    }

def QUERY_approx(ds, q):
    return {
        "budget": {
            "epsilon": 0.01,
            "delta": 0.001
        },
        "dataset": ds,
        "query": q
    }

COUNT = [
        {
            "count": {
                "column": "age",
                "mech": "Laplace"
            }
        }
    ]

COUNT_approx = [
        {
            "count": {
                "column": "age",
                "mech": "Gauss"
            }
        }
    ]

SUM = [
        {
            "sum": {
                "column": "age",
                "mech": "Laplace"
            }
        }
    ]

SUM_approx = [
        {
            "sum": {
                "column": "age",
                "mech": "Gauss"
            }
        }
    ]

MEAN = [
        {
            "mean": {
                "column": "age",
                "mech": "Laplace"
            }
        }
    ]

TUM_MIN = [
        {
            "min": {
                "column": "age"
            }
        }
    ]

TUM_MAX = [
        {
            "max": {
                "column": "age"
            }
        }
    ]

FILTER_COUNT = [
        {
            "filter": ["age > 20", "age < 60"]
        },
        {
            "count": {
                "column": "age",
                "mech": "Laplace"
            }
        }
    ]

FILTER_SUM = [
        {
            "filter": ["age > 20", "age < 60"]
        },
        {
            "sum": {
                "column": "age",
                "mech": "Laplace"
            }
        }
    ]

FILTER_MEAN = [
        {
            "filter": ["age > 20", "age < 60"]
        },
        {
            "mean": {
                "column": "age",
                "mech": "Laplace"
            }
        }
    ]

GDP_BIN_COUNT = [
        {
            "bin": {
                "age": [20,30,40,50,60]
            }
        },
        {
            "count": {
                "column": "age",
                "mech": "Laplace"
            }
        }
    ]

TUM_BIN_COUNT = [
        { 
            "bin": {
                "age": [20, 30, 40, 50, 60]
            }
        },
        { 
            "groupby":  {
                "age_binned": [30, 40, 50, 60],
            }
        },
        {
            "count": {
                "mech": "Laplace"
            }
        }
    ]

GDP_BIN_SUM = [
        {
            "bin": {
                "age": [20,30,40,50,60]
            }
        },
        {
            "sum": {
                "column": "age",
                "mech": "Laplace"
            }
        }
    ]

TUM_BIN_SUM = [
        { 
            "bin": {
                "age": [20, 30, 40, 50, 60]
            }
        },
        { 
            "groupby":  {
                "age_binned": [30, 40, 50, 60],
            }
        },
        {
            "sum": {
                "column" : "age",
                "mech": "Laplace"
            }
        }
    ]

GDP_BIN_MEAN = [
        {
            "bin": {
                "age": [20,30,40,50,60]
            }
        },
        {
            "mean": {
                "column": "age",
                "mech": "Laplace"
            }
        }
    ]

TUM_BIN_MEAN = [
        { 
            "bin": {
                "age": [20, 30, 40, 50, 60]
            }
        },
        { 
            "groupby":  {
                    "age_binned": [30, 40, 50, 60],
            }
        },
        {
            "mean": {
                "column" : "age",
                "mech": "Laplace"
            }
        }
    ]

GDP_FILTER_BIN_COUNT = [
        {
            "filter": ["age > 20", "age < 60"]
        },
        {
            "bin": {
                "age": [20,30,40,50,60]
            }
        },
        {
            "count": {
                "column": "age",
                "mech": "Laplace"
            }
        }
    ]

TUM_FILTER_BIN_COUNT = [
        { 
            "filter": ["age > 20", "age < 60"]
        },
        { 
            "bin": {
                "age": [20, 30, 40, 50, 60]
            }
        },
        { 
            "groupby":  {
                "age_binned": [30, 40, 50, 60],
            }
        },
        {
            "count": {
                "mech": "Laplace"
            }
        }
    ]

GDP_FILTER_BIN_SUM = [
        {
            "filter": ["age > 20", "age < 60"]
        },
        {
            "bin": {
                "age": [20,30,40,50,60]
            }
        },
        {
            "sum": {
                "column": "age",
                "mech": "Laplace"
            }
        }
    ]

TUM_FILTER_BIN_SUM = [
        { 
            "filter": ["age > 20", "age < 60"]
        },
        { 
            "bin": {
                "age": [20, 30, 40, 50, 60]
            }
        },
        { 
            "groupby":  {
                "age_binned": [30, 40, 50, 60],
            }
        },
        {
            "sum": {
                "column": "age",
                "mech": "Laplace"
            }
        }
    ]

GDP_FILTER_BIN_MEAN = [
        {
            "filter": ["age > 20", "age < 60"]
        },
        {
            "bin": {
                "age": [20,30,40,50,60]
            }
        },
        {
            "mean": {
                "column": "age",
                "mech": "Laplace"
            }
        }
    ]

TUM_FILTER_BIN_MEAN = [
        { 
            "filter": ["age > 20", "age < 60"] 
        },
        { 
            "bin": {
                "age": [20, 30, 40, 50, 60]
            }
        },
        { 
            "groupby":  {
                "age_binned": [30, 40, 50, 60],
            }
        },
        {
            "mean": {
                "column": "age",
                "mech": "Laplace"
            }
        }
    ]

# only supported in Tumult
TUM_GBY_MEAN = [
        { "groupby": { "job": ["Accountant", "Dentist", "High School Teacher", "Software Engineer"] } },
        { "sum": { "column": "salary" } }
    ]

# only supported in Tumult
TUM_BIN_GBY_MEAN = [
        { "bin": { "age": [18, 30, 45, 60, 75] } },
        { "groupby":  {
                "age_binned": [30, 45, 60, 75],
                "job": ["Accountant", "Dentist", "High School Teacher", "Software Engineer"]
            }
        },
        { "mean": { "column": "salary" } }
    ]