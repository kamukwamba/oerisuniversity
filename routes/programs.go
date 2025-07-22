package routes

import (
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
	var programsAvailable bool

	out := r.URL.Query().Get("out")

	admin_infor := AdminData(out)
	programdata, errout := GetAllProgramData()

	if errout != nil {
		programsAvailable = false
	} else {
		programsAvailable = true
	}

	data_out := AdminLandingData{
		Admin:         admin_infor,
		ProgramD:      programdata,
		DataAvailable: programsAvailable,
	}

	//debug failure to laod templates

	err := tpl.ExecuteTemplate(w, "programcards.html", data_out)

	if err != nil {
		log.Fatal(err)
	}
}
