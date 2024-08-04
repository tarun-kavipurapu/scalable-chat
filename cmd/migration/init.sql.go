package main

import (
	"log"
	"os"
	"strconv"
	"tarun-kavipurapu/test-go-chat/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	config.LoadConfig(".")
	dbSource := config.EnvVars.DBSource
	log.Println(dbSource)

	m, err := migrate.New("file://db/migrations", dbSource)
	if err != nil {
		log.Fatal("Init error: ", err.Error())
	}

	if len(os.Args) < 2 {
		log.Fatal("Missing migration direction argument: up, down, force, or goto")
	}

	arg := os.Args[1]

	switch arg {
	case "up":
		err = m.Up()
		if err != nil && err != migrate.ErrNoChange {
			if e, ok := err.(migrate.ErrDirty); ok {
				log.Printf("Database is dirty at version %d, forcing to the latest version...\n", e.Version)
				forceToLatestVersion(m)
			} else {
				log.Fatal("Migration error: ", err)
			}
		} else {
			log.Println("Migration up successful")
		}
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal("Migration error: ", err)
		} else {
			log.Println("Migration down successful")
		}
	case "force":
		if len(os.Args) < 3 {
			log.Fatal("Missing version argument for force")
		}
		version, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal("Invalid version argument: ", err)
		}
		if err := m.Force(version); err != nil {
			log.Fatal("Force error: ", err)
		} else {
			log.Println("Database forced to version ", version)
		}
	case "goto":
		if len(os.Args) < 3 {
			log.Fatal("Missing version argument for goto")
		}
		version, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal("Invalid version argument: ", err)
		}
		err = m.Migrate(uint(version))
		if err != nil && err != migrate.ErrNoChange {
			if e, ok := err.(migrate.ErrDirty); ok {
				log.Printf("Database is dirty at version %d, forcing to version %d...\n", e.Version, version)
				forceToVersion(m, version)
			} else {
				log.Fatal("Migration error: ", err)
			}
		} else {
			log.Printf("Migration to version %d successful\n", version)
		}
	default:
		log.Fatal("Invalid migration direction argument: ", arg)
	}
}

func forceToVersion(m *migrate.Migrate, version int) {
	if err := m.Force(version); err != nil {
		log.Fatal("Force error: ", err)
	}
	log.Printf("Database forced to version %d\n", version)

	// Retry the migration to the specified version
	if err := m.Migrate(uint(version)); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Migration retry error: ", err)
	}
	log.Printf("Migration to version %d successful after forcing\n", version)
}

func forceToLatestVersion(m *migrate.Migrate) {
	// Determine the latest version
	version, dirty, err := m.Version()
	if err != nil {
		log.Fatal("Version error: ", err)
	}

	if dirty {
		log.Printf("Database is dirty at version %d, forcing to latest version...\n", version)
	}

	if err := m.Force(int(version)); err != nil {
		log.Fatal("Force error: ", err)
	}
	log.Println("Database forced to latest version ", version)

	// Retry the migration to the latest version
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Migration retry error: ", err)
	}
	log.Println("Migration up successful after forcing to latest version")
}
