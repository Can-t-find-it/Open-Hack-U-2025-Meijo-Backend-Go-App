package main

import (
    "hacku_2025_meijo/backend/internal/database"
    "hacku_2025_meijo/backend/internal/router"
)

import (
    "fmt"
)

func main() {
    database.Connect()

    r := router.SetupRouter()
	fmt.Println("Server running on http://localhost:8080")
    r.Run(":8080")
}
