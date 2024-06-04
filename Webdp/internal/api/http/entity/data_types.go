package entity

import (
	"encoding/json"
	"fmt"

	errors "webdp/internal/api/http"
)

type DpDataType interface {
	dataType()
	GetName() string
	Valid() error
}

type DataType struct {
	Type DpDataType
}

func (d DataType) Valid() error {
	if d.Type == nil {
		return fmt.Errorf("%w: datatype is nil", errors.ErrBadType)
	}
	return d.Type.Valid()
}

func (d *DataType) UnmarshalJSON(data []byte) error {
	dp, err := UnmarshalDataType(data, json.Unmarshal)
	if err != nil {
		return err
	}
	d.Type = dp
	return nil
}

func (d DataType) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Type)
}

type rawType struct {
	Name   string    `json:"name"`
	Low    *int32    `json:"low,omitempty"`
	High   *int32    `json:"high,omitempty"`
	Labels *[]string `json:"labels,omitempty"`
}

func UnmarshalDataType(data []byte, marhalfunc func([]byte, any) error) (DpDataType, error) {
	temp := rawType{}
	err := marhalfunc(data, &temp)
	if err != nil {
		return nil, err
	}
	if temp.Name == "Int" && temp.Labels == nil && temp.Low != nil && temp.High != nil {
		it := &IntType{Low: *(temp.Low), High: *(temp.High)}
		return it, nil
	}

	if temp.Name == "Double" && temp.Labels == nil && temp.Low != nil && temp.High != nil {
		it := &DoubleType{Low: *temp.Low, High: *temp.High}
		return it, nil
	}

	if temp.Name == "Enum" && temp.Labels != nil && temp.Low == nil && temp.High == nil {
		it := &EnumType{Labels: *temp.Labels}
		return it, nil
	}

	if temp.Name == "Text" && temp.Labels == nil && temp.Low == nil && temp.High == nil {
		return &TextType{}, nil
	}

	if temp.Name == "Bool" && temp.Labels == nil && temp.Low == nil && temp.High == nil {
		return &BoolType{}, nil
	}

	return nil, fmt.Errorf("%w: unrecognized type", errors.ErrBadType)
}

type BoolType struct{}

func (b *BoolType) GetName() string {
	return "Bool"
}

func (b BoolType) Valid() error {
	return nil
}

func (b BoolType) MarshalJSON() ([]byte, error) {
	temp := make(map[string]string)
	temp["name"] = "Bool"
	return json.Marshal(&temp)
}

type TextType struct{}

func (t TextType) Valid() error {
	return nil
}

func (t TextType) GetName() string {
	return "Text"
}

func (b TextType) MarshalJSON() ([]byte, error) {
	temp := make(map[string]string)
	temp["name"] = "Text"
	return json.Marshal(&temp)
}

type IntType struct {
	Low  int32
	High int32
}

func (i IntType) Valid() error {
	if i.Low > i.High {
		return fmt.Errorf("%w: lower bound is larger than the higher bound. low: %d   high: %d", errors.ErrBadInput, i.Low, i.High)
	}
	return nil
}

func (t *IntType) GetName() string {
	return "Int"
}

func (b IntType) MarshalJSON() ([]byte, error) {
	temp := make(map[string]interface{})
	temp["name"] = "Int"
	temp["low"] = b.Low
	temp["high"] = b.High
	return json.Marshal(&temp)
}

type DoubleType struct {
	Low  int32
	High int32
}

func (i DoubleType) Valid() error {
	if i.Low > i.High {
		return fmt.Errorf("%w: lower bound is larger than the higher bound. low: %d   high: %d", errors.ErrBadInput, i.Low, i.High)
	}
	return nil
}

func (t *DoubleType) GetName() string {
	return "Double"
}

func (b DoubleType) MarshalJSON() ([]byte, error) {
	temp := make(map[string]interface{})
	temp["name"] = "Double"
	temp["low"] = b.Low
	temp["high"] = b.High
	return json.Marshal(&temp)
}

type EnumType struct {
	Labels []string
}

func (i EnumType) Valid() error {
	temp := make(map[string]int8)
	for _, label := range i.Labels {
		if _, ok := temp[label]; ok {
			return fmt.Errorf("%w: enum type must have unique labels. the label \"%s\" occurs multiple times", errors.ErrBadFormatting, label)
		} else {
			temp[label] = 0
		}
	}
	return nil
}

func (t *EnumType) GetName() string {
	return "Enum"
}

func (b EnumType) MarshalJSON() ([]byte, error) {
	temp := make(map[string]interface{})
	temp["name"] = "Enum"
	temp["labels"] = b.Labels
	return json.Marshal(&temp)
}

// func (r *rawType) dataType()    {}
func (b *BoolType) dataType()   {}
func (t *TextType) dataType()   {}
func (i *IntType) dataType()    {}
func (d *DoubleType) dataType() {}
func (e *EnumType) dataType()   {}
