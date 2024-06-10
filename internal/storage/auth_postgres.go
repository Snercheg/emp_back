package storage

import (
	"EMP_Back/internal/domain/models"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
)

type AuthPostgres struct {
	DB *sql.DB
}

func NewAuthPostgres(db *sql.DB) *AuthPostgres {
	return &AuthPostgres{DB: db}
}

// SaveUser saves a user in the database.
func (s *AuthPostgres) SaveUser(user models.User) (int64, error) {
	op := "storage.repository.SaveUser"

	stmt, err := s.DB.Prepare("INSERT INTO users (username, email, pass_hash) VALUES ($1, $2, $3) RETURNING id")
	if err != nil {
		return 0, fmt.Errorf("%s: %v", op, err)
	}

	res, err := stmt.Exec(user.Username, user.Email, user.PassHash)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // duplicate key
			return 0, fmt.Errorf("%s: %v", op, err, ErrUserExist)
		}
		return 0, fmt.Errorf("%s: %v", op, err)
	}

	// Return the UserId of the inserted row
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %v", op, err)
	}
	return id, nil
}

// GetUser returns a user by email.
func (s *AuthPostgres) GetUser(email string) (models.User, error) {
	op := "storage.repository.GetUser"
	row := s.DB.QueryRow("SELECT id, username, email, pass_hash FROM users WHERE email = $1", email)
	var u models.User
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return u, ErrUserNotFound
		}
		return u, fmt.Errorf("%s: %v", op, err)
	}
	return u, nil
}

func (s *AuthPostgres) IsAdmin(userId int64) (bool, error) {
	op := "storage.repository.IsAdmin"
	row := s.DB.QueryRow("SELECT is_admin FROM users WHERE id = $1", userId)
	var isAdmin bool
	err := row.Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return isAdmin, ErrUserNotFound
		}
		return isAdmin, fmt.Errorf("%s: %v", op, err)
	}
	return isAdmin, nil
}
