from dataclasses import dataclass


# Data types

@dataclass
class IntType:
    low: int
    high: int
    name: str = "Int"

@dataclass 
class DoubleType:
    low: int
    high: int 
    name: str = "Double"


@dataclass
class EnumType:
    labels: list[str]
    name: str = "Enum"


@dataclass
class BoolType:
    name: str = "Bool"

@dataclass 
class TextType:
    name: str = "Text"


@dataclass
class DpType:
    dptype: IntType | DoubleType | EnumType | BoolType | TextType

    def fromJson(**kwargs) -> str:
        n = kwargs.get("name")
        if n is None:
            raise Exception("failed to deserialze data type from json")
        
        elif n == "Int":
            l, h = int(kwargs["low"]), int(kwargs["high"])
            return DpType(dptype=IntType(low=l, high=h))
           
        elif n == "Double":
            l, h = int(kwargs["low"]), int(kwargs["high"])
            return DpType(dptype=DoubleType(low=l, high=h))
            
        elif n == "Enum":
            ls = kwargs["labels"]
            return DpType(dptype=EnumType(labels=ls))
          
        elif n == "Bool":
            return DpType(dptype=BoolType())
           
        elif n == "Text":
            return DpType(dptype=TextType())
          
        else:
            raise Exception("failed to deserialize data type from json")

# Budgets

@dataclass
class Budget:

    epsilon: float 
    delta: float=None

    def fromJson(**kwargs):
        e, d = kwargs.get("epsilon"), kwargs.get("delta")
        bud = None
        if e is None and d is None:
            raise Exception("failed to deserialize Budget from json")
        
        elif d is None:
            bud = Budget(epsilon=e)
        
        else:
            bud = Budget(epsilon=e, delta=d)
        
        return bud
            

    
        


