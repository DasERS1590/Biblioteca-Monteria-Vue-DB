package main

import (
	"database/sql"
	"flag"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {

	dsn := flag.String("dsn", "root:admin@tcp(localhost:3306)/biblioteca?multiStatements=true", "MySQL DSN")
	direction := flag.String("direction", "", "Migration direction: 'up' or 'down'")
	flag.Parse()

	
	if *direction != "up" && *direction != "down" {
		log.Fatal("Invalid direction. Use  -direction=up or  -direction=down ")
	}

	db, err := sql.Open("mysql", *dsn)
	if err != nil {

		log.Fatal(err)
	}
	defer db.Close()

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations",
		"mysq",
		driver)

	if err != nil {
		log.Fatalf("Error >: %v", err)

	}

	// Ejecutar la migración según el argumento proporcionado
	switch *direction {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}
}
