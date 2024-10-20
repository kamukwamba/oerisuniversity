package routes

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/kamukwamba/oerisuniversity/dbcode"
	"github.com/kamukwamba/oerisuniversity/encription"
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
	Section  string
	Question string
	Answer   string
}

type QuestionsEntered struct {
	UUID           string
	Program_Name   string
	Cource_Name    string
	Question_Count int
	Section_A      []QuestionStruct
	Section_B      []QuestionStruct
	Time           int
	Cource_Code    string
	Total_Marks    string
	Date           string
}

func Listify(question_a, question_b string) ([]QuestionStruct, []QuestionStruct) {

	var question_list_a []QuestionStruct //section a questiions
	var question_list_b []QuestionStruct

	question_a = strings.Trim(question_a, "[ ]")
	question_b = strings.Trim(question_b, "[ ]")

	section_A := strings.Split(question_a, "}")
	section_B := strings.Split(question_b, "}")

	fmt.Println(question_a, question_b)

	if len(question_a) > 0 {
		var count_out int
		for _, item := range section_A {

			count_out++

			if count_out != len(section_A) {
				the_result := strings.Trim(item, "{")
				the_result_out := strings.Split(the_result, ":")
				questions_out := QuestionStruct{
					Section:  "A",
					Question: the_result_out[0],
					Answer:   the_result_out[1],
				}

				question_list_a = append(question_list_a, questions_out)

				fmt.Println(question_list_a)

			} else {
				fmt.Println("Out of bounds")
			}

		}
	}

	//Note This code is for presenting the code
	if len(question_b) > 0 {

		var count_out int
		for _, item := range section_B {

			if len(item) > 0 {

				the_result := strings.Trim(item, "{ }")
				fmt.Println(the_result)
				count_out++

				if count_out != len(section_B) {
					the_result := strings.Trim(item, "{")
					questions_out := QuestionStruct{
						Section:  "B",
						Question: the_result,
					}

					question_list_b = append(question_list_b, questions_out)

				} else {
					fmt.Println("Out of bounds")
				}
			}

		}

	}

	return question_list_a, question_list_b
}

func Read_Exam(cource_name string) QuestionsEntered {

	get_exam := dbcode.SqlRead().DB

	var result_out QuestionsEntered

	stmt, err := get_exam.Prepare("select uuid,program_name, cource_name,question_count,section_a,section_b,cource_code,time, date from exam_questions where cource_name = ?")

	if err != nil {
		log.Fatal()
	}

	var uuid string
	var program_name string
	var cource_name_out string
	var question_count string
	var section_a string
	var section_b string
	var cource_code string
	var time_out string
	var date string

	defer stmt.Close()

	err = stmt.QueryRow(cource_name).Scan(&uuid, &program_name, &cource_name_out, &question_count, &section_a, &section_b, &cource_code, &time_out, &date)

	section_a_out, section_b_out := Listify(section_a, section_b)

	question_count_out, _ := strconv.Atoi(question_count)
	time_out_int, _ := strconv.Atoi(time_out)
	result_out = QuestionsEntered{
		UUID:           uuid,
		Program_Name:   program_name,
		Cource_Name:    cource_name,
		Question_Count: question_count_out,
		Section_A:      section_a_out,
		Section_B:      section_b_out,
		Cource_Code:    cource_code,
		Time:           time_out_int,
		Date:           date,
	}

	if err != nil {
		log.Fatal(err)
	}

	return result_out

}

func Create_Exam(question_in QuestionsEntered) bool {
	result := true
	create_exam := dbcode.SqlRead().DB

	stmt, err := create_exam.Prepare("insert into exam_questions(uuid, program_name,cource_name, question_count,section_a,section_b,cource_code,time, date) values(?,?,?,?,?,?,?,?,?)")

	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close()

	uuid := encription.Generateuudi()
	date := time.Now()
	section_a := fmt.Sprintf("%s", question_in.Section_A)
	section_b := fmt.Sprintf("%s", question_in.Section_B)

	_, err = stmt.Exec(uuid, question_in.Program_Name, question_in.Cource_Name, question_in.Question_Count, section_a, section_b, question_in.Cource_Code, question_in.Time, date)

	if err != nil {
		log.Fatal(err)
	}

	return result
}

func Update_Exam(w http.ResponseWriter, r *http.Request) {

}

func Delete_Exam(w http.ResponseWriter, r *http.Request) {

	uuid := r.URL.Query().Get("uuid")
	delete := dbcode.SqlRead().DB

	stmt, err := delete.Prepare("DELETE FROM exam_questions WHERE  uuid = ?")

	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(uuid)

	if err != nil {
		log.Fatal(err)
	}

}

// func Read_Exam_Taken(w http.ResponseWriter, r *http.Request) {

// }

// func Take_Exam(w http.ResponseWriter, r *http.Request) {

// }

// func Update_Exam_Taken(w http.ResponseWriter, r *http.Request) {

// }

// func Delete_Exam_Taken(w http.ResponseWriter, r *http.Request) {

