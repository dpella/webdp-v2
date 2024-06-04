package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
	errors "webdp/internal/api/http"
	"webdp/internal/api/http/entity"
)

/*
	This is a client that sends requests to the DP engines

*/

type result struct {
	Result []byte
	Error  error
}

type DPClient struct {
	m             map[string]entity.WebDPClientTarget
	timeout       time.Duration
	datasetURL    string
	DefaultEngine string
}

type ValidateResponse struct {
	Valid  bool   `json:"valid"`
	Status string `json:"status"`
}

/*
checks if default engine is set and exists in the engines array otherwise sets first engine in array to default
if timeout is set to nil, then the timeout limit will be 2 minutes
*/
func NewDPClient(enginesConfig entity.EnginesConfig, datasetUrl string, timeout *time.Duration) *DPClient {
	m := make(map[string]entity.WebDPClientTarget)

	for _, client := range enginesConfig.Engines {
		name := strings.ToLower(client.Name)
		m[name] = client
	}

	defaultEngine := strings.ToLower(enginesConfig.Default)

	if _, ok := m[defaultEngine]; !ok {
		for k := range m {
			defaultEngine = k
			break
		}
	}

	var to time.Duration
	if timeout == nil {
		to = time.Minute * 2
	} else {
		to = *timeout
	}
	return &DPClient{m: m, datasetURL: datasetUrl, timeout: to, DefaultEngine: defaultEngine}
}

func (c DPClient) RemoveDatasetFromEngineCache(dataset int64) error {
	w8 := sync.WaitGroup{}
	w8.Add(len(c.m))
	for _, targ := range c.m {
		url := fmt.Sprintf("%s/%d", targ.EndpointClearSingleCache, dataset)
		cont, cancel := context.WithTimeout(context.Background(), c.timeout)
		defer cancel()
		req, _ := http.NewRequestWithContext(cont, "DELETE", url, nil)

		go func() {
			defer w8.Done()
			cli := http.Client{}
			resp, err := cli.Do(req)
			if err != nil {
				logToFile(fmt.Sprintf("request to engine %s failed with error %s", req.URL, err.Error()))
			} else if resp.StatusCode != 204 {
				logToFile(fmt.Sprintf("the clearing of cache at %s exited with status code %d (expected is 204)", req.URL, resp.StatusCode))
			} else {
				logToFile(fmt.Sprintf("successful clearing of cache on engine %s for dataset %d", req.URL, dataset))
			}
		}()
	}

	w8.Wait()

	return nil
}

func logToFile(message string) {
	file, err := os.OpenFile("logfile.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	log.SetOutput(file)

	log.SetPrefix(time.Now().Format("2006-01-02 15:04:05"))

	log.Println(message)

}

/*
Returns the names of the engines to which the client has endpoint urls
*/
func (cl *DPClient) GetAvailableDPEngines() []string {
	out := make([]string, 0)
	for k := range cl.m {
		out = append(out, k)
	}
	return out
}

/*
gets the readme from documentation endpoint in the engine
*/
func (cl *DPClient) GetDocumentation(engine string) ([]byte, error) {
	targ, ok := cl.m[strings.ToLower(engine)]
	if !ok {
		return nil, fmt.Errorf("%w: unknown dp engine: %s", errors.ErrBadRequest, engine)
	}

	return cl.doRequestWithTimeout(targ.EndpointDocs, "GET", nil)
}

func newTimeoutContext(t time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), t)
}

func (cl DPClient) doRequestWithTimeout(url string, method string, requestBody []byte, headerKwargs ...string) ([]byte, error) {
	ctx, cancel := newTimeoutContext(cl.timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(requestBody))

	if err != nil {
		return nil, fmt.Errorf("%w: io: %s", errors.ErrUnexpected, err.Error())
	}

	if len(headerKwargs)%2 == 0 {
		for i := 0; i < len(headerKwargs); i += 2 {
			req.Header.Set(headerKwargs[i], headerKwargs[i+1])
		}
	}

	cli := http.Client{}

	response, err := cli.Do(req)

	if err != nil {
		if err == context.DeadlineExceeded {
			return nil, errors.ErrTimeout
		}
		return nil, fmt.Errorf("%w: %s", errors.ErrUnexpected, err.Error())
	}
	defer response.Body.Close()

	readResponse, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, fmt.Errorf("%w: io: %s", errors.ErrUnexpected, err.Error())
	}

	if response.StatusCode >= 500 {
		return nil, fmt.Errorf("%w: %s", errors.ErrUnexpected, string(readResponse))
	} else if response.StatusCode >= 400 {
		return nil, fmt.Errorf("%w: %s", errors.ErrBadRequest, string(readResponse))
	} else if response.StatusCode == 204 {
		return nil, nil
	} else {
		return readResponse, nil
	}
}

