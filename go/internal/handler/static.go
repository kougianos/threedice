package handler

import (
	"io"
	"io/fs"
	"net/http"
	"path"
)

// staticHandler serves the embedded frontend the way Spring serves
// src/main/resources/static.
//
// http.FileServer is deliberately not used: it 301s /index.html to / and writes
// a plain-text "404 page not found", where Spring serves /index.html directly
// and reports misses as RFC 7807.
func staticHandler(static fs.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Clean first so "/../secret" cannot escape the tree; fs.FS rejects the
		// rest by refusing any path that is not slash-separated and unrooted.
		name := path.Clean(r.URL.Path)[1:]
		if name == "" {
			name = "index.html"
		}

		f, err := static.Open(name)
		if err != nil {
			NotFound(w, r)
			return
		}
		defer f.Close()

		info, err := f.Stat()
		if err != nil || info.IsDir() {
			NotFound(w, r)
			return
		}

		seeker, ok := f.(io.ReadSeeker)
		if !ok {
			NotFound(w, r)
			return
		}

		// ServeContent picks the Content-Type from the extension and handles
		// conditional and range requests.
		http.ServeContent(w, r, info.Name(), info.ModTime(), seeker)
	}
}
