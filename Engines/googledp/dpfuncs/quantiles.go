package dpfuncs

import (
	"fmt"
	"googledp/entities"
	"strconv"

	"github.com/google/differential-privacy/go/v3/dpagg"
)

func quantile(step entities.QueryStep, schema []entities.Column, data [][]string, budget entities.Budget) ([]float64, error) {

	index, colType, err := getIndexAndTypeFromSchema(schema, step.GetColumn())

	if err != nil {
		return nil, err
	}

	opts, err := getQuantilesOpts(step, colType, budget)

	if err != nil {
		return nil, err
	}

	dp_quan, err := dpagg.NewBoundedQuantiles(opts)

	if err != nil {
		return nil, fmt.Errorf("something went wrong with initilizing the dp mean")
	}

	for _, value := range data {
		if f, err := strconv.ParseFloat(value[index], 64); err == nil {
			dp_quan.Add(f)
		} else {
			fmt.Println(value)
			return nil, fmt.Errorf("error in data")
		}

	}

	result := []float64{}

	sthep := 1.0 / 100
	value := 0.0

	for i := 0; i < 99; i++ {
		value += sthep

		intermediateRes, err := dp_quan.Result(value)

		if err != nil {
			return nil, err
		}

		result = append(result, intermediateRes)
	}

	return result, nil
}

func getQuantilesOpts(step entities.QueryStep, typeSpec entities.ColType, allocatedBudget entities.Budget) (*dpagg.BoundedQuantilesOptions, error) {
	if !checkIfNumber(typeSpec.GetName()) {
		return nil, fmt.Errorf("can't do mean measurment as column type is not a number")
	}

	lowBound := typeSpec.GetLow()
	highBound := typeSpec.GetHigh()

	opts := dpagg.BoundedQuantilesOptions{
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
