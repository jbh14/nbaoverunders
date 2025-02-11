package main

import (
	"errors"
	"fmt"
	//"html/template"
	"net/http"
	"strconv"

	"github.com/jbh14/nbaoverunders/internal/models"
)

// home is where "all entries" will show up
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	entries, err := app.entries.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	for _, entry := range entries {
		fmt.Fprintf(w, "%+v\n", entry)
	}

	// will uncomment later, just getting entry view working right now
	//files := []string{
	//	"./ui/html/base.tmpl",
	//	"./ui/html/partials/nav.tmpl",
	//	"./ui/html/pages/home.tmpl",
	//}

	//ts, err := template.ParseFiles(files...)
	//if err != nil {
	//	app.serverError(w, r, err)
	//	return
	//}

	//err = ts.ExecuteTemplate(w, "base", nil)
	//if err != nil {
	//	app.serverError(w, r, err)
	//}
}

func (app *application) entryView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	
	entry, err := app.entries.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	// Write the entry data as a plain-text HTTP response body.
	fmt.Fprintf(w, "%+v", entry)
}

func (app *application) entryCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	// Create some variables holding dummy data. We'll remove these later on
	// during the build.
	playername := "Dan Phelan"
	year := 2025

	// Logging for troubleshooting
	app.logger.Info(fmt.Sprintf("Attempting to insert entry: playername=%s, year=%d", playername, year))

	// Pass the data to the EntryModel.Insert() method, receiving the
	// ID of the new record back.
	id, err := app.entries.Insert(playername, year)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.logger.Info(fmt.Sprintf("Successfully inserted entry, assigned id=%d", id))
	http.Redirect(w, r, fmt.Sprintf("/entry/view?id=%d", id), http.StatusSeeOther)
}