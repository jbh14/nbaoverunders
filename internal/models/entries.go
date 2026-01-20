package models

import (
	"database/sql"
	"errors"
	"log"
	"time"
)

type Entry struct {
	ID         int
	PlayerName string
	Year       int
	Points     float32
	Created    time.Time
	Picks      []Pick
	Lock       string
}

type Pick struct {
	TeamSeasonID        int
	TeamName            string
	WinsActual          int
	LossesActual        int
	WinsLine            float64
	WinsProjected       int
	LossesProjected     int
	OverSelected        bool
	LockSelected        bool
	NosweatLockSelected bool
	Points              float32
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

	// Step 2: Fetch the associated picks, ordered by team name (ascending)
	stmt2 := `SELECT 
    			p.teamseason_id, 
    			t.teamname, 
				COALESCE(s.wins_actual, 0) AS wins_actual,
				COALESCE(s.losses_actual, 0) AS losses_actual,
				COALESCE(s.wins_line, 0) AS wins_line,
				COALESCE(s.wins_projected, 0) AS wins_projected,
				COALESCE(s.losses_projected, 0) AS losses_projected,
				p.over_selected, 
				p.lock_selected,
				p.nosweat_lock_selected
			FROM nbaoverunders.picks p
			INNER JOIN teamseasons s ON p.teamseason_id = s.id
			INNER JOIN teams t ON s.team_id = t.id
			WHERE s.season_start_year = 2025
		     AND p.entry = ? 
			 ORDER BY t.teamname ASC;`

	pickrows, err2 := m.DB.Query(stmt2, id)
	if err2 != nil {
		return Entry{}, err2
	}

	if err2 != nil {
		return Entry{}, err
	}
	defer pickrows.Close()

	// empty slice to hold picks
	var picks []Pick

	// Step 3: Iterate over the rows and populate the Picks slice
	totalPoints := 0.0 // initialize totalPoints to 0
	for pickrows.Next() {

		// Create a pointer to a new zeroed Pick struct.
		var pick Pick

		// Use rows.Scan() to copy the values from each field in the row to the
		// new Entry object that we created
		// arguments to row.Scan() must be pointers to the place you want to copy the data into
		err = pickrows.Scan(&pick.TeamSeasonID, &pick.TeamName, &pick.WinsActual, &pick.LossesActual, &pick.WinsLine, &pick.WinsProjected, &pick.LossesProjected, &pick.OverSelected, &pick.LockSelected, &pick.NosweatLockSelected)
		if err != nil {
			return Entry{}, err
		}
		// multiplier for over/under and lock
		overUnderMultiplier := 1
		if !pick.OverSelected {
			overUnderMultiplier = -1
		}
		lockMultiplier := 1
		if pick.LockSelected || pick.NosweatLockSelected {
			lockMultiplier = 2
		}
		// Calculate the points based on the wins and losses, over/under selection, and lock
		points := (float32(pick.WinsActual) - float32(pick.WinsLine)) * float32(overUnderMultiplier) * float32(lockMultiplier)

		// Check if this is a no-sweat lock with negative points
		if pick.NosweatLockSelected && points < 0 {
			pick.Points = 0 // Display as 0 for nosweat lock with negative outcome
			// Don't add to totalPoints
		} else {
			pick.Points = points           // Assign the calculated points
			totalPoints += float64(points) // Add to total
		}

		// Append it to the slice of picks
		picks = append(picks, pick)
	}

	// When the rows.Next() loop has finished, call rows.Err() to retrieve any
	// error that was encountered during the iteration
	// (don't assume that a successful iteration was completed over the whole resultset)
	if err = pickrows.Err(); err != nil {
		return Entry{}, err
	}

	// put the picks and Points into the Entry struct
	e.Points = float32(totalPoints)
	e.Picks = picks

	// Step 4: Calculate the total points for the entry and update entry in the database
	stmt3 := `UPDATE entries SET points = ? WHERE id = ?;`

	// use Exec() method on the embedded connection pool to execute the statement
	result, err := m.DB.Exec(stmt3, totalPoints, id)
	if err != nil {
		return Entry{}, err
	}

	// check for errors or no rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	if rowsAffected == 0 {
		log.Println("no entries updated with current points")
	}

	// If everything went OK, then return the filled Entry struct.
	return e, nil
}

// This will return the 10 most recently created entries.
func (m *EntryModel) Latest() ([]Entry, error) {
	stmt := `SELECT id, playername, points, year, created FROM nbaoverunders.entries
	WHERE year = 2025 ORDER BY points DESC LIMIT 20`
	// @TODO - year should be dynamic, defaulting to current year

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
		err = rows.Scan(&e.ID, &e.PlayerName, &e.Points, &e.Year, &e.Created)
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

	// to add - subquery on picks to get the number of picks for each entry and lock

	// If everything went OK then return the Entries slice
	return entries, nil
}
