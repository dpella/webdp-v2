# Engines Module

## Endpoints
Each Engine must follow the interface set below.
- POST   - /evaluate
- POST   - /accuracy
- DELETE - /cache/{id}
- POST   - /validate
- GET    - /functions
- GET    - /documentation

### Evaluate
The evaluate endpoint takes a json post request in the format: 
```
{
    "budget": {
        "epsilon": 1,
        "delta": 1 (optional)
    },
    "query": [
        querysteps here
    ],
    "dataset": id
    "schema": [
        column schema here
    ],
    "privacy_notion": "",
    "url": "" (url to where to get the data)
}
```

budget is an object
```
{
    "epsilon": 1,
    "delta": 1
}
```
where delta is optional and the values of both epsilon and delta are floats

query is an array of query steps. See x for examples of each query step. which of these to implement depends on the availability of the engine.

dataset is the id of the dataset it's an integer

schema is an array of schema which is provided from the dataset, see WebDP for information

privacy_notion is either "ApproxDP" or "PureDP" this is set in the datasetinfo

url is a callback url where the csv data can be grabbed for the dataset. Here you have a choice as an implementer, either for every query on the the dataset you can choose to either get the data every time or you can cache the data in your engine.


### Accuracy
TODO

### Cache/{id}
The cache endpoint is used to delete data that has been cached in the engine. This should be implemented if you choose to cache the data in the engine.

### Validate
Validate is used to validate if a query will run on the engine
Validate responses must conform to the below structure
If valid
```
{
    "valid": True,
    "status": "message here"
}
```
if not valid
```
{
    "valid": False,
    "status": "message here"
}
```

### Documentation
Documentation sends back a markdown page with details of the implemented functions

## Functions
Functions is a JSON object which specifies which functions are implemented in the engine
The JSON should be in the format:
```
{
    "function_1": {
        "enabled": true,
        "required_fields": {
            "field1": "what to put",
            "field2": "what to put",
            "field3": "what to put"
        }
    }
}
```

"function_1" is the name of the function, for example "select", "count", etc
"enabled" boolean if the function is available in the engine
"required_fields" object of the required fields, for example in "count" the "mechanism" is a required field  

The naming convention regarding the query features follows that of which is included in the /static folders, in the 'functions.json' file. This JSON object contains all the known query features which could be seen in an incoming request from the WebDP server. To avoid conflicts, an engine that interprets the incoming request should use this JSON object to determine if the requested query step is supported or not.

## Query Steps
### Select
```
{
    "select": [ array of columns to select ]
}
```
### Rename
```
{
    "rename": {
        "columnToRename": "RenameToThis"
    }
}
```

### Filter
```
{
    "filter": [ conditions to filter by: example "age > 18"]
}
```
### Map
```
{
    "map": {
        
    }
}
```
### Bin
```
{
    "bin": {
        "columnName": [array of bins]
    }
}
```
### GroupBy
```
{
    "groupby": {
        "columnToGroup1": [array of values to group],
        "columnToGroup2": [array of values to group],
    }
}
```
### Count
```
{
    "count": {
        "mechanism": 
    }
}
```
### Min
```
{
    "min": {
        "column"
    }
}
```
### Max
```
{
    "max": {
        "column"
    }
}
```
### Sum
```
{
    "sum": {
        "column":
        "mechanism":
    }
}
```
### Mean
```
{
    "mean": {
        "column":
        "mechanism":
    }
}
```