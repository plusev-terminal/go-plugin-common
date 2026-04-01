package utils

import (
	"bytes"
	"fmt"
	"text/template"
)

// TemplateFileToString parses a Go text/template file and executes it with the given data.
// Returns the rendered string.
func TemplateFileToString(templatePath string, data map[string]any) (string, error) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to parse template file %q: %w", templatePath, err)
	}

	var buf bytes.Buffer

	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// TemplateStringToString parses a Go text/template string and executes it with the given data.
// Returns the rendered string.
func TemplateStringToString(templateStr string, data map[string]any) (string, error) {
	tmpl, err := template.New("tpl").Parse(templateStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse template string: %w", err)
	}

	var buf bytes.Buffer

	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
