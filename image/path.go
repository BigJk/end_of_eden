package image

var searchPaths []string

func init() {
	ResetSearchPaths()
}

// AddSearchPaths adds additional search paths for images.
func AddSearchPaths(paths ...string) {
	searchPaths = append(searchPaths, paths...)
}

// ResetSearchPaths resets the search paths to the default.
func ResetSearchPaths() {
	searchPaths = []string{"./assets/images/"}
}
