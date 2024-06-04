package entity

import (
	"encoding/json"
	"fmt"
	errors "webdp/internal/api/http"
	"webdp/internal/api/http/utils"
)

// foobar interface so we can place the
// query steps in the same list

type Query struct {
	QuerySteps []QueryStep
}

type QueryStep interface {
	isQuery()
}

type QueryResult = map[string]interface{}

type QueryEvaluate struct {
	Dataset int64  `json:"dataset"`
	Budget  Budget `json:"budget"`
	Query   Query  `json:"query"`
}

type QueryCustom struct {
	Dataset int64  `json:"dataset"`
	Budget  Budget `json:"budget"`
	Query   string `json:"query" dpvalidation:"non-empty-string"`
}

type QueryAccuracy struct {
	Dataset    int64   `json:"dataset"`
	Budget     Budget  `json:"budget"`
	Query      Query   `json:"query"`
	Confidence float64 `json:"confidence"`
}

type QueryAccuracyResult struct {
	Accuracy []float64      `json:"accuracy"`
	Schema   []ColumnSchema `json:"schema"`
}

type QueryFromClientEvaluate struct {
	Budget        Budget         `json:"budget"`
	Query         Query          `json:"query"`
	Data          int64          `json:"dataset"`
	Schema        []ColumnSchema `json:"schema"`
	PrivacyNotion string         `json:"privacy_notion"`
	CallbackUrl   string         `json:"url"`
}

type QueryFromClientAccuracy struct {
	Budget        Budget         `json:"budget"`
	Query         Query          `json:"query"`
	Data          int64          `json:"dataset"`
	Schema        []ColumnSchema `json:"schema"`
	PrivacyNotion string         `json:"privacy_notion"`
	CallbackUrl   string         `json:"url"`
	Confidence    float64        `json:"confidence"`
}

type GroupByPartition struct {
	Grouping map[string]([]interface{}) `json:"groupby"`
}

type MeasurementParams struct {
	Column *string `json:"column,omitempty"`
	Mech   *string `json:"mech,omitempty"`
	Budget *Budget `json:"budget,omitempty"`
}

type ColumnMapping struct {
	Fun    string         `json:"fun"`
	Schema []ColumnSchema `json:"schema"`
}

func (q QueryEvaluate) Valid() error {
	err := q.Budget.Valid()
	if err != nil {
		return err
	}
	err = queryBudgetValidation(q.Query.QuerySteps, q.Budget)
	if err != nil {
		return err
	}
	return nil
}

func (q QueryCustom) Valid() error {
	err := utils.ValidateNonEmptyString(q)
	if err != nil {
		return err
	}
	return q.Budget.Valid()
}

func (q QueryAccuracy) Valid() error {
	if q.Confidence < 0 || q.Confidence > 1 {
		return fmt.Errorf("%w: confidence parameter is out of range. accepted range is [0,1] but given confidence was: %f", errors.ErrBadInput, q.Confidence)
	}
	return q.Budget.Valid()
}

func queryBudgetValidation(qs []QueryStep, b Budget) error {
	ms := utils.Map[QueryStep, MeasurementParams](
		utils.Filter(
			qs,
			func(qs QueryStep) bool {
				_, ok := qs.(Aggregate)
				return ok
			},
		),
		func(qs QueryStep) MeasurementParams {
			val, _ := qs.(Aggregate)
			return val.getParams()
		},
	)

	for _, m := range ms {
		err := nilOrValid(m.Budget)
		if err != nil {
			return err
		}
	}

	if atLeastOneBudgetNotNil(ms) {
		totBud := utils.Reduce[MeasurementParams, Budget](coalesceBudget(&Budget{}), ms,
			func(mp MeasurementParams, f Budget) Budget {
				bud := coalesceBudget(mp.Budget)
				del := *bud.Delta + *f.Delta
				return Budget{Epsilon: bud.Epsilon + f.Epsilon, Delta: &del}
			},
		)
		b = coalesceBudget(&b)
		if totBud.Epsilon != b.Epsilon || *totBud.Delta != *b.Delta {
			return fmt.Errorf("%w: query budgets doesn't add upp to the given total budget", errors.ErrBadInput)
		}
	}

	return nil
}

func nilOrValid(b *Budget) error {
	if b == nil {
		return nil
	}
	return b.Valid()
}

func coalesceBudget(b *Budget) Budget {
	if b == nil {
		var del float64
		return Budget{Epsilon: 0, Delta: &del}
	}
	if b.Delta == nil {
		var del float64
		return Budget{Epsilon: b.Epsilon, Delta: &del}
	}
	return *b
}

