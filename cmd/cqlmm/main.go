package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"github.com/nesv/cqlmm"
	"github.com/nesv/cqlmm/config"
)

func init() {
	log.SetFlags(0)
}

func main() {
	keyspace := flag.String("k", "", "Apply the migrations to a different keyspace")
	migrationDir := flag.String("c", "db", "Migrations directory path")
	flag.Usage = usage
	flag.Parse()

	configPath := filepath.Join(*migrationDir, "cqlmm.json")

	// We want to catch the "init" subcommand early on.
	switch cmd := flag.Arg(0); cmd {
	case "init":
		if err := cqlmm.InitializeMigrationDir(configPath); err != nil {
			log.Fatalln(err)
		}
		return
	}

	config, err := config.Load(configPath)
	if err != nil {
		log.Fatalln(err)
	}

	if *keyspace == "" {
		keyspace = &config.Keyspace
	}

	// TODO Connect to the database.

	// Take the first, remaining argument as the subcommand to invoke.
	switch cmd := flag.Arg(0); cmd {
	case "up":
		if err := cqlmm.Upgrade("", *migrationDir); err != nil {
			log.Fatalln(err)
		}

	case "down":
		if err := cqlmm.Downgrade("", *migrationDir); err != nil {
			log.Fatalln(err)
		}

	case "create":
		name := flag.Arg(1)
		if name == "" {
			log.Fatalln("no migration name specified")
		}

		switch mtyp := flag.Arg(2); mtyp {
		case "cql":
			pth, err := cqlmm.CreateCQLMigration(*migrationDir, name)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println(pth)

		case "":
			log.Fatalf("unknown migration type %q", mtyp)
		}

	case "":
		// Print usage message
		usage()
		log.Fatalln("")
	default:
		log.Fatalf("ERROR", "unknown subcommand %q", cmd)
	}
}

func usage() {
	log.Println("cqlmm - CQL Migration Manager\n")
	log.Println("Commands:")

	cmds := map[string]string{
		"create": "Create a new migration",
		"down":   "Downgrade the database",
		"init":   "Initialize the migrations directory",
		"up":     "Upgrade the database",
	}
	for cmd, desc := range cmds {
		log.Printf("\t%s\t%s", cmd, desc)
	}

	log.Println("\nFlags:")
	opts := map[string]string{
		"-k=KEYSPACE":   "Change the keyspace the migrations are applied to (overrides config)",
		"-c=CONFIG_DIR": "Path to the migrations directory",
	}
	for opt, desc := range opts {
		log.Printf("\t%s\n\t\t%s", opt, desc)
	}
}
