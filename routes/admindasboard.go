package routes

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/kamukwamba/oerisuniversity/dbcode"
)

type AdminLogData struct {
	Email    string
	Password string
}
type AdminInfo struct {
	ID       string
	Name     string
	Email    string
	Password string
}

type AdminPage struct {
	Admin     AdminInfo
	AcamsData []ProgramStruct
}
type DashData struct {
	ACAMSTotal int
	Admin      AdminInfo
}

type StudentProgramList struct {
	Program_Name string
}

func DeleteStudentInfo(uuid string) bool{

	var examtable = fmt.Sprintf("st%s", uuid)
	var deleted = true
	var tables = []string{"studentdata", "studentcridentials"}
	dbconn := dbcode.SqlRead().DB
	
	for item :=  range(tables){
	
		var preparestmt = fmt.Sprintf("DELETE FROM %s WHERE uuid = ?", item)
	
		stmt, err := dbconn.Prepare(preparestmt)
	
		if err != nil {
			fmt.Println("failed to prepare delete: ", err)
			deleted =  false
		
		}
		
		defer stmt.Close()
		
		_, err = stmt.Exec(uuid)
		
		if err != nil {
			fmt.Println("failed to delete: ", err)
			deleted = false
		}
		
	
	}
	
	
	query := fmt.Sprintf("DROP TABLE IF EXISTS %s;", examtable)

	// Execute the SQL query
	_, err := dbconn.Exec(query)
	if err != nil {
		fmt.Errorf("error dropping table %q: %w", examtable, err)
		deleted = false
	}
	
	return deleted


}


func DeleteStudent(w http.ResponseWriter, r *http.Request){

	uuid := r.URL.Query().Get(".uuid")
	
	
	
	
	success := DeleteStudentInfo(uuid)
	
	
	if success {
		w.WriteHeader(http.StatusOK)
		
	}else{
		http.Error(w, "Failed to delete item", http.StatusInternalServerError)
		return
	
	}

}

func AdminData(uuid string) AdminInfo {
	var admin_info AdminInfo
	dbconn := dbcode.SqlRead().DB

	stmt, err := dbconn.Prepare("select uuid, first_name, last_name, email from admin where uuid = ?")

	if err != nil {
		fmt.Println("Failed to load admin prepare", err)
	}

	defer stmt.Close()

	var first_name string
	var last_name string
	var id string
	var email string

	err = stmt.QueryRow(uuid).Scan(&id, &first_name, &last_name, &email)

	name := fmt.Sprintf("%s %s", first_name, last_name)

	admin_info = AdminInfo{
		ID:    id,
		Name:  name,
		Email: email,
	}

	if err != nil {
		fmt.Println("Failed to QueryRow: ", err)
	}

	return admin_info
}

func AdminAuth(data AdminLogData, dataList []dbcode.AdminInfo) (bool, AdminInfo) {

	var result bool
	var admin_data AdminInfo

	for _, admin_info := range dataList {
		name_out := fmt.Sprintf("%s %s", admin_info.First_Name, admin_info.Last_Name)

		id := admin_info.ID
		name := name_out
		email := admin_info.Email
		password := admin_info.Password
		
		matchPassword := CheckPassword(password, data.Password)
		
		

		if matchPassword == true && data.Email == email {
			fmt.Println("A match was found")
			admin_data = AdminInfo{
				ID:       id,
				Name:     name,
				Email:    email,
				Password: password,
			}
			result = true
		}
	}
	return result, admin_data
}

func AdminDashboard(w http.ResponseWriter, r *http.Request) {
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	if r.Method == "POST" {
		r.ParseForm()
		acamscount := ACAMSCount()
		adminList := dbcode.AdminGet()

		email := r.PostFormValue("email")
		password := r.PostFormValue("password")

		authget := AdminLogData{
			Email:    email,
			Password: password,
		}

		check, admin_dataout := AdminAuth(authget, adminList)

		toshow := DashData{
			ACAMSTotal: acamscount,
			Admin:      admin_dataout,
		}

		if check {
			fmt.Println("redirecting")
			err := tpl.ExecuteTemplate(w, "admindasboard.html", toshow)

			if err != nil {
				log.Fatal(err)
			}
		} else {
			err := tpl.ExecuteTemplate(w, "adminloginerror.html", nil)

			if err != nil {
				log.Fatal(err)
			}
		}

	} else {
		err := tpl.ExecuteTemplate(w, "adminloginerror.html", nil)

		if err != nil {
			log.Fatal(err)
		}

	}
	// err := tpl.ExecuteTemplate(w, "admindasboard.html", nil)

	// if err != nil {
	// 	log.Fatal(err)
	// }

}

