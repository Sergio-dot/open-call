package render

import (
	"bytes"
	"fmt"
	"github.com/Sergio-dot/open-call/pkg/config"
	"github.com/Sergio-dot/open-call/pkg/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var functions = template.FuncMap{}

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {

	return td
}

// Template renders templates using html/template
func Template(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template

	if app.UseCache {
		// get the template cache from the app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	// get the specified template from template cache
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("could not get template from cache")
	}

	// creates a buffer
	buff := new(bytes.Buffer)

	// add default data to template
	td = AddDefaultData(td)

	// execute the template and store the value in the buffer
	_ = t.Execute(buff, td)

	// write buffer content to response writer
	_, err := buff.WriteTo(w)
	if err != nil {
		fmt.Println("error writing template to browser:", err)
	}
}

// CreateTemplateCache parses all templates, including layouts, and store them
// in a cache of type map[string]*template.Template.
// The string key value represents the name of the template.
func CreateTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// look for all files that include ".page.tmpl" in the name
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return cache, err
	}

	for _, page := range pages {
		// unwrap the file path and take the file name
		name := filepath.Base(page)

		// create the template set and enable func map
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return cache, err
		}

		// look for all files that include ".layout.tmpl" in the name
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return cache, err
		}

		if len(matches) > 0 {
			// parse the layouts
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return cache, err
			}
		}

		// add the template to the cache
		cache[name] = ts
	}

	return cache, nil
}
