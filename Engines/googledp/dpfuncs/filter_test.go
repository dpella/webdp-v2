package dpfuncs

import (
	"fmt"
	"googledp/entities"
	"testing"
)

// generates test data ages [0, 99]
func genTestData() [][]string {
	var testData [][]string
	for i := range 100 {
		testData = append(testData, []string{fmt.Sprintf("david%d", i), fmt.Sprintf("%d", i)})
	}
	return testData
}

func getTestSchema() []entities.Column {
	return []entities.Column{
		{
			Name: "name",
			Type: entities.StringType{
				Name: "Text",
			},
		},
		{
			Name: "age",
			Type: entities.IntType{
				Name: "Int",
				Low:  0,
				High: 100,
			},
		},
	}
}

func TestFilterData0(t *testing.T) {

	testFilter := []entities.Filter{
		{
			Column:   "age",
			Operator: "<",
			Type:     "number",
			ToFilter: "0",
		},
	}

	res, err := filterData(genTestData(), testFilter, getTestSchema())

	if err != nil {
		t.Fatalf("test failed due to error")
	}

	if len(res) != 0 {
		t.Fatalf("test failed")
	}

}

func TestFilterData1(t *testing.T) {

	testFilter := []entities.Filter{
		{
			Column:   "age",
			Operator: "<",
			Type:     "number",
			ToFilter: "1",
		},
	}

	res, err := filterData(genTestData(), testFilter, getTestSchema())

	if err != nil {
		t.Fatalf("test failed due to error")
	}

	if len(res) != 1 {
		t.Fatalf("test failed")
	}

}

func TestFilterData2(t *testing.T) {

	testFilter := []entities.Filter{
		{
			Column:   "age",
			Operator: "<",
			Type:     "number",
			ToFilter: "100",
		},
	}

	res, err := filterData(genTestData(), testFilter, getTestSchema())

	if err != nil {
		t.Fatalf("test failed due to error")
	}

	if len(res) != 100 {
		t.Fatalf("test failed")
	}

}

func TestFilterData3(t *testing.T) {

	testFilter := []entities.Filter{
		{
			Column:   "age",
			Operator: ">",
			Type:     "number",
			ToFilter: "10",
		},
		{
			Column:   "age",
			Operator: "<",
			Type:     "number",
			ToFilter: "20",
		},
	}

	res, err := filterData(genTestData(), testFilter, getTestSchema())

	if err != nil {
		t.Fatalf("test failed due to error")
	}

	if len(res) != 9 {
		t.Fatalf("test failed")
	}

}

func TestFilterData4(t *testing.T) {

	testFilter := []entities.Filter{
		{
			Column:   "age",
			Operator: "<",
			Type:     "number",
			ToFilter: "10",
		},
		{
			Column:   "age",
			Operator: ">",
			Type:     "number",
			ToFilter: "20",
		},
	}

	res, err := filterData(genTestData(), testFilter, getTestSchema())

	if err != nil {
		t.Fatalf("test failed due to error")
	}

	if len(res) != 0 {
		t.Fatalf("test failed")
	}

}

func TestFilterData5(t *testing.T) {

	testFilter := []entities.Filter{
		{
			Column:   "age",
			Operator: ">",
			Type:     "number",
			ToFilter: "100",
		},
	}

	res, err := filterData(genTestData(), testFilter, getTestSchema())

	if err != nil {
		t.Fatalf("test failed due to error")
	}

	if len(res) != 0 {
		t.Fatalf("test failed")
	}

}

func TestFilterData6(t *testing.T) {

	testFilter := []entities.Filter{
		{
			Column:   "age",
			Operator: ">",
			Type:     "number",
			ToFilter: "0",
		},
	}

	res, err := filterData(genTestData(), testFilter, getTestSchema())

	if err != nil {
		t.Fatalf("test failed due to error")
	}

	if len(res) != 99 {
		t.Fatalf("test failed")
	}

}

func TestFilterData7(t *testing.T) {

	testFilter := []entities.Filter{
		{
			Column:   "age",
			Operator: ">=",
			Type:     "number",
			ToFilter: "0",
		},
	}

	res, err := filterData(genTestData(), testFilter, getTestSchema())

	if err != nil {
		t.Fatalf("test failed due to error")
	}

	if len(res) != 100 {
		t.Fatalf("test failed")
	}

}

func TestFilterData8(t *testing.T) {

	testFilter := []entities.Filter{
		{
			Column:   "age",
			Operator: "<=",
			Type:     "number",
			ToFilter: "50",
		},
	}

	res, err := filterData(genTestData(), testFilter, getTestSchema())

	if err != nil {
		t.Fatalf("test failed due to error")
	}

	if len(res) != 51 {
		t.Fatalf("test failed")
	}

}

func TestFilterData9(t *testing.T) {

	testFilter := []entities.Filter{
		{
			Column:   "name",
			Operator: "==",
			Type:     "string",
			ToFilter: "david1",
		},
	}

	res, err := filterData(genTestData(), testFilter, getTestSchema())

	if err != nil {
		t.Fatalf("test failed due to error")
	}

	if len(res) != 1 {
		t.Fatalf("test failed")
	}

}

