# Tests
#
# TUM01 - RENAME - SUCCESS
# TUM02 - RENAME - FAILURE
# TUM03 - FILTER - SUCCESS
# TUM04 - FILTER - FAILURE
# TUM05 - SELECT - SUCCESS
# TUM06 - SELECT - FAILURE
# TUM07 - MAP - SUCCESS
# TUM08 - MAP - FAILURE
# TUM09 - BIN - SUCCESS
# TUM10 - BIN - FAILURE
# TUM11 - GROUPBY - SUCCESS
# TUM12 - GROUPBY - FAILURE
# TUM13 - COUNT - SUCCESS
# TUM14 - COUNT - FAILURE
# TUM15 - MIN - SUCCESS
# TUM17 - MIN - FAILURE
# TUM18 - MAX - SUCCESS
# TUM19 - MAX - FAILURE
# TUM20 - SUM - SUCCESS
# TUM21 - SUM - FAILURE
# TUM22 - MEAN - SUCCESS
# TUM23 - MEAN - FAILURE


import random
from models.evaluate_request import (ColumnSchema, 
                                     DataType, 
                                     Budget, 
                                     QueryStep, 
                                     ColumnMapping, 
                                     MeasurementParams, 
                                     PrivacyNotion, 
                                     NoiseMechanism)
from models.tumult import eval_from_csv

from faker import Faker
from faker.providers import DynamicProvider

job_provider = DynamicProvider(
    provider_name="job",
    elements=["Accountant", "Dentist",
              "High School Teacher", "Software Engineer"]
)

fake = Faker()
fake.add_provider(job_provider)

FIRST_ROW = "name,age,job,salary\n"


def generate_row() -> str:
    name = fake.unique.name()
    age = random.randint(18, 100)
    job = fake.job()
    salary = random.randint(0, 100000)
    return f"{name},{age},{job},{salary}"


ROWS = 1000


def generate_test_data() -> str:
    return FIRST_ROW + '\n'.join([generate_row() for _ in range(ROWS)])


GAUSS = NoiseMechanism.GAUSS
LAPLACE = NoiseMechanism.LAPLACE

PUREDP = PrivacyNotion.PUREDP
APPROXDP = PrivacyNotion.APPROXDP

BUDGET = Budget(epsilon=1, delta=0.1)

SCHEMA = [
    ColumnSchema(name="name", type=DataType(name="Text")),
    ColumnSchema(name="age", type=DataType(name="Int", low=18, high=100)),
    ColumnSchema(name="job", type=DataType(name="Enum", labels=[
                 "Accountant", "Dentist", "High School Teacher", "Software Engineer"])),
    ColumnSchema(name="salary", type=DataType(name="Int", low=0, high=100000))
]


# def eval_from_csv(id, data, query, budget, privacy_notion, schema):

DATA = generate_test_data()

def query(q): return eval_from_csv(1, DATA, q, BUDGET, APPROXDP, SCHEMA)


def assert_query(run_query, q):
    try:
        print(run_query(q))
        return True
    except BaseException as err:
        print(err)
        return False


