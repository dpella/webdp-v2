import json, copy

from typing import List, Union

from models.evaluate_request import (
    QueryEvaluate200Response,
    PrivacyNotion,
    ColumnSchema,
    QueryStep,
    NoiseMechanism,
    ErrorMessage,
    EvalRequestWithCallBack)
from simpleeval import EvalWithCompoundTypes

import tmlt.analytics.session as tmlt_session
import tmlt.analytics.privacy_budget as tmlt_privacy_budget
import tmlt.analytics.protected_change as tmlt_protected_change
import tmlt.analytics.query_builder as tmlt_query_builder
import tmlt.analytics.query_expr as tmlt_query_expr
import tmlt.analytics.binning_spec as tmlt_binning_spec
import tmlt.analytics.keyset as tmlt_keyset
import tmlt.analytics._schema as tmlt_schema

import pyspark.sql.types as pyspark_sql_types

def query_evaluate_with_data(req: EvalRequestWithCallBack, tmlt_sess):
    """
    Evaluate a composed differentially-private query over an specified
    dataset consuming the indicated budget
    """

    # Unpack request
    query = req.query
    budget = req.budget
    privacy_notion = req.privacy_notion
    dataset_id = req.dataset
    schema = req.schema

    # Budget validation
    try:
        # Try to evaluate the query
        df = evaluateTumultQuery(dataset_id, tmlt_sess, query, budget, privacy_notion, schema)

    except (
        # Exceptions that TumultAnalytics throws on invalid queries
        QueryError,
        ValueError,
        KeyError,
        AttributeError,
        IndexError,
        RuntimeError,
    ) as error:
        return status400("Invalid query: {error}".format(error=error))
    except (
        # Any other exception is considered a 'crash'
        Exception
    ) as error:
        return status400(
            "Tumult Analytics crashed unexpectedly {ty}:\n{error}".format(
                ty=type(error).__name__, error=error
            )
        )

    # Transform the result dataframe into a JSON array of objects
    rows = []
    for row in df.toJSON().collect():
        rows.append(json.loads(row))

    # Prepare the response
    response = QueryEvaluate200Response(rows=rows)

    return status200(response)

def evaluateTumultQuery(id, tmlt_sess, query, budget, privacy_notion, schema):
    """
    Evaluate a query inside of its corresponding Spark session
    """

    # Evaluate the query inside the session
    res = tmlt_sess.evaluate(
        query_expr=to_tmlt_query(id, schema, query),
        privacy_budget=to_tmlt_budget(privacy_notion, budget),
    )

    return res

def create_tmlt_session(dataset_id, pyspark_session, privacy_notion, budget):

    return tmlt_session.Session.from_dataframe(
        source_id=to_tmlt_source_id(dataset_id),
        dataframe=pyspark_session,
        privacy_budget=to_tmlt_budget(privacy_notion, budget),
        protected_change=tmlt_protected_change.AddOneRow(),
    )

def to_tmlt_budget(privacy_notion, budget):
    """
    Transform a budget into a Tumult Analytics one
    """
    if privacy_notion == PrivacyNotion.PUREDP:
        return tmlt_privacy_budget.PureDPBudget(budget.epsilon)
    elif privacy_notion == PrivacyNotion.APPROXDP:
        return tmlt_privacy_budget.ApproxDPBudget(budget.epsilon, budget.delta)
    else:
        raise Exception("Could not transform budget: {budget}".format(budget=budget))

def dict_key(dict):
    """
    Return the name of the first key of a dictionary that is not associated with
    a None value.
    """
    for key in dict.keys():
        if dict[key] is not None:
            return key

    raise Exception("Input dictionary only contains None values: {dict}".format(dict=dict))

