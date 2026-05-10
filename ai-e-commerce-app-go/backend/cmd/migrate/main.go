package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"ai-e-commerce-app-go/backend/internal/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: go run ./cmd/migrate [up|down]")
	}

	cfg := config.Load()

	m, err := migrate.New("file://migrations", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("create migrator: %v", err)
	}
	defer m.Close()

	switch os.Args[1] {
	case "up":
		err = m.Up()
	case "down":
		err = m.Steps(-1)
	default:
		log.Fatalf("unknown migration command %q; use up or down", os.Args[1])
	}

	if errors.Is(err, migrate.ErrNoChange) {
		fmt.Println("No migration changes to apply.")
		return
	}

	if err != nil {
		log.Fatalf("run migrations: %v", err)
	}

	fmt.Println("Migrations completed.")
}
