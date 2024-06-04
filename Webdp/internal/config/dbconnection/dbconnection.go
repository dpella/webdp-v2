package dbconnection

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func ConnectPostgresDB(name string, user string, pw string, host string, port string) (*sql.DB, error) {
	connstring := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=disable", user, name, pw, host, port)
	db, err := sql.Open("postgres", connstring)

	if err != nil {
		return nil, err
	}

	return db, nil
}
