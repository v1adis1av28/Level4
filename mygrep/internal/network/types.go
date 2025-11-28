package network

type GrepRequest struct {
	Pattern string    `json:"pattern"`
	Data    string    `json:"data"`
	Flags   GrepFlags `json:"flags"`
}

type GrepFlags struct {
	Ignore      bool `json:"ignore"`
	Invert      bool `json:"invert"`
	Strict      bool `json:"strict"`
	LineNumbers bool `json:"line_numbers"`
	CountOnly   bool `json:"count_only"`
	After       int  `json:"after"`
	Before      int  `json:"before"`
	Around      int  `json:"around"`
}

type GrepResponse struct {
	Lines []string `json:"lines"`
	Count int      `json:"count"`
	Error string   `json:"error"`
}
