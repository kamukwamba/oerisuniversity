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

type Displayed struct {
	Section        string
	Question       string
	Quesion_number int
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

type DisplayExam struct {
	AlreadTaken    bool
	ExamData       QuestionsEntered
	Exam_Questions []Displayed
}

type ExamStructOut struct {
	Student_UUID   string
	Program_Name   string
	Cource_Name    string
	Cource_Code    string
	Time           string
	Attempt_Number string
	Question_List  []QuestionStruct
}

type ExamComplete struct {
	Student_UUID string
	Program_Name string
	Cource_Name  string
	Cource_Code  string
	Grade        string
}

type ExamTakenStruct struct {
	UUID            string
	Cource_Name     string
	Student_UUID    string
	Attemp_Number   string
	First_Attempted string
	Open_Period     string
	Answers         string
	Grade           string
	Date            string
	Comment         string
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

func Create_Exam(question_in Questions_Construct) bool {
	result := true
	create_exam := dbcode.SqlRead().DB

	stmt, err := create_exam.Prepare("insert into exam_questions(uuid, cource_uuid,cource_name, section_A_q,section_b_q,section_a_a, question_number) values(?,?,?,?,?,?,?,?,?)")

	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close()

	uuid := encription.Generateuudi()
	var cource_uuid string
	var cource_name string
	var answer string
	var question_a string
	var question_b string
	var question_number int

	section_enter := question_in.Section

	switch section_enter {
	case "A":
		cource_name = question_in.Cource_Name
		cource_uuid = question_in.Cource_UUID
		answer = question_in.Answer
		question_a = question_in.Question
		question_b = ""
		question_number = question_in.Question_Number

	case "B":
		cource_name = question_in.Cource_Name
		cource_uuid = question_in.Cource_UUID
		answer = question_in.Answer
		question_a = ""
		question_b = question_in.Question
		question_number = question_in.Question_Number

	}

	_, err = stmt.Exec(uuid, cource_uuid, cource_name, question_a, question_b, answer, question_number)

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

func Read_Exam_Taken(uuid string) (bool, ExamTakenStruct) {

	exam_taken := true

	var examout ExamTakenStruct

	check_if_exam_taken := dbcode.SqlRead().DB

	stmt, err := check_if_exam_taken.Prepare("select uuid, student_uuid, attemp_number, first_attempted, open_period,answers,grade,comment,date from write_exam where student_uuid = ?")

	if err != nil {
		exam_taken = false
	}

	defer stmt.Close()

	err = stmt.QueryRow(uuid).Scan(&examout.UUID, &examout.Student_UUID, &examout.Attemp_Number, &examout.First_Attempted, &examout.Open_Period, &examout.Answers, &examout.Grade, &examout.Comment, &examout.Date)

	if err != nil {
		exam_taken = false
	}

	return exam_taken, examout

}

// func Take_Exam(w http.ResponseWriter, r *http.Request) {

// }

// func Update_Exam_Taken(w http.ResponseWriter, r *http.Request) {

// }

// func Delete_Exam_Taken(w http.ResponseWriter, r *http.Request) {

// }

type TakeExamStruct struct {
	Taken         bool
	ExamTaken     QuestionsEntered
	TakensResults ExamTakenStruct
}

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

type Complition struct {
	Logo    string
	Message string
}

func TakeExam(w http.ResponseWriter, r *http.Request) {

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	var template_name string
	var display_Exam DisplayExam
	var question_out Displayed
	var question_displayed []Displayed
	var open_time_over bool
	cource_name := r.URL.Query().Get("cource_name")
	uuid := r.URL.Query().Get("uuid")

	read_exam, result_out := Read_Exam_Taken(uuid)

	attemped_number := result_out.Attemp_Number
	comment := result_out.Comment
	open_period := result_out.Open_Period
	// first_attempted := result_out.First_Attempted

	get_exam_questions := Read_Exam(cource_name)

	section_a := get_exam_questions.Section_A
	section_b := get_exam_questions.Section_B

	questions_count := 0

	template_name = "exam_complete"

	fmt.Println(template_name)

	//CREATE LIST OF QUESTION FOR STUDENTS

	year := time.Now().Year()
	month := time.Now().Month()
	day := time.Now().Day()
	hour := time.Now().Hour()
	min := time.Now().Minute()

	first_attemp := fmt.Sprintf("%s,%s,%s,%s,%s", year, month, day, hour, min)

	fmt.Println(first_attemp)
	t := time.Now()
	u := time.Date(2020, 5, 16, 20, 45, 34, 0, time.UTC)
	is_valid := t.Sub(u)

	is_valid.Hours()

	if is_valid > 7 {
		open_time_over = true
	} else {
		open_time_over = false
	}

	fmt.Println(open_time_over)

	for _, item := range section_a {
		questions_count++

		question_out = Displayed{
			Section:        "A",
			Question:       item.Question,
			Quesion_number: questions_count,
		}

		question_displayed = append(question_displayed, question_out)

	}

	for _, item := range section_b {
		questions_count++
		question_out = Displayed{
			Section:        "B",
			Question:       item.Question,
			Quesion_number: questions_count,
		}

		question_displayed = append(question_displayed, question_out)
	}

	//CULCULATE THE TIME REMANING FROM THE FIRST ATTEMPET
	attemp_number_out, err := strconv.Atoi(attemped_number)
	open_period_out, err := strconv.Atoi(open_period)

	date_out := time.Now()
	// open_period_valid := date_out - open_period_out

	fmt.Println(open_period_out, attemp_number_out, date_out)

	if err != nil {
		fmt.Print("String Conv: ", err)

	}

	if read_exam {
		var print_out_message string
		var page_data Complition

		if comment == "true" {

			print_out_message = "Congratulations you completed the cources succesfully"
			page_data = Complition{
				Logo:    "/assets",
				Message: print_out_message,
			}

		} else {
			if open_time_over {
				print_out_message = "The grace period allocated for you to retake the exam has expired to retake the exam please contact us for further information"

				page_data = Complition{
					Logo:    "/assets",
					Message: print_out_message,
				}

			}
			if attemp_number_out > 3 {
				print_out_message = "Sorry you have exided the attempts given to you to take the exam. For further assistance please contact us fo more information"
				page_data = Complition{
					Logo:    "/assets",
					Message: print_out_message,
				}

			}

			fmt.Println(page_data)

			fmt.Println(print_out_message)

		}

	} else {
		display_Exam = DisplayExam{
			AlreadTaken:    false,
			ExamData:       get_exam_questions,
			Exam_Questions: question_displayed,
		}

		err = tpl.ExecuteTemplate(w, "exam_code.html", display_Exam)

		if err != nil {
			log.Fatal(err)
		}
	}

}

func CheckFoRCource(uuid, question string) {

}

type Questions_Out struct {
	Section_A []QuestionStruct
	Section_B []QuestionStruct
}

type Questions_Construct struct {
	UUID            string
	Section         string
	Cource_UUID     string
	Cource_Name     string
	Question        string
	Answer          string
	Question_Number int
}

func Question_Count(uuid string) (bool, string) {

	dbcounter := dbcode.SqlRead().DB
	present := true
	var number string

	stmt, err := dbcounter.Prepare("select question_number  from exam_questions where cource_uuid = ?")

	if err != nil {
		fmt.Println("failed to work properly")
		present = false
	}

	defer stmt.Close()

	err = stmt.QueryRow(uuid).Scan(number)

	if err != nil {
		fmt.Println("check Query row 548 exam_go")
	}

	return present, number

}

func AddExam(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	question_section := r.URL.Query().Get("section")

	// when the page loads a uuid will be created which be usedd to track the questions created if the uuid is in the table updates wiil be made to the question list

	var questions_display_out []Questions_Construct

	program_name := r.URL.Query().Get("program_name")
	cource_name := r.URL.Query().Get("cource_name")
	cource_uuid := r.URL.Query().Get("uuid")
	exam_time := r.FormValue("exam_time")
	exam_code := r.FormValue("exam_code")

	exam_time_out, _ := strconv.Atoi(exam_time)

	fmt.Println(program_name, exam_code, exam_time_out)

	switch question_section {
	case "A":
		is_present, number := Question_Count(cource_uuid)
		question_a := r.FormValue("question_a")
		answers := r.FormValue("answer")
		var question_content Questions_Construct

		if is_present {

			number_out, _ := strconv.Atoi(number)
			number_out++
			question_content = Questions_Construct{
				Cource_UUID:     cource_uuid,
				Section:         "A",
				Cource_Name:     cource_name,
				Question:        question_a,
				Answer:          answers,
				Question_Number: number_out,
			}

			Create_Exam(question_content)

			questions_display_out = append(questions_display_out, question_content)

		} else {

			number_out := 1
			question_content = Questions_Construct{
				Cource_UUID:     cource_uuid,
				Cource_Name:     cource_name,
				Section:         "A",
				Question:        question_a,
				Answer:          answers,
				Question_Number: number_out,
			}

			Create_Exam(question_content)

			questions_display_out = append(questions_display_out, question_content)

		}

	case "B":
		is_present, number := Question_Count(cource_uuid)
		question_b := r.FormValue("question_b")
		var question_content Questions_Construct

		if is_present {

			number_out, _ := strconv.Atoi(number)
			number_out++
			question_content = Questions_Construct{
				Cource_UUID:     cource_uuid,
				Section:         "B",
				Cource_Name:     cource_name,
				Question:        question_b,
				Question_Number: number_out,
			}

			Create_Exam(question_content)

			questions_display_out = append(questions_display_out, question_content)

		} else {

			number_out := 1
			question_content = Questions_Construct{
				Cource_UUID:     cource_uuid,
				Cource_Name:     cource_name,
				Section:         "A",
				Question:        question_b,
				Question_Number: number_out,
			}

			Create_Exam(question_content)
			questions_display_out = append(questions_display_out, question_content)

		}

	case "C":

	}

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	err := tpl.ExecuteTemplate(w, "questions_out", nil)

	if err != nil {
		log.Fatal(err)
	}
}

type ExamDetails struct {
	UUID          string
	Program_Name  string
	Cource_Name   string
	Cource_Code   string
	Exam_Duration int
	Total_Marks   string
}

func LoadExamTable() {

	exam_code := dbcode.SqlRead().DB

	defer exam_code.Close()

	exam_details := `create table if not exits exam_details(
		uuid blob not null,
		program_name text,
		cource_name text,
		cource_code  text,
		duration int,
		total string
	)`

	_, exam_details_error := exam_code.Exec(exam_details)
	if exam_details_error != nil {
		log.Printf("%q: %s\n", exam_details_error, exam_details)
	}

	create_exam := `
		create table if not exists exam_questions(
		uuid blob not null,
		cource_uuid text,
		cource_name text,
		section_a_q text,
		section_b_q text,
		section_a_a text,
		question_number int)`

	_, create_exam_error := exam_code.Exec(create_exam)
	if create_exam_error != nil {
		log.Printf("%q: %s\n", create_exam_error, create_exam)
	}

	take_exam := `
		create table if not exists write_exam(
		uuid blob not null,
		cource_name text,
		student_uuid text,
		attemp_number int,
		first_attempted int,
		open_period int,
		answers text,
		grade text,
		comment text,
		completed  text,
		date text
		)`

	_, take_exam_error := exam_code.Exec(take_exam)
	if create_exam_error != nil {
		log.Printf("%q: %s\n", take_exam_error, take_exam)
	}

}