func atLeastOneBudgetNotNil(ms []MeasurementParams) bool {
	for _, val := range ms {
		if val.Budget != nil {
			return true
		}
	}
	return false
}

func (q Query) MarshalJSON() ([]byte, error) {
	stepsJSON := make([]json.RawMessage, len(q.QuerySteps))

	for i, step := range q.QuerySteps {
		stepJSON, err := json.Marshal(step)
		if err != nil {
			return nil, err
		}
		stepsJSON[i] = stepJSON
	}
	return json.Marshal(stepsJSON)
}

func (q *Query) UnmarshalJSON(data []byte) error {
	qs := make([]QueryStep, 0)
	var rawq []json.RawMessage
	err := json.Unmarshal(data, &rawq)
	if err != nil {
		return err
	}
	for _, r := range rawq {
		var m map[string]interface{}
		err := json.Unmarshal(r, &m)
		if err != nil {
			return err
		}
		// a query step has exactly one key
		if len(m) == 0 || len(m) > 1 {
			return fmt.Errorf("%w: unexpected query step: %v", errors.ErrBadFormatting, m)
		}
		for key := range m {
			mf := getUnmarshalFunc(key)
			queryStep, err := mf(r)
			if err != nil {
				return err
			}
			qs = append(qs, queryStep)
		}

	}

	q.QuerySteps = qs

	return nil
}

func getUnmarshalFunc(typ string) func([]byte) (QueryStep, error) {
	return func(data []byte) (QueryStep, error) {
		switch typ {
		case "select":
			var q SelectTransformation
			return q, marshalHelper(&q)(data)
		case "filter":
			var q FilterTransformation
			return q, marshalHelper(&q)(data)
		case "rename":
			var q RenameTransformation
			return q, marshalHelper(&q)(data)
		case "map":
			var q MapTransformation
			return q, marshalHelper(&q)(data)
		case "bin":
			var q BinTransformation
			return q, marshalHelper(&q)(data)
		case "count":
			var q CountMeasurement
			return q, marshalHelper(&q)(data)
		case "min":
			var q MinMeasurement
			return q, marshalHelper(&q)(data)
		case "max":
			var q MaxMeasurement
			return q, marshalHelper(&q)(data)
		case "mean":
			var q MeanMeasurement
			return q, marshalHelper(&q)(data)
		case "sum":
			var q SumMeasurement
			return q, marshalHelper(&q)(data)
		case "groupby":
			var q GroupByPartition
			return q, marshalHelper(&q)(data)
		default:
			return nil, fmt.Errorf("%w: could not find unmarshalfunc for %s", errors.ErrBadRequest, typ)
		}
	}
}

func marshalHelper[T QueryStep](obj *T) func([]byte) error {
	return func(data []byte) error {
		return json.Unmarshal(data, obj)
	}
}

type SelectTransformation struct {
	Columns []string `json:"select"`
}

type FilterTransformation struct {
	Filters []string `json:"filter"`
}

type RenameTransformation struct {
	Mapping map[string]string `json:"rename"`
}

type MapTransformation struct {
	Mapping ColumnMapping `json:"map"`
}

type BinTransformation struct {
	Bins map[string]([]interface{}) `json:"bin"`
}

type Aggregate interface {
	getParams() MeasurementParams
}

type CountMeasurement struct {
	Params MeasurementParams `json:"count"`
}

type MinMeasurement struct {
	Params MeasurementParams `json:"min"`
}

type MaxMeasurement struct {
	Params MeasurementParams `json:"max"`
}

type MeanMeasurement struct {
	Params MeasurementParams `json:"mean"`
}

type SumMeasurement struct {
	Params MeasurementParams `json:"sum"`
}

func (s SumMeasurement) getParams() MeasurementParams {
	return s.Params
}
func (s MeanMeasurement) getParams() MeasurementParams {
	return s.Params
}
func (s MaxMeasurement) getParams() MeasurementParams {
	return s.Params
}
func (s MinMeasurement) getParams() MeasurementParams {
	return s.Params
}
func (s CountMeasurement) getParams() MeasurementParams {
	return s.Params
}

// dummy implementation of QueryStep

func (s SelectTransformation) isQuery() {}
func (f FilterTransformation) isQuery() {}
func (r RenameTransformation) isQuery() {}
func (b BinTransformation) isQuery()    {}
func (c CountMeasurement) isQuery()     {}
func (m MinMeasurement) isQuery()       {}
func (m MaxMeasurement) isQuery()       {}
func (m MeanMeasurement) isQuery()      {}
func (s SumMeasurement) isQuery()       {}
func (g GroupByPartition) isQuery()     {}
func (m MapTransformation) isQuery()    {}
