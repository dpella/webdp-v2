package postgres

import (
	"context"
	"database/sql"
	"time"
	"webdp/internal/api/http/entity"

	"github.com/lib/pq"

	errors "webdp/internal/api/http"
)

type UserPostgres struct {
	db *sql.DB
}

func NewUserPostgres(conn *sql.DB) UserPostgres {
	return UserPostgres{db: conn}
}

func (u UserPostgres) GetUsers() ([]entity.UserResponse, error) {
	tx, err := u.db.BeginTx(context.Background(), nil)
	if err != nil {
		return []entity.UserResponse{}, err
	}
	defer dfun(err, tx)

	q := "SELECT handle, name, roles, created_time, updated_time FROM GetUsers"

	rows, err := tx.Query(q)
	if err != nil {
		return []entity.UserResponse{}, err
	}

	out := make([]entity.UserResponse, 0)
	for rows.Next() {
		usr := entity.UserResponse{}
		err = rows.Scan(&usr.Handle, &usr.Name, pq.Array(&usr.Roles), &usr.CreatedOn, &usr.UpdatedOn)
		if err != nil {
			rows.Close()
			return []entity.UserResponse{}, err
		}
		out = append(out, usr)
	}
	rows.Close()
	tx.Commit()
	return out, nil
}

func (u UserPostgres) GetUser(handle string) (entity.UserResponse, error) {
	tx, err := u.db.BeginTx(context.Background(), nil)
	if err != nil {
		return entity.UserResponse{}, err
	}

	defer dfun(err, tx)
	q := "SELECT handle, name, roles, created_time, updated_time FROM GetUsers WHERE handle = $1"
	row := tx.QueryRow(q, handle)
	var user entity.UserResponse
	err = row.Scan(&user.Handle, &user.Name, pq.Array(&user.Roles), &user.CreatedOn, &user.UpdatedOn)
	if err != nil {
		return entity.UserResponse{}, errors.ErrNotFound
	}
	tx.Commit()
	return user, nil
}

func (u UserPostgres) GetByHandle(handle string) (entity.LoginRequest, error) {
	tx, err := u.db.BeginTx(context.Background(), nil)
	if err != nil {
		return entity.LoginRequest{}, err
	}

	defer dfun(err, tx)
	q := "SELECT handle, pwd FROM Users WHERE handle = $1"
	row := tx.QueryRow(q, handle)
	var user entity.LoginRequest
	err = row.Scan(&user.Username, &user.PWD)
	if err != nil {
		return entity.LoginRequest{}, errors.ErrNotFound
	}
	tx.Commit()
	return user, nil

}

func (u UserPostgres) CreateUser(post entity.UserPost) (string, error) {
	tx, err := u.db.BeginTx(context.Background(), nil)
	if err != nil {
		return "", err
	}
	defer dfun(err, tx)
	created := time.Now().UTC()
	q1 := "INSERT INTO Users (handle, pwd, name, created_time, updated_time) VALUES ($1, $2, $3, $4, $5)"
	q2 := "INSERT INTO UserRoles (username, role) VALUES ($1, $2)"
	_, err = tx.Exec(q1, post.Handle, post.PWD, post.Name, created, created)
	if err != nil {
		return "", err
	}
	for _, role := range post.Roles {
		_, err = tx.Exec(q2, post.Handle, role)
		if err != nil {
			return "", err
		}
	}
	tx.Commit()
	return post.Handle, nil

}

func (u UserPostgres) UpdateUser(handle string, patch entity.UserPatch) (string, error) {
	tx, err := u.db.BeginTx(context.Background(), nil)
	if err != nil {
		return "", err
	}
	defer dfun(err, tx)

	updated := time.Now().UTC()

	// Supports partial updates. Lets admin update roles without interfering
	// with user's password.

	if patch.Name != "" {
		q1 := "UPDATE Users SET name = $1, updated_time = $2 WHERE handle = $3"
		_, err = tx.Exec(q1, patch.Name, updated, handle)
		if err != nil {
			return "", err
		}
	}

	if patch.PWD != "" {
		q1 := "UPDATE Users SET pwd = $1, updated_time = $2 WHERE handle = $3"
		_, err = tx.Exec(q1, patch.PWD, updated, handle)
		if err != nil {
			return "", err
		}
	}

	if patch.Roles != nil {
		d := "DELETE FROM UserRoles WHERE username = $1"
		_, err = tx.Exec(d, handle)
		if err != nil {
			return "", err
		}

		q := "INSERT INTO UserRoles (username, role) VALUES ($1, $2)"
		for _, role := range patch.Roles {
			_, err = tx.Exec(q, handle, role)
			if err != nil {
				return "", err
			}
		}
	}

	tx.Commit()
	return handle, nil
}

func (u UserPostgres) DeleteUser(handle string) error {
	tx, err := u.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer dfun(err, tx)
	q := "DELETE FROM Users WHERE handle = $1"
	_, err = tx.Exec(q, handle)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func dfun(e error, tx *sql.Tx) {
	if e != nil {
		tx.Rollback()
	}
}
