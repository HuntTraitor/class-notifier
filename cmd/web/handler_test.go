package main

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/hunttraitor/class-notifier/internal/assert"
)

func TestPing(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")

	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, body, "OK")
}

func TestHome(t *testing.T) {
	app := newTestApplication(t)

	//test for unauthenticated home page
	unauthTS := newTestServer(t, app.routes())
	defer unauthTS.Close()

	t.Run("Unauthenticated Home Page Visit", func(t *testing.T) {
		code, _, body := unauthTS.get(t, "/")
		assert.Equal(t, code, http.StatusOK)
		assert.StringContains(t, body, `<h1>Log in to view notifications</h1>`)
	})

	//test for authenticated home page
	authTS := newTestServer(t, app.sessionManager.LoadAndSave(app.mockAuthentication(app.routes())))
	defer authTS.Close()
	t.Run("Authenticated Home Page Visit", func(t *testing.T) {
		code, _, body := authTS.get(t, "/")
		assert.Equal(t, code, http.StatusOK)
		assert.StringContains(t, body, `<h1 class="notification-title">Notification List</h1>`)
	})
}

func TestViewClass(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid ID",
			urlPath:  "/class/view/1",
			wantCode: http.StatusOK,
			wantBody: "{ClassID:1 Name:Mock class Link:www.testclass.com Professor:TestProfessor}",
		},
		{
			name:     "Non-existent ID",
			urlPath:  "/class/view/2",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Negative ID",
			urlPath:  "/class/view/-1",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Decimal ID",
			urlPath:  "/class/view/1.34",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "String ID",
			urlPath:  "/class/view/foo",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Empty ID",
			urlPath:  "/class/view",
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.urlPath)
			assert.Equal(t, code, tt.wantCode)

			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}
		})
	}
}

func TestAddNotification(t *testing.T) {
	app := newTestApplication(t)
	//New test server where every route authenticates the user automatically
	ts := newTestServer(t, app.sessionManager.LoadAndSave(app.mockAuthentication(app.routes())))
	defer ts.Close()

	_, _, body := ts.get(t, "/")
	validCSRFToken := extractCSRFToken(t, body)

	const (
		existingNotification = "1"
		validNotification    = "3"
	)

	tests := []struct {
		name         string
		notification string
		csrfToken    string
		wantCode     int
	}{
		{
			name:         "Successful submission",
			notification: validNotification,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusSeeOther,
		},
		{
			name:         "Invalid CSRFtoken",
			notification: validNotification,
			csrfToken:    "wrongToken",
			wantCode:     http.StatusBadRequest,
		},
		{
			name:         "Duplicate submission",
			notification: existingNotification,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusSeeOther,
		},
		{
			name:         "Invlid Submission",
			notification: "2",
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("classId", tt.notification)
			form.Add("csrf_token", tt.csrfToken)

			code, _, _ := ts.postForm(t, "/notification/add", form)
			assert.Equal(t, code, tt.wantCode)
		})
	}
}

func TestDeleteNotification(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.sessionManager.LoadAndSave(app.mockAuthentication(app.routes())))
	defer ts.Close()

	_, _, body := ts.get(t, "/")
	validCSRFToken := extractCSRFToken(t, body)

	const (
		validNotification = "1"
	)

	tests := []struct {
		name            string
		notificationURL string
		csrfToken       string
		wantCode        int
		wantBody        string
	}{
		{
			name:            "Successful deletion",
			notificationURL: "/notification/delete/1",
			csrfToken:       validCSRFToken,
			wantCode:        http.StatusOK,
		},
		{
			name:            "Invalid CSRFToken",
			notificationURL: "/notification/delete/1",
			csrfToken:       "wrongToken",
			wantCode:        http.StatusBadRequest,
			wantBody:        "Bad Request",
		},
		{
			name:            "Invalid Deletion",
			notificationURL: "/notification/delete/2",
			csrfToken:       validCSRFToken,
			wantCode:        http.StatusInternalServerError,
			wantBody:        "Internal Server Error",
		},
		{
			name:            "Invalid ID",
			notificationURL: "/notification/delete/InvalidID",
			csrfToken:       validCSRFToken,
			wantCode:        http.StatusNotFound,
			wantBody:        "Not Found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.delete(t, tt.notificationURL, tt.csrfToken)
			assert.Equal(t, code, tt.wantCode)
			assert.Equal(t, body, tt.wantBody)
		})
	}
}

