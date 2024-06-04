package entity

import (
	"fmt"
	"time"
	errors "webdp/internal/api/http"
	"webdp/internal/api/http/utils"
)

const (
	PURE   = "PureDP"
	APPROX = "ApproxDP"
)

type DatasetInfo struct {
	Id            int64          `json:"id"`
	Name          string         `json:"name"`
	Owner         string         `json:"owner"`
	Schema        []ColumnSchema `json:"schema"`
	PrivacyNotion string         `json:"privacy_notion"`
	TotalBudget   Budget         `json:"total_budget"`
	Loaded        bool           `json:"loaded"`
	CreatedOn     time.Time      `json:"created_time,omitempty"`
	UpdatedOn     time.Time      `json:"updated_time,omitempty"`
	LoadedOn      time.Time      `json:"loaded_time,omitempty"`
}

type DatasetCreate struct {
	Name          string         `json:"name" dpvalidation:"non-empty-string"`
	Owner         string         `json:"owner" dpvalidation:"non-empty-string"`
	Schema        []ColumnSchema `json:"schema"`
	PrivacyNotion string         `json:"privacy_notion" dpvalidation:"non-empty-string"`
	TotalBudget   Budget         `json:"total_budget"`
}

type DatasetPatch struct {
	Name        string `json:"name" dpvalidation:"non-empty-string"`
	Owner       string `json:"owner" dpvalidation:"non-empty-string"`
	TotalBudget Budget `json:"total_budget"`
}

type ColumnSchema struct {
	Name string   `json:"name" dpvalidation:"non-empty-string"`
	Type DataType `json:"type"`
}

func (c ColumnSchema) Valid() error {
	err := utils.ValidateNonEmptyString(c)

	if err != nil {
		return err
	}
	return c.Type.Valid()
}

func (d DatasetCreate) Valid() error {
	err := utils.ValidateNonEmptyString(d)
	if err != nil {
		return err
	}

	if d.PrivacyNotion != PURE && d.PrivacyNotion != APPROX {
		return fmt.Errorf("%w: privacy notion should be either \"%s\" or \"%s\"", errors.ErrBadFormatting, PURE, APPROX)
	}

	err = d.TotalBudget.Valid()
	if err != nil {
		return err
	}

	for _, cs := range d.Schema {
		err := cs.Valid()
		if err != nil {
			return err
		}
	}
	return nil
}

func (d DatasetPatch) Valid() error {
	err := utils.ValidateNonEmptyString(d)
	if err != nil {
		return err
	}
	return d.TotalBudget.Valid()
}
