package main

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
)

func main() {
    connStr := "postgresql://postgres:printer0101@db.npoagmbfmexqcelbnsng.supabase.co:5432/postgres"
    
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        fmt.Printf("Error opening database: %v\n", err)
        return
    }
    defer db.Close()

    err = db.Ping()
    if err != nil {
        fmt.Printf("Error connecting to the database: %v\n", err)
        return
    }

    fmt.Println("Successfully connected to database!")
}