package test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
	"webdp/internal/api/http/entity"
)

func TestDatasetMarshaling(t *testing.T) {
	stamp := time.Now().UTC()
	dinf := entity.DatasetInfo{
		Id:    1,
		Name:  "My data",
		Owner: "Gustav Vasa",
		Schema: []entity.ColumnSchema{
			{Name: "col1", Type: entity.DataType{Type: &entity.IntType{Low: 1, High: 10}}},
			{Name: "col2", Type: entity.DataType{Type: &entity.EnumType{Labels: []string{"a", "b"}}}},
			{Name: "col3", Type: entity.DataType{Type: &entity.DoubleType{Low: 12, High: 1234}}},
			{Name: "col4", Type: entity.DataType{Type: &entity.TextType{}}},
		},
		PrivacyNotion: "PureDP",
		TotalBudget:   entity.Budget{Epsilon: 0.1234},
		Loaded:        false,

		CreatedOn: stamp,
		UpdatedOn: stamp,
	}

	res, err := json.Marshal(dinf)
	if err != nil {
		t.Error(err)
	}

	var temp entity.DatasetInfo
	err = json.Unmarshal(res, &temp)
	if err != nil {
		t.Error(err)
	}
	res, err = json.Marshal(temp)
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("\nDatasetInfo\n%v\n\n\n", temp)
	fmt.Printf("Unmarshaling of DatasetInfo\n%v\n", string(res))

}

// X
// var temp X

func TestUnmDataType(t *testing.T) {
	js := "{\"name\":\"Enum\",\"labels\":[\"a\",\"b\"]}"

	dp, err := entity.UnmarshalDataType([]byte(js), json.Unmarshal)
	if err != nil {
		t.Error(err)
	}

	var temp entity.DataType
	err = json.Unmarshal([]byte(js), &temp)

	if err != nil {
		t.Error(err)
	}

	fmt.Printf("DP: %v\n", dp)

	fmt.Printf("DP2: %v\n", temp.Type)

	back, _ := json.Marshal(temp)

	fmt.Println(string(back))

}
