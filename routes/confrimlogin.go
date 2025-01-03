package routes

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/kamukwamba/oerisuniversity/dbcode"
)

type ProgramAvailable struct {
	Available bool
}

type ADMS struct {
	UUID                                    string
	Student_UUID                            string
	First_Name                              string
	Last_Name                               string
	Email                                   string
	Payment_Method                          string
	Paid                                    string
	Accepted                                string
	Student_Results                         string
	Complete                                string
	Creative_Writing                        string
	Understanding_Miracles                  string
	Channeling_skills                       string
	Enneagram                               string
	Mythology_on_Gods_and_Goddess           string
	Herbs                                   string
	Meditation_skills                       string
	Mantras_and_Mudras                      string
	Divinations                             string
	Archetypes                              string
	Basics_in_Research                      string
	Understanding_Propaganda                string
	Great_Spiritual_Teachers                string
	Reprogramming                           string
	Shamanism                               string
	Mystery_Schools_in_the_world            string
	Law_and_Ethics_in_Metaphysical_Sciences string
	Non_Violet_Communication                string
}

type ABDMS struct {
	UUID                             string
	Student_UUID                     string
	First_Name                       string
	Last_Name                        string
	Email                            string
	Payment_Method                   string
	Paid                             string
	Accepted                         string
	Student_Results                  string
	Complete                         string
	Cause_and_Core_Issues_in_Beliefs string
	Emotional_Well_Being             string
	The_Art_of_Breathing             string
	Spiritual_symbols_and_colours    string
	Psychic_Skills                   string
	Shadow_Work                      string
	The_Craft                        string
	Hypnosis_and_Beyond              string
	Mysterious_experiences           string
	Manifestation_skills             string
	Unlocking_Creativity             string
	Transpersonal_counselling        string
	African_Healing_Arts             string
	Ceremonies_of_the_World          string
	Mother_Earth                     string
	The_Art_of_Placement             string
	Chakras_and_Auras                string
	Transforming_personalities       string
	Mayan_Calendar                   string
	Polarity_Therapy                 string
	Introduction_To_Meditation       string
	Health_and_Nutrition             string
	Setting_up_a_business            string
}
type AllStudentCources struct {
	Student_dat StudentInfo
	All         []AllCourceData
}

type AllCourceData struct {
	ProgramStruct
	Cource_Struct []CourceStruct
}

type ProgramStruct struct {
	UUID           string
	Student_UUID   string
	Program_Name   string
	First_Name     string
	Last_Name      string
	Email          string
	Payment_Method string
	Paid           string
	Approved       bool
	Applied        bool
	Completed      bool
	Date           string
	Admin_ID       string
}

type CourceStruct struct {
	UUID                  string
	Student_UUID          string
	Cource_Name           string
	Book                  string
	Module                string
	Video                 string
	Applied               bool
	Approved              bool
	Examined              bool
	Continuorse_Assesment string
	Completed             bool
	Date                  string
}

type StudentCourse struct {
	Available        bool
	StInfo           StudentInfo
	AllCourceDataOut []AllCourceData
	Admin            AdminInfo
}

func ValidateSudent(email_in, password_in string) (bool, string) {
	isstudent := true
	dbread := dbcode.SqlRead()
	stmt, err := dbread.DB.Prepare("select uuid, student_uuid, email, password from studentcridentials where email = ?")

	if err != nil {
		isstudent = false
		fmt.Println("First err")
		log.Fatal(err)
	}

	defer stmt.Close()

	var uuid string
	var student_uuid string
	var email string
	var password string

	err = stmt.QueryRow(email_in).Scan(&uuid, &student_uuid, &email, &password)

	if err != nil {
		fmt.Println("Second err")
		// log.Fatal(err)
		isstudent = false
	}

	if password_in != password {
		isstudent = false
	}

	return isstudent, student_uuid

}

func GetFromADMS(student_uuid string) {

}

func GetFromABDMS(student_uuid string) {

}