func TestFilterData10(t *testing.T) {

	testFilter := []entities.Filter{
		{
			Column:   "name",
			Operator: "!=",
			Type:     "string",
			ToFilter: "david1",
		},
	}

	res, err := filterData(genTestData(), testFilter, getTestSchema())

	if err != nil {
		t.Fatalf("test failed due to error")
	}

	if len(res) != 99 {
		t.Fatalf("test failed")
	}

}

func TestFilterData11(t *testing.T) {

	testFilter := []entities.Filter{
		{
			Column:   "age",
			Operator: "==",
			Type:     "number",
			ToFilter: "50",
		},
	}

	res, err := filterData(genTestData(), testFilter, getTestSchema())

	if err != nil {
		t.Fatalf("test failed due to error")
	}

	if len(res) != 1 {
		t.Fatalf("test failed")
	}

}

func TestFilterData12(t *testing.T) {

	testFilter := []entities.Filter{
		{
			Column:   "age",
			Operator: "!=",
			Type:     "number",
			ToFilter: "50",
		},
	}

	res, err := filterData(genTestData(), testFilter, getTestSchema())

	if err != nil {
		t.Fatalf("test failed due to error")
	}

	if len(res) != 99 {
		t.Fatalf("test failed")
	}

}

func TestFilterData13(t *testing.T) {

	testFilter := []entities.Filter{
		{
			Column:   "age",
			Operator: ">",
			Type:     "number",
			ToFilter: "10",
		},
		{
			Column:   "age",
			Operator: "<",
			Type:     "number",
			ToFilter: "20",
		},
		{
			Column:   "name",
			Operator: "!=",
			Type:     "string",
			ToFilter: "david14",
		},
	}

	res, err := filterData(genTestData(), testFilter, getTestSchema())

	if err != nil {
		t.Fatalf("test failed due to error")
	}

	if len(res) != 8 {
		t.Fatalf("test failed")
	}

}

func TestFilterData14(t *testing.T) {

	testFilter := []entities.Filter{
		{
			Column:   "age",
			Operator: ">",
			Type:     "number",
			ToFilter: "10",
		},
		{
			Column:   "age",
			Operator: "<",
			Type:     "number",
			ToFilter: "20",
		},
		{
			Column:   "name",
			Operator: "!=",
			Type:     "string",
			ToFilter: "david50",
		},
	}

	res, err := filterData(genTestData(), testFilter, getTestSchema())

	if err != nil {
		t.Fatalf("test failed due to error")
	}

	if len(res) != 9 {
		t.Fatalf("test failed")
	}

}

func TestFilterData15(t *testing.T) {

	testFilter := []entities.Filter{
		{
			Column:   "age",
			Operator: ">",
			Type:     "number",
			ToFilter: "10",
		},
		{
			Column:   "age",
			Operator: "<",
			Type:     "number",
			ToFilter: "20",
		},
		{
			Column:   "name",
			Operator: "==",
			Type:     "string",
			ToFilter: "david14",
		},
	}

	res, err := filterData(genTestData(), testFilter, getTestSchema())

	if err != nil {
		t.Fatalf("test failed due to error")
	}

	if len(res) != 1 {
		t.Fatalf("test failed")
	}

}

func TestFilterData16(t *testing.T) {

	testFilter := []entities.Filter{
		{
			Column:   "age",
			Operator: ">",
			Type:     "number",
			ToFilter: "10",
		},
		{
			Column:   "age",
			Operator: "<",
			Type:     "number",
			ToFilter: "20",
		},
		{
			Column:   "name",
			Operator: "==",
			Type:     "string",
			ToFilter: "david50",
		},
	}

	res, err := filterData(genTestData(), testFilter, getTestSchema())

	if err != nil {
		t.Fatalf("test failed due to error")
	}

	if len(res) != 0 {
		t.Fatalf("test failed")
	}

}

func TestFilterData17(t *testing.T) {

	testFilter := []entities.Filter{
		{
			Column:   "age",
			Operator: "<>",
			Type:     "number",
			ToFilter: "10",
		},
	}

	_, err := filterData(genTestData(), testFilter, getTestSchema())

	if err == nil {
		t.Fatalf("test failed")
	}
}

func TestFilterData18(t *testing.T) {

	testFilter := []entities.Filter{
		{
			Column:   "name",
			Operator: ">",
			Type:     "string",
			ToFilter: "10",
		},
	}

	_, err := filterData(genTestData(), testFilter, getTestSchema())

	if err == nil {
		t.Fatalf("test failed")
	}
}

func TestFilterData19(t *testing.T) {

	testFilter := []entities.Filter{
		{
			Column:   "age",
			Operator: ">",
			Type:     "number",
			ToFilter: "notnumber",
		},
	}

	_, err := filterData(genTestData(), testFilter, getTestSchema())

	if err == nil {
		t.Fatalf("test failed")
	}
}

func TestFilterData20(t *testing.T) {

	testFilter := []entities.Filter{
		{
			Column:   "notexist",
			Operator: ">",
			Type:     "number",
			ToFilter: "10",
		},
	}

	_, err := filterData(genTestData(), testFilter, getTestSchema())

	if err == nil {
		t.Fatalf("test failed")
	}
}
