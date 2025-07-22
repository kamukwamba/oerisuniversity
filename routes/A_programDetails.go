package routes

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/kamukwamba/oerisuniversity/dbcode"

	_ "github.com/mattn/go-sqlite3"
)

type Course_Name struct {
	ID    int
	Code  string
	Name  string
	PCode string
}

type ProgramDataEntry struct {
	ID          int
	Name        string
	Code        string
	CourseNames []Course_Name
}

func CreateProgramDB() {
	db := dbcode.SqlRead().DB

	defer db.Close()
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS ProgramData (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		programName TEXT UNIQUE,
		programCode TEXT UNIQUE,
	)`)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

}

func CreateNewProgramR(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	program_name := r.FormValue("program_name")
	program_code := r.FormValue("programcode")

	err := ConfirmProgramDataExists(program_name, program_code)

	if err != nil {
		fmt.Println("Program Exist With that Program Code Or Program Name")
	} else {
		err = CreateProgramEntry(program_name, program_code)
		if err != nil {
			fmt.Printf("Failed to Create Program Entry:: %s", err)
		} else {
			fmt.Println("Program Entry Created Sucesfully")
		}
	}

}

func ProgramCourseDataPage(w http.ResponseWriter, r *http.Request) {

}

func ConfirmProgramDataExists(programName, programCode string) error {
	db := dbcode.SqlRead().DB
	defer db.Close()

	var program_name, program_code string
	err := db.QueryRow("SELECT courseList FROM ProgramDataList WHERE programName = ? OR programCode = ? ", programName, programCode).Scan(&program_name, program_code)

	if err != nil {
		fmt.Printf("Program Name Already Exists:: %s", err)
	}

	return err
}

func CreateProgramEntry(programName, programCode string) error {

	db := dbcode.SqlRead().DB

	defer db.Close()

	_, err := db.Exec("INSERT INTO ProgramDataList (programName, programCode) VALUES (?,?)",
		programName, programCode)

	return err
}

func GetAllProgramData() ([]ProgramDataEntry, error) {

	db := dbcode.SqlRead().DB
	var programData ProgramDataEntry
	var programDataListOut []ProgramDataEntry

	rows, err := db.Query("SELECT programName, programCode FROM ProgramData")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resultList []ProgramDataEntry

	for rows.Next() {
		var program_name string
		var program_code string
		var course_names []Course_Name

		err := rows.Scan(&program_name, &program_code)
		if err != nil {
			return nil, err
		} else {
			course_names_check, errCourses := GetProgramCourses(program_code)

			if errCourses != nil {

				fmt.Printf("Failed to Get Program Courses: %s", err)
			} else {
				course_names = course_names_check
			}
		}

		programData = ProgramDataEntry{
			Name:        program_name,
			Code:        program_code,
			CourseNames: course_names,
		}

		resultList = append(programDataListOut, programData)

	}

	return resultList, nil
}

///COURSES CRUD

func GetPorgamCourseR(w http.ResponseWriter, r *http.Request) {

	program_code := r.URL.Query().Get("programcode")

	fmt.Println(program_code)

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	err := tpl.ExecuteTemplate(w, "programs.html", nil)

	if err != nil {
		log.Fatal(err)
	}

}

func CreateProgramCourseR(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	dbread := dbcode.SqlRead().DB
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	program_code := r.URL.Query().Get("programcode")
	course_name := r.FormValue("course_name")
	course_code := r.FormValue("course_code")

	var setTemplatName string

	err := CheckCourseInDataBase(program_code, course_name, course_code)

	if err != nil {
		fmt.Printf("Failed To Create New Course Because Course With Either The Same Code Or Name Already Exists In The Data Base")
		setTemplatName = "failedToCreateCourse"
	}

	stmt, err := dbread.Prepare("INSERT INTO CourseNames(courseCode,courseCode, programCode) values(?,?,?)")

	if err != nil {
		fmt.Println("Failed to get the email and password", err)
		setTemplatName = "failedToCreateCourse"

	}

	defer stmt.Close()

	_, err = stmt.Exec(course_name, course_code, program_code)

	if err != nil {
		fmt.Printf("Failed to execute db command create:: %s", err)
		setTemplatName = "failedToCreateCourse"

	}

	err = tpl.ExecuteTemplate(w, setTemplatName, nil)

	if err != nil {
		log.Fatal(err)
	}

}

func CreateCourseMaterial(program_code, course_name, course_code string) error {

	dbread := dbcode.SqlRead().DB

	stmt, err := dbread.Prepare("INSERT INTO CourseMaterial(courseCode,courseCode, programCode) values(?,?,?)")

	if err != nil {
		fmt.Println("Failed to get the email and password", err)

	}

	defer stmt.Close()

	_, err = stmt.Exec(course_name, course_code, program_code)

	if err != nil {
		fmt.Printf("Failed to execute db command create:: %s", err)

	}

	return nil

}

func CheckCourseInDataBase(program_code, course_name, course_code string) error {

	dbread := dbcode.SqlRead().DB
	var course_data Course_Name

	statement, err := dbread.Prepare("SELECT courseName, courseCode,programCode FROM CourseNames = ? WHERE courseName = ? OR courseCode = ?")

	if err != nil {
		return err
	}

	defer statement.Close()

	err = statement.QueryRow(course_name, course_code).Scan(
		&course_data.Name,
		&course_data.Code,
		&course_data.PCode,
	)

	if err != nil {
		return err
	}

	return nil
}

func GetProgramCourses(program_code string) ([]Course_Name, error) {

	db := dbcode.SqlRead().DB
	var courseData Course_Name
	var courseDataListOut []Course_Name

	rows, err := db.Query("SELECT courseName, courseCode FROM CourseNames WHERE programCode = ?", program_code)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var course_name string
		var course_code string

		err := rows.Scan(&course_name, &course_code)
		if err != nil {
			return nil, err
		} else {
		}

		courseData = Course_Name{
			Name: course_name,
			Code: course_code,
		}

		courseDataListOut = append(courseDataListOut, courseData)

	}

	return courseDataListOut, nil
}

func CreateCourseDB() {
	db := dbcode.SqlRead().DB

	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS CourseNames (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		courseName TEXT UNIQUE,
		courseCode TEXT UNIQUE,
		programCode TEXT UNIQUE,
	)`)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

}
