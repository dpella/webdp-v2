from dataclasses import dataclass
from typing import Dict, List
from enum import Enum
from dp_types import Budget, DpType
import base64


# Query steps

@dataclass
class SelectTransformation:
    columns: List[str]

    def fromJson(self, **kwargs):
        pass

@dataclass
class FilterTransformation:
    filters: List[str]

    def fromJson(**kwargs):
        fs = kwargs.get("filter")
        if fs is None:
            raise Exception(f"could not deserialise FilterTransformation: {kwargs}")
        return FilterTransformation(filters=fs)

@dataclass 
class RenameTransformation:
    mapping: Dict[str, str]

    def fromJson(**kwargs):
        rs = kwargs.get("rename")
        if rs is None:
            raise Exception(f"could not deserialise RenameTransformation: {kwargs}")
        return RenameTransformation(mapping=rs)

@dataclass 
class ColumnSchema:
    column_name: str 
    dp_type: DpType

    def fromJson(**kwargs):
        cn = kwargs.get("name")
        if cn is None:
            raise Exception("no field \"name\" in column schema")
        
        d_type = DpType.fromJson(**kwargs.get("type"))

        if d_type is None:
            raise Exception("type is None in column schema")
        
        return ColumnSchema(column_name=cn, dp_type=d_type)
        


@dataclass
class ColumnMapping:
    func: str 
    schema: List[ColumnSchema]

    def fromJson(self, **kwargs):
        pass

@dataclass
class Mech(Enum):
    GAUSS   = "Gauss"
    LAPLACE = "Laplace"


@dataclass
class MeasurementParams:
    column: str = None
    mech: Mech = None
    budget: Budget = None

    def fromJson(**kwargs):
        column = kwargs.get("column")
        m = kwargs.get("mech")
        mech = None
        if m == Mech.LAPLACE.value:
            mech = Mech.LAPLACE
        elif m == Mech.GAUSS.value:
            mech = Mech.GAUSS
        
        budget = Budget.fromJson(**kwargs.get("budget")) if kwargs.get("budget") is not None else None
        return MeasurementParams(column=column, mech=mech, budget=budget)

    
   

@dataclass 
class MapTransformation:
    mapping: ColumnMapping
    def fromJson(self, **kwargs):
        pass


@dataclass 
class BinTransformation:
    bins: Dict[str, any]
    def fromJson(self, **kwargs):
        pass

# aggregate funcs
    
@dataclass 
class CountMeasurement:
    params: MeasurementParams
    
    def __str__(self) -> str:
        return "count"

@dataclass 
class MinMeasurement:
    params: MeasurementParams
    

@dataclass 
class MaxMeasurement:
    params: MeasurementParams

   

@dataclass 
class MeanMeasurement:
    params: MeasurementParams

    def __str__(self):
        return "mean"
    
@dataclass 
class SumMeasurement:
    params: MeasurementParams

    def __str__(self):
        return "sum"


@dataclass 
class GroupByPartition:
    grouping: Dict[str, any]

    def fromJson(self, **kwargs):
        pass


# query


_query_names = ["select", "filter", "rename", "map", "bin", "count", "min", "max", "sum", "mean", "groupby"]

@dataclass
class QueryStep:
    step: SelectTransformation | FilterTransformation | RenameTransformation | MapTransformation | BinTransformation | CountMeasurement | MinMeasurement | MaxMeasurement | MeanMeasurement | SumMeasurement | GroupByPartition

    def fromJson(kwargs: Dict[str, any]):
        q = list(filter(lambda x: kwargs.get(x) is not None, _query_names))
        if len(q) != 1:
            raise Exception("Unknown query type or combination of types")
        elif q[0] == "select":
            columns = kwargs["select"]
            return QueryStep(step = SelectTransformation(columns=columns))
            
        elif q[0] == "filter":
            fs = kwargs["filter"]
            step = FilterTransformation(filters=fs)
            return QueryStep(step=step)

        elif q[0] == "map":
            raise Exception("the map transformation is not supported")

        elif q[0] == "bin":
            bs = kwargs["bin"]
            step = BinTransformation(bins=bs)
            return QueryStep(step=step)

        elif q[0] == "count":
            mp_js = kwargs["count"]
            mp = MeasurementParams.fromJson(**mp_js)
            step = CountMeasurement(params=mp)
            return QueryStep(step=step)
            
        elif q[0] == "min":
            raise Exception("the min measurement is not supported")
            
        elif q[0] == "max":
            raise Exception("the max measurement is not supported")
        
        elif q[0] == "sum":
            mp_js = kwargs["sum"]
            mp = MeasurementParams.fromJson(**mp_js)
            step = SumMeasurement(params=mp)
            return QueryStep(step=step)
            
        elif q[0] == "mean":
            raise Exception("the mean measurement is not supported")
             
        elif q[0] == "groupby":
            raise Exception("the groupby partition is not supported")
        
        elif q[0] == "rename":
            r_mappings = kwargs["rename"]
            step = RenameTransformation(mapping=r_mappings)
            return QueryStep(step=step)
        
        else:
            raise Exception(f"could not deserialise query step: {q[0]}")

class PrivacyNotion(Enum):
    PureDP = "PureDP"
    ApproxDP = "ApproxDP"


_query_fields = ["budget", "query", "privacy_notion", "dataset", "schema"]



@dataclass
class QueryRequest:
    budget: Budget
    query: List[QueryStep]
    privacy_notion: PrivacyNotion
    data: str
    csv_header: List[str]
    schema: List[ColumnSchema]
    dataLoc: str
    datasetId: int

    def fromJson(**kwargs):
        fields = list(filter(
            lambda f: kwargs.get(f) is not None,
            _query_fields
        ))
        if len(fields) != len(_query_fields):
            raise Exception("ill formatted query")
        
        b = Budget.fromJson(**kwargs["budget"])

        qs = list(map(
            lambda q: QueryStep.fromJson(q),
            kwargs["query"]
        ))
        
        pn = kwargs["privacy_notion"]

        if pn == "PureDP":
            pn = PrivacyNotion.PureDP
        elif pn == "ApproxDP":
            pn = PrivacyNotion.ApproxDP


        schema = list(map(
            lambda c: ColumnSchema.fromJson(**c),
            kwargs["schema"]
        ))

        datasetId = int(kwargs["dataset"])

        return QueryRequest(
            budget=b,
            query=qs,
            privacy_notion=pn,
            schema=schema,
            dataLoc=kwargs["url"],
            data = "",
            csv_header= "",
            datasetId=datasetId
        )
        

@dataclass 
class AccuracyRequest:
    qr: QueryRequest
    confidence: float 

    def fromJson(**kwargs):
        qr = QueryRequest.fromJson(**kwargs)

        conf = kwargs.get('confidence')

        if conf is None:
            raise Exception("confidence field can't be null in accuracy request")
        
        if not isinstance(conf, float):
            raise Exception("confidence has wrong type. expected float")
        
        if not (conf > 0.0 and conf < 1.0):
            raise Exception("the confidence level has to be on the interval (0, 1)")
        return AccuracyRequest(qr=qr, confidence=conf)