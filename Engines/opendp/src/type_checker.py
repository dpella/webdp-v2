



from io import StringIO
from typing import Dict, List, Tuple
from dp_query import ColumnSchema, Mech, PrivacyNotion
from dp_types import Budget, DpType
import pandas as pd


class DpTypechecker:

    _col_schema: Dict[str, DpType]

    _df: pd.DataFrame

    _state = 0

    _t_state = 0
    _c_state = 1
    _n_state = 2
    _n_added = 3

    _bin: Tuple[str, List[any]]
    _mech: Mech

    _pn: PrivacyNotion

    def __init__(self, col_schema: List[ColumnSchema], privacy_notion: PrivacyNotion) -> None:
        self._col_schema = DpTypechecker._make_map(col_schema)
        self._pn = privacy_notion 

        csv_head  = ",".join(self._col_schema.keys())
        self._df  = pd.read_csv(StringIO(csv_head))

    
    def select(self, cols: List[str]):
        if self._state != self._t_state:
            raise DpErr(f"Bad sequence: cannot select columns: {cols}. SelectTransformaiton cannot come after aggregate functions")
        
        new_c = {}

        for c in cols:
            typ = self._col_schema.get(c)
            if typ is None:
                raise DpErr(f"Bad SelectTransformation: Column {c} from select {cols} is not in column schema: {self._col_schema}")
            new_c[c] = typ
        
        self._col_schema = new_c
        try:
            self._df = self._df[cols]
        except:
            raise DpErr(f"Bad SeletTransformation: could not apply selection of columns: {cols}")
        return self
    

    def rename(self, map: Dict[str, str]):

        if self._state != self._t_state and self._state != self._n_added:
            raise DpErr(f"Bad sequence: cannot apply rename mapping {map} after aggregation but before noise has been added")
        
        
        for old in map:
            new = map[old]
            typ = self._col_schema.get(old)
            if typ is None:
                raise DpErr(f"Bad RenameTransformation: cannot apply rename mapping: {old} -> {new}. Key: {old} is not in column schema: {self._col_schema.keys()}")
            
            self._col_schema.pop(old)
            self._col_schema[new] = typ
        

        try:
            self._df.rename(map)
        except:
            raise DpErr(f"Bad RenameTransformation: could not apply rename mappings: {map}")
        return self 
    

    def filter(self, filters: List[str]):

        if self._state != self._t_state:
            raise DpErr(f"Bad sequence: cannot apply filter: {filter}. FilterTransformaiton cannot come after aggregate functions")

        csv_head = ",".join(self._col_schema.keys())
        frame    = pd.read_csv(StringIO(csv_head))

        try:
            for filter in filters:
                frame = frame.query(filter)
        except:
            raise DpErr(f"Bad FilterTransformation: cannot apply filter: {filter}")
        return self
    
    def count(self, column: str, mech: Mech):

        if self._state not in [self._t_state, self._c_state]:
            raise DpErr(f"Bad sequence: a CountMeasurement must come first or follow a transformation")
        
        typ = self._col_schema.get(column)

        if typ is None:
            raise DpErr(f"Bad CountMeasurement: column: {column} does not exist in the column schema: {self._col_schema.keys()}")
        
        if self._state == self._c_state:
            b_col = self._bin[0]
            if b_col != column:
                raise DpErr(f"Bad CountMeasurement: cannot count column {column}. Expected column: {b_col} from the preceding BinTransformation")
            

        self._state = self._n_state
        self._mech  = mech
        return self
    
    def bin(self, column: str, bins: List[any]):

        if self._state != self._t_state:
            raise DpErr(f"Bad sequence: a BinTransformation cannot come after a measurement")

        typ = self._col_schema.get(column)

        if typ is None:
            raise DpErr(f"Bad BinTransformation: column: {column} does not exist in the column schema: {self._col_schema.keys()}")
        
        name = typ.dptype.name
        
        if not DpTypechecker._unique(bins):
            raise DpErr(f"Bad BinTransformation: expected unique elements in bin mapping: {column} -> {bins}")

        if name == "Int":

            if not len(bins) > 1:
                raise DpErr(f"Bad BinTransformation: expected an upper and a lower bound in the bin mapping, but got: {column} -> {bins}")

            if not DpTypechecker._sorted(bins):
                raise DpErr(f"Bad BinTransformation: expected an ordered list in bin mapping: {column} -> {bins}")

            if not DpTypechecker._typecheck_list(bins, int):
                raise DpErr(f"Bad BinTransformation: expected all elements of {bins} to be of type int")
        
        elif name == "Bool":

            if not DpTypechecker._typecheck_list(bins, bool):
                raise DpErr(f"Bad BinTransformation: expected all elements of {bins} to be of type bool")
            
        elif name == "Text" or name == "Enum":

            if not DpTypechecker._typecheck_list(bins, str):
                raise DpErr(f"Bad BinTransformation: expected all elements of {bins} to be of type str")
        
        else:
            raise DpErr(f"Bad BinTransformation: column type of: {column} is not supported. ")
        
        self._state = self._c_state
        self._bin   = (column, bins)
        return self
    
    def sum(self, column: str, mech: Mech):

        if self._state != self._t_state:
            raise DpErr(f"Bad SumMeasurement: a measurement has to follow 0 or more transformations - not other measurements.")
        
        typ = self._col_schema.get(column)

        if typ is None:
            raise DpErr(f"Bad SumMeasurement: column: {column} does not exist in the column schema: {self._col_schema.keys()}")
        
        name = typ.dptype.name

        if name != "Int" and name != "Double":
            raise DpErr(f"Bad SumMeasurement: column: {column} is of type: {name}. Supported types are Int and Double.")
        
        self._mech = mech
        self._state = self._n_state
        return self
    
    def noise(self, budget: Budget):

        if self._pn is PrivacyNotion.PureDP and _c(budget.delta) > 0:
            raise DpErr(f"Bad query: privacy notion is PureDP but delta value in budget was positive")
        
        elif _c(budget.delta) > 0 and self._mech is Mech.LAPLACE:
            raise DpErr(f"Bad query: the Laplace mechanism is compatible with delta = 0, but delta was {budget.delta}")
        
        elif _c(budget.delta) == 0.0 and self._mech is Mech.GAUSS:
            raise DpErr(f"Bad query: the Gauss mechanism is not compatible with zero valued delta values")
        
        self._state = self._n_added
        return self
    
    def is_valid(self):
        if self._state == self._n_added:
            return 
        else:
            raise DpErr(f"Bad Query: The minimum requirement of a query is that it contains a CountMeasurement or a SumMeasurement")
       
        

    def _typecheck_list(ls: List[any], typ):
        return all(map(lambda x: isinstance(x, typ), ls))
    
    def _sorted(ls: List[any]):
        x = ls[0]
        for l in ls:
            if l < x:
                return False 
            x = l 
        return True
    
    def _unique(ls: List[any]):
        s = set(ls)
        return len(ls) == len(s)
    
    def _make_map(ls: List[ColumnSchema]):
        s = {}
        for l in ls:
            s[l.column_name] = l.dp_type
        return s
    
 

def _c(x):
    if x is None:
        return 0.0
    else: return x
    

    
    
    

    


class DpErr(Exception):
    def __init__(self, *args: object) -> None:
        super().__init__(*args)

    
    