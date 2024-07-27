// Package templates  provides functionality to parse HTML templates.
package templates

import (
	_ "embed"
	"html/template"
)

//go:embed htmlTemplates.tpl
var templateFile string

// ParseTemplates reads and parses the embedded HTML template file.
// It returns a pointer to a template.Template instance and any error encountered.
func ParseTemplates() (*template.Template, error) {
	// Parse the template content into a new template named "metrics"
	t, err := template.New("metrics").Parse(templateFile)
	if err != nil {
		return nil, err
	}
	return t, nil
}
