

from io import StringIO
from query_builder import *
from dp_types import *
import pandas as pd
import random as rd


_cols = "name,age,salary,job,is_cool"

_names = ["John", "Michael", "Stacy", "Olivia", "Ken"]
_jobs = ["Divorce lawyer", "Carpenter", "CIA field agent", "Doctor"]

col_schema = [
    ColumnSchema("name", DpType(TextType())),
    ColumnSchema("age", DpType(IntType(18, 90))),
    ColumnSchema("salary", DpType(DoubleType(15000.0, 200000.0))),
    ColumnSchema("job", DpType(EnumType(_jobs))),
    ColumnSchema("is_cool", DpType(BoolType()))
]

def _gen_csv(rows: int) -> str:
    return _cols + "\n" + "\n".join(map(lambda _: _gen_row(), range(0, rows)))

def _gen_row():
    n = _names[rd.randint(0, len(_names) - 1)]
    j = _jobs[rd.randint(0, len(_jobs) - 1)]
    a = rd.randint(18, 90)
    s = float(rd.randint(15000, 200000))
    b = rd.random() >= 0.5
    return ",".join(map(str, [n,a,s,j,b]))


data = pd.read_csv(StringIO(_gen_csv(100000)))

print(data)

pure = lambda: DPQueryBuilder(PrivacyNotion.PureDP, col_schema, data)
approx = lambda: DPQueryBuilder(PrivacyNotion.ApproxDP, col_schema, data)

# === PASSING ====
def test_qb_1():
    try:
        qb = DPQueryBuilder(PrivacyNotion.PureDP, col_schema, data)
        qb.apply_select(["age", "salary"])\
            .apply_rename({"age" : "ålder"})\
            .apply_filters(["ålder > 25"])\
            .apply_filters(["ålder < 65"])\
            .apply_rename({"ålder" : "the_age", "salary" : "the_salary"})\
            .make_count(column="the_age", mechanism=Mech.LAPLACE)\
            .add_noise(Budget(0.1), discrete=True).evaluate()
        
        assert True
    except Exception as e:
        print(e)
        assert 1 == 0

def test_qb_2():
    try:
        qb = DPQueryBuilder(PrivacyNotion.PureDP, column_schema=col_schema, data=data)
        qb.make_count(column="salary").add_noise(Budget(0.1), discrete=True).evaluate()
        assert True
    except Exception as e:
        print(e)
        assert 1 == 0


def test_qb_3():
    try:
        qb = DPQueryBuilder(PrivacyNotion.PureDP, col_schema, data)
        qb.make_sum(column="salary", mech=Mech.LAPLACE).add_noise(budget=Budget(0.1), discrete=False).evaluate()
        assert True
    except Exception as e:
        print(e)
        assert 1 == 0


def test_qb_4():
    try:
        qb = DPQueryBuilder(PrivacyNotion.PureDP, col_schema, data)
        qb.add_bin("age", [10, 45]).make_count("age", Mech.LAPLACE).add_noise(Budget(0.1), True).evaluate()
        assert True
    except Exception as e:
        print(e)
        assert False 

def test_qb_5():
    try:
        qb = pure()
        qb.apply_filters(["is_cool", "salary > 33000", "age < 45", "job == \"Carpenter\"", "name == \"John\""])\
            .make_count(column="name", mechanism=Mech.LAPLACE).add_noise(Budget(0.1), True).evaluate()
        assert True
    except Exception as e:
        print(e)
        assert False

def test_qb_6():
    try:
        qb = approx()
        qb.apply_select(["is_cool"]).apply_filters(["is_cool"]).apply_rename({"is_cool" : "cool_guys"})\
            .make_count("cool_guys", Mech.LAPLACE).add_noise(Budget(0.3, 0.0), True).evaluate()
        assert True
    except Exception as e:
        print(e)
        assert False 


def test_qb_7():
    try:
        approx().make_count("is_cool", Mech.GAUSS).add_noise(Budget(0.4, 0.001), True).evaluate()
        assert True 
    except Exception as e:
        print(e)
        assert False 

def test_qb_8():
    try:
        res = approx().add_bin("is_cool", ["true", "false"]).make_count("is_cool", Mech.GAUSS).add_noise(Budget(1.0,1.0), discrete=True).evaluate()
        print(res)
        assert True 
    except Exception as e:
        print(e)
        assert False 

def test_qb_9():
    try:
        res = approx().add_bin("is_cool", ["true", "false"])\
            .make_count("is_cool", Mech.GAUSS)\
            .add_noise(Budget(1.0,1.0), discrete=True)\
            .apply_rename({"is_cool" : "cool_histogram"})\
            .evaluate()
        print(res)
        assert True 
    except Exception as e:
        print(e)
        assert False 


# === invalid === 

def test_f_qb_1():
    try:
        res = approx().make_count("name", Mech.GAUSS).add_noise(Budget(0.1, 0.0), True).evaluate()
        print(res)
        assert False
    except Exception as e:
        assert True 


def test_f_qb_2():
    try:
        res = approx().make_count("name", Mech.GAUSS).apply_filters(["age > 10"]).add_noise(Budget(1.0, 0.4), True).evaluate()
        print(res)
        assert False 
    except Exception as e:
        assert True

def test_f_qb_3():
    try:
        res = approx().make_count("name", Mech.GAUSS).apply_rename({"age" : "ageee"}).add_noise(Budget(1.0, 0.4), True).evaluate()
        print(res)
        assert False 
    except Exception as e:
        assert True

def test_f_qb_4():
    try:
        res = approx().add_bin("salary", [1000.0, 25000.0]).make_count("salary", Mech.GAUSS).add_noise(Budget(0.1, 0.1), True).evaluate()
        print(res)
        assert False 
    except Exception as e:
        assert True