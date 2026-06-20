package server

import (
    "log"
    "net/http"

    "ascii-art-web/handlers"
)

func Start(addr string) error {
    mux := http.NewServeMux()
    handlers.RegisterRoutes(mux)
    mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

    log.Printf("Starting server on %s", addr)
    return http.ListenAndServe(addr, mux)
}

