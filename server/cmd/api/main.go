package main

import (
    "flag"
    "database/sql"
    "log"
    "net/http"
    _ "github.com/go-sql-driver/mysql"
)

type config struct {
    port string
    env  string
    db   struct {
        dns string
    }
}

type application struct {
    config config
   // db     *sql.DB
}

func main() {
    var con config

    flag.StringVar(&con.port, "addr", "4000", "HTTP network address port for API")
    flag.StringVar(&con.db.dns, "dns", "root:admin@(127.0.0.1:3306)/biblioteca?parseTime=true", "MySQL data source name")

    flag.Parse()

    db, err := openDB(con)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    log.Println("Database connection pool established")

    app := &application{
        config: con,
    //    db:     db,
    }

    svr := &http.Server{
        Addr:    ":" + app.config.port,
        Handler: app.routes(),
    }

    err = svr.ListenAndServe()
    if err != nil {
        log.Panic(err)
    }
}

func openDB(con config) (*sql.DB, error) {
    db, err := sql.Open("mysql", con.db.dns)
    if err != nil {
        return nil, err
    }

    if err = db.Ping(); err != nil {
        return nil, err
    }

    return db, nil
}

