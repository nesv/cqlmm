package cqlmm

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const cqlmmSectionPrefix = "-- +cql "

type Direction uint8

const (
	Up Direction = iota + 1
	Down
)

type MigrationState uint8

const (
	Pending MigrationState = iota + 1
	Partial
	Applied
)

type Migration struct {
	Name     string
	Version  int64
	Next     int64
	Previous int64
	Source   string
	State    MigrationState

	stmts map[Direction][]Stmt
}

func (m *Migration) Stmts(d Direction) []Stmt {
	return m.stmts[d]
}

func ParseMigration(pth string) (*Migration, error) {
	// Open the migration file.
	f, err := os.Open(pth)
	if err != nil {
		return nil, fmt.Errorf("cqlmm: failed to open file %q: %v", pth, err)
	}
	defer f.Close()

	// Parse the statements within the migration file.
	statements, err := parseStatements(f)

	// Parse the filename to get the version, and the name of the
	// migration.
	parts := strings.Split(filepath.Base(pth), "_")
	version, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("cqlmm: %v", err)
	}

	name := strings.Split(parts[1], ".")[0]

	migration := Migration{
		Name:    name,
		Version: version,
		Source:  pth,
		stmts:   statements,
	}

	return &migration, nil
}

func LoadMigrations(basedir string) ([]*Migration, error) {
	matches, err := filepath.Glob(filepath.Join(basedir, "migrations", "*.cql"))
	if err != nil {
		return nil, fmt.Errorf("cqlmm: %v", err)
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("no migrations to run")
	}

	if !sort.StringsAreSorted(matches) {
		sort.Strings(matches)
	}

	migrations := make([]*Migration, len(matches))
	for i, match := range matches {
		log.Println("loading migration:", filepath.Base(match))
		migration, err := ParseMigration(match)
		if err != nil {
			return nil, fmt.Errorf("cqlmm: %v", err)
		}

		migrations[i] = migration
	}

	// Set the Next and Previous values for each migration.
	for i := 0; i < len(migrations); i++ {
		if i+1 != len(migrations) {
			migrations[i].Next = migrations[i+1].Version
		}
		migrations[i].Next = 0

		if i == 0 {
			continue
		}
		migrations[i].Previous = migrations[i-1].Version
	}

	return migrations, nil
}

func RunMigrations(m []*Migration, d Direction) error {
	switch d {
	case Up:
		if len(m) == 1 {
			log.Printf("applying migration %d", m[0].Version)
		} else if len(m) > 1 {
			log.Printf("applying migrations %d..%d", m[0].Version, m[len(m)-1].Version)
		}

		for _, migration := range m {
			log.Printf("%d\t%s", migration.Version, migration.Name)
			for i, stmt := range migration.Stmts(d) {
				log.Printf("\t%d\t%s", i, stmt)
			}
		}

	case Down:
		// TODO
		// This just naively pulls the latest migration script off of
		// the pile, and runs its "down" section.
		//
		// What it *should* do is hit the database, and see which was
		// the most-recent migration that was applied, then run the
		// "down" section from that migration.
		migration := m[len(m)-1]
		log.Printf("rolling back migration %d", migration.Version)
		log.Printf("%d\t%s", migration.Version, migration.Name)
		for i, stmt := range migration.Stmts(d) {
			log.Printf("\t%d\t%s", i, stmt)
		}
	}

	return nil
}
