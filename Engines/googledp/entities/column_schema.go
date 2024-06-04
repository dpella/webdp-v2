package entities

import (
	"encoding/json"
	"fmt"
)

type Column struct {
	Name string  `json:"name"`
	Type ColType `json:"type"`
}

func (c *Column) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Name     string          `json:"name"`
		TypeData json.RawMessage `json:"type"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	c.Name = tmp.Name

	var typeInfo struct {
		TypeName string `json:"name"`
	}

	if err := json.Unmarshal(tmp.TypeData, &typeInfo); err != nil {
		return err
	}

	switch typeInfo.TypeName {
	case "Int":
		var intType IntType
		if err := json.Unmarshal(tmp.TypeData, &intType); err != nil {
			return err
		}
		c.Type = &intType
	case "Double":
		var doubleType DoubleType
		if err := json.Unmarshal(tmp.TypeData, &doubleType); err != nil {
			return err
		}
		c.Type = &doubleType
	case "Enum":
		var enumType EnumType
		if err := json.Unmarshal(tmp.TypeData, &enumType); err != nil {
			return err
		}
		c.Type = &enumType
	case "Text":
		var textType StringType
		if err := json.Unmarshal(tmp.TypeData, &textType); err != nil {
			return err
		}
		c.Type = &textType
	default:
		return fmt.Errorf("unknown type: %s", tmp.Name)
	}

	return nil
}

type ColType interface {
	GetName() string
	GetLow() int64
	GetHigh() int64
	GetLabels() []string
}

type IntType struct {
	Name string `json:"name"`
	Low  int64  `json:"low"`
	High int64  `json:"high"`
}

func (i IntType) GetName() string {
	return i.Name
}

func (i IntType) GetLow() int64 {
	return i.Low
}

func (i IntType) GetHigh() int64 {
	return i.High
}

func (i IntType) GetLabels() []string {
	return nil
}

type DoubleType struct {
	Name string `json:"name"`
	Low  int64  `json:"low"`
	High int64  `json:"high"`
}

func (i DoubleType) GetName() string {
	return i.Name
}

func (i DoubleType) GetLow() int64 {
	return i.Low
}

func (i DoubleType) GetHigh() int64 {
	return i.High
}

func (i DoubleType) GetLabels() []string {
	return nil
}

type EnumType struct {
	Name   string   `json:"name"`
	Labels []string `json:"labels"`
}

func (i EnumType) GetName() string {
	return i.Name
}

func (i EnumType) GetLow() int64 {
	return 0
}

func (i EnumType) GetHigh() int64 {
	return 0
}

func (i EnumType) GetLabels() []string {
	return nil
}

type StringType struct {
	Name string `json:"name"`
}

func (i StringType) GetName() string {
	return i.Name
}

func (i StringType) GetLow() int64 {
	return 0
}

func (i StringType) GetHigh() int64 {
	return 0
}

func (i StringType) GetLabels() []string {
	return nil
}
