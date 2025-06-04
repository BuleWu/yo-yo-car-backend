package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"zavrsni/yo-yo-car/runtimebag"
	"zavrsni/yo-yo-car/shared/constants"
)

var (
	host     = runtimebag.GetEnvString(constants.DatabaseHost, "")
	port     = runtimebag.GetEnvString(constants.DatabasePort, "")
	user     = runtimebag.GetEnvString(constants.DatabaseUser, "")
	password = runtimebag.GetEnvString(constants.DatabasePassword, "")
	dbname   = runtimebag.GetEnvString(constants.DatabaseName, "")
)

func Connect() {
	fmt.Println("Env string: ", host)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open(runtimebag.GetEnvString(constants.DatabaseDriver, ""), psqlInfo)

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
