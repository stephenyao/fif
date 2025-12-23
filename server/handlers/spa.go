package handlers

import (
	"io"
	"io/fs"
	"net/http"
	"strings"
)

// SPAHandler returns an http.HandlerFunc that serves static files from the provided
// filesystem and falls back to index.html for any routes that aren't found.
func SPAHandler(staticFS fs.FS) http.HandlerFunc {
	fileServer := http.FileServer(http.FS(staticFS))

	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// If the file exists in the FS, serve it
		f, err := staticFS.Open(strings.TrimPrefix(path, "/"))
		if err == nil {
			f.Close()
			fileServer.ServeHTTP(w, r)
			return
		}

		// Fallback to index.html for SPA routes
		index, err := staticFS.Open("index.html")
		if err != nil {
			http.NotFound(w, r)
			return
		}
		defer index.Close()

		stat, err := index.Stat()
		if err != nil {
			http.NotFound(w, r)
			return
		}

		http.ServeContent(w, r, "index.html", stat.ModTime(), index.(io.ReadSeeker))
	}
}
