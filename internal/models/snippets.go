package models

import (
	"database/sql"
	"time"
)

// Define a Snippet type to hald the data for an individual snippet.
type Snippet struct {
	ID		int
	Title	string
	Content time.Time
	Created time.Time
	Expires time.Time
}

// Wraps s sql.DB connection pool.
type SnippetModel struct {
	DB *sql.DB
}

// insert a new snippet
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	q := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(q, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Convert the id type from int64 to int
	return int(id), nil
}

// return a specific snippet
func (m *SnippetModel) Get(id int) (Snippet, error) {
	return Snippet{}, nil
}

// return 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]Snippet, error) {
	return nil, nil
}
