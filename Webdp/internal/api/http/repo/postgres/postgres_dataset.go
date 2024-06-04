package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"webdp/internal/api/http/entity"

	errors "webdp/internal/api/http"

	"github.com/lib/pq"
)

type DatasetPostgres struct {
	db *sql.DB
}

type rawtype struct {
	name   string
	low    sql.NullInt32
	high   sql.NullInt32
	labels []string
}

func NewDatasetPostgres(conn *sql.DB) DatasetPostgres {
	return DatasetPostgres{db: conn}
}

func (d DatasetPostgres) GetDataset(datasetId int64) (entity.DatasetInfo, error) {
	tx, err := d.db.BeginTx(context.Background(), nil)
	if err != nil {
		return entity.DatasetInfo{}, err
	}

	defer dfun(err, tx)
	q := "SELECT id, name, owner, privacy_notion, total_epsilon, total_delta, loaded, created_time, updated_time, loaded_time FROM LoadedDatasets WHERE id = $1"

	row := tx.QueryRow(q, datasetId)
	var id int64
	var name, owner, privn string
	var ct, ut time.Time
	var lt pq.NullTime
	var loaded bool
	var eps, del float64
	err = row.Scan(&id, &name, &owner, &privn, &eps, &del, &loaded, &ct, &ut, &lt)

	if err != nil {
		return entity.DatasetInfo{}, errors.ErrNotFound
	}
	dinfo := &entity.DatasetInfo{
		Id:            id,
		Name:          name,
		Owner:         owner,
		PrivacyNotion: privn,
		TotalBudget:   entity.Budget{Epsilon: eps, Delta: &del},
		CreatedOn:     ct,
		UpdatedOn:     ut,
		Loaded:        loaded,
		LoadedOn:      lt.Time,
	}

	cs, err := getColumnSchema(tx, datasetId)
	if err != nil {
		return entity.DatasetInfo{}, err
	}

	dinfo.Schema = cs
	tx.Commit()
	return *dinfo, nil

}

func (d DatasetPostgres) GetDatasets() ([]entity.DatasetInfo, error) {
	tx, err := d.db.BeginTx(context.Background(), nil)
	if err != nil {
		return []entity.DatasetInfo{}, err
	}

	defer dfun(err, tx)
	q := "SELECT id, name, owner, privacy_notion, total_epsilon, total_delta, loaded, created_time, updated_time, loaded_time FROM LoadedDatasets"
	rs, err := tx.Query(q)
	if err != nil {
		return []entity.DatasetInfo{}, err
	}
	temp := make([]*entity.DatasetInfo, 0)
	for rs.Next() {
		var id int64
		var name, owner, privn string
		var ct, ut time.Time
		var lt pq.NullTime
		var loaded bool
		var eps, del float64

		err = rs.Scan(&id, &name, &owner, &privn, &eps, &del, &loaded, &ct, &ut, &lt)
		if err != nil {
			return []entity.DatasetInfo{}, err
		}
		temp = append(temp, &entity.DatasetInfo{
			Id:            id,
			Name:          name,
			Owner:         owner,
			PrivacyNotion: privn,
			TotalBudget:   entity.Budget{Epsilon: eps, Delta: &del},
			CreatedOn:     ct,
			UpdatedOn:     ut,
			Loaded:        loaded,
			LoadedOn:      lt.Time,
		})
	}
	rs.Close()
	out := make([]entity.DatasetInfo, 0)
	for _, di := range temp {
		cs, err := getColumnSchema(tx, di.Id)
		if err != nil {
			return []entity.DatasetInfo{}, err
		}
		di.Schema = cs
		out = append(out, *di)
	}

	if err := tx.Commit(); err != nil {
		return []entity.DatasetInfo{}, err
	}

	return out, nil
}

func (d DatasetPostgres) CreateDataset(dataset entity.DatasetCreate) (int64, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return 0, err
	}
	defer dfun(err, tx)

	var id int64

	created := time.Now().UTC()
	q := "INSERT INTO Dataset (name, owner, privacy_notion, total_epsilon, total_delta, created_time, updated_time) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
	err = tx.QueryRow(q,
		dataset.Name,
		dataset.Owner,
		dataset.PrivacyNotion,
		dataset.TotalBudget.Epsilon,
		dataset.TotalBudget.Delta,
		created,
		created).Scan(&id)

	if err != nil {
		return 0, err
	}

	for _, cs := range dataset.Schema {
		err = insertColumnSchema(tx, id, cs)
		if err != nil {
			return 0, err
		}
	}

	tx.Commit()
	return id, nil
}

