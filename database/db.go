package database

import (
    "database/sql"
    "fmt"

    _ "github.com/go-sql-driver/mysql"
)

func New() (*sql.DB, error) {
    dsn := "root:alek1234@tcp(localhost:3306)/go_chi_1"
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }

    err = db.Ping()
    if err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }

    return db, nil
}