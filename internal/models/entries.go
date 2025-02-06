package models

import (
	"database/sql"
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
}

// This will return a specific Entry based on its id.
func (m *EntryModel) Get(id int) (Entry, error) {
	return nil, nil
}

// This will return the 10 most recently created entries.
func (m *EntryModel) Latest() ([]Entry, error) {
	return nil, nil
}