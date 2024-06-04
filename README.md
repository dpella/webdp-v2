<img src="./logos/png/logo-transparent.png" width="700">

# WebDP - Reworked

[WebDP](https://github.com/dpella/WebDP) is an open API for differential privacy frameworks built as a proof of concept by DPella. WebDP - Reworked aims to build upon the work done by the DPella team and provide a new platform with similar API specification to the original with the aim to easily be able to connect several open-source differential privacy engines.

## Project Description

This project is part of a bachelor's thesis at [Chalmers University of Technology](https://chalmers.se) and [University of Gothenburg](https://gu.se), and aims to improve [WebDP](https://github.com/dpella/WebDP), a software application serving APIs for differential privacy operations on datasets, by extending its functionality to other DP engines and frameworks. 

## Run Instructions

To run webdp first make sure docker is running on your machine then in the root folder run
```
docker-compose up
```
This will spin up the postgres database, webdp api and engines.

A user with credentials
```
username: root, password: 123
```
will be created which you can use to login, when logging in you will be sent a token which will be required to access all other endpoints

## Setup

Before opening up your instance of WebDP, it is strongly recommended to change or review *at least* the following default values:
* **.env - ROOT_PASSWORD**: Password for the root user.
* **.env - D_PASS**: Password for the database root user.
* **.env - AUTH_SIGN_KEY**: Key for signing login tokens.

## Engine configuration

* **deployment/dp-engines-config.json**: The list of active and available DP engines for which users can use.

## API Specification

There is an endpoint for viewing the API specification:

| Specification | Version 2 |
| ------------- | --------- |
| GET           | /v2/spec  |

Once docker is running, you may view the API specification wrapped in a Swagger UI by visiting

http://localhost:8000/v2/spec/index.html

The Version 1 API specification offers the same interface as the DPella Webdp.

| Type     | HTTP method | Version 1                                      | Version 2                                        |
|----------|-------------|------------------------------------------------|--------------------------------------------------|
| Auth     | POST        | /v1/login                                      | /v2/login                                        |
|          | POST        | /v1/logout                                     | /v2/logout                                       |
| Users    | GET         | /v1/users                                      | /v2/users                                        |
|          | POST        | /v1/users                                      | /v2/users                                        |
|          | GET         | /v1/user/{userHandle}                          | /v2/users/{userHandle}                           |
|          | PATCH       | /v1/user/{userHandle}                          | /v2/users/{userHandle}                           |
|          | DELETE      | /v1/user/{userHandle}                          | /v2/users/{userHandle}                           |
| Datasets | GET         | /v1/datasets                                   | /v2/datasets                                     |
|          | POST        | /v1/datasets                                   | /v2/datasets                                     |
|          | GET         | /v1/dataset/{datasetId}                        | /v2/datasets/{datasetId}                         |
|          | PATCH       | /v1/dataset/{datasetId}                        | /v2/datasets/{datasetId}                         |
|          | DELETE      | /v1/dataset/{datasetId}                        | /v2/datasets/{datasetId}                         |
|          | POST        | /v1/dataset/{datasetId}/upload                 | /v2/datasets/{datasetId}/upload                  |
| Budgets  | GET         | /v1/budget/user/{userHandle}                   | /v2/budgets/users/{userHandle}                   |
|          | GET         | /v1/budget/dataset/{datasetId}                 | /v2/budgets/datasets/{datasetId}                 |
|          | GET         | /v1/budget/allocation/{userHandle}/{datasetId} | /v2/budgets/allocations/{userHandle}/{datasetId} |
|          | POST        | /v1/budget/allocation/{userHandle}/{datasetId} | /v2/budgets/allocations/{userHandle}/{datasetId} |
|          | PATCH       | /v1/budget/allocation/{userHandle}/{datasetId} | /v2/budgets/allocations/{userHandle}/{datasetId} |
|          | DELETE      | /v1/budget/allocation/{userHandle}/{datasetId} | /v2/budgets/allocations/{userHandle}/{datasetId} |
| Queries  | POST        | /v1/query/evaluate                             | /v2/queries/evaluate?engine={engineName}         |
|          | POST        | /v1/query/accuracy                             | /v2/queries/accuracy?engine={engineName}         |
|          | POST        | /v1/query/custom                               |                                                  |
|          | POST        |                                                | /v2/queries/validate?engine={engineName}         |
|          | GET         |                                                | /v2/queries/docs?engine={engineName}             |
|          | GET         |                                                | /v2/queries/engines                              |
|          | GET         |                                                | /v2/queries/functions?engine={engineName}        |


## Key differences between v1 and v2

- As more engines will be added to webdp when doing a query you can now specify in the query parameter which engine you want to use.
- In the dataset the data curator now needs to specify a default engine that will be used unless the analyst specifies the engine in the query parameter.
- All endpoints now use plural nouns instead of mixing singular and plural.
- Three new endpoints introduced in the Queries category
  - /engines - returns a list of available engines
  - /validate - validates a query
  - /functions - if no engine parameter supplied will return which functions each engine supports. If parameter supplied will return the specified engines supported functions.
  - /docs - returns documentation markdown file

## Run API Tests

Requirements: python `requests`, `pytests` (available with `pip install`)

Enter the Webdp/ folder and run

```
pytest tests/
```
Alternatively, if you get an error that `pytest` is not recognized as a command, try the following command
```
python -m pytest tests/
```

## Demo

To run the demo you will need to install [jupyter notebook](https://jupyter.org/try).
Then navigate to the folder [demo](demo), when you are inside the demo folder run the following command.

```
jupyter notebook
```

Before running any of the demo examples ensure that the WebDP server is running by following the [run instructions](#run-instructions).

There are two demos available. The first demo is examples of different queries ran on the three available connectors. The second demo showcases a normal workflow of logging in, creating datasets, uploading data, validating queries, evaluating queries and checking budgets.

- [Demo 1](demo/demo.ipynb)
- [Demo 2](demo/presentation_demo.ipynb)

## Who we are
We are a team of students from the Department of Computer Science and Engineering at Chalmers University of Technology and University of Gothenburg. Our team consists of the following five members:
- [DAVID AL AMIRI](https://github.com/Thefriendlymoose)
- [BENJAMIN HEDE](https://github.com/hedeben)
- [ADAM NORBERG](https://github.com/Adam-Norberg)
- [SIMON PORSGAARD](https://github.com/doktorjevsky)
- SAMUEL RUNMARK THUNELL

## License

[Mozilla Public License Version 2.0](LICENSE)

