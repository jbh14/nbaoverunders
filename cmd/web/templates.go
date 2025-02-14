package main

import "github.com/jbh14/nbaoverunders/internal/models"

// Define a templateData type to act as the holding structure for
// any dynamic data that we want to pass to our HTML templates
type templateData struct {
	Entry 	models.Entry
	Entries	[]models.Entry
}