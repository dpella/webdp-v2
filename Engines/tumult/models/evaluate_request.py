from enum import Enum
from typing import List, Dict, Union
from pydantic import BaseModel


class DataType(BaseModel):
    name: str
    low: int | None = None
    high: int | None = None
    labels: List[str] | None = None

class ColumnSchema(BaseModel):
    name: str
    type: DataType

class Budget(BaseModel):
    epsilon: float
    delta: float | None = None

class ColumnMapping(BaseModel):
    fun: str
    schema: List[ColumnSchema]

class Value(BaseModel):
    #none
    pass

class NoiseMechanism(str, Enum):
    GAUSS = "Gauss"
    LAPLACE = "Laplace"

class MeasurementParams(BaseModel):
    column: str | None = None
    mech: NoiseMechanism | None = None
    budget: Budget | None = None

class QueryStep(BaseModel):
    select: List[str] | None = None
    rename: Dict[str, str] | None = None
    filter: List[str] | None = None
    map: ColumnMapping | None = None
    bin: Dict[str, List[Union[str, int, float, bool]]] | None = None
    count: MeasurementParams | None = None
    min: MeasurementParams | None = None
    max: MeasurementParams | None = None
    sum: MeasurementParams | None = None
    mean: MeasurementParams | None = None
    groupby: Dict[str, List[Union[str, int, float, bool]]] | None = None

class PrivacyNotion(str, Enum):
    PUREDP = "PureDP"
    APPROXDP = "ApproxDP"


class EvalRequest(BaseModel):
    budget: Budget
    query: List[QueryStep]
    privacy_notion: PrivacyNotion
    dataset: str
    schema: List[ColumnSchema]

class ErrorMessage(BaseModel):
    title: str = ""
    type: str = ""
    status: int = 200
    detail: str = ""

class QueryEvaluate200Response(BaseModel):
    rows: List[object]

class EvalRequestWithCallBack(BaseModel):
    budget: Budget
    query: List[QueryStep]
    privacy_notion: PrivacyNotion
    dataset: int
    schema: List[ColumnSchema]
    url: str