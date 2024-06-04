package services

import (
	"fmt"
	"strconv"
	errors "webdp/internal/api/http"
	"webdp/internal/api/http/entity"
	"webdp/internal/api/http/repo/postgres"
	"webdp/internal/api/http/utils"
)

type BudgetService struct {
	postg postgres.BudgetPostgres
}

func NewBudgetService(budgetRepo postgres.BudgetPostgres) BudgetService {
	return BudgetService{postg: budgetRepo}
}

func (b BudgetService) GetUserBudgets(userHandle string) (entity.UserBudgets, error) {
	budgets, err := b.postg.GetUserBudgets(userHandle)
	if err != nil {
		return entity.UserBudgets{}, errors.WrapDBError(err, "get user budgets", userHandle)
	}
	return budgets, nil
}

func (b BudgetService) AddConsumedBudgetToUser(user string, dataset int64, spent entity.Budget) error {
	ub, err := b.postg.GetConsumedUserBudgetOnDataset(user, dataset)
	if err != nil {
		return errors.WrapDBError(err, "get consumed budget for", user)
	}
	newConsumed := budgetAdd(ub, spent)
	if err := b.postg.UpdateUserConsumedBudget(user, dataset, newConsumed); err != nil {
		return errors.WrapDBError(err, "update consumed budget for", user)
	}
	return nil
}

func (q BudgetService) HasUserEnoughBudget(user string, dataset int64, queryBudget entity.Budget) bool {

	all, err := q.postg.GetAllocatedUserBudgetOnDataset(user, dataset)
	if err != nil {
		return false
	}

	con, err := q.postg.GetConsumedUserBudgetOnDataset(user, dataset)
	if err != nil {
		return false
	}

	diff := budSub(all, con)
	return budLeq(queryBudget, diff)
}

func (b BudgetService) GetDatasetBudget(datasetId int64) (entity.DatasetBudgetAllocationResponse, error) {
	alloc, err := b.postg.GetDatasetUserAllocations(datasetId)
	if err != nil {
		return entity.DatasetBudgetAllocationResponse{}, errors.WrapDBError(err, "get budget allocation", strconv.FormatInt(datasetId, 10))
	}
	return alloc, nil
}

func (b BudgetService) GetUserDatasetBudget(userHandle string, datasetId int64) (entity.Budget, error) {
	budget, err := b.postg.GetAllocatedUserBudgetOnDataset(userHandle, datasetId)
	if err != nil {
		return entity.Budget{}, errors.WrapDBError(err, "get user dataset budget", userHandle+" "+strconv.FormatInt(datasetId, 10))
	}
	return budget, nil
}

func (b BudgetService) PostUserDatasetBudget(userHandle string, datasetId int64, budget entity.Budget) error {
	dataset, err := b.GetDatasetBudget(datasetId)
	if err != nil {
		return err
	}

	hasallocated, err := b.UserHasDatasetBudget(userHandle, datasetId)
	if err != nil {
		return err
	}

	if hasallocated {
		return fmt.Errorf("%w: user %s already has allocated budget on dataset %d", errors.ErrBadInput, userHandle, datasetId)
	}

	if dataset.Total.Epsilon < dataset.Allocated.Epsilon+budget.Epsilon {
		return fmt.Errorf("%w: not enough epsilon budget to allocate. Total epsilon: %f, allocated: %f, new total: %f", errors.ErrBadInput, dataset.Total.Epsilon, dataset.Allocated.Epsilon, (dataset.Allocated.Epsilon + budget.Epsilon))
	}

	if coalesce(dataset.Total.Delta) < coalesce(dataset.Allocated.Delta)+coalesce(budget.Delta) {
		return fmt.Errorf("%w: not enough delta budget to allocate", errors.ErrBadInput)
	}

	if err := b.postg.CreateUserBudgetAllocation(userHandle, datasetId, budget); err != nil {
		return errors.WrapDBError(err, "allocate budget for", userHandle+" "+strconv.FormatInt(datasetId, 10))
	}

	return nil

}

func (b BudgetService) PatchUserDatasetBudget(userHandle string, datasetId int64, budget entity.Budget) error {
	// Add checks and changes to dataset allocation
	dataset, err := b.postg.GetDatasetUserAllocations(datasetId)

	if err != nil {
		return fmt.Errorf("dataset with id %d does not exist. err: %s", datasetId, err.Error())
	}

	allocated, err := b.postg.UserBudgetOnDatasetExists(userHandle, datasetId)

	if err != nil {
		return fmt.Errorf("something went wrong: %s", err.Error())
	}

	if !allocated {
		return fmt.Errorf("user with id %s does not have any budget allocated on dataset with id %d", userHandle, datasetId)
	}

	uBudget, err := b.GetUserDatasetBudget(userHandle, datasetId)

	if err != nil {
		return fmt.Errorf("something went wrong: %s", err.Error())
	}

	if dataset.Total.Epsilon < dataset.Allocated.Epsilon-uBudget.Epsilon+budget.Epsilon {
		return fmt.Errorf("dataset total epsilon allocation not enough. Total Epsilon: %f, Allocated Epsilon: %f, Previous User Epsilon: %f, New User Epsilon: %f", dataset.Total.Epsilon, dataset.Allocated.Epsilon, uBudget.Epsilon, budget.Epsilon)
	}

	if coalesce(dataset.Total.Delta) < coalesce(dataset.Allocated.Delta)-coalesce(uBudget.Delta)+coalesce(budget.Delta) {
		return fmt.Errorf("dataset total delta allocation not enough. Total Delta: %f, Allocated Delta: %f, Previous User Delta: %f, New User Delta: %f", coalesce(dataset.Total.Delta), coalesce(dataset.Allocated.Delta), coalesce(uBudget.Delta), coalesce(budget.Delta))
	}

	return b.postg.UpdateUserBudgetAllocation(userHandle, datasetId, budget)
}

func (b BudgetService) DeleteUserDatasetBudget(userHandle string, datasetId int64) error {
	if err := b.postg.DeleteUserBudgetAllocation(userHandle, datasetId); err != nil {
		return errors.WrapDBError(err, "delete user dataset budget", userHandle+" "+strconv.FormatInt(datasetId, 10))
	}
	return nil
}

func (b BudgetService) UserHasDatasetBudget(userHandle string, datasetId int64) (bool, error) {
	has, err := b.postg.UserBudgetOnDatasetExists(userHandle, datasetId)
	if err != nil {
		return false, errors.WrapDBError(err, "budget exists", userHandle+" "+strconv.FormatInt(datasetId, 10))
	}
	return has, nil
}

// helpers

func coalesce[T float32 | float64](t *T) T {
	if t == nil {
		return 0.0
	}
	return *t
}

func budgetAdd(a entity.Budget, b entity.Budget) entity.Budget {
	del := coalesce(a.Delta) + coalesce(b.Delta)
	return entity.Budget{
		Epsilon: a.Epsilon + b.Epsilon,
		Delta:   &del,
	}
}

func budSub(a entity.Budget, b entity.Budget) entity.Budget {
	del := coalesce(a.Delta) - coalesce(b.Delta)
	return entity.Budget{
		Epsilon: utils.RoundFloat(a.Epsilon, 10) - utils.RoundFloat(b.Epsilon, 10),
		Delta:   &del,
	}
}

func budLeq(a entity.Budget, b entity.Budget) bool {
	return utils.RoundFloat(a.Epsilon, 10) <= utils.RoundFloat(b.Epsilon, 10) && coalesce(a.Delta) <= coalesce(b.Delta)
}
