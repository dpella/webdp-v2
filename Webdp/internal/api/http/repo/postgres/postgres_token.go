package postgres

import "database/sql"

type TokenPostgres struct {
	db *sql.DB
}

func NewTokenPostgres(conn *sql.DB) TokenPostgres {
	return TokenPostgres{db: conn}
}

func (d TokenPostgres) SaveToken(userHandle string, stringifyedToken string) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}

	defer dfun(err, tx)

	q := "INSERT INTO UserTokens (username, token) VALUES ($1, $2)"

	_, err = tx.Exec(q, userHandle, stringifyedToken)

	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return err
}

func (d TokenPostgres) UserTokenExists(userHandle string) (bool, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return false, err
	}

	defer dfun(err, tx)

	q := "SELECT 1 FROM UserTokens WHERE username = $1"
	rows, err := tx.Query(q, userHandle)
	if err != nil {
		return false, err
	}

	ok := rows.Next()

	rows.Close()
	if err = tx.Commit(); err != nil {
		return false, err
	}

	return ok, nil

}

func (d TokenPostgres) UpdateToken(userHandle string, stringifyedToken string) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}

	defer dfun(err, tx)

	q := "UPDATE UserTokens SET token = $1 WHERE username = $2"

	_, err = tx.Exec(q, stringifyedToken, userHandle)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (d TokenPostgres) DeleteUserToken(userHandle string) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}

	defer dfun(err, tx)

	q := "DELETE FROM UserTokens WHERE username = $1"
	_, err = tx.Exec(q, userHandle)

	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (d TokenPostgres) GetUserToken(userHandle string) (string, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return "", err
	}

	defer dfun(err, tx)
	q := "SELECT token FROM UserTokens WHERE username = $1"
	row := tx.QueryRow(q, userHandle)

	var token string
	err = row.Scan(&token)

	if err != nil {
		return "", err
	}

	if err = tx.Commit(); err != nil {
		return "", err
	}

	return token, nil
}
