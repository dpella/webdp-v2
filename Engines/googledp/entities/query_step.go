package entities

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type Query []QueryStep

const (
	MEAN           = "mean"
	SUM            = "sum"
	STDEV          = "stdev"
	VARIANCE       = "variance"
	COUNT          = "count"
	QUANTILE       = "quantile"
	BIN            = "bin"
	FILTER         = "filter"
	TRANSFORMATION = "transformation"
	MEASUREMENT    = "measurement"
)

func (q *Query) UnmarshalJSON(data []byte) error {
	var stepsMap []map[string]json.RawMessage
	if err := json.Unmarshal(data, &stepsMap); err != nil {
		return err
	}

	for _, stepMap := range stepsMap {
		for key, value := range stepMap {
			var step QueryStep
			switch key {
			case MEAN:
				var meanStep MeanStep
				if err := json.Unmarshal(value, &meanStep); err != nil {
					return err
				}
				step = meanStep
			case SUM:
				var sumStep SumStep
				if err := json.Unmarshal(value, &sumStep); err != nil {
					return err
				}
				step = sumStep
			case STDEV:
				var stdevStep StdevStep
				if err := json.Unmarshal(value, &stdevStep); err != nil {
					return err
				}
				step = stdevStep
			case VARIANCE:
				var varianceStep VarianceStep
				if err := json.Unmarshal(value, &varianceStep); err != nil {
					return err
				}
				step = varianceStep
			case COUNT:
				var countStep CountStep
				if err := json.Unmarshal(value, &countStep); err != nil {
					return err
				}
				step = countStep
			case QUANTILE:
				var quantileStep QuantileStep
				if err := json.Unmarshal(value, &quantileStep); err != nil {
					return err
				}
				step = quantileStep
			case BIN:
				var temp map[string]json.RawMessage
				if err := json.Unmarshal(value, &temp); err != nil {
					return err
				}

				var binStep BinStep

				for k, v := range temp {
					binStep.Column = k
					var bins []int64

					if err := json.Unmarshal(v, &bins); err != nil {
						return err
					}
					binStep.Bins = bins
					break
				}

				step = binStep
			case FILTER:
				var temp []string
				if err := json.Unmarshal(value, &temp); err != nil {
					return err
				}

				var filterStep FilterStep

				for _, f := range temp {

					var filter Filter
					split := strings.Fields(f)
					filter.Column = split[0]
					filter.Operator = split[1]

					_, err := strconv.Atoi(split[2])
					if err != nil {
						filter.Type = "string"
					} else {
						filter.Type = "number"
					}

					filter.ToFilter = split[2]
					filterStep.Filters = append(filterStep.Filters, filter)
				}

				step = filterStep
			// Add cases for other types as needed
			default:
				return fmt.Errorf("unknown query step type: %s", key)
			}
			*q = append(*q, step)
		}
	}

	return nil
}

type QueryStep interface {
	GetOperation() string
	GetColumn() string
	GetMechanism() string
	GetBudget() *Budget
	GetBins() []int64
	GetFilters() []Filter
	GetType() string
}

type MeanStep struct {
	Column string  `json:"column"`
	Mech   string  `json:"mech"`
	Budget *Budget `json:"budget"`
}

func (s MeanStep) GetOperation() string {
	return MEAN
}

func (s MeanStep) GetColumn() string {
	return s.Column
}

func (s MeanStep) GetMechanism() string {
	return s.Mech
}

func (s MeanStep) GetBudget() *Budget {
	return s.Budget
}

func (s MeanStep) GetBins() []int64 {
	return nil
}

func (s MeanStep) GetFilters() []Filter {
	return nil
}

func (s MeanStep) GetType() string {
	return MEASUREMENT
}

type SumStep struct {
	Column string  `json:"column"`
	Mech   string  `json:"mech"`
	Budget *Budget `json:"budget"`
}

func (s SumStep) GetOperation() string {
	return SUM
}

func (s SumStep) GetColumn() string {
	return s.Column
}

func (s SumStep) GetMechanism() string {
	return s.Mech
}

