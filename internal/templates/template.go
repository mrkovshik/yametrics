package templates

import (
	"embed"
	"html/template"
	"log"
)

//go:embed htmlTemplates.tpl
var templateFile embed.FS

func ParseTemplates() (*template.Template, error) {

	content, err := templateFile.ReadFile("htmlTemplates.tpl")
	if err != nil {
		log.Fatalf("Error reading template file: %v", err)
	}
	t, err := template.New("metrics").Parse(string(content))
	if err != nil {
		return &template.Template{}, err
	}
	return t, nil
}
