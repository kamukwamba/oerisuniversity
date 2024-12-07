package routes

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/kamukwamba/oerisuniversity/dbcode"
	"github.com/kamukwamba/oerisuniversity/encription"
)

type AssesmentGrade struct {
	UUID              string
	Student_UUID      string
	Cource_Name       string
	Assesment_Title   string
	Assesment_Grade   string
	Assesment_Comment string
	Assesment_Date    string
}

type AssesmentOut struct {
	Present       bool
	AssesmentList []AssesmentGrade
}

func GetAssesmentData(student_uuid, cource_name string) (bool, []AssesmentGrade) {
	dbconn := dbcode.SqlRead().DB

	var present bool
	stmt, err := dbconn.Query("select uuid, student_uuid,cource_name,title, grade, comment, date from assesmenttable")

	if err != nil {
		fmt.Println("Failed to launch prepare statment: ", err)

	}

	defer stmt.Close()

	var assesment_out AssesmentGrade
	var assesment_out_list []AssesmentGrade

	for stmt.Next() {
		err = stmt.Scan(&assesment_out.UUID, &assesment_out.Student_UUID, &assesment_out.Cource_Name, &assesment_out.Assesment_Title, &assesment_out.Assesment_Grade, &assesment_out.Assesment_Comment, &assesment_out.Assesment_Date)

		if err != nil {
			present = false
			break
		}

		present = true

		if assesment_out.Student_UUID == student_uuid && assesment_out.Cource_Name == cource_name {
			assesment_out_list = append(assesment_out_list, assesment_out)
		}

	}

	return present, assesment_out_list

}

func HandInAssesment(w http.ResponseWriter, r *http.Request) {

	student_uuid := r.URL.Query().Get("student_uuid")
	cource_name := r.URL.Query().Get("cource_name")

	present, assesment_data := GetAssesmentData(student_uuid, cource_name)

	display_assesment := AssesmentOut{
		Present:       present,
		AssesmentList: assesment_data,
	}

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	err := tpl.ExecuteTemplate(w, "student_cource_assesment", display_assesment)

	if err != nil {
		log.Fatal(err)
	}
}

func GradeAssesment(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	uuid := encription.Generateuudi()
	student_uuid := r.URL.Query().Get("student_uuid")
	cource_name := r.URL.Query().Get("cource_name")
	assesment_title := r.FormValue("title")
	grade := r.FormValue("grade")
	comment := r.FormValue("comment")

	dbconn := dbcode.SqlRead().DB

	stmt, err := dbconn.Prepare("insert into assesmenttable (uuid, student_uuid,cource_name,title, grade, comment, date ) values(?,?,?,?,?,?,?)")

	if err != nil {
		fmt.Println("Prepare statment failed to load error: ", err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(uuid, student_uuid, cource_name, assesment_title, grade, comment)

	if err != nil {
		fmt.Println("")
	}

	display_assesment := AssesmentGrade{
		UUID:              uuid,
		Assesment_Title:   assesment_title,
		Assesment_Grade:   grade,
		Assesment_Comment: comment,
	}

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	err = tpl.ExecuteTemplate(w, "student_cource_assesment", display_assesment)

	if err != nil {
		log.Fatal(err)
	}

}

func DeleteAssesmentAdmin(w http.ResponseWriter, r *http.Request) {

	uuid := r.URL.Query().Get("uuid")
	deleteuser := dbcode.SqlRead().DB

	stmt, err := deleteuser.Prepare("delete from assesmenttable where uuid = ?")

	if err != nil {
		fmt.Println("failed to delete one")

	}
	defer stmt.Close()

	_, errde := stmt.Exec(uuid)

	if errde != nil {
		fmt.Println("failed to delete two")

	}
}

// ASSESMEENT TABLE

func LoadAssesmentTable() {
	dbconn := dbcode.SqlRead().DB

	defer dbconn.Close()

	//CREATE ACMS
	create_assesment := `
		create table if not exists assesmenttable(
			uuid blob not null,
			student_uuid text,
			cource_name text,
			title text,
			grade text,
			comment text,
			date text);`

	_, create_assesment_error := dbconn.Exec(create_assesment)
	if create_assesment_error != nil {
		log.Printf("%q: %s\n", create_assesment_error, create_assesment)
	}

}