def to_tmlt_query(dataset_id, schema, query):
    """
    Transform a query into a Tumult Analytics one
    """

    # Start with an empty QueryBuilder
    tmlt_query = tmlt_query_builder.QueryBuilder(to_tmlt_source_id(dataset_id))

    # Create a copy of the schema that might get augmented during the query
    local_schema = copy.deepcopy(schema)

    # Add each query step sequentially
    for step in query:
        # Check that an intermediate step didn't run a measurement
        if isinstance(tmlt_query, tmlt_query_expr.QueryExpr):
            raise QueryError("measurements are not allowed as intermediate steps")

        # Process the query step
        local_schema, tmlt_query = add_tmlt_query_step(local_schema, tmlt_query, step)

    # Check that we return a QueryExpr
    if not isinstance(tmlt_query, tmlt_query_expr.QueryExpr):
        raise QueryError("the last step of the query is not a measurement")

    return tmlt_query


def add_tmlt_query_step(
    schema: List[ColumnSchema],
    query: Union[
        tmlt_query_builder.QueryBuilder,
        tmlt_query_builder.GroupedQueryBuilder,
        tmlt_query_expr.QueryExpr,
    ],
    step: QueryStep,
):
    """
    Update a Tumult Analytics query with a given query step.

    NOTE: this function also returns an updated schema in case the given query
    step transforms it in any way.
    """

    # Extract the step operation
    operation = dict_key(step.model_dump())

    # =================================
    # SELECT
    if operation == "select":
        kwargs = {}
        columns = step.select

        schema = filter_schema(columns, schema)

        kwargs["columns"] = columns
        return schema, query.select(**kwargs)

    # =================================
    # RENAME
    elif operation == "rename":
        kwargs = {}
        renaming = step.rename

        for col_def in schema:
            if col_def.name in renaming:
                col_def.name = renaming[col_def.name]

        kwargs["column_mapper"] = renaming
        return schema, query.rename(**kwargs)

    # =================================
    # FILTER
    elif operation == "filter":
        filters = step.filter

        for row_filter in filters:
            kwargs = {}
            kwargs["condition"] = row_filter
            query = query.filter(**kwargs)

        return schema, query

    # =================================
    # MAP
    elif operation == "map":
        kwargs = {}
        fun = step.map.fun
        new_schema = step.map.schema

        def mapper_fun(row):
            row_names = {col_def.name: row[col_def.name] for col_def in schema}
            return EvalWithCompoundTypes(names=row_names).eval(fun)

        kwargs["f"] = mapper_fun
        kwargs["new_column_types"] = to_tmlt_schema(new_schema)
        return new_schema, query.map(**kwargs)

    # =================================
    # BIN
    elif operation == "bin":
        binning = step.bin

        for col, bins in binning.items():
            kwargs = {}
            col_def = find_column_def(col, schema)

            new_col = ColumnSchema(name="{col}_binned".format(col=col), type=col_def.type)
            schema.append(new_col)

            kwargs["column"] = col
            kwargs["spec"] = tmlt_binning_spec.BinningSpec(bins, names=bins[1:])
            kwargs["name"] = new_col.name
            query = query.bin_column(**kwargs)

        return schema, query

    # =================================
    # COUNT
    elif operation == "count":
        kwargs = {}
        params = step.count

        if params.budget is not None:
            raise QueryError("Tumult Analytics does not support specifying count budget")

        if params.mech is not None:
            kwargs["mechanism"] = to_tmlt_count_mechanism(params.mech)

        kwargs["name"] = "count"
        return schema, query.count(**kwargs)

    # =================================
    # MIN
    elif operation == "min":
        kwargs = {}
        params = step.min

        if params.budget is not None:
            raise QueryError("Tumult Analytics does not support specifying min budget")

        if params.mech is not None:
            raise QueryError("Tumult Analytics does not support specifying min mechanism")

        if params.column is None:
            raise QueryError("min measurment requires `column` parameter")

        col_def = find_column_def(params.column, schema)

        if not numeric_type(col_def.type):
            raise QueryError("min measurement only supports numeric types")

        kwargs["column"] = params.column
        kwargs["low"] = col_def.type.low
        kwargs["high"] = col_def.type.high
        kwargs["name"] = "{col}_min".format(col=params.column)
        return schema, query.min(**kwargs)

    # =================================
    # MAX
    elif operation == "max":
        kwargs = {}
        params = step.max

        if params.budget is not None:
            raise QueryError("Tumult Analytics does not support specifying max budget")

        if params.mech is not None:
            raise QueryError("Tumult Analytics does not support specifying max mechanism")

        if params.column is None:
            raise QueryError("max measurment requires `column` parameter")

        col_def = find_column_def(params.column, schema)

        if not numeric_type(col_def.type):
            raise QueryError("max measurement only supports numeric types")

        kwargs["column"] = params.column
        kwargs["low"] = col_def.type.low
        kwargs["high"] = col_def.type.high
        kwargs["name"] = "{col}_max".format(col=params.column)
        return schema, query.max(**kwargs)

    # =================================
    # SUM
    elif operation == "sum":
        kwargs = {}
        params = step.sum

        if params.budget is not None:
            raise QueryError("Tumult Analytics does not support specifying sum budget")

        if params.mech is not None:
            kwargs["mechanism"] = to_tmlt_sum_mechanism(params.mech)

        if params.column is None:
            raise QueryError("sum measurement requires `column` parameter")

        col_def = find_column_def(params.column, schema)

        if not numeric_type(col_def.type):
            raise QueryError("sum measurement only supports numeric types")

        kwargs["column"] = params.column
        kwargs["low"] = col_def.type.low
        kwargs["high"] = col_def.type.high
        kwargs["name"] = "{col}_sum".format(col=params.column)
        return schema, query.sum(**kwargs)

    # =================================
    # MEAN
    elif operation == "mean":
        kwargs = {}
        params = step.mean

        if params.budget is not None:
            raise QueryError("Tumult Analytics does not support specifying mean budget")

        if params.mech is not None:
            kwargs["mechanism"] = to_tmlt_mean_mechanism(params.mech)

        if params.column is None:
            raise QueryError("mean measurment requires `column` parameter")

        col_def = find_column_def(params.column, schema)

        if not numeric_type(col_def.type):
            raise QueryError("mean measurement only supports numeric types")

        kwargs["column"] = params.column
        kwargs["low"] = col_def.type.low
        kwargs["high"] = col_def.type.high
        kwargs["name"] = "{name}_mean".format(name=params.column)
        return schema, query.average(**kwargs)

    # =================================
    # GROUPBY
    elif operation == "groupby":
        kwargs = {}
        keys = step.groupby

        kwargs["keys"] = tmlt_keyset.KeySet.from_dict(keys)
        return schema, query.groupby(**kwargs)


