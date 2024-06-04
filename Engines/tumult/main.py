import json, requests, os

from fastapi import FastAPI, Response, Request, HTTPException
from fastapi.responses import PlainTextResponse
from models.evaluate_request import EvalRequestWithCallBack, Budget
from models.tumult import status400, query_evaluate_with_data, create_tmlt_session
from dataframe.to_pyspark_session import from_csv

app = FastAPI()

app.cache = {}

@app.middleware("http")
async def log_requests(request: Request, call_next):
    print(f"Request: {request.method} {request.url}")
    respone = await call_next(request)
    return respone

@app.post("/evaluate")
def post_evaluate_v2(query: EvalRequestWithCallBack, response: Response):
    if query.dataset not in app.cache:
        try:
            resp = requests.get(query.url, timeout=30)
            if resp.status_code == 200:
                pyspark_sess = from_csv(resp.text, query.schema)
                tmlt_sess = create_tmlt_session(query.dataset, pyspark_sess, query.privacy_notion, Budget(epsilon=float('inf'), delta=float('inf')))
                app.cache[query.dataset] = tmlt_sess
            else:
                response.status_code = 400
                return {"error": "failed to retrieve data"}
        except requests.exceptions.Timeout:
            response.status_code = 500
            return {"error": "failed to retrieve data"}
    resp, status_code = query_evaluate_with_data(query, app.cache[query.dataset])
    response.status_code = status_code
    return resp

@app.post("/accuracy")
def post_accuracy(response: Response):
    resp, status_code = status400("Tumult Analytics does not support computing query accuracy")
    response.status_code = status_code
    return resp

@app.get("/functions")
def read_help():
    with open('functions.json', 'r') as helper_file:
        return json.loads(helper_file.read())
    
@app.get("/documentation", response_class=PlainTextResponse)
def get_docs():
    file_path = "README.md"
    if not os.path.exists(file_path):
        raise HTTPException(status_code=404, detail="File not found")
    try:
        with open(file_path, "r") as file:
            content = file.read()
        return PlainTextResponse(content, media_type="text/markdown")
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"error reading file: {e}")


@app.delete("/cache/{dataset_id}")
def delete_cached_dataset(dataset_id: int, response: Response):
    if dataset_id in app.cache:
        del app.cache[dataset_id]
        response.status_code = 204
        return {"status": f"dataset with id {dataset_id} deleted from cache"}
    response.status_code = 404
    return {"error": f"dataset with id {dataset_id} not found in cache"}

@app.post("/validate")
def post_validate(query: EvalRequestWithCallBack, response: Response):
    if query.dataset not in app.cache:
        try:
            resp = requests.get(query.url, timeout=30)
            if resp.status_code == 200:
                pyspark_sess = from_csv(resp.text, query.schema)
                tmlt_sess = create_tmlt_session(query.dataset, pyspark_sess, query.privacy_notion, Budget(epsilon=float('inf'), delta=float('inf')))
                app.cache[query.dataset] = tmlt_sess
            else:
                response.status_code = 400
                return {"valid": False, "status": "failed to validate query"}
        except requests.exceptions.Timeout:
            response.status_code = 500
            return {"valid": False, "status": "failed to validate query"}
    resp, status_code = query_evaluate_with_data(query, app.cache[query.dataset])
    response.status_code = status_code
    if status_code == 200:
        return {"valid": True, "status": "query is valid in tumult"}
    else:
        return {"valid": False, "status": "failed to validate query"}
