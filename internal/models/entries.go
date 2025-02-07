package models

import (
	"database/sql"
	"errors"
	"time"
)

type Entry struct {
	ID int
	PlayerName string
	Year int
	Points float32
	Created time.Time
}


// Define a EntryModel type which wraps a sql.DB connection pool.
type EntryModel struct {
	DB *sql.DB
}

// This will insert a new player entry into the database.
func (m *EntryModel) Insert(playerName string, year int) (int, error) {
	// default Points to 0 on creation? 
	return 0, nil

	// SQL statement to execute
	stmt := `INSERT INTO entries (playername, year, points, created)
	VALUES(?, ?, 0, UTC_TIMESTAMP() )`

	
	// use Exec() method on the embedded connection pool to execute the statement
	result, err := m.DB.Exec(stmt, playerName, year)
		if err != nil {
		return 0, err
	}
	
	// Use the LastInsertId() method on the result to get the ID of our
	// newly inserted record in the entries table.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	
	// The ID returned has the type int64, so we convert it to an int type
	// before returning.
	return int(id), nil
}

// This will return a specific Entry based on its id.
func (m *EntryModel) Get(id int) (Entry, error) {
	stmt := `SELECT id, playername, year, created FROM entries
			WHERE id = ?`
	
	row := m.DB.QueryRow(stmt, id)
	
	// Initialize a new zeroed Entry struct.
	var e Entry
	
	// args to row.Scan are *pointers* to the place we want to copy the data into,
	// and the number of args must be exactly the same as the number of
	// columns returned by SQL statement
	err := row.Scan(&e.ID, &e.PlayerName, &e.Year, &e.Created)
	if err != nil {
		// return ErrNoRecord if no rows found, otherwise actual err
		if errors.Is(err, sql.ErrNoRows) {
			return Entry{}, ErrNoRecord
		} else {
			return Entry{}, err
		}
	}

	// If everything went OK, then return the filled Entry struct.
	return e, nil
}

// This will return the 10 most recently created entries.
func (m *EntryModel) Latest() ([]Entry, error) {
	return nil, nil
}