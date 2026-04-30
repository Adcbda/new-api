package common

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalizeAppBasePath(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "empty uses root", input: "", want: ""},
		{name: "slash uses root", input: "/", want: ""},
		{name: "adds leading slash", input: "new-api", want: "/new-api"},
		{name: "trims trailing slash", input: "/new-api/", want: "/new-api"},
		{name: "cleans duplicate slashes", input: "//new-api//console/", want: "/new-api/console"},
		{name: "rejects traversal", input: "/../new-api", wantErr: true},
		{name: "rejects query", input: "/new-api?x=1", wantErr: true},
		{name: "rejects fragment", input: "/new-api#top", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NormalizeAppBasePath(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestWithAppBasePathStripsPrefixAndRejectsRootPath(t *testing.T) {
	original := AppBasePath
	AppBasePath = "/new-api"
	t.Cleanup(func() {
		AppBasePath = original
	})

	called := false
	handler := WithAppBasePath(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		require.Equal(t, "/api/status", r.URL.Path)
		require.Equal(t, "x=1", r.URL.RawQuery)
		require.Equal(t, "/api/status?x=1", r.RequestURI)
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/new-api/api/status?x=1", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	require.True(t, called)
	require.Equal(t, http.StatusNoContent, rec.Code)

	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/status", nil))
	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestWithAppBasePathRedirectsExactPrefixToSlash(t *testing.T) {
	original := AppBasePath
	AppBasePath = "/new-api"
	t.Cleanup(func() {
		AppBasePath = original
	})

	handler := WithAppBasePath(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("inner handler should not be called for exact base path")
	}))

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/new-api?x=1", nil))
	require.Equal(t, http.StatusMovedPermanently, rec.Code)
	require.Equal(t, "/new-api/?x=1", rec.Header().Get("Location"))
}

func TestBuildPublicURLUsesAppBasePathWithoutDuplicatingIt(t *testing.T) {
	original := AppBasePath
	AppBasePath = "/new-api"
	t.Cleanup(func() {
		AppBasePath = original
	})

	require.Equal(
		t,
		"https://example.com/new-api/console/topup?show_history=true",
		BuildPublicURL("https://example.com", "/console/topup?show_history=true"),
	)
	require.Equal(
		t,
		"https://example.com/new-api/console/topup",
		BuildPublicURL("https://example.com/new-api", "/console/topup"),
	)
	require.Equal(
		t,
		"https://example.com/root/new-api/console/topup",
		BuildPublicURL("https://example.com/root", "/console/topup"),
	)
}

func TestBuildPublicURLRootModePreservesExistingBehavior(t *testing.T) {
	original := AppBasePath
	AppBasePath = ""
	t.Cleanup(func() {
		AppBasePath = original
	})

	require.Equal(t, "https://example.com/console/log", BuildPublicURL("https://example.com", "/console/log"))
}
