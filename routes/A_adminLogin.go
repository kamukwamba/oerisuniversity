package routes

import (
	"html/template"
	"log"
	"net/http"

	"github.com/kamukwamba/oerisuniversity/dbcode"
)

// Login PAGE FOR ADMIN LOG IN
type AdminLandingData struct {
	Admin         AdminInfo
	ProgramD      []ProgramDataEntry
	DataAvailable bool
}

func AdminLogin(w http.ResponseWriter, r *http.Request) {
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	err := tpl.ExecuteTemplate(w, "A_adminlogin.html", nil)

	if err != nil {
		log.Fatal(err)
	}
}

// ADMIN DASH BOARD
func IsHXRequest(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}
func AdminDashboard(w http.ResponseWriter, r *http.Request) {
	tpl = template.Must(template.ParseGlob("templates/*.html"))
	var cardDataAvailable bool
	if r.Method == "POST" {
		r.ParseForm()
		adminList := dbcode.AdminGet()

		email := r.PostFormValue("email")
		password := r.PostFormValue("password")

		authget := AdminLogData{
			Email:    email,
			Password: password,
		}

		check, admin_dataout := AdminAuth(authget, adminList)

		cardData, errGAPD := GetAllProgramData()

		if errGAPD != nil {
			cardDataAvailable = false
		} else {
			cardDataAvailable = true
		}

		toshow := AdminLandingData{
			Admin:         admin_dataout,
			ProgramD:      cardData,
			DataAvailable: cardDataAvailable,
		}

		if check {

			if IsHXRequest(r) {
				w.Header().Set("HX-Redirect", "/dashboard")
				w.WriteHeader(http.StatusOK)
				return
			}

			err := tpl.ExecuteTemplate(w, "A_adminDasboard.html", toshow)

			if err != nil {
				log.Fatal(err)
			}
		} else {
			err := tpl.Execute(w, "Invalid username or password")

			if err != nil {
				log.Fatal(err)
			}
		}

	} else {
		err := tpl.Execute(w, "Invalid username or password")

		if err != nil {
			log.Fatal(err)
		}

	}

}
