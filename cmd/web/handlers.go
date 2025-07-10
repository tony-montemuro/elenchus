package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/tony-montemuro/elenchus/internal/models"
	"github.com/tony-montemuro/elenchus/internal/validator"
)

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, r, http.StatusOK, "home.tmpl", data)
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = loginForm{}
	data.RangeRules = validator.RangeRules[validator.LoginForm]
	app.render(w, r, http.StatusOK, "login.tmpl", data)
}

func (app *application) signup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = signupForm{}
	data.RangeRules = validator.RangeRules[validator.SignUpForm]
	app.render(w, r, http.StatusOK, "signup.tmpl", data)
}

func (app *application) signupPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := signupForm{
		FirstName: r.PostForm.Get("first-name"),
		LastName:  r.PostForm.Get("last-name"),
		Email:     r.PostForm.Get("email"),
		Password:  r.PostForm.Get("password"),
		Password2: r.PostForm.Get("password2"),
	}

	formName := validator.SignUpForm
	rangeRules := validator.RangeRules[formName]
	errs := validator.GetRangeErrors(form, formName)
	for _, err := range errs {
		form.AddFieldError(err.Key, err.Error())
	}
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address.")
	form.CheckField(form.Password == form.Password2, "password", "Passwords do not match.")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		data.RangeRules = rangeRules
		app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}

	err = app.profiles.Insert(form.FirstName, form.LastName, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address already in use.")
			data := app.newTemplateData(r)
			data.Form = form
			data.RangeRules = rangeRules
			app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl", data)
			return
		} else {
			app.serverError(w, r, err)
		}

		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful! Please log in.")

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *application) loginPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := loginForm{
		Email:    r.PostForm.Get("email"),
		Password: r.PostForm.Get("password"),
	}

	formName := validator.LoginForm
	errs := validator.GetRangeErrors(form, formName)
	for _, err := range errs {
		form.AddFieldError(err.Key, err.Error())
	}
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address.")

	rangeRules := validator.RangeRules[formName]
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		data.RangeRules = rangeRules
		app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}

	profile, err := app.profiles.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect.")

			data := app.newTemplateData(r)
			data.Form = form
			data.RangeRules = rangeRules
			app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), authenticatedUserIdKey, profile.ID)
	app.sessionManager.Put(r.Context(), "flash", fmt.Sprintf("Login successful! Welcome, %s!", profile.FirstName))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) logoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Remove(r.Context(), authenticatedUserIdKey)
	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) quizList(w http.ResponseWriter, r *http.Request) {
	quizzes, err := app.quizzes.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Data = QuizzesPageData{
		Quizzes: quizzes,
	}
	app.render(w, r, http.StatusOK, "quizzes.tmpl", data)
}

func (app *application) create(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = createForm{}
	data.RangeRules = validator.RangeRules[validator.CreateForm]
	app.render(w, r, http.StatusOK, "create.tmpl", data)
}

