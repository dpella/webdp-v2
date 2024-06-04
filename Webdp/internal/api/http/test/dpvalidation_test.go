package test

import (
	"testing"
	"webdp/internal/api/http/entity"
	"webdp/internal/api/http/utils"
)

type valid interface {
	Valid() error
}

func TestDpValidation(t *testing.T) {

	type temp struct {
		Field1 string `dpvalidation:"non-empty-string"`
		Field2 int
	}

	v := temp{}

	err := utils.ValidateNonEmptyString(v)

	if err == nil {
		t.Error("expected error")
	}

	v.Field1 = "hello"

	err = utils.ValidateNonEmptyString(v)

	if err != nil {
		t.Error(err)
	}
}

func TestBudgets(t *testing.T) {
	bsV := validBudgets()
	for _, b := range bsV {
		testValid(b, t)
	}

	bsIv := invalidBudgets()
	for _, b := range bsIv {
		testInvalid(b, t)
	}
}

func TestDataTypes(t *testing.T) {
	v := validDataTypes()
	for _, d := range v {
		testValid(d, t)
	}

	iv := invalidDataTypes()
	for _, d := range iv {
		testInvalid(d, t)
	}
}

func TestColumnSchemas(t *testing.T) {
	v := validColumnSchemas()
	for _, c := range v {
		testValid(c, t)
	}

	iv := invalidColumnSchemas()
	for _, c := range iv {
		testInvalid(c, t)
	}
}

func TestDatasetPatch(t *testing.T) {
	v := entity.DatasetPatch{Name: "hej", Owner: "då", TotalBudget: entity.Budget{Epsilon: 1}}
	testValid(v, t)

	v = entity.DatasetPatch{Name: "", Owner: ":)", TotalBudget: entity.Budget{Epsilon: 1}}
	testInvalid(v, t)

	v = entity.DatasetPatch{Name: ":)", Owner: "", TotalBudget: entity.Budget{Epsilon: 1}}
	testInvalid(v, t)

	v = entity.DatasetPatch{Name: ":)", Owner: ":)", TotalBudget: entity.Budget{Epsilon: -1}}
	testInvalid(v, t)

	v = entity.DatasetPatch{}
	testInvalid(v, t)

}

func TestDatasetCreate(t *testing.T) {
	v := entity.DatasetCreate{Name: "name", Owner: "myowner", Schema: validColumnSchemas(), PrivacyNotion: "PureDP", TotalBudget: validBudgets()[0]}
	testValid(v, t)
	v = entity.DatasetCreate{Name: "name", Owner: "myowner", Schema: validColumnSchemas(), PrivacyNotion: "ApproxDP", TotalBudget: validBudgets()[0]}
	testValid(v, t)

	v = entity.DatasetCreate{Name: "", Owner: "myowner", Schema: validColumnSchemas(), PrivacyNotion: "PureDP", TotalBudget: validBudgets()[0]}
	testInvalid(v, t)

	v = entity.DatasetCreate{Name: "name", Owner: "", Schema: validColumnSchemas(), PrivacyNotion: "PureDP", TotalBudget: validBudgets()[0]}
	testInvalid(v, t)

	v = entity.DatasetCreate{Name: "name", Owner: "myowner", Schema: invalidColumnSchemas(), PrivacyNotion: "ApproxDP", TotalBudget: validBudgets()[0]}
	testInvalid(v, t)

	v = entity.DatasetCreate{Name: "", Owner: "myowner", Schema: validColumnSchemas(), PrivacyNotion: "Pure Deep Shit", TotalBudget: validBudgets()[0]}
	testInvalid(v, t)

	v = entity.DatasetCreate{Name: "", Owner: "myowner", Schema: validColumnSchemas(), PrivacyNotion: "PureDP", TotalBudget: invalidBudgets()[0]}
	testInvalid(v, t)

}

func TestQueryEvaluate(t *testing.T) {
	v := validQueryEvaluate()
	for _, q := range v {
		testValid(q, t)
	}
	iv := invalidQueryEvaluate()
	for _, q := range iv {
		testInvalid(q, t)
	}
}

func testValid(v valid, t *testing.T) {
	err := v.Valid()
	if err != nil {
		t.Error(err)
	}
}

func testInvalid(v valid, t *testing.T) {
	err := v.Valid()
	if err == nil {
		t.Error("expected invalid input")
	}
}

// test data

// BUDGETS

