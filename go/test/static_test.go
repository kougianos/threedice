package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The Java service serves the frontend from src/main/resources/static; these
// pin the behaviour Go's http.FileServer would otherwise get wrong.

func TestStatic_ServesIndexAtRoot(t *testing.T) {
	resp := get(t, "/")
	require.Equal(t, http.StatusOK, resp.status)
	assert.Contains(t, string(resp.body), "<title>ThreeDice</title>")
}

func TestStatic_ServesIndexHtmlDirectly(t *testing.T) {
	// http.FileServer would 301 this to "/"; Spring serves it as-is.
	resp := get(t, "/index.html")
	require.Equal(t, http.StatusOK, resp.status)
	assert.Contains(t, string(resp.body), "<title>ThreeDice</title>")
}

func TestStatic_ServesAssets(t *testing.T) {
	for _, asset := range []string{"/css/style.css", "/js/app.js", "/favicon.svg"} {
		resp := get(t, asset)
		assert.Equal(t, http.StatusOK, resp.status, "asset %s", asset)
		assert.NotEmpty(t, resp.body, "asset %s", asset)
	}
}

func TestStatic_MissingFileReturnsProblemJSON(t *testing.T) {
	// http.FileServer would write a plain-text 404 here.
	resp := get(t, "/nope.html")
	require.Equal(t, http.StatusNotFound, resp.status)
	assert.Contains(t, resp.detail(t), "No endpoint GET /nope.html")
}

func TestStatic_DoesNotEscapeEmbeddedTree(t *testing.T) {
	resp := get(t, "/../go.mod")
	assert.Equal(t, http.StatusNotFound, resp.status)
}
