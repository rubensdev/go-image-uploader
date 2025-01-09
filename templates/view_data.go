package templates

type Metadata map[string]any

type ViewData struct {
	Title string
	Lang  string
	Meta  Metadata
}
