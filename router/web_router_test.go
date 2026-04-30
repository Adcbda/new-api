package router

import (
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/stretchr/testify/require"
)

func TestPrepareIndexPageInjectsAppBasePathRuntime(t *testing.T) {
	original := common.AppBasePath
	common.AppBasePath = "/new-api"
	t.Cleanup(func() {
		common.AppBasePath = original
	})

	input := []byte(`<!doctype html><html><head><link href="/static/app.css"><script src="/static/app.js"></script></head><body><div id="root"></div></body></html>`)

	got := string(prepareIndexPage(input))

	require.Contains(t, got, `<base href="/new-api/">`)
	require.Contains(t, got, `window.__APP_BASE_PATH__="/new-api"`)
	require.Contains(t, got, `href="/new-api/static/app.css"`)
	require.Contains(t, got, `src="/new-api/static/app.js"`)
}

func TestPrepareIndexPageLeavesRootModeUntouched(t *testing.T) {
	original := common.AppBasePath
	common.AppBasePath = ""
	t.Cleanup(func() {
		common.AppBasePath = original
	})

	input := []byte(`<!doctype html><html><head><link href="/static/app.css"></head><body></body></html>`)

	require.Equal(t, string(input), string(prepareIndexPage(input)))
}
