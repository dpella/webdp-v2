package dpfuncs

import (
	"fmt"
	"googledp/entities"
	"strconv"

	"github.com/google/differential-privacy/go/v3/dpagg"
)

func stDev(step entities.QueryStep, schema []entities.Column, data [][]string, budget entities.Budget) (float64, error) {
	column := step.GetColumn()
	index, colType, err := getIndexAndTypeFromSchema(schema, column)

	if err != nil {
		return 0, err
	}

	opts, err := getStdevOpts(step, colType, budget)

	if err != nil {
		return 0, err
	}

	dp_stdev, err := dpagg.NewBoundedStandardDeviation(opts)

	if err != nil {
		return 0, fmt.Errorf("something went wrong with initilizing the dp mean")
	}

	for _, value := range data {
		if f, err := strconv.ParseFloat(value[index], 64); err == nil {
			dp_stdev.Add(f)
		} else {
			fmt.Println(value)
			return 0, fmt.Errorf("error in data")
		}

	}

	return dp_stdev.Result()

}

func getStdevOpts(step entities.QueryStep, typeSpec entities.ColType, allocatedBudget entities.Budget) (*dpagg.BoundedStandardDeviationOptions, error) {

	if !checkIfNumber(typeSpec.GetName()) {
		return nil, fmt.Errorf("can't do mean measurment as column type is not a number")
	}

	lowBound := typeSpec.GetLow()
	highBound := typeSpec.GetHigh()

	opts := dpagg.BoundedStandardDeviationOptions{
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
