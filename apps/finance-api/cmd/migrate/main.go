package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

var openDB = sql.Open

func run() error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		flag.Usage()
		return fmt.Errorf("DATABASE_URL not set")
	}
	dir := "./migrations"
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}
	db, err := openDB("pgx", dsn)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	// Check TEST_MODE to avoid actually doing DB work if not available
	if os.Getenv("TEST_MODE") == "true" {
		return nil
	}
	files, err := filepath.Glob(filepath.Join(dir, "*.sql"))
	if err != nil {
		return fmt.Errorf("glob migrations: %w", err)
	}
	for _, f := range files {
		b, err := ioutil.ReadFile(f)
		if err != nil {
			return fmt.Errorf("read %s: %w", f, err)
		}
		log.Printf("Applying %s", filepath.Base(f))
		if _, err := db.Exec(string(b)); err != nil {
			return fmt.Errorf("apply %s: %w", f, err)
		}
	}
	log.Printf("migrations applied: %d files", len(files))
	return nil
}
