package entity

import (
	"fmt"
	errors "webdp/internal/api/http"
)

type Budget struct {
	Epsilon float64  `json:"epsilon"`
	Delta   *float64 `json:"delta,omitempty"`
}

func (b Budget) Valid() error {
	if b.Delta != nil && (*b.Delta < 0 || b.Epsilon <= 0) {
		return fmt.Errorf("%w: epsilon and delta can't be negative: %s", errors.ErrBadInput, printBudget(b))
	} else if b.Epsilon <= 0 {
		return fmt.Errorf("%w: epsilon must be positive: %s", errors.ErrBadInput, printBudget(b))
	}
	return nil
}

func printBudget(b Budget) string {
	eps := b.Epsilon
	var del float64
	if b.Delta == nil {
		del = 0.0
	} else {
		del = *b.Delta
	}

	return fmt.Sprintf("(%f, %f)", eps, del)
}

type UserBudgets = []UserBudgetsResponse

type UserBudgetsResponse struct {
	Did       int64  `json:"dataset"`
	Allocated Budget `json:"allocated"`
	Consumed  Budget `json:"consumed"`
}

type DatasetBudgetAllocationResponse struct {
	Total      Budget            `json:"total"`
	Allocated  Budget            `json:"allocated"`
	Consumed   Budget            `json:"consumed"`
	Allocation []UserBudgetModel `json:"allocation"`
}

type UserBudgetModel struct {
	User      string `json:"user"`
	Allocated Budget `json:"allocated"`
	Consumed  Budget `json:"consumed"`
}
