package requests

import "googledp/entities"

type Accuracy struct {
	Budget      entities.Budget `json:"budget"`
	Query       []interface{}   `json:"query"`
	Schema      []interface{}   `json:"schema"`
	CallbackUrl string          `json:"url"`
}
