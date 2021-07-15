package postgres

import (
	"calenwu.com/snippetbox/pkg/models"
	"database/sql"
	"time"
)

// SnippetModel Define a SnippetModel type which wraps a sql.DB connection pool.
type SnippetModel struct {
	DB *sql.DB
}

// Get This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
		WHERE expires > current_timestamp AND id = $1`
	row := m.DB.QueryRow(stmt, id)
	s := &models.Snippet{}
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	return s, nil
}

// Latest This will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return nil, nil
}

// Insert This will insert a new snippet into the database.
func (m *SnippetModel) Insert(title, content string, expires int) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
		VALUES($1, $2, current_timestamp, $3)
		RETURNING id;
	`
	lastInsertedId := -1
	err := m.DB.QueryRow(stmt, title, content, time.Now().AddDate(
		0, 0, expires)).Scan(&lastInsertedId)
	if err != nil {
		return 0, err
	}
	if lastInsertedId == -1 {
		return 0, nil
	}
	return lastInsertedId, nil
}
