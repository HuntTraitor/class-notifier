package main

import (
	"bytes"
	"context"
	"html"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"
	"time"

	"github.com/hunttraitor/class-notifier/internal/models/mocks"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
)

// regex that captures csrf token value from html
var csrfTokenRX = regexp.MustCompile(`<input type="hidden" name="csrf_token" value="(.+)">`)

func extractCSRFToken(t *testing.T, body string) string {
	//returns array of matched expression and then the value of the expression
	matches := csrfTokenRX.FindStringSubmatch(body)
	if len(matches) < 2 {
		t.Fatal("no csrf token found in body")
	}
	return html.UnescapeString(string(matches[1]))
}

func newTestApplication(t *testing.T) *application {
	templateCache, err := newTemplateCache()
	if err != nil {
		t.Fatal(err)
	}

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	return &application{
		logger:         slog.New(slog.NewTextHandler(io.Discard, nil)),
		classes:        &mocks.ClassModel{},
		notifications:  &mocks.NotificationModel{},
		users:          &mocks.UserModel{},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}
}

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	ts.Client().Jar = jar

	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	return &testServer{ts}
}

func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, string) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)
	return rs.StatusCode, rs.Header, string(body)
}

func (ts *testServer) delete(t *testing.T, urlPath string, csrfToken string) (int, http.Header, string) {
	req, err := http.NewRequest(http.MethodDelete, ts.URL+urlPath, nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("X-CSRF-Token", csrfToken)
	rs, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)
	return rs.StatusCode, rs.Header, string(body)
}

func (ts *testServer) postForm(t *testing.T, urlPath string, form url.Values) (int, http.Header, string) {
	rs, err := ts.Client().PostForm(ts.URL+urlPath, form)
	if err != nil {
		t.Fatal(err)
	}

	//Read response from body
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, string(body)
}

func (app *application) mockAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), isAuthenticatedContextKey, true)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
