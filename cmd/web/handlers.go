package main

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/heisenberg8055/gosts/internal/models"
	"github.com/heisenberg8055/gosts/internal/validator"
)

type neuteredFileSystem struct {
	fs http.FileSystem
}

type snippet struct {
	Title   string
	Content string
	Expires int
	validator.Validator
}

type userLoginForm struct {
	Email    string
	Password string
	validator.Validator
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.html", data)
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = snippet{
		Expires: 365,
	}
	app.render(w, http.StatusOK, "create.html", data)
}

func (app *application) createSnippetPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := snippet{
		Title:   r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
		Expires: expires,
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedValues(form.Expires, 1, 7, 365), "expires", "This field must equals 1, 7 or 365")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.html", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", fmt.Sprintf("Snippet %d successfully created!", id))
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (app *application) viewSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, http.StatusOK, "view.html", data)
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}
	s, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}
			return nil, err
		}
	}
	return f, nil
}

type userSignupForm struct {
	Name     string
	Email    string
	Password string
	validator.Validator
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, http.StatusOK, "signup.html", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	formData := userSignupForm{
		Name:     r.PostForm.Get("name"),
		Email:    r.PostForm.Get("email"),
		Password: r.PostForm.Get("password"),
	}

	formData.CheckField(validator.NotBlank(formData.Name), "name", "This field cannot be blank")
	formData.CheckField(validator.NotBlank(formData.Email), "email", "This field cannot be blank")
	formData.CheckField(validator.Matches(formData.Email, validator.EmailRegex), "email", "This field must be a valid email address")
	formData.CheckField(validator.NotBlank(formData.Password), "password", "This field cannot be blank")
	formData.CheckField(validator.MinChars(formData.Password, 8), "password", "This field must be alteast 8 character long")
	formData.CheckField(validator.MaxChars(formData.Password, 36), "password", "This field must be less than 36 character long")

	if !formData.Valid() {
		data := app.newTemplateData(r)
		data.Form = formData
		app.render(w, http.StatusUnprocessableEntity, "signup.html", data)
		return
	}

	err = app.users.Insert(formData.Name, formData.Email, formData.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			formData.AddFieldError("email", "Email address is already in use")
			data := app.newTemplateData(r)
			data.Form = formData
			app.render(w, http.StatusUnprocessableEntity, "signup.html", data)
		} else {
			app.serverError(w, err)
		}
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "Your Signup is successful. Please log in.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, http.StatusOK, "login.html", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}
	formData := userLoginForm{
		Email:    r.PostForm.Get("email"),
		Password: r.PostForm.Get("password"),
	}
	formData.CheckField(validator.NotBlank(formData.Email), "email", "This field cannot be blank")
	formData.CheckField(validator.Matches(formData.Email, validator.EmailRegex), "email", "This field must be a valid email address")
	formData.CheckField(validator.NotBlank(formData.Password), "password", "This field cannot be blank")

	if !formData.Valid() {
		data := app.newTemplateData(r)
		data.Form = formData
		app.render(w, http.StatusUnprocessableEntity, "login.html", data)
		return
	}

	id, err := app.users.Authenticate(formData.Email, formData.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			formData.AddNonFieldError("Email or password is incorrent")
			data := app.newTemplateData(r)
			data.Form = formData
			app.render(w, http.StatusUnprocessableEntity, "login.html", data)
		} else {
			app.serverError(w, err)
		}
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	fmt.Println(id)
	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")

	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
