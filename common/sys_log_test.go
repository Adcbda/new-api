package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormatStartupAccessURLIncludesAppBasePath(t *testing.T) {
	original := AppBasePath
	AppBasePath = "/new-api"
	t.Cleanup(func() {
		AppBasePath = original
	})

	require.Equal(t, "http://localhost:3000/new-api/", formatStartupAccessURL("localhost", "3000"))
	require.Equal(t, "http://192.168.1.10:3000/new-api/", formatStartupAccessURL("192.168.1.10", "3000"))
}

func TestFormatStartupAccessURLRootModePreservesTrailingSlash(t *testing.T) {
	original := AppBasePath
	AppBasePath = ""
	t.Cleanup(func() {
		AppBasePath = original
	})

	require.Equal(t, "http://localhost:3000/", formatStartupAccessURL("localhost", "3000"))
}
