package services

import "webdp/internal/api/http/repo/postgres"

type InternalDatasetService struct {
	repo postgres.DatasetPostgres
}

func NewInternalDatasetService(repo postgres.DatasetPostgres) *InternalDatasetService {
	return &InternalDatasetService{repo: repo}
}

func (ds InternalDatasetService) GetTable(datasetid int64) ([]byte, error) {
	return ds.repo.GetUploadedData(datasetid)
}
