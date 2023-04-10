package game

type Artifact struct {
	ID          string
	Name        string
	Description string
	Order       int
	Price       int
	Callbacks   map[string]OwnedCallback
}

type ArtifactInstance struct {
	TypeID string
	GUID   string
	Owner  string
}