# TUM01 - RENAME - SUCCESS
print("RUNNING TEST: TUM01 - RENAME  - SUCCESS")
TUM01 = [QueryStep(rename={"age": "new_age"}), QueryStep(count=MeasurementParams(mech=LAPLACE))]
assert assert_query(query, TUM01)
# TUM02 - RENAME  - FAILURE
print("RUNNING TEST: TUM02 - RENAME  - FAILURE")
TUM02 = [QueryStep(rename={"no_age": "new_age"}), QueryStep(count=MeasurementParams(mech=LAPLACE))]
assert not assert_query(query, TUM02)
# TUM03 - FILTER  - SUCCESS
print("RUNNING TEST: TUM03 - FILTER  - SUCCESS")
TUM03 = [QueryStep(filter=["age >= 40", "age <= 60"]), QueryStep(count=MeasurementParams(mech=LAPLACE))]
assert assert_query(query, TUM03)
# TUM04 - FILTER  - FAILURE
print("RUNNING TEST: TUM04 - FILTER  - FAILURE")
TUM04 = [QueryStep(filter=["no_age > 40", "no_age < 60", "no_job = Dentist"]), QueryStep(count=MeasurementParams(mech=LAPLACE))]
assert not assert_query(query, TUM04)
# TUM05 - SELECT  - SUCCESS
print("RUNNING TEST: TUM05 - SELECT  - SUCCESS")
TUM05 = [QueryStep(select=["age"]), QueryStep(count=MeasurementParams(mech=LAPLACE))]
assert assert_query(query, TUM05)
# TUM06 - SELECT  - FAILURE
print("RUNNING TEST: TUM06 - SELECT  - FAILURE")
TUM06 = [QueryStep(select=["no_age"]), QueryStep(count=MeasurementParams(mech=LAPLACE))]
assert not assert_query(query, TUM06)
# TUM07 - MAP     - SUCCESS
# print("RUNNING TEST: TUM07 - MAP     - SUCCESS")
# TUM07 = []
# assert assert_query(query, TUM07)
# TUM08 - MAP     - FAILURE
# print("RUNNING TEST: TUM08 - MAP     - FAILURE")
# TUM08 = []
# assert assert_query(query, TUM08)
# TUM09 - BIN     - SUCCESS
print("RUNNING TEST: TUM09 - BIN     - SUCCESS")
TUM09 = [QueryStep(bin={"age": [10, 20, 30, 40, 50]}), QueryStep(count=MeasurementParams(mech=LAPLACE))]
assert assert_query(query, TUM09)
# TUM10 - BIN     - FAILURE
print("RUNNING TEST: TUM10 - BIN     - FAILURE")
TUM10 = [QueryStep(bin={"no_age": [10, 20, 30, 40, 50]}), QueryStep(count=MeasurementParams(mech=LAPLACE))]
assert not assert_query(query, TUM10)
# TUM11 - GROUPBY - SUCCESS
print("RUNNING TEST: TUM11 - GROUPBY - SUCCESS")
TUM11 = [QueryStep(groupby={"age": [10, 20, 30, 40, 50]}), QueryStep(count=MeasurementParams(mech=LAPLACE))]
assert assert_query(query, TUM11)
# TUM12 - GROUPBY - FAILURE
print("RUNNING TEST: TUM12 - GROUPBY - FAILURE")
TUM12 = [QueryStep(groupby={"no_age": [10, 20, 30, 40, 50]}),QueryStep(count=MeasurementParams(mech=LAPLACE))]
assert not assert_query(query, TUM12)
# TUM13 - COUNT   - SUCCESS
print("RUNNING TEST: TUM13 - COUNT   - SUCCESS")
TUM13 = [QueryStep(count=MeasurementParams(mech=LAPLACE))]
assert assert_query(query, TUM13)
# TUM14 - COUNT   - FAILURE
print("RUNNING TEST: TUM14 - COUNT   - FAILURE")
TUM14 = [QueryStep(count=MeasurementParams(column="age"))]
assert not assert_query(query, TUM14)
# TUM15 - MIN     - SUCCESS
print("RUNNING TEST: TUM15 - MIN     - SUCCESS")
TUM15 = [QueryStep(min=MeasurementParams(column="age"))]
assert assert_query(query, TUM15)
# TUM16 - MIN     - FAILURE
print("RUNNING TEST: TUM16 - MIN     - FAILURE")
TUM16 = [QueryStep(min=MeasurementParams(column="age", mech=LAPLACE))]
assert not assert_query(query, TUM16)
# TUM17 - MAX     - SUCCESS
print("RUNNING TEST: TUM17 - MAX     - SUCCESS")
TUM17 = [QueryStep(max=MeasurementParams(column="age"))]
assert assert_query(query, TUM17)
# TUM18 - MAX     - FAILURE
print("RUNNING TEST: TUM17 - MAX     - FAILURE")
TUM18 = [QueryStep(max=MeasurementParams(column="age", mech=LAPLACE))]
assert not assert_query(query, TUM18)
# TUM19 - SUM     - SUCCESS
#print("RUNNING TEST: TUM19 - SUM     - SUCCESS")
#TUM19 = [QueryStep(sum=MeasurementParams(column="salary", mech=LAPLACE))]
#assert assert_query(query, TUM19)
# TUM20 - SUM     - FAILURE
print("RUNNING TEST: TUM20 - SUM     - FAILURE")
TUM20 = [QueryStep(sum=MeasurementParams(column="name"))]
assert not assert_query(query, TUM20)
# TUM21 - MEAN    - SUCCESS
#print("RUNNING TEST: TUM21 - MEAN    - SUCCESS")
#TUM21 = [QueryStep(mean=MeasurementParams(column="age"))]
#assert assert_query(query, TUM21)
# TUM22 - MEAN    - FAILURE
TUM22 = [QueryStep(mean=MeasurementParams(column="name"))]
print("RUNNING TEST: TUM22 - MEAN    - FAILURE")
assert not assert_query(query, TUM22)