func StudentACAMSData(student_uuid string) {
	dbread := dbcode.SqlRead()

	stmt, err := dbread.DB.Prepare("select program_list from studentprogramlist where student_uuid = ?")

	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var program_list string

	err = stmt.QueryRow(student_uuid).Scan(&program_list)

	if err != nil {
		fmt.Println("FAILED TO GET STUDENT PROGRAM LIST")
	}

}

func GetACAMSStudents(admin_id string) []ProgramStruct {
	dbread := dbcode.SqlRead()
	var conuter int
	var datalist []ProgramStruct

	var cource_data ProgramStruct
	rows, err := dbread.DB.Query("select * from  acams")
	if err != nil {
		fmt.Println("Failed to get acams student data")
	}
	defer rows.Close()

	for rows.Next() {
		conuter += 1

		var uuid string
		var student_uuid string
		var program_name string
		var first_name string
		var last_name string
		var email string
		var applied bool
		var approved bool
		var payment_method string
		var paid string
		var completed bool
		var date string
		err := rows.Scan(&uuid, &student_uuid, &program_name, &first_name, &last_name, &email, &applied, &approved, &payment_method, &paid, &completed, &date)

		cource_data = ProgramStruct{
			UUID:           uuid,
			Student_UUID:   student_uuid,
			Program_Name:   program_name,
			First_Name:     first_name,
			Last_Name:      last_name,
			Email:          email,
			Applied:        applied,
			Approved:       approved,
			Payment_Method: payment_method,
			Paid:           paid,
			Completed:      completed,
			Date:           date,
			Admin_ID:       admin_id,
		}

		fmt.Println("the cource data out: ", cource_data)
		if err != nil {
			fmt.Println("Check the scan for student data admindashboard")
			log.Fatal(err)
		}

		datalist = append(datalist, cource_data)

	}

	return datalist
}

func ACAMSCount() int {
	dbread := dbcode.SqlRead()
	var counter int

	rows, err := dbread.DB.Query("select * from  acams")
	if err != nil {
		fmt.Println("Failed to get acams student data")
	}
	defer rows.Close()

	for rows.Next() {
		counter += 1

	}

	return counter
}

func GetStudentAllDetails(uuid string) StudentInfo {
	dbread := dbcode.SqlRead()

	stmt, err := dbread.DB.Prepare("select  uuid, first_name, last_name, phone, email, date_of_birth,marital_status,country,eduction_background,program,high_scholl_confirmation,grammer_comprihention,waiver,number_of_children,school_atteneded,major_studied,degree_obtained,current_occupetion,field_interested_in,mps_techqnique_Practiced,previouse_experince,purpose_of_enrollment,use_of_degree,reason_for_choice,method_of_incounter from studentdata where uuid = ?")

	if err != nil {
		fmt.Println("the student uuid", uuid)

	}
	defer stmt.Close()

	var data StudentInfo

	err = stmt.QueryRow(uuid).Scan(&data.UUID, &data.First_Name, &data.Last_Name, &data.Phone, &data.Email, &data.Date_Of_Birth, &data.Marital_Status, &data.Country, &data.Education_Background, &data.Program, &data.High_School, &data.Grammer_Confirmation, &data.Waiver, &data.Children, &data.School_Attended, &data.Major_In, &data.Degree_Obtained, &data.Current_Occupation, &data.Field_Interested, &data.Prio_Techniques, &data.Previouse_Experience, &data.Purpose_Of_Enrollment, &data.Use_Of_Knowledge, &data.Reason_For_Choice, &data.Method_Of_Encounter)

	if err != nil {
		log.Fatal(err)

	}

	// fmt.Println("student info: ", data) //PRINT STUDENT DATA

	return data
}

