package dpfuncs

import (
	"fmt"
	"googledp/entities"
	"googledp/requests"
	"strconv"
)

type ResultType struct {
	Rows []map[string]float64 `json:"rows"`
}

type binKey struct {
	FromInclude  int64
	ToNotInclude int64
}

func NewEvalQuery(req requests.Evaluate, data [][]string) (*ResultType, error) {

	allSubQueryHasBudget, err := isCorrectBudgetAndAllSubMeasHasBudget(req.Query, req.Budget, req.PrivacyNotion)

	if err != nil {
		return nil, err
	}

	subQueries, err := splitIntoSubQueries(req.Query)

	if err != nil {
		return nil, err
	}

	topLevelBudget := req.Budget
	schema := req.Schema

	// incase querystep budgets not set

	var numberOfMeasurements int64
	for _, step := range req.Query {
		if step.GetType() == entities.MEASUREMENT {
			numberOfMeasurements++
		}
	}

	epsilon := topLevelBudget.Epsilon / float64(numberOfMeasurements)
	var delta float64
	if topLevelBudget.Delta != nil {
		delta = *topLevelBudget.Delta / float64(numberOfMeasurements)
	} else {
		delta = 0
	}

	results := ResultType{
		Rows: make([]map[string]float64, 0),
	}

	for _, subQ := range subQueries {
		subQData := data
		bins := make(map[binKey][][]string)
		doBins := false
		for _, step := range subQ {
			var budget entities.Budget
			if !allSubQueryHasBudget {
				budget = entities.Budget{
					Epsilon: epsilon,
					Delta:   &delta,
				}
			} else {
				budget = *step.GetBudget()
			}

			switch step.GetOperation() {
			case entities.MEAN:
				var err error
				if doBins {
					budgetBins := getBudgetBins(budget, bins)
					_, err = doBinnedEval(bins, step, schema, budgetBins, mean, &results)
				} else {
					_, err = doEval(subQData, step, schema, budget, mean, &results)
				}
				if err != nil {
					return nil, err
				}
			case entities.SUM:
				var err error
				if doBins {
					budgetBins := getBudgetBins(budget, bins)
					_, err = doBinnedEval(bins, step, schema, budgetBins, sum, &results)
				} else {
					_, err = doEval(subQData, step, schema, budget, sum, &results)
				}
				if err != nil {
					return nil, err
				}
			case entities.STDEV:
				var err error
				if doBins {
					budgetBins := getBudgetBins(budget, bins)
					_, err = doBinnedEval(bins, step, schema, budgetBins, stDev, &results)
				} else {
					_, err = doEval(subQData, step, schema, *step.GetBudget(), stDev, &results)
				}
				if err != nil {
					return nil, err
				}
			case entities.VARIANCE:
				var err error
				if doBins {
					budgetBins := getBudgetBins(budget, bins)
					_, err = doBinnedEval(bins, step, schema, budgetBins, variance, &results)
				} else {
					_, err = doEval(subQData, step, schema, budget, variance, &results)
				}
				if err != nil {
					return nil, err
				}
			case entities.COUNT:
				var err error
				if doBins {
					budgetBins := getBudgetBins(budget, bins)
					_, err = doBinnedEval(bins, step, schema, budgetBins, count, &results)
				} else {
					_, err = doEval(subQData, step, schema, budget, count, &results)
				}
				if err != nil {
					return nil, err
				}
			case entities.QUANTILE:
				if doBins {
					tempMap := make(map[string]float64)
					budgetBins := getBudgetBins(budget, bins)
					for key, binnedData := range bins {
						result, err := quantile(step, schema, binnedData, budgetBins)

						if err != nil {
							return nil, err
						}

						tempMap := make(map[string]float64)

						for i, ires := range result {
							tempMap[fmt.Sprintf("bin_%d_%d_%s_%s", key, i, step.GetOperation(), step.GetColumn())] = ires
						}

					}
					results.Rows = append(results.Rows, tempMap)
				} else {
					result, err := quantile(step, schema, subQData, budget)

					if err != nil {
						return nil, err
					}

					tempMap := make(map[string]float64)

					for i, ires := range result {
						tempMap[fmt.Sprintf("%d_%s_%s", i, step.GetOperation(), step.GetColumn())] = ires
					}

					results.Rows = append(results.Rows, tempMap)
				}
			case entities.BIN:
				index, colType, err := getIndexAndTypeFromSchema(schema, step.GetColumn())

				if err != nil {
					return nil, err
				}

				if !checkIfNumber(colType.GetName()) {
					return nil, fmt.Errorf("can only do bin on columns with numbers")
				}

				biinz := step.GetBins()

				if !isUniqueAndSorted(biinz) {
					return nil, fmt.Errorf("bins not well formatted, bins should be unique and in ascending order")
				}

				if len(biinz) < 2 {
					return nil, fmt.Errorf("bins not well formatted, minimum 2 bins")
				}

				for i := range biinz[:len(biinz)-1] {
					bk := binKey{
						FromInclude:  biinz[i],
						ToNotInclude: biinz[i+1],
					}
					bins[bk] = make([][]string, 0)
				}

				for _, row := range subQData {
					for key := range bins {
						if f, err := strconv.ParseFloat(row[index], 64); err == nil {
							if f >= float64(key.FromInclude) && f < float64(key.ToNotInclude) {
								bins[key] = append(bins[key], row)
								break
							}
						} else {
							return nil, fmt.Errorf("type mismatch, error in data")
						}
					}
				}

				doBins = true

			case entities.FILTER:
				filters := step.GetFilters()
				filteredData, err := filterData(data, filters, schema)
				if err != nil {
					return nil, err
				}
				subQData = filteredData
			default:
				return nil, fmt.Errorf("oops something went wrong")
			}
		}
	}

	return &results, nil
}

type Binz map[binKey][][]string
type Data [][]string
type EvalFunc func(entities.QueryStep, []entities.Column, [][]string, entities.Budget) (float64, error)
type QS entities.QueryStep
type CS []entities.Column

func doBinnedEval(bins Binz, step QS, schema CS, budget entities.Budget, operation EvalFunc, res *ResultType) (*ResultType, error) {

	for key, binnedData := range bins {
		tempMap := make(map[string]float64)
		result, err := operation(step, schema, binnedData, budget)
		if err != nil {
			return nil, err
		}

		tempMap[fmt.Sprintf("%s_binned", step.GetColumn())] = float64(key.ToNotInclude)
		tempMap[step.GetOperation()] = result

		res.Rows = append(res.Rows, tempMap)
	}

	return res, nil
}

func doEval(data Data, step QS, schema CS, budget entities.Budget, operation EvalFunc, res *ResultType) (*ResultType, error) {

	result, err := operation(step, schema, data, budget)
	if err != nil {
		return nil, err
	}

	tempMap := make(map[string]float64)

	tempMap[fmt.Sprintf("%s_%s", step.GetColumn(), step.GetOperation())] = result

	res.Rows = append(res.Rows, tempMap)

	return res, nil
}

func getBudgetBins(budget entities.Budget, bins Binz) entities.Budget {
	nBins := len(bins)

	budgetBins := entities.Budget{
		Epsilon: budget.Epsilon / float64(nBins),
	}

	var del float64
	if budget.Delta != nil {
		del = *budget.Delta / float64(nBins)
	} else {
		del = 0
	}

	budgetBins.Delta = &del
	return budgetBins
}
