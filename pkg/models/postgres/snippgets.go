package postgres

import (
	"calenwu.com/snippetbox/pkg/models"
	"database/sql"
	"time"
)

// SnippetModel Define a SnippetModel type which wraps a sql.DB connection pool.
type SnippetModel struct {
	DB         *sql.DB
	InsertStmt *sql.Stmt
}
func NewSnippetModel(db *sql.DB) (*SnippetModel, error) {
	// Use the Prepare method to create a new prepared statement for the current connection pool.
	insertStmt, err := db.Prepare(
		`INSERT INTO snippets (title, content, created, expires)
		VALUES($1, $2, current_timestamp, $3)
		RETURNING id;`)
	if err != nil {
		return nil, err
	}
	// Store it in our Model object, alongside the connection poool
	return &SnippetModel{db, insertStmt}, nil
}

// InsertNewSnippet Insert a new SnippetModel into the database
func (m *SnippetModel) InsertNewSnippet(title, content string, expires int) (int, error) {
	// Call Exec directly on the prepared statement, rather than the connection pool.
	lastInsertedId := -1
	err := m.InsertStmt.QueryRow(title, content, time.Now().AddDate(
		0, 0, expires)).Scan(&lastInsertedId)
	return lastInsertedId, err
}

// Get This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
		WHERE expires > current_timestamp AND id = $1`
	s := &models.Snippet{}
	row := m.DB.QueryRow(stmt, id)
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
	stmt := `SELECT id, title, content, created, expires FROM snippets
		WHERE expires > current_timestamp 
		ORDER BY created
		DESC LIMIT 10
	`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var snippets []*models.Snippet
	for rows.Next() {
		s :=& models.Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
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

func (m *SnippetModel) ExampleTransaction() error {
	// Calling the Begin() method on the connection pool creates a new sql.Tx
	// object, which represents the in -progress database transaction.
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}
	// Call Exec() on the transaction, passin gin your statement and any parameters.
	// It's important to notice that tx.Exec() is called on the transaction object just created,
	// NOT the connection opool. Although we're using tx.Exec() here you can also use tx.Query() and tx.QueryRow()
	// in exactly the same way
	_, err = tx.Exec("INSERT INTO ...")
	if err != nil {
		// If there is any error, we call the tx.Rollback() method on the transaction.
		// This will abort the transaction and no changes will be made to the database.
		tx.Rollback()
		return err
	}
	_, err = tx.Exec("UPDATE ...")
	if err != nil {
		tx.Rollback()
		return err
	}
	// If there are no errors, the transaction can be committed. Its IMPORTANT to call either Rollback() or Commit()
	// before the function returns. Otherwise the connection will stay open and and wont be returned to the pool
	return tx.Commit()
}