//ROUTER CODE

//Approve Cources

func UpdateCource(uuid, cource_name string) bool {
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	cource_approved := true

	fmt.Println("Cource name: ", cource_name)
	fmt.Println("UUID: ", uuid)

	dbread := dbcode.SqlRead().DB

	dbupdate_statement := fmt.Sprintf(`UPDATE %s SET approved = ? WHERE student_uuid = ? `, cource_name)

	statement, err := dbread.Prepare(dbupdate_statement)

	if err != nil {
		error_text := fmt.Sprintf("line 44 error from update prepare:: %s", err)
		ErrorPrintOut("studentportal", "ApplyForCource", error_text)

		cource_approved = false
	}
	defer statement.Close()

	_, errup := statement.Exec(true, uuid)

	if errup != nil {
		error_text := fmt.Sprintf("line 50 error from update prepare:: %s", errup)
		ErrorPrintOut("studentportal", "ApplyForCource", error_text)
		cource_approved = false
	}
	return cource_approved

}

func DeleteCource(w http.ResponseWriter, r *http.Request) {

	cource_name := r.URL.Query().Get("cource_name")
	cource_uuid := r.URL.Query().Get("cource_uuid")

	dbconn := dbcode.SqlRead().DB

	stmt, err := dbconn.Prepare("Delete from cource_table where uuid = ?")

	if err != nil {
		http.Redirect(w, r, "/error", http.StatusSeeOther)
	}

	defer stmt.Close()

	_, err = stmt.Exec(cource_uuid)
	if err != nil {
		http.Redirect(w, r, "/error", http.StatusSeeOther)

	}

	if !DeleteAllQuestions(cource_name) {
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		fmt.Println("Failed to delete questions")
	}

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	err = tpl.ExecuteTemplate(w, "empty_tr", nil)

	if err != nil {
		fmt.Println("Failed to delete")
	}

}

func ApproveCourceUpdate(w http.ResponseWriter, r *http.Request) {
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	uuid := r.URL.Query().Get("user_uuid")
	cource_name := r.URL.Query().Get("cource_name")

	fmt.Println(uuid, cource_name)

	var setroute string
	result := UpdateCource(uuid, cource_name)

	if result {

		setroute = "cource_approved_admin"
	} else {
		setroute = "approval_error"
	}
	// dbread := dbcode.SqlRead().DB

	err := tpl.ExecuteTemplate(w, setroute, nil)

	if err != nil {
		log.Fatal(err)
	}

}

