package game

type Mod struct {
	Name        string `json:"name"`
	Author      string `json:"author"`
	Description string `json:"description"`
	Version     string `json:"version"`
	URL         string `json:"url"`
}