func validBudgets() []entity.Budget {
	d := 0.001
	return []entity.Budget{
		{Epsilon: 0.1},
		{Epsilon: 0.0},
		{Epsilon: 0.1, Delta: &d},
	}
}

func invalidBudgets() []entity.Budget {
	d := 0.001
	id := -d
	return []entity.Budget{
		{Epsilon: -0.1},
		{Epsilon: 0.1, Delta: &id},
		{Epsilon: -0.1, Delta: &d},
		{Epsilon: -0.1, Delta: &id},
	}
}

// DATATYPES

func validDataTypes() []entity.DataType {
	return []entity.DataType{
		{Type: &entity.IntType{Low: 1, High: 10}},
		{Type: &entity.TextType{}},
		{Type: &entity.BoolType{}},
		{Type: &entity.DoubleType{Low: -1, High: 0}},
		{Type: &entity.EnumType{Labels: []string{"a", "b", "c"}}},
	}
}

func invalidDataTypes() []entity.DataType {
	return []entity.DataType{
		{Type: &entity.DoubleType{Low: 10}},
		{Type: &entity.DoubleType{Low: -10, High: -11}},
		{Type: &entity.IntType{Low: 1}},
		{Type: &entity.IntType{Low: 10, High: 5}},
		{Type: &entity.EnumType{Labels: []string{"a", "a"}}},
		{},
	}
}

// COLUMN SCHEMAS

func validColumnSchemas() []entity.ColumnSchema {
	return []entity.ColumnSchema{
		{Name: "column", Type: entity.DataType{Type: &entity.BoolType{}}},
	}
}

func invalidColumnSchemas() []entity.ColumnSchema {
	return []entity.ColumnSchema{
		{Type: entity.DataType{Type: &entity.BoolType{}}},
		{Name: "myname"},
		{Name: ""},
	}
}

// QueryEvaluate

func validQueryEvaluate() []entity.QueryEvaluate {
	del := 1.0
	del2 := 0.3
	del3 := del - del2
	return []entity.QueryEvaluate{
		{
			Dataset: 1,
			Budget:  entity.Budget{Epsilon: 1, Delta: &del},
			Query: entity.Query{
				QuerySteps: []entity.QueryStep{
					entity.SelectTransformation{Columns: []string{"foo"}},
					entity.CountMeasurement{Params: entity.MeasurementParams{
						Budget: &entity.Budget{Epsilon: 0.5, Delta: &del2},
					}},
					entity.MeanMeasurement{Params: entity.MeasurementParams{
						Budget: &entity.Budget{Epsilon: 0.5, Delta: &del3},
					}},
				},
			},
		},
		{
			Dataset: 1,
			Budget:  entity.Budget{Epsilon: 1},
			Query: entity.Query{
				QuerySteps: []entity.QueryStep{
					entity.SelectTransformation{Columns: []string{"foo", "bar"}},
					entity.CountMeasurement{},
				},
			},
		},
	}
}

func invalidQueryEvaluate() []entity.QueryEvaluate {
	del := 0.1
	return []entity.QueryEvaluate{
		{
			Dataset: 1,
			Budget:  invalidBudgets()[0],
			Query: entity.Query{
				QuerySteps: []entity.QueryStep{
					entity.SelectTransformation{Columns: []string{"foo"}},
					entity.CountMeasurement{},
				},
			},
		},
		{
			Dataset: 1,
			Budget:  entity.Budget{Epsilon: 1},
			Query: entity.Query{
				QuerySteps: []entity.QueryStep{
					entity.SelectTransformation{Columns: []string{"foo"}},
					entity.CountMeasurement{
						Params: entity.MeasurementParams{
							Budget: &entity.Budget{Epsilon: 1},
						},
					},
					entity.CountMeasurement{
						Params: entity.MeasurementParams{
							Budget: &entity.Budget{Epsilon: 1},
						},
					},
				},
			},
		},
		{
			Dataset: 1,
			Budget:  entity.Budget{Epsilon: 1},
			Query: entity.Query{
				QuerySteps: []entity.QueryStep{
					entity.SelectTransformation{Columns: []string{"bar"}},
					entity.MaxMeasurement{
						Params: entity.MeasurementParams{
							Budget: &entity.Budget{Epsilon: 0.5},
						},
					},
					entity.MinMeasurement{
						Params: entity.MeasurementParams{
							Budget: &entity.Budget{Epsilon: 0.5, Delta: &del},
						},
					},
				},
			},
		},
	}
}
