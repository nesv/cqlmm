package cqlmm

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
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

func parseStatements(r io.Reader) (map[Direction][]string, error) {
	var buf bytes.Buffer
	scanner := bufio.NewScanner(r)
	stmts := make(map[Direction][]string, 0)

	upSections, downSections := 0, 0
	var dir Direction
	for scanner.Scan() {
		ln := scanner.Text()

		if strings.HasPrefix(ln, cqlmmSectionPrefix) {
			section := strings.TrimSpace(ln[len(cqlmmSectionPrefix):])
			switch section {
			case "up":
				if upSections > 0 {
					return nil, fmt.Errorf("cqlmm: too many up sections in migration")
				}
				upSections++
				dir = Up

			case "down":
				if downSections > 0 {
					return nil, fmt.Errorf("cqlmm: too many down sections in migration")
				}
				downSections++
				dir = Down

			default:
				return nil, fmt.Errorf("cqlmm: bad section name %q", section)
			}
			continue
		}

		if _, ok := stmts[dir]; !ok {
			stmts[dir] = make([]string, 0)
		}

		if strings.HasPrefix(ln, "--") {
			// Ignore comments.
			continue
		}

		_, err := buf.WriteString(ln + "\n")
		if err != nil {
			return nil, fmt.Errorf("cqlmm: %v", err)
		}

		// Check to see if the current line ends with a semi-colon. If
		// it does, then dump the buffer, and append it to the list of
		// statements for this section.
		if strings.HasSuffix(ln, ";\n") {
			stmt := buf.Bytes()
			stmts[dir] = append(stmts[dir], string(stmt))
			buf.Reset()
		}
	}

	if upSections == 0 {
		return nil, fmt.Errorf("cqlmm: no up section in migration file")
	}

	if downSections == 0 {
		return nil, fmt.Errorf("cqlmm: no down section in migration file")
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("cqlmm: %v", err)
	}

	return stmts, nil
}

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

	stmts map[Direction][]string
}

func (m *Migration) Stmts(d Direction) []string {
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

	return nil
}
