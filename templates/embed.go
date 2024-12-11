package templates

import (
	"embed"
	"fmt"
	"github.com/samber/do"
	"html/template"
)

//go:embed html/*
var HTMLTemplates embed.FS

func NewHtmlTemplates(_ *do.Injector) (*template.Template, error) {
	template, err := template.ParseFS(HTMLTemplates, "html/**/*.html")
	if err != nil {
		return nil, fmt.Errorf("error parsing html templates: %w", err)
	}

	return template, nil
}
