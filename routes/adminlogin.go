package routes

import (
	"html/template"
	"log"
	"net/http"
)

func AdminLogin(w http.ResponseWriter, r *http.Request) {
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	err := tpl.ExecuteTemplate(w, "adminlogin.html", nil)

	if err != nil {
		log.Fatal(err)
	}
}
