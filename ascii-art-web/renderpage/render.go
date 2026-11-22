package renderpage

import (
	"bytes"
	"html/template"
	"net/http"
)

type PageData struct {
	Result string
}

var tmpl *template.Template

func init() {
	var err error
	tmpl, err = template.ParseFiles("templates/index.html")
	if err != nil {
		panic(err)
	}
}

func RenderPage(w http.ResponseWriter, data PageData) {
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	buf.WriteTo(w)
}
