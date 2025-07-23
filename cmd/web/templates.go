package main

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"

	"github.com/justinas/nosurf"
	uicomponents "github.com/tony-montemuro/elenchus/internal/ui"
	"github.com/tony-montemuro/elenchus/internal/validator"
	"github.com/tony-montemuro/elenchus/ui"
)

type templateData struct {
	Form            any
	RangeRules      validator.FormRangeRules
	Flash           *uicomponents.Flash
	IsAuthenticated bool
	CSRFToken       string
	Script          string
	Data            any
}

var functions = template.FuncMap{
	"humanDate":  humanDate,
	"timeAgo":    timeAgo,
	"pluralize":  pluralize,
	"percentage": percentage,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"html/base.tmpl",
			"html/partials/*.tmpl",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}

func (app *application) newTemplateData(r *http.Request) templateData {
	var flash *uicomponents.Flash

	f, ok := app.sessionManager.Pop(r.Context(), "flash").(*uicomponents.Flash)
	if ok {
		flash = f
	} else {
		flash = nil
	}

	return templateData{
		Flash:           flash,
		IsAuthenticated: app.isAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
	}
}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("template %s does not exist", page)
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(status)

	err := ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
}
