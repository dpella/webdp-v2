import tempfile
import tmlt.analytics.utils as tmlt_utils

import pyspark.sql as pyspark_sql
import pyspark.sql.types as pyspark_sql_types

def from_csv(data, schema):
    spark = pyspark_sql.SparkSession.builder.config(conf=tmlt_utils.get_java_11_config()).getOrCreate()
    with tempfile.NamedTemporaryFile(delete=False, mode='w') as temp_file:
        temp_file.write(data)
        temp_file.flush()
        return spark.read.csv(
            temp_file.name,
            header=True,
            inferSchema=True,
        )

def to_pyspark_schema(schema):
    """
    Transform a schema into a PySpark one
    """
    return pyspark_sql_types.StructType(
        [
            pyspark_sql_types.StructField(column.name, to_pyspark_type(column.type))
            for column in schema
        ]
    )

def to_pyspark_type(datatype):
    """
    Transform a type into a PySpark one
    """
    if datatype.name == "Bool":
        return pyspark_sql_types.BooleanType()
    elif datatype.name == "Int":
        return pyspark_sql_types.IntegerType()
    elif datatype.name == "Double":
        return pyspark_sql_types.FloatType()
    elif datatype.name == "Text":
        return pyspark_sql_types.StringType()
    elif datatype.name == "Enum":
        return pyspark_sql_types.StringType()
    else:
        raise Exception("Could not transform column type: {ty}".format(ty=datatype))

'''
def from_base64(data, schema):
    rows = base64.b64decode(data).decode().split("\n")
    spark = pyspark_sql.SparkSession.builder.config(conf=tmlt_utils.get_java_11_config()).getOrCreate()
    return spark.read.csv(
        spark.sparkContext.parallelize(rows),
        header=True,
        inferSchema=False,
        schema=to_pyspark_schema(schema),
    )



url = "jdbc:postgresql://localhost:5432/tumult_test"
table = "test_table"
properties = {
    "user": "user",
    "password": "password",
    "driver": "org.postgresql.Driver"
} 

def from_postgres_table_csv(table_name, schema):
    with psycopg.connect("host=postgres port=5432 dbname=tumult_test user=user password=password") as conn:
        with conn.cursor() as cur:
            cur.execute(f"SELECT * FROM {table_name}")
            rows = cur.fetchall()

            col_names = [desc[0] for desc in cur.description]

            csv_data = StringIO()
            csv_writer = csv.writer(csv_data)
            print(col_names)
            csv_writer.writerow(col_names)
            csv_writer.writerows(rows)

            csv_data.seek(0)

            print(csv_data.getvalue())

    spark = pyspark_sql.SparkSession.builder.config(conf=tmlt_utils.get_java_11_config()).getOrCreate()
    return spark.read.csv(
        spark.sparkContext.parallelize(csv_data),
        header=True,
        inferSchema=False,
        schema=to_pyspark_schema(schema),
    )


def from_postgres_table_csv2(table_name, schema):
    with psycopg.connect("host=postgres port=5432 dbname=tumult_test user=user password=password") as conn:
        with conn.cursor() as cur:
            with cur.copy("COPY (SELECT * FROM test_table) TO STDOUT WITH CSV DELIMITER ','") as copy:
                print("here")
                file = StringIO()
                for row in copy.rows():
                    print(row)
                    file.write(row)
                file.seek(0)
                csv = file.read()
        
                print(csv)

    spark = pyspark_sql.SparkSession.builder.config(conf=tmlt_utils.get_java_11_config()).getOrCreate()
    return spark.read.csv(
        spark.sparkContext.parallelize(csv),
        header=True,
        inferSchema=False,
        schema=to_pyspark_schema(schema),
    ) 


def from_postgres_table(table_name, schema):
     spark = pyspark_sql.SparkSession.builder.config("spark.jars", "/app/postgresql.jar").getOrCreate()
     return spark.read.jdbc(
         url=url, 
         table=table, 
         properties=properties
         )


 '''