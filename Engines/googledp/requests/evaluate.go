package requests

import (
	"googledp/entities"
)

type Evaluate struct {
	Budget        entities.Budget   `json:"budget"`
	Query         entities.Query    `json:"query"`
	Dataset       int64             `json:"dataset"`
	Schema        []entities.Column `json:"schema"`
	PrivacyNotion string            `json:"privacy_notion"`
	CallbackUrl   string            `json:"url"`
}
