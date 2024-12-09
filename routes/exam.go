package routes

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/kamukwamba/oerisuniversity/dbcode"
	"github.com/kamukwamba/oerisuniversity/encription"
)

type ExamDetails struct {
	UUID          string
	Program_Name  string
	Cource_Name   string
	Cource_Code   string
	Exam_Duration int
	Total_Marks   string
}

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
	Section         string
	Question        string
	Question_Number string
	Answer          string
}

type QuestionsEntered struct {
	UUID            string
	Program_Name    string
	Question_Number string
	Cource_Name     string
	Question_Count  int
	Section_A       []QuestionStruct
	Section_B       []QuestionStruct
	Time            int
	Cource_Code     string
	Total_Marks     string
	Date            string
}

type Questions_Out struct {
	Section_A []QuestionStruct
	Section_B []QuestionStruct
}

type Answer_Out struct {
	Cource_UUID     string
	Student_UUID    string
	Qustion_UUID    string
	Question_Number string
	Question        string
	Answer          string
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

type CreateExamResponse struct {
	Details_Messages string
	Questions        Questions_Construct
}

type DisplayExam struct {
	AlreadyTaken   bool
	Student_UUID   string
	Cource_UUIDOut string
	Attempt_Number int ``
	ExamData       ExamDetails
	Exam_Questions []Question_Structure
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

type Exam_Details struct {
	UUID         string
	Cource_UUID  string
	Program_Name string
	Cource_Name  string
	Cource_Code  string
	Duration     string
	Total_Marks  string
}

type Question_Structure struct {
	Section         string
	Question_UUID   string
	Question_Number string
	Question        string
}

func CounterFunc(str string) int {
	var counter int

	counter = counter + 1

	return counter

}

func ToUpperCase(str string) string {

	return strings.ToUpper(str)
}

func Clean(str string) string {

	str_out := strings.Trim(str, "_")

	join_string := fmt.Sprintf(" %s ", str_out)
	capitalised := ToUpperCase(join_string)

	return capitalised
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

func Read_Exam(cource_name string) ([]Question_Structure, string) {

	get_exam := dbcode.SqlRead().DB
	var question_structure Question_Structure
	var question_structure_list []Question_Structure

	stmt, err := get_exam.Query("select uuid, section, cource_uuid,cource_name,question, question_number from exam_questions where cource_name = ?", cource_name)

	if err != nil {
		log.Fatal("Failed to get exam questions: ", err)
	}

	var uuid string
	var section string
	var cource_uuid string
	var cource_name_out string
	var question string
	var question_number string

	defer stmt.Close()

	for stmt.Next() {
		err = stmt.Scan(&uuid, &section, &cource_uuid, &cource_name_out, &question, &question_number)

		if err != nil {
			fmt.Println("failed to query row: ", err)
			break
		} else {

			if section == "A" {
				question_structure = Question_Structure{
					Section:         section,
					Question_UUID:   uuid,
					Question_Number: question_number,
					Question:        question,
				}

				question_structure_list = append(question_structure_list, question_structure)

			} else if section == "B" {
				question_structure = Question_Structure{
					Section:         section,
					Question_UUID:   uuid,
					Question_Number: question_number,
					Question:        question,
				}
				question_structure_list = append(question_structure_list, question_structure)

			}
		}
	}

	if err = stmt.Err(); err != nil {
		fmt.Println("stmt scan failed, error out: ", err)
	}

	if err != nil {
		log.Fatal(err)
	}

	return question_structure_list, cource_uuid

}
func Grade_Exam(w http.ResponseWriter, r *http.Request) {

	student_uuid := r.URL.Query().Get("student_uuid")
	cource_name := r.URL.Query().Get("cource_name")

	dbconn := dbcode.SqlRead().DB

	var cource_uuid string

	stmt, err := dbconn.Prepare("selecte cource_uuid from from questions where cource_name = ?")

	if err != nil {
		fmt.Println("Failed to initialize prepare statement: ", err)
	}

	defer stmt.Close()

	err = stmt.QueryRow(cource_name).Scan(&cource_uuid)

	if err != nil {
		fmt.Println("Failed to create prapare statement: ", err)
	}

	defer stmt.Close()

	RetrieveStudentExam(cource_uuid, student_uuid)

	fmt.Println(student_uuid, cource_name)

}
func Update_Exam_Details(details Exam_Details) (bool, string) {
	present := true
	update_msg := "Exam details updated successfully"

	// THE RESULT WILL TELL THE INSTRUCT THAT THE DETAILS ARE ALRED PRESENT
	saved_succesfully := false
	uuid := encription.Generateuudi()
	var details_out string
	dbcon := dbcode.SqlRead().DB
	stmt, err := dbcon.Prepare("select cource_uuid from exam_details where cource_uuid = ?")

	if err != nil {
		present = false
		fmt.Println("Program Details Not Present: ", err)
	}

	defer stmt.Close()

	err = stmt.QueryRow(details.Cource_UUID).Scan(&details_out)

	if err != nil {
		fmt.Println("Failed to get details out: ", err)
		present = false

	}

	if !present {
		dbcreate := dbcode.SqlRead().DB
		stmt, err := dbcreate.Prepare("insert into exam_details(uuid,cource_uuid,program_name,cource_name,cource_code,duration, total_marks) values(?,?,?,?,?,?,?)")

		if err != nil {
			fmt.Println("Failedd to save exam details: ", err)
		}

		defer stmt.Close()

		_, err = stmt.Exec(uuid, details.Cource_UUID, details.Program_Name, details.Cource_Name, details.Cource_Code, details.Duration, details.Total_Marks)

		saved_succesfully = true
		if err != nil {
			fmt.Println("Failed to create details entry: ", err)
			saved_succesfully = false
		}

		update_msg = "Exam details created successfully"
	} else if present {
		dbupdate := dbcode.SqlRead().DB
		stmt, err := dbupdate.Prepare("UPDATE exam_details SET  program_name, cource_name, cource_code, duration, total_marks) where uuid = ?")

		if err != nil {
			fmt.Println("Failed to Execute Update Query: ", err)
		}

		defer stmt.Close()

		_, err = stmt.Exec(details.Program_Name, details.Cource_Name, details.Cource_Code, details.Duration, details.Total_Marks)

		if err != nil {
			fmt.Println("Failed To Update Exam Details: ", err)
		}
	}

	return saved_succesfully, update_msg

}

func Create_Exam_Details(details Exam_Details) (bool, string) {
	present := true
	message_out := "Exam details saved succesfully"

	// THE RESULT WILL TELL THE INSTRUCT THAT THE DETAILS ARE ALRED PRESENT
	saved_succesfully := false
	uuid := encription.Generateuudi()
	var details_out string
	dbcon := dbcode.SqlRead().DB
	stmt, err := dbcon.Prepare("select cource_uuid from exam_details where cource_uuid = ?")

	if err != nil {
		present = false
		fmt.Println("Program Details Not Present: ", err)
	}

	defer stmt.Close()

	err = stmt.QueryRow(details.Cource_UUID).Scan(&details_out)

	if err != nil {
		fmt.Println("Failed to get details out: ", err)
		present = false

	}

	if !present {
		dbcreate := dbcode.SqlRead().DB
		stmt, err := dbcreate.Prepare("insert into exam_details(uuid,cource_uuid,program_name,cource_name,cource_code,duration, total_marks) values(?,?,?,?,?,?,?)")

		if err != nil {
			fmt.Println("Failedd to save exam details: ", err)
		}

		defer stmt.Close()

		_, err = stmt.Exec(uuid, details.Cource_UUID, details.Program_Name, details.Cource_Name, details.Cource_Code, details.Duration, details.Total_Marks)

		saved_succesfully = true

		if err != nil {
			fmt.Println("Failed to create details entry: ", err)
			saved_succesfully = false
		}
	} else if present {
		message_out = "Exam details present in database. If you wuld like to make changes please select the 'Update Details'."
	}

	return saved_succesfully, message_out
}

func Create_Exam(question_in Questions_Construct) bool {
	result := true
	create_exam := dbcode.SqlRead().DB

	stmt, err := create_exam.Prepare("insert into exam_questions(uuid, section,cource_uuid,cource_name,question,answer, question_number) values(?,?,?,?,?,?,?)")

	if err != nil {
		log.Fatal("Failed to entere exam", err)
	}

	defer stmt.Close()

	uuid := encription.Generateuudi()
	var section string
	var cource_uuid string
	var cource_name string
	var answer string
	var question string
	var question_number int

	section_enter := question_in.Section

	switch section_enter {
	case "A":
		cource_name = question_in.Cource_Name
		cource_uuid = question_in.Cource_UUID
		section = question_in.Section
		answer = question_in.Answer
		question = question_in.Question
		question_number = question_in.Question_Number

	case "B":

		section = question_in.Section
		cource_name = question_in.Cource_Name
		cource_uuid = question_in.Cource_UUID
		answer = question_in.Answer
		question = question_in.Question
		question_number = question_in.Question_Number

	}

	_, err = stmt.Exec(uuid, section, cource_uuid, cource_name, question, answer, question_number)

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
		ExamData    []Question_Structure
	}

	var to_show ExamOut

	cource_name_out := get_program_data.Cource_Name
	//Get Program Details to use when creating a database entry

	if exam_value == "true" {

		result_out, _ := Read_Exam(cource_name_out)

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

	funcMap := template.FuncMap{
		"ToCapital": ToUpperCase,
		"Clean":     Clean,
	}

	tpl = template.Must(template.New("exam_code.html").Funcs(funcMap).ParseGlob("templates/*.html"))

	var template_name string
	var display_Exam DisplayExam

	// GET VARIABLES FROM ROUTE
	cource_name := r.URL.Query().Get("cource_name")
	uuid := r.URL.Query().Get("uuid")

	// VERIFY IF EXAM HAS BEEN ATTEMPTED
	if_taken, result_out := Read_Exam_Taken(uuid)

	var attemped_number int
	var comment string
	var open_period int
	fmt.Println(comment)

	// CHECK IF STUDENT STILL HAS ATTEMPTS TO WRITE EXAM "START"
	if if_taken {
		attemped_number, _ = strconv.Atoi(result_out.Attemp_Number)
		comment = result_out.Comment
		open_period, _ = strconv.Atoi(result_out.Open_Period)

		fmt.Println(comment)

	} else {
		attemped_number = 0
		comment = ""
		open_period = 7
	}

	if attemped_number > 3 || open_period > 7 {
		template_name = "exam_failed"
	} else {
		template_name = "exam_code.html"

	}

	// CHECK IF STUDENT STILL HAS ATTEMPTS TO WRITE EXAM "END"

	// GET EXAM QUESTIONS "START"
	get_exam_questions, cource_uuid_out := Read_Exam(cource_name)

	// GET EXAM QUESTIONS "END"

	// GET EXAM DETAILS "START"

	fmt.Println("The Cource UUID is: ", cource_uuid_out)

	get_exam_details := GetExamDetails(cource_uuid_out)

	display_Exam = DisplayExam{
		AlreadyTaken:   false,
		Student_UUID:   uuid,
		Cource_UUIDOut: cource_uuid_out,
		Attempt_Number: attemped_number,
		ExamData:       get_exam_details,
		Exam_Questions: get_exam_questions,
	}

	// GET EXAM DETAILS "END"

	//CREATE LIST OF QUESTION FOR STUDENTS

	err := tpl.ExecuteTemplate(w, template_name, display_Exam)

	if err != nil {
		log.Fatal(err)
	}

}

func CheckFoRCource(uuid, question string) {

}

func Question_Count(uuid string) (bool, string) {

	dbcounter := dbcode.SqlRead().DB
	present := true
	counter_present := false
	var number string
	var number_list []string
	var last_number string

	fmt.Println("The UUID ENTERED: ", uuid)

	stmt, err := dbcounter.Query("SELECT question_number  FROM exam_questions WHERE cource_uuid = ?", uuid)

	if err != nil {
		fmt.Println("failed to work properly", err)
		present = false
	}

	defer stmt.Close()

	for stmt.Next() {
		err = stmt.Scan(&number)
		if err != nil {
			fmt.Println("failed to query row: ", err)
			break
		} else {
			number_list = append(number_list, number)
			counter_present = true
		}

	}

	if err = stmt.Err(); err != nil {
		fmt.Println("stmt scan failed, error out: ", err)
	}

	if counter_present {
		last_number = number_list[len(number_list)-1]
		number = last_number
	}

	fmt.Println("Length of number list: ", len(number_list))

	return present, number

}

func GetExamDetails(courceuuid string) ExamDetails {

	var exam_details ExamDetails

	fmt.Println("The Cource Code Search: ", courceuuid, ":")
	dbconn := dbcode.SqlRead().DB

	stmt, err := dbconn.Prepare("select  program_name, cource_name, cource_code, duration, total_marks from exam_details where cource_uuid = ?")

	if err != nil {
		fmt.Println("Prepare Statement failed  in 'GetExamDetails' error out : ", err)
	}

	defer stmt.Close()

	err = stmt.QueryRow(courceuuid).Scan(&exam_details.Program_Name, &exam_details.Cource_Name, &exam_details.Cource_Code, &exam_details.Exam_Duration, &exam_details.Total_Marks)

	if err != nil {
		fmt.Println("Failed to get program details: ", err)
	}

	return exam_details

}

func QuestionUUID(cource_uuid string) []string {

	dbcounter := dbcode.SqlRead().DB

	var uuid_list []string
	var question_uuid string

	fmt.Println("The UUID ENTERED: ", cource_uuid)

	stmt, err := dbcounter.Query("SELECT uuid  FROM exam_questions WHERE cource_uuid = ?", cource_uuid)

	if err != nil {
		fmt.Println("failed to work properly", err)

	}

	defer stmt.Close()

	for stmt.Next() {
		err = stmt.Scan(&question_uuid)
		if err != nil {
			fmt.Println("failed to query row: ", err)
			break
		} else {
			uuid_list = append(uuid_list, question_uuid)

		}

	}

	if err = stmt.Err(); err != nil {
		fmt.Println("stmt scan failed, error out: ", err)
	}

	return uuid_list

}

type QuestionData struct {
	Question_UUID   string
	Question        string
	Question_Number string
}

func GetQuestionData(question_uuid string) QuestionData {

	var question_data QuestionData

	dbconn := dbcode.SqlRead().DB

	stmt, err := dbconn.Prepare("select uuid, question, question_number from exam_questions where uuid = ?")

	if err != nil {
		fmt.Println("Failed to get qustin data, error text: ", err)
	}

	defer stmt.Close()

	err = stmt.QueryRow(question_uuid).Scan(&question_data.Question_UUID, &question_data.Question, &question_data.Question_Number)

	if err != nil {
		fmt.Println("Failed to get question data", err)
	}

	return question_data

}

func SubmitExam(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	student_uuid := r.URL.Query().Get("student_uuid")
	cource_uuid := r.URL.Query().Get("cource_uuid")

	fmt.Println("Cource_UUID: ", cource_uuid)

	question_uuids := QuestionUUID(cource_uuid)

	var store_answer []Answer_Out

	for _, item := range question_uuids {

		answer := r.FormValue(item)
		data_out := GetQuestionData(item)

		answer_out := Answer_Out{
			Cource_UUID:     cource_uuid,
			Student_UUID:    student_uuid,
			Qustion_UUID:    item,
			Question_Number: data_out.Question_Number,
			Question:        data_out.Question,
			Answer:          answer,
		}

		saved := RecordStudentMarks(answer_out)

		fmt.Println(saved)

	}

	fmt.Println(store_answer)

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	err := tpl.ExecuteTemplate(w, "exam_taken.html", student_uuid)

	if err != nil {
		log.Fatal(err)
	}

}

func AddExamDetails(w http.ResponseWriter, r *http.Request) {

	section_out := r.URL.Query().Get("section")

	fmt.Println(section_out)

	var exam_responce CreateExamResponse
	var save_details bool

	var template_name string
	var details_message string
	cource_uuid := r.URL.Query().Get("uuid")
	pr_name := r.FormValue("program_name")
	c_name := r.FormValue("cource_name")
	cource_code := r.FormValue("cource_code")
	exam_time := r.FormValue("exam_time")
	total_marks := r.FormValue("total_marks")

	fmt.Println("Program Details: ", pr_name, c_name, cource_code, exam_time, total_marks)

	create_exam_detaile := Exam_Details{
		Cource_UUID:  cource_uuid,
		Program_Name: pr_name,
		Cource_Name:  c_name,
		Cource_Code:  cource_code,
		Duration:     exam_time,
		Total_Marks:  total_marks,
	}

	if section_out == "save" {
		save_details, details_message = Create_Exam_Details(create_exam_detaile)

	}
	if section_out == "update" {
		save_details, details_message = Update_Exam_Details(create_exam_detaile)

	}

	if save_details {
		exam_responce = CreateExamResponse{
			Details_Messages: details_message}

		template_name = "details_saved_temp"

	} else {
		exam_responce = CreateExamResponse{
			Details_Messages: details_message}

		template_name = "details_saved_temp"

	}

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	fmt.Println(template_name)

	err := tpl.ExecuteTemplate(w, template_name, exam_responce)

	if err != nil {
		log.Fatal(err)
	}
}

func AddExam(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	question_section := r.URL.Query().Get("section")

	// when the page loads a uuid will be created which be usedd to track the questions created if the uuid is in the table updates wiil be made to the question list

	fmt.Println("Section Letter: ", question_section)
	var exam_responce CreateExamResponse

	var template_name string

	program_name := r.URL.Query().Get("program_name")
	cource_name := r.URL.Query().Get("cource_name")
	cource_uuid := r.URL.Query().Get("uuid")
	exam_time := r.FormValue("exam_time")
	exam_code := r.FormValue("exam_code")

	fmt.Println("Cource UUID: ", cource_uuid)

	exam_time_out, _ := strconv.Atoi(exam_time)

	fmt.Println(program_name, exam_code, exam_time_out)

	switch question_section {

	case "A":
		is_present, number := Question_Count(cource_uuid)

		question_a := r.FormValue("question_a")
		answers := r.FormValue("answer")

		if len(answers) < 1 {
			answers = "false"
		}
		var question_content Questions_Construct

		if is_present {

			number_out, _ := strconv.Atoi(number)
			fmt.Println("Present", number_out)

			question_content = Questions_Construct{
				Cource_UUID:     cource_uuid,
				Section:         "A",
				Cource_Name:     cource_name,
				Question:        question_a,
				Answer:          answers,
				Question_Number: number_out + 1,
			}

			Create_Exam(question_content)

			exam_responce = CreateExamResponse{
				Details_Messages: "",
				Questions:        question_content}

			template_name = "questions_out_a"

		} else {

			number_out := 1
			fmt.Println("Not Present", number_out)
			question_content = Questions_Construct{
				Cource_UUID:     cource_uuid,
				Cource_Name:     cource_name,
				Section:         "A",
				Question:        question_a,
				Answer:          answers,
				Question_Number: number_out,
			}

			Create_Exam(question_content)

			exam_responce = CreateExamResponse{
				Details_Messages: "",
				Questions:        question_content}

			template_name = "questions_out_a"

		}

	case "B":
		is_present, number := Question_Count(cource_uuid)
		question_b := r.FormValue("question_b")
		var question_content Questions_Construct

		if is_present {

			number_out, _ := strconv.Atoi(number)

			question_content = Questions_Construct{
				Cource_UUID:     cource_uuid,
				Section:         "B",
				Cource_Name:     cource_name,
				Question:        question_b,
				Question_Number: number_out + 1,
			}

			Create_Exam(question_content)

			exam_responce = CreateExamResponse{
				Details_Messages: "",
				Questions:        question_content}

			template_name = "questions_out_b"

		} else {

			number_out := 1
			question_content = Questions_Construct{
				Cource_UUID:     cource_uuid,
				Cource_Name:     cource_name,
				Section:         "B",
				Question:        question_b,
				Question_Number: number_out,
			}

			Create_Exam(question_content)

			exam_responce = CreateExamResponse{
				Details_Messages: "",
				Questions:        question_content}

			template_name = "questions_out_b"

		}

	}

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	fmt.Println(template_name)

	err := tpl.ExecuteTemplate(w, template_name, exam_responce)

	if err != nil {
		log.Fatal(err)
	}
}

func LoadExamTable() {

	exam_code := dbcode.SqlRead().DB

	defer exam_code.Close()

	exam_details := `create table if not exists exam_details(
		uuid blob not null,
		cource_uuid text,
		program_name text,
		cource_name text,
		cource_code  text,
		duration int,
		total_marks string
	)`

	_, exam_details_error := exam_code.Exec(exam_details)
	if exam_details_error != nil {
		log.Printf("%q: %s\n", exam_details_error, exam_details)
	}

	create_exam := `
		create table if not exists exam_questions(
		uuid blob not null,
		section text,
		cource_uuid text,
		cource_name text,
		question text,
		answer text,
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
