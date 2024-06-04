package dpfuncs

import (
	"fmt"
	"googledp/entities"
	"strconv"

	"github.com/google/differential-privacy/go/v3/dpagg"
)

func variance(step entities.QueryStep, schema []entities.Column, data [][]string, budget entities.Budget) (float64, error) {
	column := step.GetColumn()
	index, colType, err := getIndexAndTypeFromSchema(schema, column)

	if err != nil {
		return 0, err
	}

	opts, err := getVarianceOpts(step, colType, budget)

	if err != nil {
		return 0, err
	}

	dp_var, err := dpagg.NewBoundedVariance(opts)

	if err != nil {
		return 0, fmt.Errorf("something went wrong with initilizing the dp mean")
	}

	for _, value := range data {
		if f, err := strconv.ParseFloat(value[index], 64); err == nil {
			dp_var.Add(f)
		} else {
			fmt.Println(value)
			return 0, fmt.Errorf("error in data")
		}

	}

	return dp_var.Result()
}

func getVarianceOpts(step entities.QueryStep, typeSpec entities.ColType, allocatedBudget entities.Budget) (*dpagg.BoundedVarianceOptions, error) {
	if !checkIfNumber(typeSpec.GetName()) {
		return nil, fmt.Errorf("can't do mean measurment as column type is not a number")
	}

	lowBound := typeSpec.GetLow()
	highBound := typeSpec.GetHigh()

	opts := dpagg.BoundedVarianceOptions{
		MaxPartitionsContributed:     1,
		MaxContributionsPerPartition: 1,
		Lower:                        float64(lowBound),
		Upper:                        float64(highBound),
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
