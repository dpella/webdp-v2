package dpfuncs

import (
	"fmt"
	"googledp/entities"
	"strconv"

	"github.com/google/differential-privacy/go/v3/dpagg"
)

func sum(step entities.QueryStep, schema []entities.Column, data [][]string, budget entities.Budget) (float64, error) {
	index, colType, err := getIndexAndTypeFromSchema(schema, step.GetColumn())

	if err != nil {
		return 0, err
	}

	opts, err := getSumOpts(step, colType, budget)

	if err != nil {
		return 0, err
	}

	dp_sum, err := dpagg.NewBoundedSumFloat64(opts)

	if err != nil {
		return 0, fmt.Errorf("something went wrong with initilizing the dp mean")
	}

	for _, value := range data {
		if f, err := strconv.ParseFloat(value[index], 64); err == nil {
			dp_sum.Add(f)
		} else {
			fmt.Println(value)
			return 0, fmt.Errorf("error in data")
		}

	}

	return dp_sum.Result()
}

func getSumOpts(step entities.QueryStep, typeSpec entities.ColType, allocatedBudget entities.Budget) (*dpagg.BoundedSumFloat64Options, error) {
	if !checkIfNumber(typeSpec.GetName()) {
		return nil, fmt.Errorf("can't do mean measurment as column type is not a number")
	}

	lowBound := typeSpec.GetLow()
	highBound := typeSpec.GetHigh()

	opts := dpagg.BoundedSumFloat64Options{
		MaxPartitionsContributed: 1,
		Lower:                    float64(lowBound),
		Upper:                    float64(highBound),
	}
	mech := step.GetMechanism()

	if mech != "" {
		opts.Noise = getNoiseGenerator(mech)
		mech = "Laplace"
	}

	opts.Epsilon = allocatedBudget.Epsilon
	opts.Delta = *allocatedBudget.Delta

	return &opts, nil
}
