package requests

import "googledp/entities"

type Validate struct {
	Budget        entities.Budget `json:"budget"`
	Query         []interface{}   `json:"query"`
	Dataset       int64           `json:"dataset"`
	Schema        []interface{}   `json:"schema"`
	PrivacyNotion string          `json:"privacy_notion"`
	CallbackUrl   string          `json:"url"`
}