def to_tmlt_count_mechanism(mech):
    """
    Transform a noise mechanism to a Tumult Analytics count mechanism
    """
    if mech == NoiseMechanism.GAUSS:
        return tmlt_query_expr.CountMechanism.GAUSSIAN
    elif mech == NoiseMechanism.LAPLACE:
        return tmlt_query_expr.CountMechanism.LAPLACE
    else:
        raise Exception("Invalid count noise mechanism: {mech}".format(mech=mech))


def to_tmlt_sum_mechanism(mech):
    """
    Transform a noise mechanism to a Tumult Analytics sum mechanism
    """
    if mech == NoiseMechanism.GAUSS:
        return tmlt_query_expr.SumMechanism.GAUSSIAN
    elif mech == NoiseMechanism.LAPLACE:
        return tmlt_query_expr.SumMechanism.LAPLACE
    else:
        raise Exception("Invalid sum noise mechanism: {mech}".format(mech=mech))


def to_tmlt_mean_mechanism(mech):
    """
    Transform a noise mechanism to a Tumult Analytics mean mechanism
    """
    if mech == NoiseMechanism.GAUSS:
        return tmlt_query_expr.AverageMechanism.GAUSSIAN
    elif mech == NoiseMechanism.LAPLACE:
        return tmlt_query_expr.AverageMechanism.LAPLACE
    else:
        raise Exception("Invalid mean noise mechanism: {mech}".format(mech=mech))

