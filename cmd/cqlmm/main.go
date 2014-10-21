package main

import (
	"flag"
	"log"

	"github.com/nesv/cqlmm/config"
)

func init() {
	log.SetFlags(0)
}

func main() {
	keyspace := flag.String("k", "", "Apply the migrations to a different keyspace")
	configPath := flag.String("c", "db/cqlmm.json", "Path to config file")
	flag.Usage = usage
	flag.Parse()

	config, err := config.Load(*configPath)
	if err != nil {
		log.Fatalln(err)
	}

	// Take the first, remaining argument as the subcommand to invoke.
	switch cmd := flag.Arg(0); cmd {
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
		"-k=KEYSPACE": "Change the keyspace the migrations are applied to (overrides config)",
		"-c=CONFIG": "Path to the configuration file"
	}
	for opt, desc := range opts {
		log.Printf("\t%s\n\t\t%s", opt, desc)
	}
}
