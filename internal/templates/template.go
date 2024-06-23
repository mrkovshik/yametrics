// Package templates  provides functionality to parse HTML templates.
package templates

import (
	"embed"
	"html/template"
	"log"
)

//go:embed htmlTemplates.tpl
var templateFile embed.FS

// ParseTemplates reads and parses the embedded HTML template file.
// It returns a pointer to a template.Template instance and any error encountered.
func ParseTemplates() (*template.Template, error) {
	// Read the contents of the embedded template file
	content, err := templateFile.ReadFile("htmlTemplates.tpl")
	if err != nil {
		log.Fatalf("Error reading template file: %v", err)
	}

	// Parse the template content into a new template named "metrics"
	t, err := template.New("metrics").Parse(string(content))
	if err != nil {
		return nil, err
	}
	return t, nil
}
