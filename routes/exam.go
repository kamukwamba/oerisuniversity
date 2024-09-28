package routes

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/kamukwamba/oerisuniversity/dbcode"
)

type ExamQuestons struct {
	Section   string
	Questions []string
}

type ExamStruct struct {
	UUID         string
	Program_Name string
	Cource_Name  string
	Questions    []ExamQuestons
}

type QuestionStruct struct {
	Question string
	Answer   string
}

func GetExam(uuid string) {

}

func CreatePage(w http.ResponseWriter, r *http.Request) {
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	program_data := r.URL.Query().Get("uuid")

	get_program_data := GetProgramDetailsSingle(program_data)

	err := tpl.ExecuteTemplate(w, "create_exam_template", get_program_data)

	if err != nil {
		log.Fatal(err)
	}

}

func TakeExam(w http.ResponseWriter, r *http.Request) {

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	uuid := r.URL.Query().Get("cource_name")

	err := tpl.ExecuteTemplate(w, "write_exam.html", nil)
	fmt.Println(uuid)

	if err != nil {
		log.Fatal(err)
	}

}

func AddExam(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	fmt.Println("The One Picked")
	question := r.FormValue("question")
	answer := r.FormValue("answer_a")

	uuid := r.URL.Query().Get("uuid")
	section := r.URL.Query().Get("section")

	question_struct := QuestionStruct{
		Question: question,
		Answer:   answer,
	}
	fmt.Println(question_struct)

	fmt.Println("Question:", question, "\nAnswer:", answer, "\nUUID: ", uuid, section)

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	err := tpl.ExecuteTemplate(w, "question_section_a", nil)

	if err != nil {
		log.Fatal(err)
	}
}

func LoadExamTable() {

	exam_code := dbcode.SqlRead().DB

	defer exam_code.Close()

	create_exam := `
	create table if not exists create_exam(
			uuid blob not null,
			student_uuid text,
			program_name text,
			cource_name text,
			section_a_questions text,
			section_a_answers text,
			section_b_questions text,
			section_b_answers text

		);`

	_, create_exam_error := exam_code.Exec(create_exam)
	if create_exam_error != nil {
		log.Printf("%q: %s\n", create_exam_error, create_exam)
	}

	take_exam := `
		create table if not exists take_exam(
			uuid blob not null,
			exam_uuid blob,
			program_name text,
			cource_name text,
			student_uuid text,
			section_a_questions text,
			section_b_questions text,
			date text

		);`

	_, take_exam_error := exam_code.Exec(take_exam)
	if create_exam_error != nil {
		log.Printf("%q: %s\n", take_exam_error, take_exam)
	}

}
