package dpfuncs

import (
	"fmt"
	"googledp/entities"
	"strconv"
)

func filterData(data [][]string, filters []entities.Filter, cols []entities.Column) ([][]string, error) {
	filteredData := data
	for _, filter := range filters {
		var tempData [][]string
		col := filter.Column
		op := filter.Operator
		ty := filter.Type
		val := filter.ToFilter
		index, _, err := getIndexAndTypeFromSchema(cols, col)

		if err != nil {
			return nil, err
		}

		for _, row := range filteredData {
			err = caseHereFilter(op, ty, val, row, index, &tempData)
			if err != nil {
				return nil, err
			}
		}
		filteredData = tempData
	}
	return filteredData, nil
}

func caseHereFilter(operation string, columnType string, filterVal string, row []string, index int, tempData *[][]string) error {
	switch operation {
	case "<":
		return doComparisonNumber(operation, lt, columnType, filterVal, row, index, tempData)
	case "<=":
		return doComparisonNumber(operation, lte, columnType, filterVal, row, index, tempData)
	case ">":
		return doComparisonNumber(operation, gt, columnType, filterVal, row, index, tempData)
	case ">=":
		return doComparisonNumber(operation, gte, columnType, filterVal, row, index, tempData)
	case "==":
		if columnType == "number" {
			return doComparisonNumber(operation, eq, columnType, filterVal, row, index, tempData)
		} else {
			doComparisonString(stringEq, filterVal, row, index, tempData)
		}
	case "!=":
		if columnType == "number" {
			return doComparisonNumber(operation, neq, columnType, filterVal, row, index, tempData)
		} else {
			doComparisonString(stringNeq, filterVal, row, index, tempData)
		}
	default:
		return fmt.Errorf("%s not supported", operation)
	}

	return nil
}

func doComparisonNumber(operation string, compOp comparisonOpNumber, columnType string, filterVal string, row []string, index int, tempData *[][]string) error {
	if columnType != "number" {
		return fmt.Errorf("can't use %s operator on string", operation)
	}

	g, err := strconv.ParseFloat(filterVal, 64)

	if err != nil {
		return fmt.Errorf("provided value not a number %s", filterVal)
	}

	f, err := strconv.ParseFloat(row[index], 64)

	if err != nil {
		return fmt.Errorf("type mismatch, error in data")
	}

	if compOp(f, g) {
		*tempData = append(*tempData, row)
	}

	return nil
}

func doComparisonString(compOp comparisonOpString, filterVal string, row []string, index int, tempData *[][]string) {
	if compOp(row[index], filterVal) {
		*tempData = append(*tempData, row)
	}
}

type comparisonOpString func(string, string) bool

func stringEq(s1 string, s2 string) bool {
	return s1 == s2
}

func stringNeq(s1 string, s2 string) bool {
	return s1 != s2
}

type comparisonOpNumber func(float64, float64) bool

func lt(n1 float64, n2 float64) bool {
	return n1 < n2
}

func lte(n1 float64, n2 float64) bool {
	return n1 <= n2
}

func gt(n1 float64, n2 float64) bool {
	return n1 > n2
}

func gte(n1 float64, n2 float64) bool {
	return n1 >= n2
}

func eq(n1 float64, n2 float64) bool {
	return n1 == n2
}

func neq(n1 float64, n2 float64) bool {
	return n1 != n2
}
