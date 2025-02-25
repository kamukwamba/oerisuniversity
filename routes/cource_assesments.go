package routes

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"io"
	"path/filepath"
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
	Assesment []AssesmentGrade
	Student_UUID string
	Cource_Name string
	StInfo StudentInfo
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
	studentdata := GetStudentAllDetails(student_uuid)

	present, assesment_data := GetAssesmentData(student_uuid, cource_name)

	display_assesment := AssesmentOut{
		Present:       present,
		Assesment: assesment_data,
		StInfo:      studentdata,
		Cource_Name: cource_name,
		Student_UUID: student_uuid,

		
	}

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	err := tpl.ExecuteTemplate(w, "handInAsignments.html", display_assesment)

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

func SaveGrade(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

}


func DownloadAssesments(w  http.ResponseWriter, r *http.Request){
	student_uuid := r.URL.Query().Get("student_uuid")
	cource_name := r.URL.Query().Get("cource_name")
	pdf_filename := r.URL.Query().Get("file_name")
	
	dbFilePath := fmt.Sprintf("./assesmentFiles/%s/%s/%s.pdf", student_uuid, cource_name, pdf_filename)
	

	if _, err := os.Stat(dbFilePath); os.IsNotExist(err) {
		http.Error(w, "Database file not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(dbFilePath)))

	
	file, err := os.Open(dbFilePath)
	if err != nil {
		http.Error(w, "Unable to open the file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "Error writing file to response", http.StatusInternalServerError)
		return
	}
}

type FileDownload struct {
		FileName string
		Student_UUID string
		Cource_Name string
	}
	type HandedIn struct {
		Handed []FileDownload
		Assesment []AssesmentGrade
	}
func GradeCA(w http.ResponseWriter, r *http.Request) {

	student_uuid := r.URL.Query().Get("student_uuid")
	cource_name := r.URL.Query().Get("cource_name")
	
	theString := fmt.Sprintf("assesmentFiles/%s/%s", student_uuid, cource_name)
	filesAs, err  := listAssignments(theString)
	
	
	
	var fileData FileDownload
	
	var fileDataList []FileDownload
	
	for _,item := range filesAs{
		fileData = FileDownload {
			FileName: item,
			Student_UUID: student_uuid,
			Cource_Name: cource_name,
		}
		
		fileDataList = append(fileDataList, fileData)
		
	}
	

	fmt.Println("Route Has Been Hit")

	var data_out HandedIn
	var data_out_list []AssesmentGrade
	

	present, assesment_data := GetAssesmentData(student_uuid, cource_name)

	if present {
		data_out_list = assesment_data
	}
	
	data_out = HandedIn{
		Handed: fileDataList,
		Assesment: data_out_list,
	}

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	err = tpl.ExecuteTemplate(w, "admin_cource_assesment", data_out)

	if err != nil {
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

}

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






func listAssignments(dir string) ([]string, error) {
	var pdfFiles []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Check if the file has a .pdf extension
		if !info.IsDir() && filepath.Ext(info.Name()) == ".pdf" {
			pdfFiles = append(pdfFiles, info.Name())
		}
		return nil
	})
	return pdfFiles, err
}