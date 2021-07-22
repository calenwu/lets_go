package postgres

import (
	"database/sql"

	"calenwu.com/snippetbox/pkg/models"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	stmt := `INSERT INTO users (name, email, hashed_password, created)
	VALUES($1, $2, $3, current_timestamp)
	RETURNING id;
	`
	lastInsertedId := -1
	err = m.DB.QueryRow(stmt, name, email, string(hashedPassword)).Scan(&lastInsertedId)
	if err != nil {
		if postgresErr, ok := err.(*pq.Error); ok {
			if postgresErr.Code == "23505" {
				return models.ErrDuplicateEmail
			}
		}
	}
	return err
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte
	err := m.DB.QueryRow(
		"SELECT id, hashed_password FROM users WHERE email = $1;", email,
		).Scan(&id, &hashedPassword)
	if err == sql.ErrNoRows {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return -1, models.ErrInvalidCredentials
	} else if err != nil {
		return -1, err
	}
	return id, nil
}

func (m *UserModel) Get(id int) (*models.User, error) {
	s := &models.User{}
	stmt := `SELECT id, name, email, created FROM users WHERE id = $1`
	err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Name, &s.Email, &s.Created)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	return s, err
}
