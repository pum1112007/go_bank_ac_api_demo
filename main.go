package main

import (
	"GO_BANK_AC_API_DEMO/api"
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	api.StartServer(":"+os.Getenv("PORT"), db)
}
