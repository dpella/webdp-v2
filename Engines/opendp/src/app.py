


from io import StringIO
import threading
import time
import os
from typing import Dict, Tuple
from flask import Flask, jsonify, make_response, request, Response, send_from_directory, abort, render_template
import requests

from type_checker import DpErr
from query_builder import DPSyntaxException
from dp_query import AccuracyRequest, QueryRequest
from query_service import NotSupportedException, QueryService

import pandas as pd



app = Flask(__name__)

class Cache:
    _cache: Dict[int, Tuple[int, str]]
    _lock: threading.Lock
    
    def __init__(self) -> None:
        self._cache = {}
        self._lock  = threading.Lock()

    
    def delete_dataset(self, did: int):
        self._lock.acquire()
        try:
            self._cache.pop(did)
        except:
            print(f"{did} not in cache")
        finally:
            self._lock.release()

    
    def update_cache(self, did: int, url: str):
        self._lock.acquire()
        try:
            if not self._cache.__contains__(did):
                # fetch
                resp = requests.get(url)
                if resp.status_code != 200:
                    raise ServerException("Could not fetch data from WebDP")
                
                resp.encoding = 'utf-8'
                recieved_at = int(time.time())
                self._cache[did] = (recieved_at, resp.text)
        finally:
            self._lock.release()
    
    def get_csv(self, did: int) -> str:
        self._lock.acquire()
        try:
            return self._cache.get(did)[1]
        finally:
            self._lock.release()

class QHandler:

    service: QueryService
    cache: Cache

    def __init__(self, service: QueryService, cache: Cache):
        self.service = service
        self.cache   = cache
    
    def evaluate(self, req: QueryRequest):
        try:
            self.service.typecheck_query(
                query_steps=req.query,
                column_schema=req.schema,
                budget=req.budget,
                privacy_notion=req.privacy_notion
            )
            data = self.cache.get_csv(req.datasetId)
            app.logger.debug(req.query)
            
            dummy = pd.read_csv(StringIO(data))
            app.logger.debug(dummy.columns) 
            result = self.service.build_query_from_sequence(
                query_steps=req.query,
                column_schema=req.schema,
                budget=req.budget,
                privacy_notion=req.privacy_notion,
                dataset=data
            )
            return {"rows": [result.evaluate()]}, 200 
        except DPSyntaxException as e:
            return str(e), 400
        except NotSupportedException as e:
            return str(e), 400
        except DpErr as e:
            return str(e), 400
        except Exception as e:
            app.log_exception(e)
            return "unexpected error", 500
    
    def accuracy(self, req: QueryRequest, confidence: float):
        try:
            self.service.typecheck_query(
                query_steps=req.query,
                column_schema=req.schema,
                budget=req.budget,
                privacy_notion=req.privacy_notion
            )
            

            data = self.cache.get_csv(req.datasetId)
            
            result = self.service.build_query_from_sequence(
                query_steps=req.query,
                column_schema=req.schema,
                budget=req.budget,
                privacy_notion=req.privacy_notion,
                dataset=data
            )
            return jsonify([result.accuracy(confidence=confidence)]), 200
        except DPSyntaxException as e:
            return str(e), 400
        except NotSupportedException as e:
            return str(e), 400
        except DpErr as e:
            return str(e), 400
        except Exception as e:
            app.log_exception(e)
            return "unexpected error", 500
    
    def validate(self, req: QueryRequest):
        try:
            self.service.typecheck_query(
                query_steps=req.query,
                column_schema=req.schema,
                budget=req.budget,
                privacy_notion=req.privacy_notion
            )
            
            return {"valid": True, "status": "query is valid in OpenDP"}, 200
        
        except NotSupportedException as e:
            return {"valid": False, "status": str(e)}, 200
        except DpErr as e:
            return {"valid": False, "status": str(e)}, 200
        except Exception as e:
            app.log_exception(e)
            return {"valid": False, "status": "failed to validate query in OpenDP"}, 200
        
        


_cache = Cache()
_handler = QHandler(QueryService(), _cache)


@app.route("/evaluate", methods=["POST"])
def evaluate_query():
    data = request.get_json()
    try:
        qr = QueryRequest.fromJson(**data)
        _cache.update_cache(qr.datasetId, qr.dataLoc)
        return _handler.evaluate(qr)
    except Exception as e:
        return str(e), 500
    
@app.route("/accuracy", methods=["POST"])
def accuracy():
    data = request.get_json()
    try:
        ar = AccuracyRequest.fromJson(**data)
        url = ar.qr.dataLoc
        did = ar.qr.datasetId
        _cache.update_cache(did, url)
        return _handler.accuracy(req=ar.qr, confidence=ar.confidence)
    except:
        return "", 500
    
@app.route("/validate", methods = ['POST'])
def validate():
    data = request.get_json()
    try:
        qr = QueryRequest.fromJson(**data)
        return _handler.validate(qr)
    except Exception as e:
        app.log_exception(e)
        return {"valid": False, "status": str(e)}, 200

@app.route("/functions", methods=["GET"])
def get_functions():
    file_path = "./static/functions.json"
    if not os.path.exists(file_path):
        return abort(404, description="File not found")
    
    try:
        with open(file_path, "r") as file:
            content = file.read()
        return Response(content, mimetype="application/json")
    except Exception as e:
        return abort(500, description=f"Error reading file: {e}")

@app.route("/documentation", methods = ["GET"])
def get_documentation():
    file_path = "README.md"
    if not os.path.exists(file_path):
        return abort(404, description="File not found")
    
    try:
        with open(file_path, "r") as file:
            content = file.read()
        return Response(content, mimetype="text/markdown")
    except Exception as e:
        return abort(500, description=f"Error reading file: {e}")

@app.route('/cache/<int:dataset_id>', methods = ["DELETE"])
def delete_from_cache(dataset_id: int):
    _cache.delete_dataset(dataset_id)
    return "", 204

class ServerException(Exception):
    def __init__(self, *args: object) -> None:
        super().__init__(*args)

if __name__ == '__main__':
    app.run(debug=True, port=8000, host='0.0.0.0')


