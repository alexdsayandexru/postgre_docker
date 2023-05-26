package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func main() {
	execute_query(`DROP TABLE IF EXISTS users`)

	execute_query(`CREATE TABLE public.users
	(
		id bigint NOT NULL,
		name text,
		created date,
		PRIMARY KEY (id)
	);`)
}

const (
	host     = "localhost"
	port     = 5555
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

func execute_query(query string) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("open:", err)
		return
	}

	_, err = db.Exec(query)
	if err != nil {
		fmt.Println("exec:", err)
		return
	}
	db.Close()
	fmt.Print("success:", query)
}