func (d DatasetPostgres) UpdateDataset(dataset int64, patch entity.DatasetPatch) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}

	defer dfun(err, tx)
	updated := time.Now().UTC()
	q := "UPDATE Dataset SET name = $1, owner = $2, total_epsilon = $3, total_delta = $4, updated_time = $5 WHERE id = $6"

	_, err = tx.Exec(q, patch.Name, patch.Owner, patch.TotalBudget.Epsilon, patch.TotalBudget.Delta, updated, dataset)
	if err != nil {
		return errors.ErrNotFound
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil

}

func (d DatasetPostgres) UploadData(dataset int64, data []byte) error {
	ltime := time.Now()
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}

	defer dfun(err, tx)

	q := "INSERT INTO DataUpload (dataset, loaded_time, loaded_data) VALUES ($1, $2, $3)"
	_, err = tx.Exec(q, dataset, ltime, data)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (d DatasetPostgres) GetUploadedData(dataset int64) ([]byte, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return []byte{}, err
	}
	q := "SELECT loaded_data FROM DataUpload WHERE dataset = $1"
	row := tx.QueryRow(q, dataset)
	var out []byte
	err = row.Scan(&out)

	if err != nil {
		return []byte{}, errors.ErrNotFound
	}

	if err = tx.Commit(); err != nil {
		return []byte{}, err
	}
	return out, nil
}

func (d DatasetPostgres) DeleteDataset(dataset int64) error {
	q := "DELETE FROM Dataset WHERE id = $1"
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}

	defer dfun(err, tx)

	_, err = tx.Exec(q, dataset)
	if err != nil {
		return errors.ErrNotFound
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func insertColumnSchema(tx *sql.Tx, dataset int64, cs entity.ColumnSchema) error {
	q := ""
	var vars []any
	switch t := cs.Type.Type.(type) {
	case *entity.IntType:
		q = "INSERT INTO ColumnSchemas (dataset, column_name, data_type, low, high) VALUES ($1, $2, $3, $4, $5)"
		vars = []any{dataset, cs.Name, "Int", t.Low, t.High}
	case *entity.DoubleType:
		q = "INSERT INTO ColumnSchemas (dataset, column_name, data_type, low, high) VALUES ($1, $2, $3, $4, $5)"
		vars = []any{dataset, cs.Name, "Double", t.Low, t.High}
	case *entity.BoolType:
		q = "INSERT INTO ColumnSchemas (dataset, column_name, data_type) VALUES ($1, $2, $3)"
		vars = []any{dataset, cs.Name, "Bool"}
	case *entity.TextType:
		q = "INSERT INTO ColumnSchemas (dataset, column_name, data_type) VALUES ($1, $2, $3)"
		vars = []any{dataset, cs.Name, "Text"}
	case *entity.EnumType:
		q = "INSERT INTO ColumnSchemas (dataset, column_name, data_type, labels) VALUES ($1, $2, $3, $4)"
		vars = []any{dataset, cs.Name, "Enum", pq.Array(t.Labels)}
	default:
		return fmt.Errorf("%w: invalid data type", errors.ErrBadType)
	}
	_, err := tx.Exec(q, vars...)
	return err

}

func getColumnSchema(tx *sql.Tx, dataset int64) ([]entity.ColumnSchema, error) {
	q := "SELECT column_name, data_type, low, high, labels FROM ColumnSchemas WHERE dataset = $1"
	rs, err := tx.Query(q, dataset)
	if err != nil {
		return []entity.ColumnSchema{}, err
	}

	out := make([]entity.ColumnSchema, 0)
	for rs.Next() {
		raw := rawtype{}
		col := ""
		err = rs.Scan(&col, &raw.name, &raw.low, &raw.high, pq.Array(&raw.labels))
		if err != nil {
			return []entity.ColumnSchema{}, err
		}
		dp, err := fromRaw(raw)
		if err != nil {
			return []entity.ColumnSchema{}, err
		}
		out = append(out, entity.ColumnSchema{Name: col, Type: dp})
	}
	rs.Close()
	return out, nil
}

func fromRaw(r rawtype) (entity.DataType, error) {
	switch r.name {
	case "Int":
		return entity.DataType{Type: &entity.IntType{Low: r.low.Int32, High: r.high.Int32}}, nil
	case "Double":
		return entity.DataType{Type: &entity.DoubleType{Low: r.low.Int32, High: r.high.Int32}}, nil
	case "Bool":
		return entity.DataType{Type: &entity.BoolType{}}, nil
	case "Text":
		return entity.DataType{Type: &entity.TextType{}}, nil
	case "Enum":
		return entity.DataType{Type: &entity.EnumType{Labels: r.labels}}, nil
	default:
		return entity.DataType{}, fmt.Errorf("%w: invalid datatype", errors.ErrBadType)
	}
}
