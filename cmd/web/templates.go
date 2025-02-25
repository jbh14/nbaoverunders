package main

import (
	"html/template"
	"path/filepath"
	"time"
	"github.com/jbh14/nbaoverunders/internal/models"
)

// Define a templateData type to act as the holding structure for
// any dynamic data that we want to pass to our HTML templates
type templateData struct {
	CurrentYear int
	Entry 	models.Entry
	Entries	[]models.Entry
}

// humanDate functio - returns a formatted representation of a time.Time object.
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// initialize a template.FuncMap object and store it in a global variable
// this is essentially a string-keyed map which acts as a lookup between the names of custom template functions and the functions themselves
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache
	cache := map[string]*template.Template{}

	// Use the filepath.Glob() function to get a slice of all filepaths that
	// match the pattern "./ui/html/pages/*.tmpl"
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	// Loop through the page filepaths one-by-one
	for _, page := range pages {
		
		// extract the file name (like 'home.tmpl') from the full filepath
		name := filepath.Base(page)

		// the template.FuncMap must be registered with the template set before calling the ParseFiles() method
		// to do so, use template.New() to create an empty template set, use the Funcs() method to register the
		// template.FuncMap, and then parse the file as normal
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}

		// Call ParseGlob() *on this template set* to add any partials
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		// Call ParseFiles() *on this template set* to add the page template
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Add the template set to the map, using the name of the page as the key
		cache[name] = ts
	}

	return cache, nil
}