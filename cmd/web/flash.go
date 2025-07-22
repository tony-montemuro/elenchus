package main

import (
	"net/http"

	"github.com/tony-montemuro/elenchus/internal/ui"
)

func (app *application) addFlashToSession(flash *ui.Flash, r *http.Request) {
	app.sessionManager.Put(r.Context(), "flash", flash)
}

func (app *application) addSuccessFlashToSession(message string, r *http.Request) {
	app.addFlashToSession(ui.GetSuccessFlash(message), r)
}

func (app *application) addErrorFlashToSession(message string, r *http.Request) {
	app.addFlashToSession(ui.GetErrorFlash(message), r)
}
