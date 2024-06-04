# Webdp

## Endpoints for the front-end and engines

A connector acts as a middle man between WebDP and the DP library. It reads and handles the incoming requests from WebDP, converting them into a format that the DP library understands.

Once the docker is running, you can get a better overview of what endpoints a connector must implement by viewing the API specification at the following address:

http://localhost:8080/v2/spec/index.html

## Query endpoints
The following endpoint(s) queries to either a specified engine, or all engines. You can find format of the requests WebDP sends to the engines and the format of the response it expects back in the demos.

* **Evaluate** - For a chosen engine, asks it to calculate a DP result given a dataset, budget, and a list of steps. 
* **Accuracy** - For a chosen engine, asks for the Accuracy and Confidence to the result of a given list of steps and a budget.
* **Validate** - For one or all engines, asks for whether they can evaluate a given list of steps.
* **Functions** - For one or all engines, returns what functionality it offers, such as supported DP functions and noise mechanisms.
* **Docs** - For one or all engines, returns their engine's documentation/README.

The following endpoint(s) queries to WebDP itself.

* **Engines** - Returns a list of the currently enabled DP engines.

# Deployment

/deployment contains the deployment files for the application, including a Dockerfile, an initiation file for the database, and the engine config file.

# Internal

/internal contains the internal packages that should not be imported by the applications outside of the project. It's organized into several subdirectories, each with a specific responsibility:

api: contains the HTTP API implementation. Notice that there is a possibility to add other API implementations, like gRPC.

config: contains the configuration package. Environment variables are loaded into config and are used from there. Parts of config are passed down to lower layers.

util: contains the utility package that consists of smaller helper functions.

## Tests

Requirements: python `requests`, `pytests` (available with `pip install`)

Run from Webdp/ folder:

```
pytest tests/
```