package ludoc

type Category struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Order       int    `json:"order"`
}

type Function struct {
	Name        string   `json:"name"`
	Args        []string `json:"args"`
	Return      string   `json:"return"`
	Category    string   `json:"category"`
	Description string   `json:"description"`
}

type Global struct {
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description"`
}

type Docs struct {
	curCat string

	Categories map[string]Category `json:"categories"`
	Globals    map[string]Global   `json:"globals"`
	Functions  map[string]Function `json:"functions"`
}

func New() *Docs {
	return &Docs{
		Categories: map[string]Category{},
		Globals:    map[string]Global{},
		Functions:  map[string]Function{},
	}
}

func (d *Docs) Function(name string, description string, ret string, args ...string) {
	d.Functions[name] = Function{
		Name:        name,
		Args:        args,
		Return:      ret,
		Description: description,
		Category:    d.curCat,
	}
}

func (d *Docs) Global(name string, description string) {
	d.Globals[name] = Global{
		Name:        name,
		Description: description,
		Category:    d.curCat,
	}
}

func (d *Docs) Category(name string, description string, order int) {
	d.Categories[name] = Category{
		Name:        name,
		Description: description,
		Order:       order,
	}
	d.PushCategory(name)
}

func (d *Docs) PushCategory(name string) {
	d.curCat = name
}

func (d *Docs) PopCategory() {
	d.curCat = ""
}
