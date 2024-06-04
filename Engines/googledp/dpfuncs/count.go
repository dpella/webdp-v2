package dpfuncs

import (
	"fmt"
	"googledp/entities"

	"github.com/google/differential-privacy/go/v3/dpagg"
)

func count(step entities.QueryStep, schema []entities.Column, data [][]string, budget entities.Budget) (float64, error) {
	_, colType, err := getIndexAndTypeFromSchema(schema, step.GetColumn())

	if err != nil {
		return 0, err
	}

	opts, err := getCountOpts(step, colType, budget)

	if err != nil {
		return 0, err
	}

	dp_count, err := dpagg.NewCount(opts)

	if err != nil {
		return 0, fmt.Errorf("something went wrong with initilizing the dp mean")
	}

	err = dp_count.IncrementBy(int64(len(data)))

	if err != nil {
		return 0, fmt.Errorf("error in data")
	}

	result, err := dp_count.Result()

	return float64(result), err

}

func getCountOpts(step entities.QueryStep, typeSpec entities.ColType, allocatedBudget entities.Budget) (*dpagg.CountOptions, error) {
	if !checkIfNumber(typeSpec.GetName()) {
		return nil, fmt.Errorf("can't do mean measurment as column type is not a number")
	}

	opts := dpagg.CountOptions{
		MaxPartitionsContributed: 1,
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
