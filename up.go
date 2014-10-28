package cqlmm

func Upgrade(dbURL, migrationDir string) error {
	migrations, err := LoadMigrations(migrationDir)
	if err != nil {
		return err
	}

	return RunMigrations(migrations, Up)
}
