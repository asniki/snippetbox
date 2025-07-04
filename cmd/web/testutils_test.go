package main

import (
	"asniki/snippetbox/internal/models/mocks"
	"bytes"
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

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
)

// csrfTokenRX is a regular expression which captures the CSRF token value from the HTML for the user signup page
var csrfTokenRX = regexp.MustCompile(`<input type='hidden' name='csrf_token' value='(.+)'>`)

// newTestApplication returns an instance of the application struct containing mocked dependencies
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
		snippets:       &mocks.SnippetModel{}, // use the mock
		users:          &mocks.UserModel{},    // use the mock
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}
}

// testServer embeds a httptest.Server instance
type testServer struct {
	*httptest.Server
}

// customTransport is a custom http.RoundTripper that adds default headers
type customTransport struct {
	headers          http.Header
	defaultTransport http.RoundTripper
}

// RoundTrip adds the headers to the request (implements the http.RoundTripper)
func (t *customTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	newReq := req.Clone(req.Context())

	for key, values := range t.headers {
		for _, value := range values {
			newReq.Header.Add(key, value)
		}
	}

	return t.defaultTransport.RoundTrip(newReq)
}

// newTestServer initializes and returns a new instance of the custom testServer type
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

	defaultHeaders := http.Header{}
	defaultHeaders.Set("Origin", ts.URL)
	ts.Client().Transport = &customTransport{
		headers:          defaultHeaders,
		defaultTransport: ts.Client().Transport,
	}

	return &testServer{ts}
}

// get makes a GET request to a given url path using the test server client, and returns
// the response status code, headers and body
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
	bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, string(body)
}

// extractCSRFToken extracts the CSRF token (if one exists) from a HTML response body
func extractCSRFToken(t *testing.T, body string) string {
	matches := csrfTokenRX.FindStringSubmatch(body)
	if len(matches) < 2 {
		t.Fatal("no csrf token found in body")
	}

	return html.UnescapeString(string(matches[1]))
}

// postForm sends POST requests to the test server
func (ts *testServer) postForm(t *testing.T, urlPath string, form url.Values) (int, http.Header, string) {

	rs, err := ts.Client().PostForm(ts.URL+urlPath, form)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, string(body)
}
