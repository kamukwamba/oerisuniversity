package routes

import (
	"database/sql/driver"
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

type StudentInfo struct {
	UUID                  string
	First_Name            string
	Last_Name             string
	Phone                 string
	Email                 string
	Date_Of_Birth         string
	Gender                string
	Marital_Status        string
	Country               string
	Education_Background  string
	Program               string
	High_School           string
	Grammer_Confirmation  string
	Waiver                string
	Children              string
	School_Attended       string
	Major_In              string
	Degree_Obtained       string
	Current_Occupation    string
	Field_Interested      string
	Prio_Techniques       string
	Previouse_Experience  string
	Purpose_Of_Enrollment string
	Use_Of_Knowledge      string
	Reason_For_Choice     string
	Method_Of_Encounter   string
}

type StudentCridentials struct {
	UUID        string
	StudentUUID string
	Email       string
	Password    string
}

type ProgramListName struct {
	Program_Name []string
}

type StringSlice []string

func SendEMAIL() {
	fmt.Println("New Student Registered")
}

func Validation(email string) bool {

	result := false

	return result
}

func CreateStudentCridentials(studentdate StudentCridentials) bool {

	confirm_creation := true
	dbread := dbcode.SqlRead()
	cridentials, err := dbread.DB.Begin()
	if err != nil {
		log.Fatal()
	}

	uuid := encription.Generateuudi()
	//RETRIEVE ENCRIPTION CODE
	// student_uuid := encription.EncryptData(studentdate.StudentUUID)
	// student_email := encription.EncryptData(studentdate.Email)
	// student_password := encription.EncryptData(studentdate.Password)
	student_uuid := studentdate.StudentUUID
	student_email := studentdate.Email
	student_password := studentdate.Password

	stmt, err := cridentials.Prepare("insert into studentcridentials(uuid, student_uuid, email,password) values(?,?,?,?)")

	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close()
	_, err = stmt.Exec(uuid, student_uuid, student_email, student_password)

	if err != nil {
		log.Fatal(err)
		fmt.Println("PART 2: Failed to save to cridentials")
		confirm_creation = false
	}

	err = cridentials.Commit()
	if err != nil {
		log.Fatal(err)
		confirm_creation = false
	}

	return confirm_creation
}

func CreateStudent(data StudentInfo) bool {
	dbread := dbcode.SqlRead()
	entry, err := dbread.DB.Begin()
	var result bool = true
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := entry.Prepare("insert into studentdata(uuid, first_name, last_name, phone, email, date_of_birth,marital_status,country,eduction_background,program,high_scholl_confirmation,grammer_comprihention,waiver,number_of_children,school_atteneded,major_studied,degree_obtained,current_occupetion,field_interested_in,mps_techqnique_Practiced,previouse_experince,purpose_of_enrollment,use_of_degree,reason_for_choice,method_of_incounter) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,?, ?, ?, ?, ?)")

	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(data.UUID, data.First_Name, data.Last_Name, data.Phone, data.Email, data.Date_Of_Birth, data.Marital_Status, data.Country, data.Education_Background, data.Program, data.High_School, data.Grammer_Confirmation, data.Waiver, data.Children, data.School_Attended, data.Major_In, data.Degree_Obtained, data.Current_Occupation, data.Field_Interested, data.Prio_Techniques, data.Previouse_Experience, data.Purpose_Of_Enrollment, data.Use_Of_Knowledge, data.Reason_For_Choice, data.Method_Of_Encounter)
	if err != nil {
		ErrorPrintOut("enrollment", "CreateStudent", fmt.Sprintf("%s", err))
		fmt.Println("PART 2: Failed to execute")
		result = false
	}

	err = entry.Commit()
	if err != nil {
		log.Fatal(err)
		result = false
	}

	return result

}

func FindStudent(email string) bool {
	dbread := dbcode.SqlRead()
	var result bool = true

	stmt, err := dbread.DB.Prepare("select first_name from studentdata where email = ?")
	if err != nil {
		log.Fatal(err)
		result = false
	}
	defer stmt.Close()

	var first_name string
	err = stmt.QueryRow(email).Scan(&first_name)

	if err != nil {
		result = false

	}

	return result
}

func (stringSlice StringSlice) Value() (driver.Value, error) {
	var quotedStrings []string
	for _, str := range stringSlice {
		quotedStrings = append(quotedStrings, strconv.Quote(str))
	}
	value := fmt.Sprintf("[%s]", strings.Join(quotedStrings, ","))
	return value, nil
}

func Enrollment(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	if r.Method == "POST" {
		fmt.Println("Form obtained")
	}
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	//debug failure to laod templates

	err := tpl.ExecuteTemplate(w, "enrollstudent.html", nil)

	if err != nil {
		log.Fatal(err)
	}
}

func ConfirmEnrollment(w http.ResponseWriter, r *http.Request) {
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	r.ParseForm()

	//convert on to true
	program_name := "ACAMS"
	var studentsdatain StudentInfo
	uuid := encription.Generateuudi()
	first_name := r.FormValue("first_name")
	last_name := r.FormValue("last_name")
	phone := r.FormValue("phone_number")
	email := r.FormValue("email")
	dateofbirth := r.FormValue("date_of_birth")
	gender := r.FormValue("studentsgender")
	payment_type := r.FormValue("payment")
	maritalstatus := r.FormValue("marital_status")
	country := r.FormValue("country")
	educationlevel := r.FormValue("education_level")
	program := r.FormValue("program")
	confirmdiplomer := r.FormValue("ucms_diplomer")
	languagecomprihension := r.FormValue("launguage_comprihension")
	waiver := r.FormValue("waiver")
	chidrencount := r.FormValue("children_count")
	collegename := r.FormValue("trade_college_name")
	collegemajor := r.FormValue("college_major")
	collegediplomer := r.FormValue("college_diplomers")
	currentoccupation := r.FormValue("current_occupation")
	fieldofinterest := r.FormValue("field_of_interest")
	priorexperience := r.FormValue("prior_experience")
	priorknowledge := r.FormValue("prior_knowledge")
	purposeofenrolling := r.FormValue("purpose_of_enrolling")
	applicationofknowledge := r.FormValue("application_of_knowledge")
	reasonforchossingucms := r.FormValue("reason_for_chossing_ucms")
	methodofknowledge := r.FormValue("method_of_knowledge")

	studentsdatain = StudentInfo{
		UUID:                  uuid,
		First_Name:            first_name,
		Last_Name:             last_name,
		Phone:                 phone,
		Email:                 email,
		Date_Of_Birth:         dateofbirth,
		Gender:                gender,
		Marital_Status:        maritalstatus,
		Country:               country,
		Education_Background:  educationlevel,
		Program:               program,
		High_School:           confirmdiplomer,
		Grammer_Confirmation:  languagecomprihension,
		Waiver:                waiver,
		Children:              chidrencount,
		School_Attended:       collegename,
		Major_In:              collegemajor,
		Degree_Obtained:       collegediplomer,
		Current_Occupation:    currentoccupation,
		Field_Interested:      fieldofinterest,
		Prio_Techniques:       priorexperience,
		Previouse_Experience:  priorknowledge,
		Purpose_Of_Enrollment: purposeofenrolling,
		Use_Of_Knowledge:      applicationofknowledge,
		Reason_For_Choice:     reasonforchossingucms,
		Method_Of_Encounter:   methodofknowledge,
	}

	chaeck_user := FindStudent(email)

	type Inuse struct {
		Present bool
	}

	if chaeck_user {
		fmt.Println("User With eamil already in database")

		user_email_in_use := Inuse{
			Present: true,
		}

		err := tpl.ExecuteTemplate(w, "enrollstudent.html", user_email_in_use)

		if err != nil {
			log.Fatal(err)
		}

	} else {
		fmt.Println("Create New User")
		result := CreateStudent(studentsdatain)

		if result {

			current_time := time.Now().Local()

			date_applied := fmt.Sprintf("%s", current_time)

			acams_data := ACAMS{
				Student_UUID:   uuid,
				First_Name:     first_name,
				Last_Name:      last_name,
				Email:          email,
				Applied:        true,
				Approved:       false,
				Payment_Method: payment_type,
				Paid:           "pending",
				Completed:      false,
				Date:           date_applied,
			}

			addedtoacams := CreateACAMS(acams_data, payment_type)

			if addedtoacams {
				SendEMAIL()
				studentcridentials := StudentCridentials{
					StudentUUID: uuid,
					Email:       email,
					Password:    email,
				}
				CreateStudentCridentials(studentcridentials)
				AddStudentPrograms(uuid, program_name)

			} else {
				fmt.Println("Problem with adding student to acams")
			}
		} else {
			fmt.Println("FAILED TO CREATE NEW USER")
		}

		err := tpl.ExecuteTemplate(w, "confirmenroll.html", nil)

		if err != nil {
			log.Fatal(err)
		}
	}

}