func TestUserSignup(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	//make get request and extract the token from it
	_, _, body := ts.get(t, "/user/signup")
	validCSRFToken := extractCSRFToken(t, body)

	const (
		validName            = "Hunter"
		validPassword        = "validPa$$word"
		validConfirmPassword = "validPa$$word"
		validEmail           = "htratar@ucsc.edu"
		formTag              = `<form action="/user/signup" method="POST" novalidate>`
	)

	tests := []struct {
		name                string
		userName            string
		userEmail           string
		userPassword        string
		userConfirmPassword string
		csrfToken           string
		wantCode            int
		wantFormTag         string
	}{
		{
			name:                "Valid submission",
			userName:            validName,
			userEmail:           validEmail,
			userPassword:        validPassword,
			userConfirmPassword: validConfirmPassword,
			csrfToken:           validCSRFToken,
			wantCode:            http.StatusSeeOther,
		},
		{
			name:                "Invalid CSRF Token",
			userName:            validName,
			userEmail:           validEmail,
			userPassword:        validPassword,
			userConfirmPassword: validConfirmPassword,
			csrfToken:           "wrongToken",
			wantCode:            http.StatusBadRequest,
		},
		{
			name:                "Empty name",
			userName:            "",
			userEmail:           validEmail,
			userPassword:        validPassword,
			userConfirmPassword: validConfirmPassword,
			csrfToken:           validCSRFToken,
			wantCode:            http.StatusUnprocessableEntity,
			wantFormTag:         formTag,
		},
		{
			name:                "Empty email",
			userName:            validName,
			userEmail:           "",
			userPassword:        validPassword,
			userConfirmPassword: validConfirmPassword,
			csrfToken:           validCSRFToken,
			wantCode:            http.StatusUnprocessableEntity,
			wantFormTag:         formTag,
		},
		{
			name:                "Empty password",
			userName:            validName,
			userEmail:           validEmail,
			userPassword:        "",
			userConfirmPassword: validConfirmPassword,
			csrfToken:           validCSRFToken,
			wantCode:            http.StatusUnprocessableEntity,
			wantFormTag:         formTag,
		},
		{
			name:                "Empty confirmation password",
			userName:            validName,
			userEmail:           validEmail,
			userPassword:        validPassword,
			userConfirmPassword: "",
			csrfToken:           validCSRFToken,
			wantCode:            http.StatusUnprocessableEntity,
			wantFormTag:         formTag,
		},
		{
			name:                "Invlid email",
			userName:            validName,
			userEmail:           "htratar@ucsc.",
			userPassword:        validPassword,
			userConfirmPassword: validConfirmPassword,
			csrfToken:           validCSRFToken,
			wantCode:            http.StatusUnprocessableEntity,
			wantFormTag:         formTag,
		},
		{
			name:                "Short password",
			userName:            validName,
			userEmail:           validEmail,
			userPassword:        "pa$$",
			userConfirmPassword: "pa$$",
			csrfToken:           validCSRFToken,
			wantCode:            http.StatusUnprocessableEntity,
			wantFormTag:         formTag,
		},
		{
			name:                "Passwords don't match",
			userName:            validName,
			userEmail:           validEmail,
			userPassword:        validPassword,
			userConfirmPassword: "validpa$$$word",
			csrfToken:           validCSRFToken,
			wantCode:            http.StatusUnprocessableEntity,
			wantFormTag:         formTag,
		},
		{
			name:                "Duplicate email",
			userName:            validName,
			userEmail:           "dupe@example.com",
			userPassword:        validPassword,
			userConfirmPassword: validConfirmPassword,
			csrfToken:           validCSRFToken,
			wantCode:            http.StatusUnprocessableEntity,
			wantFormTag:         formTag,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("name", tt.userName)
			form.Add("email", tt.userEmail)
			form.Add("password", tt.userPassword)
			form.Add("confirmedpassword", tt.userConfirmPassword)
			form.Add("csrf_token", tt.csrfToken)

			code, _, body := ts.postForm(t, "/user/signup", form)
			assert.Equal(t, code, tt.wantCode)

			if tt.wantFormTag != "" {
				assert.StringContains(t, body, tt.wantFormTag)
			}
		})
	}
}
