package main

import (
    "log"

    "ascii-art-web/server"
)

func main() {
    if err := server.Start(":8080"); err != nil {
        log.Fatal(err)
    }
}
