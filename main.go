package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServerHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join(
			"templates",
			t.filename,
		)))

		t.templ.Execute(w, nil)
	})
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
            <html>
                <head>
                    <title>Chat</title>
                </head>
                <body>
                    Let's chat!
                </body>
            </html>
        `))
	})

	if err := http.ListenAndServe("127.0.0.1:8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
