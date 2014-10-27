package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

var configTmpl = `{
        "hosts": ["cql://user:password@127.0.0.1:9042"],
        "keyspace": "my_keyspace"
}
`

func InitializeMigrationDir(cfgPath string) error {
	dir := filepath.Dir(cfgPath)
	if dir == "" {
		return errors.New("init: migration directory is not specified")
	}

	// Create the migrations directory.
	migDir := filepath.Join(dir, "migrations")
	if err := os.MkdirAll(migDir, 0755); err != nil {
		return fmt.Errorf("init: failed to create directories: %v", err)
	}

	// Generate a starter cqlmm.json config, from a template.
	tmpl, err := template.New("config").Parse(configTmpl)
	if err != nil {
		return fmt.Errorf("init: failed to parse config template: %v", err)
	}

	f, err := os.Create(cfgPath)
	if err != nil {
		return fmt.Errorf("init: failed to open file %s for writing: %v", cfgPath, err)
	}
	defer f.Close()

	if err := tmpl.ExecuteTemplate(f, "config", nil); err != nil {
		return fmt.Errorf("init: failed to execute config template: %v", err)
	}

	return nil
}
