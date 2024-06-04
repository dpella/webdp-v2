package entities

type Budget struct {
	Epsilon float64  `json:"epsilon"`
	Delta   *float64 `json:"delta,omitempty"`
}
