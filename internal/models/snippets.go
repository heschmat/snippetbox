package models

import (
	"database/sql"
	"errors"
	"time"
)

// Define a Snippet type to hald the data for an individual snippet.
type Snippet struct {
	ID		int
	Title	string
	Content string
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
	q := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND id = ?`

	// This returns a `sql.Row` object which holds the result from the database.
	row := m.DB.QueryRow(q, id)

	var s Snippet

	// row.Scan copies the values from each field in sql.Row to the corresponding field
	// in the Snippet struct, s;
	// N.B. row.Scan requires *pointers* as its argument.
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		// row.Scan() returns sql.ErrNoRows as error if the query returns no rows.
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}
	// If everything went well, retunr the filled Snippet struct.
	return s, nil
}

// return 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]Snippet, error) {
	return nil, nil
}
