


from dp_types import BoolType, DoubleType, IntType, TextType, EnumType
from dp_query import ColumnSchema
from type_checker import*




col_schema = [
    ColumnSchema("name", DpType(TextType())),
    ColumnSchema("age", DpType(IntType(18, 90))),
    ColumnSchema("salary", DpType(DoubleType(15000.0, 200000.0))),
    ColumnSchema("job", DpType(EnumType(["A", "B", "C"]))),
    ColumnSchema("is_cool", DpType(BoolType()))
]


p = lambda: DpTypechecker(col_schema, PrivacyNotion.PureDP)
ap = lambda: DpTypechecker(col_schema, PrivacyNotion.ApproxDP)

# Good queries

def test_tc_1():
    t = p()
    try: 
        t.select(["name", "age", "salary"])\
            .filter(["age > 10"])\
            .bin("age", [10, 30])\
            .count("age", Mech.LAPLACE)\
            .noise(Budget(0.2))

        t.is_valid()
        assert True
    except:
        assert False



def test_tc_2():
    t = p()
    try:
        t.count("is_cool", Mech.LAPLACE).noise(Budget(1))
        t.is_valid()
        assert True
    except:
        assert False



def test_tc_3():
    t = p()
    try:
        t.sum("age", Mech.LAPLACE).noise(Budget(1))
        t.is_valid()
        assert True
    except:
        assert False 


def test_tc_4():
    t = p()
    try:
        t.bin("is_cool", [True, False]).count("is_cool", Mech.LAPLACE).noise(Budget(0.1)).is_valid()
        assert True
    except:
        assert False

def test_tc_5():
    t = p()
    try:
        t.bin("name", ["John", "Mary"]).count("name", Mech.LAPLACE).noise(Budget(0.1)).is_valid()
        assert True
    except:
        assert False 


def test_tc_6():
    t = ap()
    try:
        t.rename({"age" : "baby"})\
            .bin("baby", [1,2,3,4,5,6,7])\
            .count("baby", Mech.LAPLACE)\
            .noise(Budget(0.1))\
            .rename({"baby" : "age"})
        t.is_valid()
        assert True
    except:
        assert False 
    
def test_tc_7():
    t = ap()
    try:
        t.rename({"age" : "baby"})\
            .bin("baby", [1,2,3,4,5,6,7])\
            .count("baby", Mech.GAUSS)\
            .noise(Budget(0.1, 0.0001))\
            .rename({"baby" : "age"})
        t.is_valid()
        assert True
    except:
        assert False 

def test_tc_8():
    t = p()
    try:
        t.rename({"salary" : "sal"}).bin("age", [10, 50]).count("age", Mech.LAPLACE).noise(Budget(1))
        t.is_valid()
        assert True
    except:
        assert False 


# bad queries 


def test_tcfail_1():
    t = p()
    try:
        t.rename({"age" : "blah"})
        t.is_valid()
        assert False 
    except:
        assert True


def test_tcfail_2():
    t = p()
    try:
        t.select(["age", "boys"])
        assert False 
    except:
        assert True


def test_tcfail_3():
    t = p()
    try:
        t.rename({"not_here" : "wont be there"})
        assert False 
    except:
        assert True 


def test_tcfail_4():
    t = ap()
    try:
        t.rename({"age" : "ok_rename"})\
            .select(["ok_rename", "salary"])\
            .bin("salary", [1000, 1000000])\
            .count("salary", Mech.GAUSS).noise(Budget(1,1))
        assert False 
    except:
        assert True 

def test_tcfail_5():
    t = p()
    try:
        t.bin("age", [100, 10])
        assert False 
    except:
        assert True 


def test_tcfail_6():
    t = p()
    try:
        t.bin("age", [10])
        assert False 
    except:
        assert True 
    

def test_tcfail_7():
    t = p()
    try:
        t.count("age", Mech.LAPLACE).select(["age"]).noise(Budget(1))
        assert False 
    except:
        assert True 


def test_tcfail_8():
    t = p()
    try:
        t.count("age", Mech.GAUSS).noise(Budget(1, 1))
        assert False 
    except:
        assert True 

def test_tcfail_9():
    t = ap()
    try:
        t.count("age", Mech.GAUSS).noise(Budget(1, 0))
        assert False 
    except:
        assert True 


def test_tc_fail_10():
    t = p()
    try:
        t.count("age", Mech.LAPLACE).noise(Budget(1, 1))
        assert False 
    except:
        assert True 


def test_tc_fail_11():
    t = p()
    try:
        t.filter(["askdfaksfd"])
        assert False 
    except:
        assert True 

def test_tc_fail_12():
    t = p()
    try:
        t.select(["age"]).filter(["salary > 0"])
        assert False 
    except:
        assert True 


def test_tc_fail_13():
    t = p() 
    try:
        t.sum("is_cool", Mech.LAPLACE).noise(Budget(1))
        assert False 
    except:
        assert True


t = p()

t.filter(["age > 50"])