package test

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"testing"
	"webdp/internal/api/http/entity"
)

func TestQueryBudgetsInvalid(t *testing.T) {

	a := readQs("./bad_queries.json")
	for i, q := range a {
		err := q.Valid()
		if err == nil {
			t.Errorf("expected error for q: %d", i)
		}
	}
}

func TestQueryBudgetIsValid(t *testing.T) {
	a := readQs("./good_queries.json")
	for _, q := range a {
		err := q.Valid()
		if err != nil {
			t.Errorf("expected no error but got: %v", err)
		}
	}
}

func readQs(path string) []entity.QueryEvaluate {
	fs, err := os.Open(path)
	if err != nil {
		log.Fatal("Failed to open file")
	}
	defer fs.Close()

	byteValues, err := io.ReadAll(fs)

	if err != nil {
		log.Fatal("Failed to read file")
	}

	var qs []entity.QueryEvaluate

	err = json.Unmarshal(byteValues, &qs)

	if err != nil {
		log.Fatal("Failed to unmarshal file")
	}
	return qs
}
