package structs

type MESSAGE struct {
	Channel   int64     `json:"channel"`
	Message   string    `json:"message"`
	Image     string    `json:"image"`
	EmbedRows [][]EMBED `json:"embedRows"`
	Markdown  string    `json:"markdown"`
}

type EMBED struct {
	Name   string `json:"name"`
	Action string `json:"action"`
	Type   int    `json:"type"`
}