func GetStudentProgramDataAdmin(programlist []string, students_uuid string) ([]AllCourceData, bool) {

	var programdata ProgramStruct
	var courcedata []CourceStruct
	var allcourcedataout AllCourceData
	var allcourcedataoutlist []AllCourceData

	var programsavailable ProgramAvailable
	var available []bool

	for _, program := range programlist {

		fmt.Println(program)

		switch program {

		case "ACAMS":

			is_present, dataout, _ := GetACAMSAdmin(students_uuid, "one")

			if is_present {

				var programdataacams ACAMS = dataout
				fmt.Println("is present")
				programdata.UUID = programdataacams.UUID
				programdata.Student_UUID = programdataacams.Student_UUID
				programdata.Program_Name = programdataacams.Program_Name
				programdata.First_Name = programdataacams.First_Name
				programdata.Last_Name = programdataacams.Last_Name
				programdata.Email = programdataacams.Email
				programdata.Payment_Method = programdataacams.Payment_Method
				programdata.Paid = programdataacams.Paid
				programdata.Approved = programdataacams.Approved
				programdata.Applied = programdataacams.Applied
				programdata.Completed = programdataacams.Completed
				programdata.Date = programdataacams.Date

				courcedata = GetFromACAMSCources(students_uuid)

				allcourcedataout.ProgramStruct = programdata
				allcourcedataout.Cource_Struct = courcedata

				allcourcedataoutlist = append(allcourcedataoutlist, allcourcedataout)

				available = append(available, true)
			} else {
				available = append(available, false)
			}

		case "ACMS":
			is_present, data_out, _ := GetACMSAdmin(students_uuid, "one")

			var programdataacms ACMS = data_out
			if is_present {
				programdata.UUID = programdataacms.UUID
				programdata.Student_UUID = programdataacms.Student_UUID
				programdata.Program_Name = programdataacms.Program_Name
				programdata.First_Name = programdataacms.First_Name
				programdata.Last_Name = programdataacms.Last_Name
				programdata.Email = programdataacms.Email
				programdata.Payment_Method = programdataacms.Payment_Method
				programdata.Paid = programdataacms.Paid
				programdata.Approved = programdataacms.Approved
				programdata.Applied = programdataacms.Applied
				programdata.Completed = programdataacms.Completed
				programdata.Date = programdataacms.Date

				courcedata = GetFromACMSCources(students_uuid)

				allcourcedataout.ProgramStruct = programdata
				allcourcedataout.Cource_Struct = courcedata

				allcourcedataoutlist = append(allcourcedataoutlist, allcourcedataout)

				available = append(available, true)

				fmt.Println("Program Out: ", allcourcedataout.ProgramStruct, data_out)
			} else {
				available = append(available, false)
			}
		case "ADMS":
			is_present, data_out, _ := GetADMSAdmin(students_uuid, "one")

			var programdataacms ProgramStruct = data_out
			if is_present {
				programdata.UUID = programdataacms.UUID
				programdata.Student_UUID = programdataacms.Student_UUID
				programdata.Program_Name = programdataacms.Program_Name
				programdata.First_Name = programdataacms.First_Name
				programdata.Last_Name = programdataacms.Last_Name
				programdata.Email = programdataacms.Email
				programdata.Payment_Method = programdataacms.Payment_Method
				programdata.Paid = programdataacms.Paid
				programdata.Approved = programdataacms.Approved
				programdata.Applied = programdataacms.Applied
				programdata.Completed = programdataacms.Completed
				programdata.Date = programdataacms.Date

				courcedata = GetFromADMSCoures(students_uuid)

				allcourcedataout.ProgramStruct = programdata
				allcourcedataout.Cource_Struct = courcedata

				allcourcedataoutlist = append(allcourcedataoutlist, allcourcedataout)

				available = append(available, true)

				fmt.Println("Program Out: ", allcourcedataout.ProgramStruct, data_out)
			} else {
				available = append(available, false)
			}

		case "ABDMS":
			is_present, data_out, _ := GetABDMSAdmin(students_uuid, "one")

			var programdataacms ProgramStruct = data_out
			if is_present {
				programdata.UUID = programdataacms.UUID
				programdata.Student_UUID = programdataacms.Student_UUID
				programdata.Program_Name = programdataacms.Program_Name
				programdata.First_Name = programdataacms.First_Name
				programdata.Last_Name = programdataacms.Last_Name
				programdata.Email = programdataacms.Email
				programdata.Payment_Method = programdataacms.Payment_Method
				programdata.Paid = programdataacms.Paid
				programdata.Approved = programdataacms.Approved
				programdata.Applied = programdataacms.Applied
				programdata.Completed = programdataacms.Completed
				programdata.Date = programdataacms.Date

				courcedata = GetFromABDMSCources(students_uuid)

				allcourcedataout.ProgramStruct = programdata
				allcourcedataout.Cource_Struct = courcedata

				allcourcedataoutlist = append(allcourcedataoutlist, allcourcedataout)

				available = append(available, true)

				fmt.Println("Program Out: ", allcourcedataout.ProgramStruct, data_out)
			} else {
				available = append(available, false)
			}

		default:
			available = append(available, false)

		}

	}

	programsavailable.Available = available[0]

	return allcourcedataoutlist, available[0]

}

func ReadMessageAdmin(w http.ResponseWriter, r *http.Request) {

	student_uuid := r.URL.Query().Get("student_uuid")
	message_seen := r.URL.Query().Get("message_seen")

	switch message_seen {
	case "not":
		UpdateMessages(student_uuid)

	}

	messages := ReadMessage(student_uuid)

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	err := tpl.ExecuteTemplate(w, "readMessagesAdmin", messages)

	if err != nil {
		log.Fatal(err)
	}

}

