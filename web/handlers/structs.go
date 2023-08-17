package handlers

type Error struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

type ProjectInfo struct {
	Name   string    `json:"name"`
	Author string    `json:"author"`
	Build  BuildInfo `json:"build"`
}

type BuildInfo struct {
	Build       string `json:"build"`
	Commit      string `json:"commit"`
	Branch      string `json:"branch"`
	Environment string `json:"environment"`
	Url         string `json:"url"`
}
