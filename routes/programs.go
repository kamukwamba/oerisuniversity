package routes

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func Programs(w http.ResponseWriter, r *http.Request) {

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	//debug failure to laod templates

	err := tpl.ExecuteTemplate(w, "programs.html", nil)

	if err != nil {
		log.Fatal(err)
	}
}

func Programcards(w http.ResponseWriter, r *http.Request) {

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	out := r.URL.Query().Get("out")

	fmt.Println("Admin ID: ", out)

	admin_infor := AdminData(out)

	data_out := AdminPage{
		Admin: admin_infor,
	}

	//debug failure to laod templates

	err := tpl.ExecuteTemplate(w, "programcards.html", data_out)

	if err != nil {
		log.Fatal(err)
	}
}