func UpdateProgram(student_uuid, table_name string) bool {
	var programcompleted bool

	dbread := dbcode.SqlRead()

	// .dbread"UPDATE artist_t SET check_s = ? WHERE artist_n = ?", "2021-05-20", 42

	approval_string := fmt.Sprintf("UPDATE %s SET completed = ? WHERE student_uuid = ?", table_name)
	stmt, err := dbread.DB.Prepare(approval_string)

	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close()

	_, erre := stmt.Exec(true, student_uuid)

	if erre != nil {
		log.Fatal(err)
	}

	// _, err := dbread.DB.Exec("update  acams set accepted = ? where student_uuid = ?", "true", student_uuid)
	return programcompleted

}

func Update(student_uuid, table_name string) bool {

	dbread := dbcode.SqlRead()
	updated := true

	approval_string := fmt.Sprintf("UPDATE %s SET approved = ? WHERE student_uuid = ?", table_name)

	stmt, err := dbread.DB.Prepare(approval_string)

	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, erre := stmt.Exec(true, student_uuid)

	if erre != nil {
		log.Fatal(err)
	}

	// _, err := dbread.DB.Exec("update  acams set accepted = ? where student_uuid = ?", "true", student_uuid)

	return updated

}

func GetFromADMSOne(student_uuid string) (bool, ADMS) {
	var result bool

	var dataout ADMS

	return result, dataout
}

func GetFromADBMSOne(student_uuid string) (bool, ABDMS) {
	var result bool

	var dataout ABDMS

	return result, dataout
}

func GetStudentPrograms(student_uuid string) []string {
	dbread := dbcode.SqlRead()

	var listout []string

	stmt, err := dbread.DB.Prepare("select program_list from studentprogramlist where student_uuid = ?")

	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var program_list string

	err = stmt.QueryRow(student_uuid).Scan(&program_list)

	trimedlist := strings.Trim(program_list, "[]")
	list_out := strings.Split(trimedlist, ",")

	for _, item := range list_out {
		trimedlistone := strings.Trim(item, "\"")
		trimedlisttwo := strings.Trim(trimedlistone, " \"")

		if len(trimedlisttwo) > 1 {
			listout = append(listout, trimedlisttwo)

		} else {
			continue
		}
	}

	if err != nil {
		fmt.Println("FAILED TO GET STUDENT PROGRAM LIST")
	}

	fmt.Println("listout", listout)

	return listout
}

