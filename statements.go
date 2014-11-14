package cqlmm

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"fmt"
	"hash"
	"io"
	"strings"
)

type Stmt struct {
	Id       string
	Order    int
	HashFunc hash.Hash
	hashed   string
	raw      string
}

func (s *Stmt) Hash() (string, error) {
	if s.hashed != "" {
		return s.hashed, nil
	}

	if s.HashFunc == nil {
		s.HashFunc = sha1.New()
	}

	if _, err := io.WriteString(s.HashFunc, s.raw); err != nil {
		return "", err
	}

	s.hashed = fmt.Sprintf("%x", s.HashFunc.Sum(nil))
	return s.hashed, nil
}

func (s Stmt) String() string {
	return s.raw
}

func parseStatements(r io.Reader) (map[Direction][]Stmt, error) {
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

		// Ignore comments.
		if strings.HasPrefix(ln, "--") {
			continue
		}

		// Skip blank lines.
		if ln == "" {
			continue
		}

		// Write the current line to the buffer, and re-add the newline
		// that was stripped by the earlier call to
		// strings.TrimSpace().
		_, err := buf.WriteString(ln + "\n")
		if err != nil {
			return nil, fmt.Errorf("cqlmm: %v", err)
		}

		// Check to see if the current line ends with a semi-colon. If
		// it does, then dump the buffer, and append it to the list of
		// statements for this section.
		if strings.HasSuffix(ln, ";") {
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

	// Now, create Stmt objects out of the raw statements we have.
	statements := make(map[Direction][]Stmt, 0)
	for dir, rawStmts := range stmts {
		if _, ok := statements[dir]; !ok {
			statements[dir] = make([]Stmt, len(rawStmts))
		}

		for i, s := range rawStmts {
			stmt := Stmt{
				raw:      s,
				Order:    i,
				HashFunc: sha1.New(),
			}
			statements[dir][i] = stmt
		}
	}

	return statements, nil
}
