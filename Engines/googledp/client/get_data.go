package client

import (
	"encoding/csv"
	"net/http"
)

func GetCSVData(url string) ([][]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	r := csv.NewReader(resp.Body)

	records, err := r.ReadAll()

	if err != nil {
		return nil, err
	}

	return records, nil

}
