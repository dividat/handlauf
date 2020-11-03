package main

import (
	"mime"
	"net/http"
	"path"
	"strings"
)

//go:generate go run github.com/256dpi/embed -strings -pkg main -out ui_files.go debug/build

func uiHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check method
		if r.Method != "GET" {
			http.Error(w, "", http.StatusMethodNotAllowed)
			return
		}

		// remove leading slash
		pth := strings.TrimPrefix(r.URL.Path, "/")

		// get content
		content, ok := files[pth]
		if !ok {
			pth = "index.html"
			content, _ = files[pth]
		}

		// get mime type
		mimeType := mime.TypeByExtension(path.Ext(pth))

		// set content type
		w.Header().Set("Content-Type", mimeType)

		// write file
		_, _ = w.Write([]byte(content))
	})
}