func (s SumStep) GetBudget() *Budget {
	return s.Budget
}

func (s SumStep) GetBins() []int64 {
	return nil
}

func (s SumStep) GetFilters() []Filter {
	return nil
}

func (s SumStep) GetType() string {
	return MEASUREMENT
}

type StdevStep struct {
	Column string  `json:"column"`
	Mech   string  `json:"mech"`
	Budget *Budget `json:"budget"`
}

func (s StdevStep) GetOperation() string {
	return STDEV
}

func (s StdevStep) GetColumn() string {
	return s.Column
}

func (s StdevStep) GetMechanism() string {
	return s.Mech
}

func (s StdevStep) GetBudget() *Budget {
	return s.Budget
}

func (s StdevStep) GetBins() []int64 {
	return nil
}

func (s StdevStep) GetFilters() []Filter {
	return nil
}

func (s StdevStep) GetType() string {
	return MEASUREMENT
}

type VarianceStep struct {
	Column string  `json:"column"`
	Mech   string  `json:"mech"`
	Budget *Budget `json:"budget"`
}

func (s VarianceStep) GetOperation() string {
	return VARIANCE
}

func (s VarianceStep) GetColumn() string {
	return s.Column
}

func (s VarianceStep) GetMechanism() string {
	return s.Mech
}

func (s VarianceStep) GetBudget() *Budget {
	return s.Budget
}

func (s VarianceStep) GetBins() []int64 {
	return nil
}

func (s VarianceStep) GetFilters() []Filter {
	return nil
}

func (s VarianceStep) GetType() string {
	return MEASUREMENT
}

type CountStep struct {
	Column string  `json:"column"`
	Mech   string  `json:"mech"`
	Budget *Budget `json:"budget"`
}

func (s CountStep) GetOperation() string {
	return COUNT
}

func (s CountStep) GetColumn() string {
	return s.Column
}

func (s CountStep) GetMechanism() string {
	return s.Mech
}

func (s CountStep) GetBudget() *Budget {
	return s.Budget
}

func (s CountStep) GetBins() []int64 {
	return nil
}

func (s CountStep) GetFilters() []Filter {
	return nil
}

func (s CountStep) GetType() string {
	return MEASUREMENT
}

type QuantileStep struct {
	Column string  `json:"column"`
	Mech   string  `json:"mech"`
	Budget *Budget `json:"budget"`
}

func (s QuantileStep) GetOperation() string {
	return QUANTILE
}

func (s QuantileStep) GetColumn() string {
	return s.Column
}

func (s QuantileStep) GetMechanism() string {
	return s.Mech
}

func (s QuantileStep) GetBudget() *Budget {
	return s.Budget
}

func (s QuantileStep) GetBins() []int64 {
	return nil
}

func (s QuantileStep) GetFilters() []Filter {
	return nil
}

func (s QuantileStep) GetType() string {
	return MEASUREMENT
}

type BinStep struct {
	Column string  `json:"column"`
	Bins   []int64 `json:"bins"`
}

func (s BinStep) GetOperation() string {
	return BIN
}

func (s BinStep) GetColumn() string {
	return s.Column
}

func (s BinStep) GetMechanism() string {
	return ""
}

func (s BinStep) GetBudget() *Budget {
	return nil
}

func (s BinStep) GetBins() []int64 {
	return s.Bins
}

func (s BinStep) GetFilters() []Filter {
	return nil
}

func (s BinStep) GetType() string {
	return TRANSFORMATION
}

type Filter struct {
	Column   string
	Operator string
	Type     string
	ToFilter string
}

type FilterStep struct {
	Filters []Filter
}

func (s FilterStep) GetOperation() string {
	return FILTER
}

func (s FilterStep) GetColumn() string {
	return ""
}

func (s FilterStep) GetMechanism() string {
	return ""
}

func (s FilterStep) GetBudget() *Budget {
	return nil
}

func (s FilterStep) GetBins() []int64 {
	return nil
}

func (s FilterStep) GetFilters() []Filter {
	return s.Filters
}

func (s FilterStep) GetType() string {
	return TRANSFORMATION
}
