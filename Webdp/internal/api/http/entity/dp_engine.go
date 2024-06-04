package entity

type URL = string

type EnginesConfig struct {
	Default string              `json:"default"`
	Engines []WebDPClientTarget `json:"engines"`
}

type WebDPClientTarget struct {
	Name                     string `json:"name"`
	EndpointEvaluate         URL    `json:"evaluate_url"`
	EndpointAccuracy         URL    `json:"accuracy_url"`
	EndpointClearSingleCache URL    `json:"delete_url"`
	EndpointValidate         URL    `json:"validation_url"`
	EndpointFunctions        URL    `json:"functions_url"`
	EndpointDocs             URL    `json:"documentation_url"`
}
