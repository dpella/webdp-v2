package dpfuncs

import (
	"fmt"
	"googledp/dpfuncs/validator"
	"googledp/entities"

	"github.com/google/differential-privacy/go/v3/noise"
)

func isCorrectBudgetAndAllSubMeasHasBudget(query entities.Query, topLevelBudget entities.Budget, privacyNotation string) (bool, error) {

	allSubQueryHasBudget := true

	var epsilon float64
	var delta float64

	for _, step := range query {
		if step.GetType() == entities.MEASUREMENT {
			budget := step.GetBudget()
			if budget != nil {
				e, d, err := correctBudget(*budget, privacyNotation)

				if err != nil {
					return false, err
				}

				epsilon += e
				delta += d

				allSubQueryHasBudget = allSubQueryHasBudget && true
			}
			allSubQueryHasBudget = allSubQueryHasBudget && false
		}
	}

	return allSubQueryHasBudget && topLevelBudget.Epsilon == epsilon && topLevelBudget.Delta == &delta, nil
}

func correctBudget(budget entities.Budget, privacyNotation string) (float64, float64, error) {
	if privacyNotation == "ApproxDP" {
		if budget.Epsilon > 0 && budget.Delta != nil {
			return budget.Epsilon, 0, nil
		}
		return 0, 0, fmt.Errorf("approxDP budget incorrectly formatted")
	} else if privacyNotation == "PureDP" {
		if budget.Epsilon > 0 && budget.Delta == nil {
			return budget.Epsilon, *budget.Delta, nil
		}
		return 0, 0, fmt.Errorf("puredp budget incorrectly formatted")
	}

	return 0, 0, fmt.Errorf("invalid privacy notation %s", privacyNotation)
}

func isValidShapeSM(input [][]string) bool {
	smValidator := validator.NewSMValidator()

	for _, subQ := range input {
		if !smValidator.VerifyInputs(subQ) {
			return false
		}
		smValidator.Reset()
	}

	return true
}

func splitIntoSubQueries(query entities.Query) ([][]entities.QueryStep, error) {
	// Allowed Query format
	// Filter -> Bins -> Measurement -> Result
	// Filter -> Measurement -> Result
	// Bins -> Measurement -> Result
	// Measurement -> Result
	stepsOps := make([]string, 0)
	blablaValidation := make([]string, 0)
	for _, step := range query {
		if step.GetType() == entities.MEASUREMENT {
			blablaValidation = append(blablaValidation, step.GetType())
		} else {
			blablaValidation = append(blablaValidation, step.GetOperation())
		}
		stepsOps = append(stepsOps, step.GetType())
	}

	formattedQuery := make([][]entities.QueryStep, 0)
	forValidation := make([][]string, 0)

	var tempQS []entities.QueryStep
	var tempFV []string

	for i, op := range stepsOps {
		tempQS = append(tempQS, query[i])
		tempFV = append(tempFV, blablaValidation[i])
		if op == entities.MEASUREMENT {
			formattedQuery = append(formattedQuery, tempQS)
			forValidation = append(forValidation, tempFV)
			tempQS = nil
			tempFV = nil
		}
	}

	if tempQS != nil {
		formattedQuery = append(formattedQuery, tempQS)
		forValidation = append(forValidation, tempFV)
	}

	if isValidShapeSM(forValidation) {
		return formattedQuery, nil
	}

	return nil, fmt.Errorf("query is not wellformed, filter -> bin -> measurement in this order")

}

func getIndexAndTypeFromSchema(columns []entities.Column, columnName string) (int, entities.ColType, error) {
	for i, col := range columns {
		if col.Name == columnName {
			return i, col.Type, nil
		}
	}
	return -1, nil, fmt.Errorf("column with name %s not in schema", columnName)
}

func checkIfNumber(name string) bool {
	if name == "Int" || name == "Double" {
		return true
	}
	return false
}

func getNoiseGenerator(name string) noise.Noise {
	if name == "Laplace" {
		return noise.Laplace()
	} else if name == "Gaussian" {
		return noise.Gaussian()
	}
	return noise.Laplace()
}

func isUniqueAndSorted(arr []int64) bool {
	if len(arr) <= 1 {
		return true
	}

	for i := 1; i < len(arr); i++ {
		if arr[i] <= arr[i-1] {
			return false
		}
	}

	return true
}
