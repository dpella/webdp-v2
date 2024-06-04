package test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"
	"webdp/internal/api/http/client"
	"webdp/internal/api/http/entity"

	"github.com/gorilla/mux"
)

func TestEvaluate1(t *testing.T) {
	go setupMockService("12345")

	dp := entity.WebDPClientTarget{
		Name:             "opendp",
		EndpointEvaluate: "http://localhost:12345/evaluate",
		EndpointAccuracy: "http://localhost:12345/accuracy",
	}

	config := entity.EnginesConfig{
		Default: "opendp",
		Engines: []entity.WebDPClientTarget{dp},
	}

	cli := client.NewDPClient(config, "http://webdp-api:8001/datasets", nil)

	// so we have time to spin up the server
	fmt.Printf("getting some sleep ... \n")
	time.Sleep(time.Second * 10)
	fmt.Printf("resuming testing ... \n")

	// ok request

	req := entity.QueryFromClientEvaluate{
		Data:   1,
		Budget: entity.Budget{Epsilon: 0.2},
		Query:  entity.Query{QuerySteps: []entity.QueryStep{entity.SelectTransformation{Columns: []string{"my", "columns"}}}},
	}
	_, err := cli.EvaluateQuery("opendp", req)
	if err != nil {
		t.Error(err)
	}

}

func setupMockService(port string) {

	r := mux.NewRouter()

	r.HandleFunc("/evaluate", func(w http.ResponseWriter, r *http.Request) {
		var temp entity.QueryFromClientEvaluate
		err := json.NewDecoder(r.Body).Decode(&temp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		time.Sleep(time.Second)

		m := make(entity.QueryResult, 0)
		m["nice"] = "result!"
		m["of"] = "your query!"

		js, _ := json.Marshal(m)
		w.WriteHeader(200)
		w.Write(js)

	})

	r.HandleFunc("/accuracy", func(w http.ResponseWriter, r *http.Request) {
		var temp entity.QueryFromClientAccuracy
		err := json.NewDecoder(r.Body).Decode(&temp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		time.Sleep(time.Second)

		ret := []float64{0.1}
		js, _ := json.Marshal(ret)

		w.WriteHeader(200)
		w.Write(js)
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}

func TestClientTimeout(t *testing.T) {
	go setupMockService("12345")

	time.Sleep(time.Second * 7)

	dp := entity.WebDPClientTarget{
		Name:             "opendp",
		EndpointEvaluate: "http://localhost:12345/evaluate",
		EndpointAccuracy: "http://localhost:12345/accuracy",
	}

	config := entity.EnginesConfig{
		Default: "opendp",
		Engines: []entity.WebDPClientTarget{dp},
	}

	timeout := time.Millisecond * 20
	cli := client.NewDPClient(config, "http://webdp-api:8001/datasets", &timeout)

	req := entity.QueryFromClientEvaluate{
		Data:   1,
		Budget: entity.Budget{Epsilon: 0.2},
		Query:  entity.Query{QuerySteps: []entity.QueryStep{entity.SelectTransformation{Columns: []string{"my", "columns"}}}},
	}
	res, err := cli.EvaluateQuery("opendp", req)
	if err == nil {
		t.Errorf("expected timeout error but got: %v", res)
	}

	fmt.Printf("%s\n", err)

}
