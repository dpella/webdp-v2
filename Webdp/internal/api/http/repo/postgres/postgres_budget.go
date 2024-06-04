package postgres

import (
	"database/sql"
	"fmt"
	errors "webdp/internal/api/http"
	"webdp/internal/api/http/entity"
	"webdp/internal/api/http/utils"
)

type BudgetPostgres struct {
	db *sql.DB
}

func NewBudgetPostgres(conn *sql.DB) BudgetPostgres {
	return BudgetPostgres{db: conn}
}

func (b BudgetPostgres) GetUserBudgets(userHandle string) ([]entity.UserBudgetsResponse, error) {
	tx, err := b.db.Begin()
	if err != nil {
		return []entity.UserBudgetsResponse{}, err
	}
	defer dfun(err, tx)
	q := "SELECT dataset, all_epsilon, all_delta, con_epsilon, con_delta FROM UserBudgetAllocation WHERE userid = $1"

	rows, err := tx.Query(q, userHandle)
	if err != nil {
		return []entity.UserBudgetsResponse{}, err
	}

	out := make([]entity.UserBudgetsResponse, 0)
	for rows.Next() {
		var did int64
		var aeps float64
		var adel, ceps, cdel sql.NullFloat64
		err = rows.Scan(&did, &aeps, &adel, &ceps, &cdel)
		if err != nil {
			rows.Close()
			return []entity.UserBudgetsResponse{}, err
		}
		out = append(out, entity.UserBudgetsResponse{
			Did:       did,
			Allocated: entity.Budget{Epsilon: aeps, Delta: &adel.Float64},
			Consumed:  entity.Budget{Epsilon: ceps.Float64, Delta: &cdel.Float64},
		})
	}
	rows.Close()
	if err = tx.Commit(); err != nil {
		return []entity.UserBudgetsResponse{}, err
	}

	return out, nil
}

func (b BudgetPostgres) GetDatasetUserAllocations(datasetId int64) (entity.DatasetBudgetAllocationResponse, error) {
	tx, err := b.db.Begin()
	if err != nil {
		return entity.DatasetBudgetAllocationResponse{}, err
	}
	defer dfun(err, tx)

	q := "SELECT total_epsilon, total_delta, all_epsilon, all_delta, con_epsilon, con_delta FROM DatasetAllocatedConsumed WHERE id = $1"

	row := tx.QueryRow(q, datasetId)
	var tote, totd, ae, ad, ce, cd float64
	err = row.Scan(&tote, &totd, &ae, &ad, &ce, &cd)

	if err != nil {
		return entity.DatasetBudgetAllocationResponse{}, err
	}

	budmodel, err := getUserAllocations(datasetId, tx)
	if err != nil {
		return entity.DatasetBudgetAllocationResponse{}, err
	}

	out := entity.DatasetBudgetAllocationResponse{
		Total:      entity.Budget{Epsilon: tote, Delta: &totd},
		Allocated:  entity.Budget{Epsilon: ae, Delta: &ad},
		Consumed:   entity.Budget{Epsilon: ce, Delta: &cd},
		Allocation: budmodel,
	}

	if err = tx.Commit(); err != nil {
		return entity.DatasetBudgetAllocationResponse{}, err
	}
	return out, nil
}

func getUserAllocations(datasetId int64, tx *sql.Tx) ([]entity.UserBudgetModel, error) {
	q := "SELECT userid, all_epsilon, all_delta, con_epsilon, con_delta FROM UserBudgetAllocation WHERE dataset = $1"
	rows, err := tx.Query(q, datasetId)
	if err != nil {
		return []entity.UserBudgetModel{}, err
	}

	out := make([]entity.UserBudgetModel, 0)

	for rows.Next() {
		var user string
		var aeps float64
		var adel, ceps, cdel sql.NullFloat64
		err = rows.Scan(&user, &aeps, &adel, &ceps, &cdel)
		if err != nil {
			rows.Close()
			return []entity.UserBudgetModel{}, err
		}
		out = append(out, entity.UserBudgetModel{
			User:      user,
			Allocated: entity.Budget{Epsilon: aeps, Delta: &adel.Float64},
			Consumed:  entity.Budget{Epsilon: ceps.Float64, Delta: &cdel.Float64},
		})
	}
	rows.Close()
	return out, nil
}