func (app *application) createPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := createForm{
		Notes: r.PostForm.Get("notes"),
	}

	formName := validator.CreateForm
	errs := validator.GetRangeErrors(form, formName)
	for _, err := range errs {
		form.AddFieldError(err.Key, err.Error())
	}

	rangeRules := validator.RangeRules[formName]
	data := app.newTemplateData(r)
	data.Form = form
	data.RangeRules = rangeRules

	if !form.Valid() {
		app.render(w, r, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	quiz, err := app.generateQuiz(form.Notes, r.Context())
	if err != nil {
		if errors.Is(err, ErrGenerationRefusal) {
			form.AddFieldError("notes", err.Error())
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "create.tmpl", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	profileID, err := app.getProfileID(r)
	if err != nil {
		app.serverError(w, r, errors.New("user attempted to create quiz without proper authorization!"))
		return
	}
	id, err := app.quizzesService.UploadQuiz(quiz, *profileID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data.Flash = "Quiz created!"

	http.Redirect(w, r, fmt.Sprintf("/quizzes/%d", id), http.StatusSeeOther)
}

func (app *application) quiz(w http.ResponseWriter, r *http.Request) {
	quizID, err := strconv.Atoi(r.PathValue("quizID"))
	if err != nil {
		app.redirectNotFound(w, r, "user attempted to access a quiz that does not exist", err)
		return
	}

	profileID, _ := app.getProfileID(r)
	quiz, err := app.quizzesService.GetQuizByID(quizID, profileID)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.redirectNotFound(w, r, "user attempted to access a quiz that does not exist", err)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Data = QuizPageData{
		Quiz: quiz,
	}
	app.render(w, r, http.StatusOK, "quiz.tmpl", data)
}

func (app *application) quizPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%v", r.PathValue("quizID"))
}

func (app *application) profile(w http.ResponseWriter, r *http.Request) {
	profileID, err := app.getProfileID(r)
	if err != nil {
		app.redirectNotFound(w, r, "user attempted to access profile page without proper authorization!", err)
		return
	}

	published, err := app.quizzes.GetPublishedQuizzesByProfile(profileID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	unpublished, err := app.quizzes.GetUnpublishedQuizzesByProfile(profileID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Data = ProfilePageData{
		Published:   published,
		Unpublished: unpublished,
	}
	data.Script = "profile.js"

	app.render(w, r, http.StatusOK, "profile.tmpl", data)
}

func (app *application) edit(w http.ResponseWriter, r *http.Request) {
	quizID, err := strconv.Atoi(r.PathValue("quizID"))
	if err != nil {
		app.redirectNotFound(w, r, "user attempted to edit a quiz that does not exist", err)
		return
	}

	quiz, err := app.getQuizByID(quizID, r)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.redirectNotFound(w, r, "user attempted to access a quiz that does not exist", err)
		} else if errors.Is(err, ErrNotEditable) {
			app.redirectHome(w, r, "Quiz cannot be edited.", err.Error(), err)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Form = editForm{}
	data.Data = QuizPageData{
		Quiz: quiz,
	}
	data.RangeRules = validator.RangeRules[validator.EditForm]

	app.render(w, r, http.StatusOK, "edit.tmpl", data)
}

func (app *application) editPost(w http.ResponseWriter, r *http.Request) {
	quizID, err := strconv.Atoi(r.PathValue("quizID"))
	if err != nil {
		app.redirectNotFound(w, r, "user attempted to edit a quiz that does not exist", err)
		return
	}

	quiz, err := app.getQuizByID(quizID, r)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.redirectNotFound(w, r, "user attempted to access a quiz that does not exist", err)
		} else if errors.Is(err, ErrNotEditable) {
			app.redirectHome(w, r, "Quiz cannot be edited.", err.Error(), err)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	err = r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form, err := newEditForm(r.PostForm)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	formName := validator.EditForm
	errs := validator.GetRangeErrors(form, formName)
	errs = append(errs, validator.GetAggregateFieldRangeErrors(form.serializeQuestionContent(), formName, "question")...)
	errs = append(errs, validator.GetAggregateFieldRangeErrors(form.serializeAnswerContent(), formName, "answer")...)
	for _, err := range errs {
		form.AddFieldError(err.Key, err.Error())
	}
	for key, points := range form.serializeQuestionPoints() {
		form.CheckField(validator.Gte(points, 1), key, "This field cannot be less than 1.")
		form.CheckField(validator.Lte(points, 1000), key, "This field cannot be more than 1000.")
	}

	newQuiz, err := app.buildNewQuizPublic(quiz, form)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	rangeRules := validator.RangeRules[formName]
	data := app.newTemplateData(r)
	data.Data = QuizPageData{
		Quiz: newQuiz,
	}
	data.RangeRules = rangeRules
	data.Form = form

	if !form.Valid() {
		app.render(w, r, http.StatusUnprocessableEntity, "edit.tmpl", data)
		return
	}

	app.render(w, r, http.StatusOK, "edit.tmpl", data)
}
