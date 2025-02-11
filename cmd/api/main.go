package main

import (
    "log"
    "net/http"
    
    "github.com/nurdamiron/print-automation/internal/config"
    "github.com/nurdamiron/print-automation/internal/api"
    "github.com/nurdamiron/print-automation/internal/db"
)

func main() {
    cfg, err := config.Load()
    if err != nil {
        log.Fatal(err)
    }

    db, err := db.Connect(cfg.DatabaseURL)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    server := api.NewServer(cfg, db)
    log.Fatal(http.ListenAndServe(cfg.ServerAddr, server.Router))
}