func (b BudgetPostgres) GetAllocatedUserBudgetOnDataset(userHandle string, datasetId int64) (entity.Budget, error) {
	return b.budgetHelper("all_epsilon", "all_delta", userHandle, datasetId)
}

func (b BudgetPostgres) GetConsumedUserBudgetOnDataset(userhandle string, datasetId int64) (entity.Budget, error) {
	return b.budgetHelper("con_epsilon", "con_delta", userhandle, datasetId)
}

func (b BudgetPostgres) budgetHelper(eps string, delta string, userhandle string, datasetId int64) (entity.Budget, error) {
	tx, err := b.db.Begin()
	if err != nil {
		return entity.Budget{}, err
	}
	defer dfun(err, tx)

	q := fmt.Sprintf("SELECT COALESCE(%s, 0), COALESCE(%s, 0) FROM UserBudgetAllocation WHERE userid = $1 AND dataset = $2", eps, delta)
	row := tx.QueryRow(q, userhandle, datasetId)
	var e float64
	var d sql.NullFloat64
	err = row.Scan(&e, &d)
	if err != nil {
		return entity.Budget{}, errors.ErrNotFound
	}
	if err = tx.Commit(); err != nil {
		return entity.Budget{}, err
	}
	return entity.Budget{Epsilon: e, Delta: &d.Float64}, nil
}

func (b BudgetPostgres) CreateUserBudgetAllocation(userhandle string, datasetId int64, allocation entity.Budget) error {
	tx, err := b.db.Begin()
	if err != nil {
		return err
	}
	defer dfun(err, tx)
	q := "INSERT INTO UserBudgetAllocation (dataset, userid, all_epsilon, all_delta) VALUES ($1, $2, $3, $4)"
	_, err = tx.Exec(q, datasetId, userhandle, allocation.Epsilon, allocation.Delta)
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (b BudgetPostgres) UpdateUserBudgetAllocation(userHandle string, datasetId int64, allocation entity.Budget) error {
	tx, err := b.db.Begin()
	if err != nil {
		return err
	}

	defer dfun(err, tx)

	q := "UPDATE UserBudgetAllocation SET all_epsilon = $1, all_delta = $2 WHERE dataset = $3 AND userid = $4"
	_, err = tx.Exec(q, utils.RoundFloat(allocation.Epsilon, 10), allocation.Delta, datasetId, userHandle)
	if err != nil {
		return errors.ErrNotFound
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (b BudgetPostgres) UpdateUserConsumedBudget(userHandle string, datasetId int64, newConsumed entity.Budget) error {
	tx, err := b.db.Begin()
	if err != nil {
		return err
	}

	defer dfun(err, tx)

	q := "UPDATE UserBudgetAllocation SET con_epsilon = $1, con_delta = $2 WHERE dataset = $3 AND userid = $4"

	_, err = tx.Exec(q, newConsumed.Epsilon, newConsumed.Delta, datasetId, userHandle)
	if err != nil {
		return errors.ErrNotFound
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (b BudgetPostgres) DeleteUserBudgetAllocation(userHandle string, datasetId int64) error {
	tx, err := b.db.Begin()
	if err != nil {
		return err
	}

	defer dfun(err, tx)

	q := "DELETE FROM UserBudgetAllocation WHERE userid = $1 AND dataset = $2"
	_, err = tx.Exec(q, userHandle, datasetId)
	if err != nil {
		return errors.ErrNotFound
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (b BudgetPostgres) UserBudgetOnDatasetExists(handle string, dataset int64) (bool, error) {
	tx, err := b.db.Begin()
	if err != nil {
		return false, err
	}

	defer dfun(err, tx)
	q := "SELECT 1 FROM UserBudgetAllocation WHERE dataset = $1 AND userid = $2"
	rows, err := tx.Query(q, dataset, handle)
	if err != nil {
		return false, err
	}

	ok := rows.Next()
	rows.Close()
	err = tx.Commit()
	if err != nil {
		return false, err
	}

	return ok, nil
}