// }

func CreatePage(w http.ResponseWriter, r *http.Request) {
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	program_data := r.URL.Query().Get("uuid")
	exam_value := r.URL.Query().Get("exam_present")
	get_program_data := GetProgramDetailsSingle(program_data)

	type ExamOut struct {
		Present     bool
		Cource_Data CourceDataStruct
		ExamData    QuestionsEntered
	}

	var to_show ExamOut

	cource_name_out := get_program_data.Cource_Name
	//Get Program Details to use when creating a database entry

	if exam_value == "true" {

		result_out := Read_Exam(cource_name_out)
		to_show = ExamOut{
			Present:     true,
			Cource_Data: get_program_data,
			ExamData:    result_out,
		}

	} else {
		to_show = ExamOut{
			Present:     false,
			Cource_Data: get_program_data,
		}

	}

	err := tpl.ExecuteTemplate(w, "create_exam_template", to_show)

	if err != nil {
		log.Fatal(err)
	}

}

func TakeExam(w http.ResponseWriter, r *http.Request) {

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	uuid := r.URL.Query().Get("")

	err := tpl.ExecuteTemplate(w, "exam_code.html", nil)
	fmt.Println(uuid)

	if err != nil {
		log.Fatal(err)
	}

}

func CheckFoRCource(uuid, question string) {

}

type Questions_Out struct {
	Section_A []QuestionStruct
	Section_B []QuestionStruct
}

func AddExam(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	var question_count int

	// when the page loads a uuid will be created which be usedd to track the questions created if the uuid is in the table updates wiil be made to the question list

	program_name := r.URL.Query().Get("program_name")
	cource_name := r.URL.Query().Get("cource_name")
	cource_uuid := r.URL.Query().Get("uuid")
	question_a := r.FormValue("question_a")
	question_b := r.FormValue("question_b")
	exam_time := r.FormValue("exam_time")
	exam_code := r.FormValue("exam_code")

	exam_time_out, _ := strconv.Atoi(exam_time)

	var question_list_a []QuestionStruct
	var question_list_b []QuestionStruct

	section_A := strings.Split(question_a, "}")
	section_B := strings.Split(question_b, "}")

	if len(question_a) > 0 {
		var count_out int
		for _, item := range section_A {

			count_out++
			fmt.Println(1)

			if count_out != len(section_A) {
				the_result := strings.Trim(item, "{")
				the_result_out := strings.Split(the_result, ":")
				questions_out := QuestionStruct{
					Section:  "A",
					Question: the_result_out[0],
					Answer:   the_result_out[1],
				}

				question_list_a = append(question_list_a, questions_out)

				fmt.Println(question_list_a)

				question_count++
			} else {
				fmt.Println("Out of bounds")
			}

		}
	}

	//Note This code is for presenting the code
	if len(question_b) > 0 {

		var count_out int
		for _, item := range section_B {

			if len(item) > 0 {

				the_result := strings.Trim(item, "{ }")
				fmt.Println(the_result)
				count_out++

				if count_out != len(section_B) {
					the_result := strings.Trim(item, "{")
					questions_out := QuestionStruct{
						Section:  "B",
						Question: the_result,
					}

					question_list_b = append(question_list_b, questions_out)

					question_count++
				} else {
					fmt.Println("Out of bounds")
				}
			}

		}

	}

	//{question one: true}{question two: true}{question three: true}{question four: true}
	//{question five}{question six}{question seven}{question eight}

	question_data := Questions_Out{
		Section_A: question_list_a,
		Section_B: question_list_b,
	}

	total_questions := QuestionsEntered{
		Cource_Name:    cource_name,
		Program_Name:   program_name,
		Question_Count: question_count,
		Section_A:      question_list_a,
		Section_B:      question_list_b,
		Cource_Code:    exam_code,
		Time:           exam_time_out,
	}

	result_out := Create_Exam(total_questions)

	if !result_out {
		fmt.Println("failed to save")

	} else {
		ExamTrue(cource_uuid)

	}

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	err := tpl.ExecuteTemplate(w, "questions_out", question_data)

	if err != nil {
		log.Fatal(err)
	}
}

func LoadExamTable() {

	exam_code := dbcode.SqlRead().DB

	defer exam_code.Close()

	create_exam := `
		create table if not exists exam_questions(
		uuid blob not null,
		program_name text,
		cource_name text,
		question_count int,
		section_a text,
		section_b text,
		cource_code text,
		time int,
		date text)`

	_, create_exam_error := exam_code.Exec(create_exam)
	if create_exam_error != nil {
		log.Printf("%q: %s\n", create_exam_error, create_exam)
	}

	take_exam := `
		create table if not exists write_exam(
		uuid blob not null,
		cource_name text,
		attemp_number int,
		first_attempted int,
		open_period int,
		answers text,
		grade text,
		comment text,
		date text
		)`

	_, take_exam_error := exam_code.Exec(take_exam)
	if create_exam_error != nil {
		log.Printf("%q: %s\n", take_exam_error, take_exam)
	}

}
