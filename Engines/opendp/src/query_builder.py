

from dataclasses import dataclass
from typing import Dict, List, Tuple
import pandas as pd
from opendp.transformations import make_split_dataframe, make_select_column, then_cast, then_impute_constant, then_clamp, then_sum, make_count_by_categories, then_count
from opendp.measurements import then_base_discrete_gaussian, then_base_discrete_laplace, then_base_laplace, then_base_gaussian
from opendp.prelude import atom_domain, vector_domain, L1Distance, L2Distance, symmetric_distance, binary_search_param, gaussian_scale_to_accuracy, laplacian_scale_to_accuracy, discrete_laplacian_scale_to_accuracy, discrete_gaussian_scale_to_accuracy
from opendp.combinators import make_fix_delta, make_zCDP_to_approxDP
from opendp.mod import enable_features

from dp_query import ColumnSchema, Mech, PrivacyNotion, Budget, DpType



class DPQueryBuilder:

    """
    This class is a wrapper around the OpenDP library for making differentially private queries.
    It uses pandas to extend OpenDP with dynamc filters and rename operations on columns - a trivial 
    extension as they have the output stability of 1 and therefore does not affect the noise.

    The grammar of the builder can be represented as the regex:
        (R + F)*(BC + C + S)NR*
    where:
        R = rename
        F = filter
        B = bin transformation
        C = count measurement
        S = sum measurement
    
    Each function returns the object itself to enable streams-like chaining of operations
    """

    def __init__(self, privacy_notion: PrivacyNotion, column_schema: List[ColumnSchema], data: pd.DataFrame) -> None:
        """
            Args:
                privacy_notion (PrivacyNotion): either PureDP or ApproxDP
                column_schema (List[ColumnSchema]): the columns and types of the columns of the data
                data (pandas.DataFrame): the dataset
        """
        self._ACCEPT_FRBSC: int                          = 0 # starting state (accepts all but noise)
        self._ACCEPT_BC: int                             = 1 # accepts counts that are histograms
        self._ACCEPT_N: int                              = 2 # accept noise
        self._NOISE_ADDED: int                           = 3 # final state 
        self._compiled: bool                             = False 
        self._state: int                                 = self._ACCEPT_FRBSC

        self._privacy_notion: PrivacyNotion              = privacy_notion
        self._dataset_column_schema: List[ColumnSchema]  = column_schema
        self._data: pd.DataFrame                         = data

        self._bin_map:Tuple[str,List[any]]               = None
        self._hist_labels                                = None


        self._dp_query                                   = None
        self._col_mapping: Dict[str, DpType]             = self._make_column_map(column_schema)
        self._current_column                             = ""
        self._mech                                       = None

        self._result: DpResult                           = None

    
    def apply_select(self, columns: List[str]):

        if self._state != self._ACCEPT_FRBSC:
            raise DPSyntaxException("Bad select transformation: A select transformation cannot follow aggregate functions or bin mappings")
        
        new_col_schema = {}
        for c in columns:
            dptyp = self._col_mapping.get(c)
            if dptyp is None:
                raise DPSyntaxException(f"Bad select transformation: The column: {c} does not exist")
            new_col_schema[c] = dptyp
        
        new_data_frame = self._data[columns]
        self._data = new_data_frame
        self._col_mapping = new_col_schema
        return self
    
    
    
    def apply_filters(self, filters: List[str]):
        """
            Filters the data using the filters given as arguments
            Args:
                filters (List[str]): a list of pandas-compatible filters
            Returns:
                DPQueryBuilder
        """

        if self._state != self._ACCEPT_FRBSC:
            raise DPSyntaxException("A filter transformation cannot follow aggregate functions or bin mappings")
        
        for filter in filters:
            try:
                self._data = self._data.query(filter)
            except Exception as e:
                raise DPSyntaxException(f"Your filter: {filter} seems to be badly formated")
            
        return self
    
    def apply_rename(self, rename_mappings: Dict[str, str]):
        """
            Renames the columns according to the given map
            Args:
                rename_mappings (Dict[str, str]): map of old_name -> new_name
            Returns:
                DPQueryBuilder
        
        """

        if self._state != self._ACCEPT_FRBSC and self._state != self._NOISE_ADDED:
            raise DPSyntaxException("A rename mapping cannot follow aggregate functions or bin transformations. You can rename a result after having added noise or before aggregating")
        
        try:
            self._data = self._data.rename(columns=rename_mappings)
            for old in rename_mappings:
                if not self._col_mapping.__contains__(old):
                    raise DPSyntaxException(f"The column: {old} does not exist in the dataset. Error was found in rename mapping: {rename_mappings}")
                new_col = rename_mappings[old]
                col_type = self._col_mapping.pop(old)
                self._col_mapping[new_col] = col_type
                if self._current_column == old:
                    self._current_column = new_col

        except DPSyntaxException as dp_e:
            raise dp_e
        
        except Exception as e:
            raise DPSyntaxException(f"Your rename mapping: {rename_mappings} seems to be badly formated")
        
        return self
    
    def make_sum(self, column: str, mech:Mech = Mech.LAPLACE):
        """
            Makes a dp sum on a column
            
            Args:
                column (str): the column to sum
                mech (Mech): either Laplace or Gauss. Default is Laplace
            
            Returns: 
                DPQueryBuilder
        """

        self._mech = mech

        if self._state == self._ACCEPT_BC:
            raise DPSyntaxException("Expected a count measurement but instead got a sum measurement")
        
        if self._state == self._ACCEPT_N:
            raise DPSyntaxException("Expected noise to be added but instead go a sum measurement")
        
        if self._state == self._NOISE_ADDED:
            raise DPSyntaxException("Noise has already been added to this query. There is nothing more you can do but evaluate it.")
        
        col_type = self._col_mapping.get(column)
        if col_type is None:
            col_names = list(self._col_mapping.keys())
            raise DPSyntaxException(f"The column: {column} seems not to exist in the dataset with column names: {col_names}")
        
        enable_features("contrib")
    
        
        dp_q = (
            make_split_dataframe(separator=',', col_names=self._data.columns.to_list()) >> 
            make_select_column(key=column, TOA=str)
        )

        if col_type.dptype.name == "Int":
            self._current_column = column
            bounds = (col_type.dptype.low, col_type.dptype.high)
            dp_q = (
                dp_q >> 
                then_cast(TOA=int) >> 
                then_impute_constant(0) >> 
                then_clamp(bounds=bounds) >> 
                then_sum()
            )
            self._dp_query = dp_q
            self._state = self._ACCEPT_N
            return self
          
        elif col_type.dptype.name == "Double":
            self._current_column = column
            bounds = (float(col_type.dptype.low), float(col_type.dptype.high))
            dp_q = (
                dp_q >> 
                then_cast(TOA=float) >> 
                then_impute_constant(0.0) >> 
                then_clamp(bounds=bounds) >> 
                then_sum()
            )
            self._dp_query = dp_q
            self._state = self._ACCEPT_N
            return self
        else:
            raise DPSyntaxException(f"The column {column} has type {col_type.dptype.name} which cannot be summed")
         
    def add_bin(self, column: str, bins: List[any]):
        """
            Adds instructions to make a histogram to the sequence

            Args:
                column (str): the column to bin
                bins (List[any]): the categories to sum the column by
            
            Returns:
                DPQueryBuilder
        """
        if self._state == self._ACCEPT_BC:
            raise DPSyntaxException("You can only make one histogram (bin transformation) per query")
        if self._state == self._ACCEPT_N:
            raise DPSyntaxException("Expected noise to be added but was given a bin transformation")
        if self._state == self._NOISE_ADDED:
            raise DPSyntaxException("Noise has been added to this query: no more transformations can be made")
        # only ok state left
        if self._col_mapping.get(column).dptype.name == "Bool":
            bs = []
            for b in bins:
                if b == True:
                    bs.append("True")
                elif b == False:
                    bs.append("False")
                elif b == "true":
                    bs.append("True")
                elif b == "false":
                    bs.append("False")
            bins = bs
        self._bin_map = (column, bins)
        self._state = self._ACCEPT_BC # wait for count
        return self
    
    def make_count(self, column: str, mechanism=Mech.LAPLACE):
        """
        Makes a dp count on a column

        Args:
            column (str): the column to sum
            mech (Mech): either Laplace or Gauss. Default is Laplace
            
        Returns: 
            DPQueryBuilder
        
        """
        enable_features("contrib")
        self._dp_query = (
            make_split_dataframe(separator=',', col_names=self._data.columns.to_list()) >> 
            make_select_column(key=column, TOA=str)
        )
        self._current_column = column

        self._mech = mechanism

        if self._state == self._ACCEPT_BC:
            # make histogram
            col_type = self._col_mapping.get(column)
            if col_type is None:
                raise DPSyntaxException(f"The column: {column} does not exist in the dataset")
            
            if self._bin_map[0] != column:
                raise DPSyntaxException(f"Expected a count of the column: {self._bin_map[0]} but got: {column}")
            
            categories = None 
            if col_type.dptype.name == "Int":
                if not (all(map(lambda x: isinstance(x, int), self._bin_map[1]))):
                    raise DPTypeError(f"The values of your bins: {self._bin_map[1]} are not of type Int, which is the column type of column: {column}")
                
                categories = list(map(str, list(range(min(*self._bin_map[1]), max(*self._bin_map[1])+1))))
            elif col_type.dptype.name == "Double":
                raise DPTypeError(f"You cannot make histograms for columns of type Double")
            else:
                categories = self._bin_map[1]
            
            self._hist_labels = categories
            self._dp_query = self._dp_query >> make_count_by_categories(
                input_domain=vector_domain(atom_domain=atom_domain(T=str)),
                input_metric=symmetric_distance(),
                null_category=False,
                categories=categories,
                MO = L1Distance[int] if mechanism is Mech.LAPLACE else L2Distance[int]
            )
            

        elif self._state == self._ACCEPT_FRBSC:
            # make row count 
            self._dp_query = self._dp_query >> then_count()
        
        elif self._state == self._ACCEPT_N:
            raise DPSyntaxException("Expected noise to be added but was given a count measurement.")
        
        elif self._state == self._NOISE_ADDED:
            raise DPSyntaxException(f"Noise has been added to the query: you cannot make another count on column: {column}.")
        
        self._state = self._ACCEPT_N
        return self

    def add_noise(self, budget: Budget, discrete: bool):
        """
        Adds noise to the aggregation
        Args:
            budget (Budget): how much (epsilon, delta) to spend on the query
            discrete (bool): whether to use a discrete noise distribution or not
        Returns:
            DPQueryBuilder
        
        """
        budget = Budget(float(budget.epsilon), float(_coalesce(budget.delta)))
        mechanism = self._mech
        if mechanism is Mech.GAUSS and self._privacy_notion is PrivacyNotion.PureDP:
            raise DPSyntaxException(f"The Gauss mechanism is not compatible with {self._privacy_notion.value}")
        
        if mechanism is Mech.LAPLACE and _coalesce(budget.delta) != 0.0:
            raise DPSyntaxException(f"The Laplace mechanism is not compatible with postitive or negative values of delta")
        
        if mechanism is Mech.GAUSS and _coalesce(budget.delta) == 0.0:
            raise DPSyntaxException(f"The Gauss mechanism is not defined for delta = {_coalesce(budget.delta)}")
        
        if self._state != self._ACCEPT_N:
            raise DPSyntaxException("Noise cannot be added at this stage: you probably haven't chained an aggregate function yet")
        
        gauss = lambda scale: then_base_discrete_gaussian(scale) if discrete else then_base_gaussian(scale)
        lap   = lambda scale: then_base_discrete_laplace(scale) if discrete else then_base_laplace(scale)

        noise = None 
        if mechanism is Mech.GAUSS:
            noise = lambda scale: self._dp_query >> gauss(scale)
        elif mechanism is Mech.LAPLACE:
            noise = lambda scale: self._dp_query >> lap(scale)
        else:
            raise DPSyntaxException(f"unrecognized noise mechanism: {mechanism}")
        
        if self._privacy_notion is PrivacyNotion.ApproxDP and _coalesce(budget.delta) != 0.0:
            delta = budget.delta
            eps = budget.epsilon
            q = lambda scale: make_fix_delta(make_zCDP_to_approxDP(
                measurement=noise(scale)
            ), delta=delta)

            scale = binary_search_param(
                make_chain=q,
                d_in=1, d_out=(eps, delta)
            )
            self._result = DpResult(discrete=discrete, scale=scale, function=q(scale), column=self._current_column, mech=mechanism)
        else:
            eps = budget.epsilon
            scale = binary_search_param(
                make_chain=noise,
                d_in=1, d_out=eps
            )
            self._result = DpResult(discrete=discrete,scale=scale, function=noise(scale), column=self._current_column, mech=mechanism)
        
        self._state = self._NOISE_ADDED
        return self
    
    def evaluate(self) -> Dict[str, any]:
        if self._state != self._NOISE_ADDED:
            raise DPSyntaxException("Noise has not been added to the query: can't evaluate")
        
        f = self._result.function
        csv = self._data.to_csv(index=False, header=False)
        res = None 
        if self._hist_labels is not None:
            res = self._make_histogram(f(csv))
        else:
            res = f(csv)
        return {self._current_column: res}
    
    def accuracy(self, confidence:float):
        if self._state != self._NOISE_ADDED:
            raise DPSyntaxException("Noise has not been added to the query: can't calculate accuracy")
        scale = self._result.scale
        mech  = self._result.mech
        disc = self._result.discrete
        if disc and mech is Mech.LAPLACE:
            return discrete_laplacian_scale_to_accuracy(scale=scale,alpha=1-confidence)
        elif (not disc) and mech is Mech.LAPLACE:
            return laplacian_scale_to_accuracy(scale=scale,alpha=1-confidence)
        elif disc:
            return discrete_gaussian_scale_to_accuracy(scale=scale,alpha=1-confidence)
        else:
            return gaussian_scale_to_accuracy(scale=scale,alpha=1-confidence)

    # evaluate will throw an exception if the query is not valid
    def validate(self):
        self.evaluate() 
    
    def _make_histogram(self, results: List[int]) -> Dict[any, int]:
        typ = self._col_mapping.get(self._current_column)
        f = lambda x: int(x) if typ.dptype.name == "Int" else x
        out = {}
        for col, res in zip(self._hist_labels, results):
            out[f(col)] = res 
        return out
    
    def _make_column_map(self, schema: List[ColumnSchema]) -> Dict[str, DpType]:
        out = {}
        for cs in schema:
            out[cs.column_name] = cs.dp_type
        return out
    
   


class DPSyntaxException(Exception):
    def __init__(self, *args: object) -> None:
        super().__init__(*args)


class DPTypeError(Exception):
    def __init__(self, *args: object) -> None:
        super().__init__(*args)

@dataclass
class DpResult:
    scale: any
    function: any
    mech: Mech
    column: any 
    discrete: bool



def _coalesce(x):
    if x is None:
        return 0.0
    return x