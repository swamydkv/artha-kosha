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
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		fmt.Println("DATABASE_URL not set")
		flag.Usage()
		os.Exit(1)
	}
	dir := "./migrations"
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	files, err := filepath.Glob(filepath.Join(dir, "*.sql"))
	if err != nil {
		log.Fatalf("glob migrations: %v", err)
	}
	// sort by name (Glob already returns sorted on most systems)
	for _, f := range files {
		b, err := ioutil.ReadFile(f)
		if err != nil {
			log.Fatalf("read %s: %v", f, err)
		}
		log.Printf("Applying %s", filepath.Base(f))
		if _, err := db.Exec(string(b)); err != nil {
			log.Fatalf("apply %s: %v", f, err)
		}
	}
	log.Printf("migrations applied: %d files", len(files))
}
