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
func (m *EntryModel) Insert(playername string, year int) (int, error) {

	// SQL statement to execute - default to 0 on creation
	stmt := `INSERT INTO nbaoverunders.entries (playername, year, created)
	VALUES(?, ?, UTC_TIMESTAMP())`

	
	// use Exec() method on the embedded connection pool to execute the statement
	result, err := m.DB.Exec(stmt, playername, year)
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
	stmt := `SELECT id, playername, year, created FROM nbaoverunders.entries
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
	stmt := `SELECT id, playername, year, points, created FROM nbaoverunders.entries
	WHERE year >= 2025 ORDER BY id DESC LIMIT 20`
	
	// Query() method on the connection pool returns a sql.Rows resultset containing the query result
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// defer rows.Close() to ensure the sql.Rows resultset is
	// always properly closed before the Latest() method returns
	// this defer statement should come after checking for an error from the Query() method
	defer rows.Close()

	// Initialize empty slice to hold the Entry structs
	var entries []Entry

	// Use rows.Next to iterate through the rows in the resultset. This
	// prepares the first (and then each subsequent) row to be acted on by the
	// rows.Scan() method. If iteration over all the rows completes then the
	// resultset automatically closes itself and frees-up the underlying
	// database connection
	for rows.Next() {
		
		// Create a pointer to a new zeroed Entry struct.
		var e Entry
	
		// Use rows.Scan() to copy the values from each field in the row to the
		// new Entry object that we created
		// arguments to row.Scan() must be pointers to the place you want to copy the data into
		err = rows.Scan(&e.ID, &e.PlayerName, &e.Year, &e.Points, &e.Created)
		if err != nil {
			return nil, err
		}
		
		// Append it to the slice of entries
		entries = append(entries, e)
	}

	// When the rows.Next() loop has finished, call rows.Err() to retrieve any
	// error that was encountered during the iteration
	// (don't assume that a successful iteration was completed over the whole resultset)
	if err = rows.Err(); err != nil {
		return nil, err
	}
	
	// If everything went OK then return the Entries slice
	return entries, nil
}