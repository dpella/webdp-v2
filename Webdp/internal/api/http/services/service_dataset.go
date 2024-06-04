package services

import (
	"fmt"
	"strconv"
	errors "webdp/internal/api/http"
	"webdp/internal/api/http/entity"
	"webdp/internal/api/http/repo/postgres"
)

type DatasetService struct {
	postg      postgres.DatasetPostgres
	budgetRepo postgres.BudgetPostgres
}

func NewDatasetService(datasetRepo postgres.DatasetPostgres, budgetRepo postgres.BudgetPostgres) DatasetService {
	return DatasetService{postg: datasetRepo, budgetRepo: budgetRepo}
}

func (d DatasetService) GetAllDatasets() ([]entity.DatasetInfo, error) {
	data, err := d.postg.GetDatasets()
	if err != nil {
		return []entity.DatasetInfo{}, errors.WrapDBError(err, "get datasets", "all")
	}
	return data, nil
}

func (d DatasetService) GetDataset(id int64) (entity.DatasetInfo, error) {
	dataset, err := d.postg.GetDataset(id)
	if err != nil {
		return entity.DatasetInfo{}, errors.WrapDBError(err, "get", strconv.FormatInt(id, 10))
	}
	return dataset, nil
}

func (d DatasetService) GetDatasetOwner(id int64) (string, error) {
	dataset, err := d.GetDataset(id)
	if err != nil {
		return "", err
	}
	return dataset.Owner, nil
}

func (d DatasetService) CreateDataset(ds entity.DatasetCreate) (int64, error) {
	id, err := d.postg.CreateDataset(ds)
	if err != nil {
		return 0, errors.WrapDBError(err, "create", strconv.FormatInt(id, 10))
	}
	return id, nil
}

func (d DatasetService) UpdateDataset(datasetId int64, patch entity.DatasetPatch) error {
	allocs, err := d.budgetRepo.GetDatasetUserAllocations(datasetId)
	if err != nil {
		return errors.WrapDBError(err, "update", strconv.FormatInt(datasetId, 10))
	}

	if allocs.Allocated.Epsilon > patch.TotalBudget.Epsilon || coalesce[float64](allocs.Allocated.Delta) > coalesce[float64](patch.TotalBudget.Delta) {
		return fmt.Errorf("%w: you cannot set a lower budget than what has already been allocated", errors.ErrBadRequest)
	}

	if err := d.postg.UpdateDataset(datasetId, patch); err != nil {
		return errors.WrapDBError(err, "update", strconv.FormatInt(datasetId, 10))
	}
	return nil
}

func (d DatasetService) DeleteDataset(id int64) error {
	if err := d.postg.DeleteDataset(id); err != nil {
		return errors.WrapDBError(err, "delete", strconv.FormatInt(id, 10))
	}
	return nil
}

func (d DatasetService) UploadData(datasetid int64, data []byte) error {
	err := d.postg.UploadData(datasetid, data)
	if err != nil {
		return errors.WrapDBError(err, "upload", strconv.FormatInt(datasetid, 10))
	}
	return nil
}
