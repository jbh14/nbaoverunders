-- Create a new UTF-8 `nbaoverunders` database.
-- CREATE DATABASE nbaoverunders CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE DATABASE IF NOT EXISTS nbaoverunders CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Switch to using the `nbaoverunders` database.
USE nbaoverunders;

-- Create an `entries` table.
CREATE TABLE entries (
	id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
	playername VARCHAR(100) NOT NULL,
	year INTEGER NOT NULL,
	points DECIMAL(10, 2),
	created DATETIME NOT NULL
);
-- Add an index on the created column.
CREATE INDEX idx_entries_created ON entries(created);

-- Add some dummy records
INSERT INTO entries (playername, year, points, created) VALUES (
	'Brendan Heinz',
	2024,
	20.50,
	UTC_TIMESTAMP()
);
INSERT INTO entries (playername, year, points, created) VALUES (
	'Thomas Bruch',
	2024,
	17.00,
	UTC_TIMESTAMP()
);
INSERT INTO entries (playername, year, points, created) VALUES (
	'Andy Heinz',
	2024,
	23.50,
	UTC_TIMESTAMP()
);

-- create a "teams" table with an entry for each team
CREATE TABLE teams (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    teamname VARCHAR(100) NOT NULL UNIQUE
);

-- create a "team seasons" table with an entry for each team for each year
CREATE TABLE teamseasons (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    team_id INTEGER NOT NULL,
    season_start_year INTEGER NOT NULL,
	wins_actual INTEGER,
	losses_actual INTEGER,
	wins_line DECIMAL(4,1),
	losses_line DECIMAL(4,1),
	wins_projected INTEGER,
	losses_projected INTEGER,
	projected_over BOOLEAN,
	confirmed_over BOOLEAN,
    FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE
);

-- create team records and team seasons for all teams for 2024-2025 season
INSERT INTO teams (teamname) VALUES
    ('Atlanta Hawks'),
    ('Boston Celtics'),
    ('Brooklyn Nets'),
    ('Charlotte Hornets'),
    ('Chicago Bulls'),
    ('Cleveland Cavaliers'),
    ('Dallas Mavericks'),
    ('Denver Nuggets'),
    ('Detroit Pistons'),
    ('Golden State Warriors'),
    ('Houston Rockets'),
    ('Indiana Pacers'),
    ('Los Angeles Clippers'),
    ('Los Angeles Lakers'),
    ('Memphis Grizzlies'),
    ('Miami Heat'),
    ('Milwaukee Bucks'),
    ('Minnesota Timberwolves'),
    ('New Orleans Pelicans'),
    ('New York Knicks'),
    ('Oklahoma City Thunder'),
    ('Orlando Magic'),
    ('Philadelphia 76ers'),
    ('Phoenix Suns'),
    ('Portland Trail Blazers'),
    ('Sacramento Kings'),
    ('San Antonio Spurs'),
    ('Toronto Raptors'),
    ('Utah Jazz'),
    ('Washington Wizards');

-- Insert team seasons for the 2024-2025 season
INSERT INTO teamseasons (team_id, season_start_year)
SELECT id, 2024 FROM teams;

-- update Hawks season with wins and losses
UPDATE teamseasons
SET wins_actual = 36, losses_actual = 40, wins_line = 36.5
WHERE season_start_year = 2024
AND team_id = (SELECT ID FROM teams WHERE teamname = 'Atlanta Hawks');

-- update Cavs season with wins and losses
UPDATE teamseasons
SET wins_actual = 61, losses_actual = 15, wins_line = 48.5
WHERE season_start_year = 2024
AND team_id = (SELECT ID FROM teams WHERE teamname = 'Cleveland Cavaliers');

-- update Clippers season with wins and losses
UPDATE teamseasons
SET wins_actual = 44, losses_actual = 32, wins_line = 35.5
WHERE season_start_year = 2024
AND team_id = (SELECT ID FROM teams WHERE teamname = 'Los Angeles Clippers');

-- update Pacers season with wins and losses
UPDATE teamseasons
SET wins_actual = 45, losses_actual = 31, wins_line = 46.5
WHERE season_start_year = 2024
AND team_id = (SELECT ID FROM teams WHERE teamname = 'Indiana Pacers');


-- create "picks" table
CREATE TABLE picks (
	id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
	entry INTEGER NOT NULL,
	teamseason_id INTEGER NOT NULL,
	over_selected BOOLEAN,
	lock_selected BOOLEAN,
	FOREIGN KEY (teamseason_id) REFERENCES teamseasons(id) ON DELETE CASCADE,
	FOREIGN KEY (entry) REFERENCES entries(id) ON DELETE CASCADE
);

-- create a few "picks" entries here, but most will be added by the user
-- Brendan over on the Hawks
INSERT INTO picks (entry, teamseason_id, over_selected, lock_selected)
SELECT 
    e.id, 
    ts.id, 
    TRUE, 
    FALSE
FROM entries e
JOIN teamseasons ts ON ts.team_id = (SELECT id FROM teams WHERE teamname = 'Atlanta Hawks') 
    AND ts.season_start_year = 2024
WHERE e.playername = 'Brendan Heinz' AND e.year = 2024;

-- Brendan over on the Cavs
INSERT INTO picks (entry, teamseason_id, over_selected, lock_selected)
SELECT 
    e.id, 
    ts.id, 
    TRUE, 
    FALSE
FROM entries e
JOIN teamseasons ts ON ts.team_id = (SELECT id FROM teams WHERE teamname = 'Cleveland Cavaliers') 
    AND ts.season_start_year = 2024
WHERE e.playername = 'Brendan Heinz' AND e.year = 2024;

-- Brendan under on the Clippers
INSERT INTO picks (entry, teamseason_id, over_selected, lock_selected)
SELECT 
    e.id, 
    ts.id, 
    FALSE, 
    FALSE
FROM entries e
JOIN teamseasons ts ON ts.team_id = (SELECT id FROM teams WHERE teamname = 'Los Angeles Clippers') 
    AND ts.season_start_year = 2024
WHERE e.playername = 'Brendan Heinz' AND e.year = 2024;

-- Brendan lock over on the Pacers
INSERT INTO picks (entry, teamseason_id, over_selected, lock_selected)
SELECT 
    e.id, 
    ts.id, 
    TRUE, 
    TRUE
FROM entries e
JOIN teamseasons ts ON ts.team_id = (SELECT id FROM teams WHERE teamname = 'Indiana Pacers') 
    AND ts.season_start_year = 2024
WHERE e.playername = 'Brendan Heinz' AND e.year = 2024;