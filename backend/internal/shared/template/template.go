package template

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sync"
)

type TemplateParser struct {
	templatesDir string
	cache        map[string]*template.Template
	mu           sync.RWMutex
}

func NewTemplateParser(templatesDir string) *TemplateParser {
	return &TemplateParser{
		templatesDir: templatesDir,
		cache:        make(map[string]*template.Template),
	}
}

// Parse loads and parses a template file, then executes it with the given context
func (tp *TemplateParser) Parse(templateName string, context map[string]interface{}) (string, error) {
	tmpl, err := tp.getTemplate(templateName)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, context); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", templateName, err)
	}

	return buf.String(), nil
}

func (tp *TemplateParser) getTemplate(templateName string) (*template.Template, error) {
	tp.mu.RLock()
	tmpl, exists := tp.cache[templateName]
	tp.mu.RUnlock()

	if exists {
		return tmpl, nil
	}

	return tp.loadTemplate(templateName)
}

func (tp *TemplateParser) loadTemplate(templateName string) (*template.Template, error) {
	tp.mu.Lock()
	defer tp.mu.Unlock()

	// Double-check after acquiring write lock
	if tmpl, exists := tp.cache[templateName]; exists {
		return tmpl, nil
	}

	templatePath := filepath.Join(tp.templatesDir, templateName)

	content, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template %s: %w", templateName, err)
	}

	tmpl, err := template.New(templateName).Parse(string(content))
	if err != nil {
		return nil, fmt.Errorf("failed to parse template %s: %w", templateName, err)
	}

	tp.cache[templateName] = tmpl
	return tmpl, nil
}

func (tp *TemplateParser) ClearCache() {
	tp.mu.Lock()
	defer tp.mu.Unlock()
	tp.cache = make(map[string]*template.Template)
}
