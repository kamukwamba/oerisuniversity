package routes

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"ucmps/dbcode"
	"ucmps/encription"
)

type CourceDataStruct struct {
	UUID             string
	Program_Name     string
	Cource_Name      string
	Cource_Aseesment string
	Video_List       string
	Module           string
	Book             string
	Exam             bool
}

type ProgramDataOut struct {
	Present      bool
	Program_Name string
	ProgramData  []CourceDataStruct
}

type CourceDataUpdate struct {
	Update bool
	Data   CourceDataStruct
}

func UpdateProgramDetails(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	uuid := r.URL.Query().Get("uuid")

	program_name := r.FormValue("program_name")
	cource_name := r.FormValue("cource_name")

	book_link := r.FormValue("book_link")
	module_link := r.FormValue("module_link")
	video_link := r.FormValue("video_link")
	assesment_link := r.FormValue("assesment_link")

	update := dbcode.SqlRead().DB

	stmt, err := update.Prepare("UPDATE cource_table SET(program_name = ?, cource_name = ? , cource_assesment, video_list,module,recomended_book, exam_file) where uuid = ? ")
	if err != nil {
		fmt.Println(err)
	}

	defer stmt.Close()

	_, erre := stmt.Exec(program_name, cource_name, assesment_link, video_link, module_link, book_link, uuid)

	if erre != nil {
		log.Fatal(err)
	}

	err_out := tpl.ExecuteTemplate(w, "cource_data_saved", nil)

	if err_out != nil {
		log.Fatal(err)
	}

}
func GetProgramDetailsSingle(uuid_out string) CourceDataStruct {
	var cource_data_out CourceDataStruct
	get_one := dbcode.SqlRead().DB

	stmt, err := get_one.Prepare("select uuid, program_name, cource_name, cource_assesment, video_list,module,recomended_book, exam_file from cource_table where uuid = ?")

	if err != nil {
		fmt.Println("Error One", err)
	}

	defer stmt.Close()

	err = stmt.QueryRow(uuid_out).Scan(&cource_data_out.UUID, &cource_data_out.Program_Name, &cource_data_out.Cource_Name, &cource_data_out.Cource_Aseesment, &cource_data_out.Video_List, &cource_data_out.Module, &cource_data_out.Book, &cource_data_out.Exam)

	if err != nil {
		fmt.Println("UUID: ", uuid_out)

		fmt.Println("Error two", err)
	}

	return cource_data_out
}

func GetProgramDetails(program_name string) []CourceDataStruct {
	var cuorce_data_out_list []CourceDataStruct
	var cource_data_out CourceDataStruct

	get_cource_data := dbcode.SqlRead().DB

	statement, err := get_cource_data.Query("select * from cource_table")

	if err != nil {
		log.Fatal(err)
	}
	defer statement.Close()

	fmt.Println("The Program Name:: ", program_name)

	for statement.Next() {
		err := statement.Scan(
			&cource_data_out.UUID, &cource_data_out.Program_Name, &cource_data_out.Cource_Name, &cource_data_out.Cource_Aseesment, &cource_data_out.Video_List, &cource_data_out.Module, &cource_data_out.Book, &cource_data_out.Exam,
		)
		if err != nil {
			log.Fatal(err)
		}

		if cource_data_out.Program_Name == program_name {
			cuorce_data_out_list = append(cuorce_data_out_list, cource_data_out)
		} else {
			continue
		}

	}

	return cuorce_data_out_list
}

func ProgramDetails(w http.ResponseWriter, r *http.Request) {

	path := r.PathValue("id")
	fmt.Println("The Path Value", path)
	var program_data ProgramDataOut

	result := GetProgramDetails(path)

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	if len(result) >= 1 {

		program_data = ProgramDataOut{
			Present:      true,
			Program_Name: path,
			ProgramData:  result,
		}

	} else {
		program_data = ProgramDataOut{
			Present: false,
		}
	}

	fmt.Println(program_data)
	err := tpl.ExecuteTemplate(w, "programedetails.html", program_data)

	if err != nil {
		log.Fatal(err)
	}
}

func CreateCourseData(w http.ResponseWriter, r *http.Request) {
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	parameter_in := r.URL.Query().Get("parameter")

	var data_out CourceDataUpdate
	fmt.Println(parameter_in)

	if parameter_in == "update" {
		uuid := r.URL.Query().Get("uuid")
		data_out.Update = true
		data_out.Data = GetProgramDetailsSingle(uuid)

	} else if parameter_in == "create" {
		data_out.Update = false
	}

	err := tpl.ExecuteTemplate(w, "create_cource_data", data_out)

	if err != nil {
		log.Fatal(err)
	}

}

func CloseCreateCourseData(w http.ResponseWriter, r *http.Request) {

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	fmt.Println("Close Cource Create")

	err := tpl.ExecuteTemplate(w, "cource_data_close", nil)

	if err != nil {
		log.Fatal(err)
	}
}

func AddCourceData(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	program_name := r.FormValue("program_name")
	cource_name := r.FormValue("cource_name")

	book_link := r.FormValue("book_link")
	module_link := r.FormValue("module_link")
	video_link := r.FormValue("video_link")
	assesment_link := r.FormValue("assesment_link")

	create_uuid := encription.Generateuudi()
	exam := false

	fmt.Println(program_name, cource_name, book_link, module_link, video_link, assesment_link)

	create_cource := dbcode.SqlRead().DB
	statment, err := create_cource.Prepare("insert into cource_table(uuid, program_name, cource_name, cource_assesment, video_list,module,recomended_book, exam_file) values(?,?,?,?,?,?,?,?)")

	if err != nil {
		log.Fatal(err)
	}

	defer statment.Close()

	_, err = statment.Exec(create_uuid, program_name, cource_name, assesment_link, video_link, module_link, book_link, exam)

	if err != nil {
		log.Fatal(err)
	}

	err_out := tpl.ExecuteTemplate(w, "cource_data_saved", nil)

	if err_out != nil {
		log.Fatal(err)
	}

}