/*
Should return true if the client recognizes the name passed as argument
*/
func (cl *DPClient) IsAvailable(engine string) bool {
	engine = strings.ToLower(engine)
	_, ok := cl.m[engine]
	return ok
}

func (cl DPClient) makeUrl(dataset int64) string {
	s := fmt.Sprintf("%s/%d", cl.datasetURL, dataset)
	return s
}

/*
Sends a http request to the engine's evaluation endpoint. Will fail if:
  - engine doesn't exist
  - unexpected errors (like json serialisation)
  - the request takes more than `timeout` to complete
*/
func (cl *DPClient) EvaluateQuery(engine string, query entity.QueryFromClientEvaluate) (entity.QueryResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cl.timeout)
	defer cancel()
	query.CallbackUrl = cl.makeUrl(query.Data)
	var res entity.QueryResult
	js, err := json.Marshal(query)
	if err != nil {
		return res, fmt.Errorf("%w: marshal: %s", errors.ErrBadFormatting, err.Error())
	}

	engine = strings.ToLower(engine)
	cli, ok := cl.m[engine]

	if !ok {
		return res, fmt.Errorf("%w: unknown dp engine: %s", errors.ErrBadRequest, engine)
	}

	url := cli.EndpointEvaluate
	if url == "" {
		return res, fmt.Errorf("%w: engine: %s does not support evaluation of queries. there is no known endpoint of evaluating queries", errors.ErrBadRequest, cli.Name)
	}
	respChan := make(chan result, 1)

	go func() {
		defer close(respChan)
		r, err := doRequest(url, js)
		select {
		case <-ctx.Done(): // did the parent already return?
			return
		default:
			respChan <- result{Result: r, Error: err}
			return
		}
	}()

	for {
		select {
		case response := <-respChan:
			if response.Error != nil {
				return res, response.Error // TODO here
			}
			err := json.Unmarshal(response.Result, &res)
			if err != nil {
				return entity.QueryResult{}, fmt.Errorf("%w: unmarshal: %s", errors.ErrBadFormatting, err.Error())
			}
			return res, nil

		case <-ctx.Done():
			return res, fmt.Errorf("%w: query failed due to timeout", errors.ErrTimeout)
		}
	}
}

/*
Sends a Validate Query to all Engines with a Validate endpoint.
Returns a map with the result for each Engine.
*/
func (cl *DPClient) ValidateQueryAll(query entity.QueryFromClientEvaluate) (map[string]interface{}, error) {
	query.CallbackUrl = cl.makeUrl(query.Data)

	validateResult := make(map[string]interface{})
	var wg sync.WaitGroup
	wg.Add(len(cl.m))
	for engine := range cl.m {
		go func(eng string, que entity.QueryFromClientEvaluate) {
			defer wg.Done()
			resp, err := cl.ValidateQuery(eng, que)
			if err != nil {
				validateResult[eng] = ValidateResponse{Valid: false, Status: err.Error()}
				return
			}
			validateResult[eng] = resp
		}(engine, query)
	}
	wg.Wait()
	return validateResult, nil
}

/*
ValidateQuery sends a query to an engine to see if it can be executed, but doesn't apply the full data.
Will fail if:
  - engine doesn't exist
  - unexpected errors (like json serialisation)
  - the request takes more than `timeout` to complete
*/
func (cl *DPClient) ValidateQuery(engine string, query entity.QueryFromClientEvaluate) (interface{}, error) {
	// Send the engine the URL to retrieve the data
	query.CallbackUrl = cl.makeUrl(query.Data)

	// Check the format of the query to send the engine
	js, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("%w: marshal: %s", errors.ErrBadFormatting, err.Error())
	}

	engine = strings.ToLower(engine)
	cli, ok := cl.m[engine]

	if !ok {
		return nil, fmt.Errorf("%w: unknown dp engine: %s", errors.ErrBadRequest, engine)
	}

	url := cli.EndpointValidate
	if url == "" {
		return nil, fmt.Errorf("%w: engine: %s does not support validation of queries. there is no known endpoint of validating queries", errors.ErrBadRequest, cli.Name)
	}

	r, err := cl.doRequestWithTimeout(url, "POST", js, "Content-Type", "application/json")

	if err != nil {
		return nil, fmt.Errorf("%w: %s failed to validate query", errors.ErrUnexpected, engine)
	}

	var resp ValidateResponse
	err = json.Unmarshal(r, &resp)

	if err != nil {
		return nil, fmt.Errorf("%w: %s failed to validate query", errors.ErrUnexpected, engine)
	}

	return resp, nil

}

