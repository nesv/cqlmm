package cqlmm

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

var migrationTmpl = `-- +cql up
-- CQL commands in section "up" are executed when this migration is applied


-- +cql down
-- CQL commands in section "down" are executed when this migration is rolled back

`

func CreateCQLMigration(migrationDir, name string) (string, error) {
	// Get the current time, and prepend it to the name of the migration
	// file we are going to create.
	now := time.Now()
	filename := fmt.Sprintf("%d_%s.cql", now.Unix(), name)
	pth := filepath.Join(migrationDir, "migrations", filename)

	// Now, create a new migration file, and fill it in with our template.
	f, err := os.Create(pth)
	if err != nil {
		return pth, fmt.Errorf("create: failed to create migration file %s: %v", pth, err)
	}
	defer f.Close()

	tmpl, err := template.New("migration").Parse(migrationTmpl)
	if err != nil {
		return pth, fmt.Errorf("create: failed to parse migration template: %v", err)
	}

	err = tmpl.ExecuteTemplate(f, "migration", nil)
	return pth, nil
}
