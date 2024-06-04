




from io import StringIO
from typing import List

from type_checker import DpTypechecker
from query_builder import DPQueryBuilder
from dp_types import Budget
from dp_query import BinTransformation, ColumnSchema, CountMeasurement, FilterTransformation, MeasurementParams, PrivacyNotion, QueryStep, RenameTransformation, SelectTransformation, SumMeasurement
import pandas as pd


class QueryService:

    def build_query_from_sequence(
            self, 
            query_steps: List[QueryStep], 
            column_schema: List[ColumnSchema], 
            budget: Budget, 
            privacy_notion: PrivacyNotion, 
            dataset: str
            ) -> DPQueryBuilder:
        
        df = pd.read_csv(StringIO(dataset))

        qb = DPQueryBuilder(
            privacy_notion=privacy_notion,
            column_schema=column_schema,
            data=df
            )
        
        for step in query_steps:
            qstep = step.step
            if isinstance(qstep, FilterTransformation):
                qb.apply_filters(qstep.filters)

            elif isinstance(qstep, RenameTransformation):
                qb.apply_rename(qstep.mapping)

            elif isinstance(qstep, BinTransformation):
                self._is_valid_bin_transformation(qstep)
                for k in qstep.bins:
                    bin_list = qstep.bins[k]
                    qb.add_bin(k, bin_list)
                
            elif isinstance(qstep, CountMeasurement):
                pars = qstep.params
                self._is_valid_params(pars)
                qb.make_count(column=pars.column, mechanism=pars.mech)\
                    .add_noise(budget=budget, discrete=True)
                
            elif isinstance(qstep, SumMeasurement):
                pars = qstep.params
                self._is_valid_params(pars)
                qb.make_sum(column=pars.column, mech=pars.mech)\
                    .add_noise(budget=budget, discrete=self._get_column_type(pars.column, column_schema) == "Int")
            
            elif isinstance(qstep, SelectTransformation):
                qb.apply_select(qstep.columns)
                
            else:
                raise NotSupportedException(f"The query step: {step.step} is not supported")

        return qb
    

    def typecheck_query(
            self,
            query_steps: List[QueryStep],
            budget: Budget,
            column_schema: List[ColumnSchema],
            privacy_notion: PrivacyNotion
    ):
        
        t = DpTypechecker(col_schema=column_schema, privacy_notion=privacy_notion)
        for step in query_steps:
            qs = step.step
            if isinstance(qs, SelectTransformation):
                t.select(qs.columns)

            elif isinstance(qs, FilterTransformation):
                t.filter(qs.filters)

            elif isinstance(qs, RenameTransformation):
                t.rename(qs.mapping)

            elif isinstance(qs, BinTransformation):
                if len(qs.bins.keys()) != 1:
                    raise NotSupportedException(f"Bad BinTransformation: cannot process bin mapping: {qs.bins} -- map has to have size 1")
                else:
                    for b in qs.bins:
                        t.bin(b, qs.bins[b])

            elif isinstance(qs, SumMeasurement):
                par = qs.params
                self._is_valid_params(par)

                t.sum(column=par.column, mech=par.mech)
                t.noise(budget=budget)

            elif isinstance(qs, CountMeasurement):
                par = qs.params
                self._is_valid_params(par)

                t.count(par.column, par.mech)
                t.noise(budget=budget)
            
            else:
                raise NotSupportedException(f"Unknown query step: {qs}")

        return t.is_valid()


    # valid if the list of values is homogenously typed
    # valid if there is exactly one mapping
    def _is_valid_bin_transformation(self, bin: BinTransformation):
        if len(bin.bins) != 1:
            raise NotSupportedException(f"Your BinTransformation is badly formatted: the number of mappings has to be exactly 1")
        
        for k in bin.bins.keys():
            bin_list = bin.bins[k]
            if not (
                all(map(lambda x: isinstance(x, str), bin_list)) or 
                all(map(lambda x: isinstance(x, int), bin_list)) or 
                all(map(lambda x: isinstance(x, bool), bin_list)) or 
                all(map(lambda x: isinstance(x, float), bin_list))
                ):
                raise NotSupportedException(f"The bin values are not of the same type: {bin_list}")
            
    

    def _is_valid_params(self, m_params: MeasurementParams):
        if m_params.mech is None:
            raise NotSupportedException("Bad MeasurementParams: The noise mechanism field of an aggregate cannot be null")
        elif m_params.column is None:
            raise NotSupportedException("Bad MeasurementParams: The column field cannot be null")
        
    
    def _get_column_type(self, col: str, cs: List[ColumnSchema]) -> str:
        for s in cs:
            if s.column_name == col:
                return s.dp_type.dptype.name
        return ""


class NotSupportedException(Exception):
    def __init__(self, *args: object) -> None:
        super().__init__(*args)