/*
Sends a http request to the engine's accuracy endpoint. Will fail if:
  - engine doesn't exist
  - unexpected errors (like json serialisation)
  - the request takes more than `timeout` to complete
  - the engine matching is case insensitive
*/
func (cl *DPClient) GetQueryAccuracy(engine string, query entity.QueryFromClientAccuracy) ([]float64, error) {

	query.CallbackUrl = cl.makeUrl(query.Data)

	// marshaling
	js, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("%w: marshal: %s", errors.ErrBadFormatting, err.Error())
	}

	// does the engine exist?
	engine = strings.ToLower(engine)
	cli, ok := cl.m[engine]
	if !ok {
		return nil, fmt.Errorf("%w: unsupported engine: %s", errors.ErrBadRequest, engine)
	}

	// get url and setup means of communicating result
	url := cli.EndpointAccuracy
	if url == "" {
		return nil, fmt.Errorf("%w: engine: %s does not support accuracy measurements of queries. there is no known endpoint of calculating accuracy", errors.ErrBadRequest, cli.Name)
	}

	resp, err := cl.doRequestWithTimeout(url, "POST", js, "Content-Type", "application/json")
	if err != nil {

		return nil, fmt.Errorf("%w: call to dp engine failed: %s", errors.ErrUnexpected, err.Error())
	}

	var res []float64

	err = json.Unmarshal(resp, &res)

	if err != nil {
		return nil, fmt.Errorf("%w: io: %s", errors.ErrUnexpected, err.Error())
	}

	return res, nil

}

func (cl *DPClient) GetSingleEngineFunctions(engine string) (interface{}, error) {
	engine = strings.ToLower(engine)
	dpClient, ok := cl.m[engine]

	if !ok {
		return nil, fmt.Errorf("unknown dp engine: %s", engine)
	}

	if dpClient.EndpointFunctions == "" {
		return nil, fmt.Errorf("engine %s has not implemented help", engine)
	}

	resp, err := http.Get(dpClient.EndpointFunctions)

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve help")
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to read response body")
	}

	var data interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return data, nil

}

func (cl *DPClient) GetAllEngineFunctions() (map[string]interface{}, error) {
	engineData := make(map[string]interface{})

	for engine, clStruct := range cl.m {
		if clStruct.EndpointFunctions != "" {
			resp, err := http.Get(clStruct.EndpointFunctions)
			if err != nil {
				return nil, err
			}

			body, err := io.ReadAll(resp.Body)

			if err != nil {
				return nil, err
			}

			var data interface{}
			if err := json.Unmarshal(body, &data); err != nil {
				return nil, err
			}

			engineData[engine] = data

		} else {
			engineData[engine] = "engine has not implemented features"
		}
	}

	return engineData, nil

}

/*
The URL must match an endpoint as defined in the DP engines configuration file.
*/
func doRequest(url string, request []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(request))
	if err != nil {
		return []byte{}, err
	}

	req.Header.Set("Content-Type", "application/json")

	cli := http.Client{}

	dpResp, err := cli.Do(req)
	if err != nil {
		return []byte{}, fmt.Errorf("%w: client: %s", errors.ErrUnexpected, err.Error())
	}

	if dpResp.StatusCode > 400 {

		resp, err := io.ReadAll(dpResp.Body)
		if err != nil {
			panic(err)
		}
		return []byte{}, fmt.Errorf("%w: computation of query failed, %s", errors.ErrBadRequest, string(resp))
	}

	ret, err := io.ReadAll(dpResp.Body)

	if err != nil {
		return []byte{}, fmt.Errorf("%w: io: %s", errors.ErrUnexpected, err.Error())
	}

	if dpResp.StatusCode >= 400 {
		return []byte{}, fmt.Errorf("%w: %s", errors.ErrUnexpected, string(ret))
	}

	defer dpResp.Body.Close()

	return ret, nil
}
