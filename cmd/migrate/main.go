package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rizkiar00/homework/pkg/config"
)

const defaultMigrationPath = "migrations/db"

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if len(os.Args) < 2 {
		printUsage()
		return nil
	}

	command := os.Args[1]
	flags := flag.NewFlagSet(command, flag.ExitOnError)
	migrationPath := flags.String("path", defaultMigrationPath, "migration files path")
	steps := flags.Int("steps", 0, "number of migration steps for up/down")
	version := flags.Int("version", 0, "version used by force")
	if err := flags.Parse(os.Args[2:]); err != nil {
		return err
	}

	cfg, err := config.New()
	if err != nil {
		return err
	}

	if !cfg.Database.IsConfigured() {
		return errors.New("database config is not complete")
	}

	migrator, err := migrate.New("file://"+*migrationPath, cfg.Database.GenerateConnectionString())
	if err != nil {
		return err
	}
	defer migrator.Close()

	switch command {
	case "up":
		return runUp(migrator, *steps)
	case "down":
		return runDown(migrator, *steps)
	case "force":
		return runForce(migrator, *version)
	case "version":
		return runVersion(migrator)
	default:
		printUsage()
		return fmt.Errorf("unknown migration command: %s", command)
	}
}

func runUp(migrator *migrate.Migrate, steps int) error {
	var err error
	if steps <= 0 {
		err = migrator.Up()
	} else {
		err = migrator.Steps(steps)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		fmt.Println("migration already up to date")
		return nil
	}
	if err != nil {
		return err
	}

	fmt.Println("migration up completed")
	return nil
}

func runDown(migrator *migrate.Migrate, steps int) error {
	if steps <= 0 {
		steps = 1
	}

	var err error
	err = migrator.Steps(-steps)

	if errors.Is(err, migrate.ErrNoChange) {
		fmt.Println("migration already at base version")
		return nil
	}
	if err != nil {
		return err
	}

	fmt.Println("migration down completed")
	return nil
}

func runForce(migrator *migrate.Migrate, version int) error {
	if version < 0 {
		return errors.New("version must be greater than or equal to 0")
	}

	if err := migrator.Force(version); err != nil {
		return err
	}

	fmt.Printf("migration forced to version %d\n", version)
	return nil
}

func runVersion(migrator *migrate.Migrate) error {
	version, dirty, err := migrator.Version()
	if errors.Is(err, migrate.ErrNilVersion) {
		fmt.Println("migration version: none")
		return nil
	}
	if err != nil {
		return err
	}

	fmt.Printf("migration version: %d dirty: %t\n", version, dirty)
	return nil
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "usage:")
	fmt.Fprintln(os.Stderr, "  go run ./cmd/migrate up")
	fmt.Fprintln(os.Stderr, "  go run ./cmd/migrate up -steps=1")
	fmt.Fprintln(os.Stderr, "  go run ./cmd/migrate down")
	fmt.Fprintln(os.Stderr, "  go run ./cmd/migrate down -steps=2")
	fmt.Fprintln(os.Stderr, "  go run ./cmd/migrate force -version=1")
	fmt.Fprintln(os.Stderr, "  go run ./cmd/migrate version")
}
