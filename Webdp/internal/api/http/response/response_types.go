package response

/*
This is mostly for OAS generation.
*/

type Error struct {
	Title  string `json:"title"`
	Detail string `json:"detail"`
	Status int    `json:"status"`
	Type   string `json:"type"`
}

type Id struct {
	Id int64 `json:"id"`
}

type AllFunctions struct {
	Engine1 EngineFunctions `json:"engine1"`
	Engine2 EngineFunctions `json:"engine2"`
}

type EngineFunctions struct {
	Feature1 Function `json:"feature1"`
	Feature2 Function `json:"feature2"`
	Feature3 Function `json:"feature3"`
}

type Function struct {
	Enabled bool `json:"enabled"`
	Opt     struct {
		Optoptions string `json:"options"`
	} `json:"optional_fields"`
	Req struct {
		Reqoptions string `json:"options"`
	} `json:"required_fields"`
}