func StudentProfileData(w http.ResponseWriter, r *http.Request) {

	studentuuid := r.URL.Query().Get("student_uuid")

	admin_id := r.URL.Query().Get("out")

	admin_infor := AdminData(admin_id)

	fmt.Println("Admin Id: ", admin_id, "Admin Infor: ", admin_infor)

	studentdataout := GetStudentAllDetails(studentuuid)
	listout := GetStudentPrograms(studentuuid)

	programdataout, present := GetStudentProgramDataAdmin(listout, studentuuid)

	studentdataoutadmin := StudentCourse{
		Available:        present,
		StInfo:           studentdataout,
		AllCourceDataOut: programdataout,
		Admin:            admin_infor,
	}
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	err := tpl.ExecuteTemplate(w, "studentdetailstemplate.html", studentdataoutadmin)

	if err != nil {
		log.Fatal(err)
	}

}

func Example(w http.ResponseWriter, r *http.Request) {
	fmt.Println("The Code Is Working Out")
}

func CloseAdmintDiv(w http.ResponseWriter, r *http.Request) {
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	err := tpl.ExecuteTemplate(w, "student_cource_assesment", nil)

	if err != nil {
		log.Fatal(err)
	}
}

func ACAMSStudentData(w http.ResponseWriter, r *http.Request) {
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	admin_id := r.URL.Query().Get("out")
	acamsstudents := GetACAMSStudents(admin_id)

	fmt.Println("ADMIN ID::::", admin_id)
	fmt.Println("The Route Has Been Hi")
	admin_infor := AdminData(admin_id)

	data_out := AdminPage{
		Admin:     admin_infor,
		AcamsData: acamsstudents,
	}

	err := tpl.ExecuteTemplate(w, "studentdataadmin.html", data_out)

	if err != nil {
		log.Fatal(err)
	}
}

type CreatNews struct {
	Admin AdminInfo
	AllNews []NewsStruct
}

func AdminNews(w http.ResponseWriter, r *http.Request) {

	_, all := ReadNews("o", "many")
	admin_id := r.URL.Query().Get("out")
	admin_infor := AdminData(admin_id)

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	display_data := CreatNews{
		Admin: admin_infor,
		AllNews: all,
	}

	err := tpl.ExecuteTemplate(w, "NewsAdminCreate.html", display_data)
	
	fmt.Println("All The News: ", all)

	if err != nil {
		log.Fatal(err)
	}
}

func DeleteMessageRouter(w http.ResponseWriter, r *http.Request) {
	uuid := r.URL.Query().Get("uuid")
	deleted := DeleteMessage(uuid)
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	if deleted {
		fmt.Println("Deleted")
	}
	err := tpl.ExecuteTemplate(w, "deleted_replacement", nil)

	if err != nil {
		log.Fatal(err)
	}

}

func IsItWorking(w http.ResponseWriter, r *http.Request) {

	fmt.Println("the check is working")
}

func ApproveProgram(w http.ResponseWriter, r *http.Request) {
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	user_uuid := r.URL.Query().Get("user_uuid")
	program := r.URL.Query().Get("program")

	fmt.Println(user_uuid, program)

	switch program {
	case "acams":
		if Update(user_uuid, "acams") {
			fmt.Println("Update was succesfull")

		}

	case "acms":
		if Update(user_uuid, "acms") {
			fmt.Println("Update was succesfull")
		}
	case "adms":
		if Update(user_uuid, "adms") {
			fmt.Println("Update was succesfull")
		}
	case "abdms":
		if Update(user_uuid, "abdms") {
			fmt.Println("Update was succesfull")
		}

	}

	// datastring := fmt.Sprintf("The querris are %s ", dataout)
	// fmt.Fprint(w, datastring)

	// keys, ok := r.URL.Query()["id"]

	// if ok {
	// 	fmt.Println(keys)
	// }

	err := tpl.ExecuteTemplate(w, "approved", nil)

	if err != nil {
		log.Fatal(err)
	}
}

func AdminMessagesGet(w http.ResponseWriter, r *http.Request) {

}

func AdminMessagesSend(w http.ResponseWriter, r *http.Request) {

}