func StudentPortal(w http.ResponseWriter, r *http.Request) {

	studentuuid := r.URL.Query().Get("student")

	tpl = template.Must(template.ParseGlob("templates/*.html"))
	var programdataout []AllCourceData
	var studentinfo StudentInfo

	studentprogramlist := GetStudentPrograms(studentuuid)

	programdataout, present := GetStudentProgramData(studentprogramlist, studentuuid)
	studentinfo = GetStudentAllDetails(studentuuid)

	students_data := StudentCourse{
		Available:        present,
		StInfo:           studentinfo,
		AllCourceDataOut: programdataout,
	}

	err := tpl.ExecuteTemplate(w, "studentportal.html", students_data)

	if err != nil {
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

}

func ConfirmStudentLogin(w http.ResponseWriter, r *http.Request) {

	// var students_data_acams ACAMS

	//WORK ON STUDENT VALIDAION AND SECURITY CHECK

	tpl = template.Must(template.ParseGlob("templates/*.html"))
	var programdataout []AllCourceData
	var studentinfo StudentInfo

	var present bool
	var setroute string
	r.ParseForm()

	if r.Method == "POST" {

		studentemail := r.FormValue("studentemail")
		password := r.FormValue("studentpassword")

		confirm, studentuuid := ValidateSudent(studentemail, password)
		fmt.Println("Confirmed")

		if confirm {

			// Redirect with query parameters

			http.Redirect(w, r, fmt.Sprintf("/studentportal?student=%s", studentuuid), http.StatusSeeOther)

			return

		} else {
			setroute = "loginerror.html"
		}

	}

	students_data := StudentCourse{
		Available:        present,
		StInfo:           studentinfo,
		AllCourceDataOut: programdataout,
	}

	tpl.ExecuteTemplate(w, setroute, students_data)

}

func StudentProcced(w http.ResponseWriter, r *http.Request) {

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	var setroute string
	r.ParseForm()

	student_uuid := r.PathValue("id")

	dbread := dbcode.SqlRead()
	result := GetStudentPrograms(student_uuid)

	var completed_out bool
	var count int

	for _, item := range result {
		count += 1
		fmt.Println(item)

	}
	var complete bool
	var result_lenth int

	result_lenth = len(result) - 1

	current_program := result[result_lenth]

	fmt.Println("The Current Program is : ", current_program)

	switch current_program {
	case "ACAMS":
		stmt, err := dbread.DB.Prepare("select completed from acams where student_uuid = ?")
		if err != nil {
			log.Fatal(err)
		}

		defer stmt.Close()

		err = stmt.QueryRow(student_uuid).Scan(&complete)

		if err != nil {
			log.Fatal(err)
		}

		if !complete {
			completed_out = false
			setroute = "programdeniedapproval.html"

		} else {

			checkACMS := CheckACMS(student_uuid)

			if checkACMS {
				setroute = "programdenied.html"
			} else {
				completed_out = true
				CreateACMS(student_uuid)
				AddToProgramList("acms", student_uuid)
				setroute = "programad.html"

			}

		}

	case "ACMS":
		stmt, err := dbread.DB.Prepare("select completed from acms where student_uuid = ?")
		if err != nil {
			log.Fatal(err)
		}

		defer stmt.Close()

		err = stmt.QueryRow(student_uuid).Scan(&complete)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("The Student UUID: ", student_uuid)

		if !complete {
			completed_out = false
			setroute = "programdeniedapproval.html"

		} else {

			checkADMS := CheckADMS(student_uuid)

			fmt.Println("Check ADMS: ", checkADMS)

			if checkADMS {
				setroute = "programdenied.html"

			} else {
				completed_out = true
				setroute = "programad.html"
				AddToProgramList("adms", student_uuid)
				CreateADMS(student_uuid)

			}
			fmt.Println("The Current Route is: ", setroute)

		}
	case "ADMS":
		stmt, err := dbread.DB.Prepare("select completed from adms where student_uuid = ?")
		if err != nil {
			log.Fatal(err)
		}

		defer stmt.Close()

		err = stmt.QueryRow(student_uuid).Scan(&complete)

		if err != nil {
			log.Fatal(err)
		}

		if !complete {
			completed_out = false
			setroute = "programdeniedapproval.html"

		} else {

			fmt.Println("Completed")

			checkACMS := CheckABDMS(student_uuid)

			if checkACMS {
				setroute = "programdenied.html"
			} else {
				completed_out = true
				setroute = "programad.html"
				CreateABDMS(student_uuid)
				AddToProgramList("abdms", student_uuid)

			}

		}
	case "ABDMS":
		stmt, err := dbread.DB.Prepare("select completed from abdms where student_uuid = ?")
		if err != nil {
			log.Fatal(err)
		}

		defer stmt.Close()

		err = stmt.QueryRow(student_uuid).Scan(&complete)

		if err != nil {
			log.Fatal(err)
		}

		if !complete {
			completed_out = false
			setroute = "programdeniedapproval.html"

		} else {

			fmt.Println("Completed")

			checkACMS := CheckADMS(student_uuid)

			if checkACMS {
				setroute = "programdenied.html"
			} else {
				completed_out = true
				CreateABDMS(student_uuid)
				AddToProgramList("adms", student_uuid)
				setroute = "programad.html"

			}

		}

		// default:
		// 	setroute = "programacompleted.html"

	}

	tpl.ExecuteTemplate(w, setroute, completed_out)

}
