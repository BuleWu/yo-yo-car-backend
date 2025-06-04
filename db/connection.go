package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func Connect(host, port, user, password, database, driver string) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=disable",
		host, port, user, password, database)

	db, err := sql.Open(driver, psqlInfo)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
}
