package postgres

import (
	"EMP_Back/internal/domain/models"
	"EMP_Back/internal/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"os"
)

type Storage struct {
	db *sql.DB
}

func New() (*Storage, error) {
	op := "storage.postgres.New"

	db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)
	return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(ctx context.Context, username string, email string, passHash []byte) (int64, error) {
	op := "storage.postgres.SaveUser"

	_, err := s.db.Prepare("INSERT INTO users (username, email, pass_hash) VALUES ($1, $2, $3) RETURNING id")
	if err != nil {
		return 0, fmt.Errorf("%s: %v", op, err)
	}

	res, err := s.db.ExecContext(ctx, username, email, passHash)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // duplicate key
			return 0, fmt.Errorf("%s: %v", op, err, storage.ErrUserExist)
		}
		return 0, fmt.Errorf("%s: %v", op, err)
	}

	// Return the ID of the inserted row.

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %v", op, err)
	}
	return id, nil
}

// User returns a user by email.
func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	op := "storage.postgres.User"
	row := s.db.QueryRowContext(ctx, "SELECT id, username, email, pass_hash FROM users WHERE email = $1", email)
	var u models.User
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PassHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return u, storage.ErrUserNotFound
		}
		return u, fmt.Errorf("%s: %v", op, err)
	}
	return u, nil
}

// IsAdmin returns true if the user is an admin.
func (s *Storage) IsAdmin(ctx context.Context, id int64) (bool, error) {
	op := "storage.postgres.IsAdmin"
	row := s.db.QueryRowContext(ctx, "SELECT is_admin FROM users WHERE id = $1", id)
	var isAdmin bool
	err := row.Scan(&isAdmin)
	if err != nil {
		if err == sql.ErrNoRows {
			return isAdmin, storage.ErrAppNotFound
		}
		return isAdmin, fmt.Errorf("%s: %v", op, err)
	}
	return isAdmin, nil
}