def to_tmlt_type(datatype):
    """
    Transform a type into a Tumult Analytics one

    NOTE: Tumult Analytics does not support `Bool`!!!
    """
    if datatype.name == "Int":
        return tmlt_schema.ColumnType.INTEGER
    elif datatype.name == "Double":
        return tmlt_schema.ColumnType.DECIMAL
    elif datatype.name == "Text":
        return tmlt_schema.ColumnType.VARCHAR
    elif datatype.name == "Enum":
        return tmlt_schema.ColumnType.VARCHAR
    else:
        raise Exception("Could not transform column type: {ty}".format(ty=datatype))

def to_tmlt_schema(schema):
    """
    Transform a schema into a Tumult Analytics one
    """
    return {column.name: to_tmlt_type(column.type) for column in schema}

def find_column_def(column: str, schema: List[ColumnSchema]):
    """
    Find a column in a given schema, raising a QueryException if the column does not exist
    """
    col_def = next(filter(lambda col_def: col_def.name == column, schema), None)

    if col_def is None:
        raise QueryError("column {col} does not exist in the dataset schema".format(col=column))

    return col_def


def filter_schema(columns: List[str], schema: List[ColumnSchema]):
    """
    Filter the columns of a given schema, raising a QueryException if a column does not exist
    """
    new_schema = []
    for column in columns:
        new_schema.append(find_column_def(column, schema))
    return new_schema

def numeric_type(type):
    """
    Check that a type is numeric
    """
    return type.name in ["Int", "Double"]

class QueryError(BaseException):
    """
    Raised when a request passes an invalid query
    """

    pass

def to_tmlt_source_id(dataset_id):
    """
    Assign a string name to a dataset id
    """
    return "dataset_{num}".format(num=dataset_id)

def status200(body=None):
    """
    return HTTP 200 (Ok)
    """
    return body, 200

# 4xx return codes


def status400(error_msg):
    """
    return HTTP 400 (Bad request)
    """
    return ErrorMessage(title="Bad request", detail=error_msg, status=400), 400


'''
Not used functions

def eval_from_csv(id, data, query, budget, privacy_notion, schema):
    """
    Evaluate a query inside of its corresponding Spark session
    """

    # Retrieve the dataset Tumult Analytics session
    session = create_session_from_csv(id, data, schema, privacy_notion, budget)

    # Evaluate the query inside the session
    res = session.evaluate(
        query_expr=to_tmlt_query(id, schema, query),
        privacy_budget=to_tmlt_budget(privacy_notion, budget),
    )

    return res



def query_evaluate(query_evaluate_request: EvalRequest):
    """
    Evaluate a composed differentially-private query over an specified
    dataset consuming the indicated budget
    """

    # Unpack request
    data = query_evaluate_request.dataset
    query = query_evaluate_request.query
    budget = query_evaluate_request.budget
    privacy_notion = query_evaluate_request.privacy_notion
    dataset_id = "tempidhere"
    schema = query_evaluate_request.schema

    # Budget validation
    try:
        # Try to evaluate the query
        df = evaluateTumultQuery(dataset_id, data, query, budget, privacy_notion, schema)

    except (
        # Exceptions that TumultAnalytics throws on invalid queries
        QueryError,
        ValueError,
        KeyError,
        AttributeError,
        IndexError,
        RuntimeError,
    ) as error:
        return status400("Invalid query: {error}".format(error=error))
    except (
        # Any other exception is considered a 'crash'
        Exception
    ) as error:
        return status400(
            "Tumult Analytics crashed unexpectedly {ty}:\n{error}".format(
                ty=type(error).__name__, error=error
            )
        )

    # Transform the result dataframe into a JSON array of objects
    rows = []
    for row in df.toJSON().collect():
        rows.append(json.loads(row))

    # Prepare the response
    response = QueryEvaluate200Response(rows=rows)

    # Log the request
    print("[WebDP] query_evaluate:")
    print(data, [dict_key(step.model_dump()) for step in query], budget)

    return status200(response)


def create_session(dataset_id, data, schema, privacy_notion, budget):
    df = from_base64(data, schema)

    return tmlt_session.Session.from_dataframe(
        source_id=to_tmlt_source_id(dataset_id),
        dataframe=df,
        privacy_budget=to_tmlt_budget(privacy_notion, budget),
        protected_change=tmlt_protected_change.AddOneRow(),
    )





'''