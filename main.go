package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	config := loadConfig()
	connStr := getConnectionString(config)

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println(`
		1 : DROP TABLE
		2 : CREATE TABLE
		3 : GENERATE USERS
		4 : PRINT USERS
		5 : DELETE USERS
		`)
		fmt.Print("Enter command: ")
		text, _ := reader.ReadString('\n')
		switch text {
		case "1\n":
			executeSingleQuery(connStr, `DROP TABLE IF EXISTS users`)
		case "2\n":
			executeSingleQuery(connStr, `CREATE TABLE public.users
			(
				id bigint NOT NULL,
				name text,
				created date,
				PRIMARY KEY (id)
			);`)
		case "3\n":
			generateUsers(connStr, 100)
		case "4\n":
			printData(connStr)
		case "5\n":
			executeSingleQuery(connStr, `DELETE FROM users`)
		default:
			return
		}
	}
}

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

func loadConfig() *Config {
	var config Config
	raw, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	json.Unmarshal(raw, &config)
	return &config
}

func getConnectionString(c *Config) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Port, c.User, c.Password, c.Database)
}

func generateUsers(connStr string, count int) {
	for i := 1; i <= count; i++ {
		t := time.Now().AddDate(0, 0, i)
		query := fmt.Sprintf("INSERT INTO users (id, name, created) "+
			"VALUES(%d, %s, %s)",
			i, fmt.Sprintf("'User %d'", i), fmt.Sprintf("'%s'", t.Format("2006-01-02")))

		executeSingleQuery(connStr, query)
	}
}

func executeSingleQuery(connStr string, query string) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("sql.Open:", connStr, err)
		return
	}

	defer db.Close()

	_, err = db.Exec(query)
	if err != nil {
		fmt.Println("db.Exec:", err)
	}
}

func printData(connStr string) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("sql.Open:", connStr, err)
		return
	}

	defer db.Close()

	rows, err := db.Query("SELECT name, created, id FROM users")
	if err != nil {
		fmt.Println("select:", connStr, err)
		return
	}

	for rows.Next() {
		var id int
		var name string
		var created time.Time
		err = rows.Scan(&name, &created, &id)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(id, name, created)
		}
	}
}
