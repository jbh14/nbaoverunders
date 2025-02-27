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

-- removing this for now as part of docker compose
-- create a web user - NO need to do this, already did as part of snippetbox
-- CREATE USER 'web'@'localhost';
-- GRANT SELECT, INSERT, UPDATE, DELETE ON entries.* TO 'web'@'localhost';
-- Important: Make sure to swap 'pass' with a password of your own choosing.
-- ALTER USER 'web'@'localhost' IDENTIFIED BY 'pass';