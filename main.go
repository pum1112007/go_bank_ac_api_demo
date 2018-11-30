package main

import (
	"database/sql"
	"go_bank_ac_api_demo/userapi"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {

	//db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	db, err := sql.Open("postgres", "postgres://panuwatsg:@localhost/panuwatsg?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	createTable := `
	CREATE TABLE IF NOT EXISTS secrets (
		id SERIAL PRIMARY KEY,
		key TEXT
	);
	CREATE TABLE IF NOT EXISTS Users (
		id SERIAL PRIMARY KEY,
		first_name TEXT,
		last_name TEXT
	);
	CREATE TABLE IF NOT EXISTS BankAccount (
		id SERIAL PRIMARY KEY,
		user_id INT,
		number INT,
		name TEXT,
		balance INT
		 
	);
	`
	if _, err := db.Exec(createTable); err != nil {
		log.Fatal(err)
	}
	userapi.StartServer(":"+os.Getenv("PORT"), db)
}
