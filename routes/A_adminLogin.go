package routes

import (
	"html/template"
	"log"
	"net/http"
	
	"github.com/kamukwamba/oerisuniversity/dbcode"
)


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



func AdminDashboard(w http.ResponseWriter, r *http.Request) {
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	if r.Method == "POST" {
		r.ParseForm()
		

		adminList := dbcode.AdminGet()
		var cardDataAvailable bool

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
			
			err := tpl.ExecuteTemplate(w, "A_adminDashboard.html", toshow)

			CreateCookie(admin_dataout.First_Name,admin_dataout.ID, w,r)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			err := tpl.ExecuteTemplate(w, "A_adminLoginError.html", nil)

			if err != nil {
				log.Fatal(err)
			}
		}

	} else {

		
		err := tpl.ExecuteTemplate(w, "A_adminLoginError.html", nil)

		if err != nil {
			log.Fatal(err)
		}

	}


}