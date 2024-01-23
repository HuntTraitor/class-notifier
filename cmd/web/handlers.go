package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hunttraitor/class-notifier/internal/models"
	"github.com/hunttraitor/class-notifier/internal/models/validator"

	"github.com/julienschmidt/httprouter"
)

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	ConfirmedPassword   string `form:"confirmedpassword"`
	validator.Validator `form:"-"`
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	//fetches classes from database for every call to home even if user is not authenticated
	classes, err := app.classes.Classlist()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	if !app.isAuthenticated(r) {
		data := app.newTemplateData(r)
		data.Classes = classes
		app.render(w, r, http.StatusOK, "home.html", data)
		return
	}

	//fetches email from session
	email := app.sessionManager.GetString(r.Context(), "authenticatedUserEmail")

	//Fetches notifications from database for user
	notifications, err := app.notifications.NotificationList(email)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Classes = classes
	data.Notifications = notifications
	app.render(w, r, http.StatusOK, "home.html", data)
}

// Sends a post request to add the class
func (app *application) addClass(w http.ResponseWriter, r *http.Request) {

	classid := 31139
	name := "CSE 102 - 01   Introduction to Analysis of Algorithms"
	link := "https://pisa.ucsc.edu/class_search/index.php?action=detail&class_data=YToyOntzOjU6IjpTVFJNIjtzOjQ6IjIyNDAiO3M6MTA6IjpDTEFTU19OQlIiO3M6NToiMzExMzkiO30%253D"
	professor := "Chatziafratis,E."

	id, err := app.classes.Insert(classid, name, link, professor)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/class/view?id=%d", id), http.StatusSeeOther)
}

// redircts to the link for the class
func (app *application) viewClass(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.notFound(w)
		return
	}

	class, err := app.classes.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	fmt.Fprintf(w, "%+v", class)
}

func (app *application) addNotification(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	classid, err := strconv.Atoi(r.PostFormValue("classId"))
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	expires := 7

	//gets email from user session
	email := app.sessionManager.GetString(r.Context(), "authenticatedUserEmail")

	err = app.notifications.Insert(email, classid, expires)
	//check if error violates unique constraint
	if err != nil {
		if errors.Is(err, models.ErrDuplicateNotification) {
			app.sessionManager.Put(r.Context(), "flash", DuplicateNotificationFlashError)
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Class notification successfully added!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) deleteNotification(w http.ResponseWriter, r *http.Request) {

	params := httprouter.ParamsFromContext(r.Context())

	notificationid, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.notFound(w)
		return
	}

	err = app.notifications.Delete(notificationid)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Class notification successfully deleted!")

	//redirect using HTMX
	app.redirect(w)
}

// will add a page later that displays these nicely!
func (app *application) viewNotifications(w http.ResponseWriter, r *http.Request) {

	classes, err := app.notifications.NotificationList("htratar@ucsc.edu")

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	for _, class := range classes {
		fmt.Fprintf(w, "%+v\n", class)
	}
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, r, http.StatusOK, "signup.html", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm

	//decodes form to check if its a valid form
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	//setting up validators for form fields
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field msut be at least 8 characters long")
	form.CheckField(validator.NotBlank(form.ConfirmedPassword), "confirmedpassword", "This field cannot be blank")
	form.CheckField(validator.Equal(form.Password, form.ConfirmedPassword), "confirmedpassword", "Password must match")

	//If there is a validation error post form with errors
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "signup.html", data)
		return
	}

	err = app.users.Insert(form.Name, form.Email, form.Password)
	//if the error violates the unique constraint in the DB
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "signup.html", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, r, http.StatusOK, "login.html", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "login.html", data)
		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "login.html", data)
			return
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	//RenewToken() to change the session ID. We want to generate a new session any time a user does privilege changes i.e. login/logout
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	//add a authenticatedUserID field to the session to save login
	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)
	app.sessionManager.Put(r.Context(), "authenticatedUserEmail", form.Email)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	//Renew token on current session to change session ID
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	//remove authenticatedUserID on logout
	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	app.sessionManager.Remove(r.Context(), "authenticatedUserEmail")
	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
