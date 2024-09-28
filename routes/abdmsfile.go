package routes

import (
	"fmt"
	"log"
	"time"
	"ucmps/dbcode"
	"ucmps/encription"
)

func CheckABDMS(student_uuid string) bool {

	confirmacms := true
	dbread := dbcode.SqlRead().DB
	var acms_data_out ACMS

	statement, err := dbread.Prepare("select uuid, student_uuid,program_name ,first_name, last_name, email, applied, approved, payment_method, paid, completed, date from abdms where student_uuid = ?")

	if err != nil {
		error_out := fmt.Sprintf("%s", err)
		ErrorPrintOut("abdms 1", "GetABDMS", error_out)
		confirmacms = false
	}

	defer statement.Close()

	err = statement.QueryRow(student_uuid).Scan(
		&acms_data_out.UUID,
		&acms_data_out.Student_UUID,
		&acms_data_out.Program_Name,
		&acms_data_out.First_Name,
		&acms_data_out.Last_Name,
		&acms_data_out.Email,
		&acms_data_out.Applied,
		&acms_data_out.Approved,
		&acms_data_out.Payment_Method,
		&acms_data_out.Paid,
		&acms_data_out.Completed,
		&acms_data_out.Date,
	)

	if err != nil {
		fmt.Println("Not in the ABDMS")
		confirmacms = false
	}

	return confirmacms

}

func AddToABDMSCources(student_uuid, date_in, payment_type string) bool {
	dbcreate := dbcode.SqlRead().DB

	added_to_cource_table := true

	cource_table := []string{
		"creative_writing",
		"understanding_miracles",
		"channeling_skills",
		"enneagram",
		"mythology_on_gods_and_goddess",
		"herbs",
		"meditation_skills",
		"mantras_and_mudras",
		"divinations",
		"archetypes",
		"basics_in_research",
		"understanding_propaganda",
		"great_Spiritual_teachers",
		"reprogramming",
		"shamanism",
		"mystery_schools_in_the_world",
		"law_and_ethics_in_metaphysical_sciences",
		"non_violet_communication"}

	for _, item := range cource_table {

		uuid := encription.Generateuudi()

		cource_create, err := dbcreate.Begin()

		if err != nil {
			error_out := fmt.Sprintf("AddToACAMSCources: %s", err)
			ErrorPrintOut("abdmsfile", "AddToABDMSCources", error_out)
		}

		insert_String := fmt.Sprintf(`insert into %s(
			uuid,
			student_uuid,
			cource_name,
			book,
			module,
			video,
			applied,
			approved,
			examined,
			continuorse_assesment,
			completed,
			date) values(?,?,?,?,?,?,?,?,?,?,?,?)`, item)

		statment, err := cource_create.Prepare(insert_String)

		if err != nil {
			error_out := fmt.Sprintf("failed to insert student in: %s,error: %s", item, err)

			ErrorPrintOut("abdmsfile", "AddToABDMSCources", error_out)
		}

		defer statment.Close()

		var applied_value bool

		if payment_type == "lamp_sum" {
			applied_value = true
		} else {

			applied_value = false
		}
		var cource_name = item
		var book string = fmt.Sprintf("Book for: %s", item)
		var video string = fmt.Sprintf("video for: %s", item)
		var module string = fmt.Sprintf("video for: %s", item)
		var applied bool = applied_value
		var approved bool = false
		var examinde bool = false
		var continuorse_assesment string = "0"
		var completed bool = false

		_, err = statment.Exec(
			uuid,
			student_uuid,
			cource_name,
			book,
			video,
			module,
			applied,
			approved,
			examinde,
			continuorse_assesment,
			completed,
			date_in,
		)

		if err != nil {
			error_out := fmt.Sprintf("failde to aad to cource table: %s, error: %s", item, err)
			ErrorPrintOut("abdmsfile", "AddToABDMSCources", error_out)

		}

		err = cource_create.Commit()

		if err != nil {
			error_out := fmt.Sprintf("failde to commit to cource table: %s, error: %s", item, err)
			ErrorPrintOut("abdmsfile", "AddToABDMSCources", error_out)

		}

	}

	return added_to_cource_table
}

func GetFromABDMSCources(student_uuid string) []CourceStruct {
	var cource_data_out CourceStruct

	var cource_data_out_list []CourceStruct

	cource_table := []string{
		"creative_writing",
		"understanding_miracles",
		"channeling_skills",
		"enneagram",
		"mythology_on_gods_and_goddess",
		"herbs",
		"meditation_skills",
		"mantras_and_mudras",
		"divinations",
		"archetypes",
		"basics_in_research",
		"understanding_propaganda",
		"great_Spiritual_teachers",
		"reprogramming",
		"shamanism",
		"mystery_schools_in_the_world",
		"law_and_ethics_in_metaphysical_sciences",
		"non_violet_communication"}

	for _, item := range cource_table {
		dbread := dbcode.SqlRead().DB

		data_query_string := fmt.Sprintf("select uuid, student_uuid, cource_name, book, module,video, applied, approved, examined, continuorse_assesment,completed, date from %s where student_uuid = ?", item)

		stmt, err := dbread.Prepare(data_query_string)

		if err != nil {
			error_out := fmt.Sprintf("getting from cource data: %s", err)
			ErrorPrintOut("abdmsfile", "GetFromABDMSCource", error_out)
			log.Fatal(err)
		}

		defer stmt.Close()

		err = stmt.QueryRow(student_uuid).Scan(&cource_data_out.UUID,
			&cource_data_out.Student_UUID,
			&cource_data_out.Cource_Name,
			&cource_data_out.Book,
			&cource_data_out.Module,
			&cource_data_out.Video,
			&cource_data_out.Applied,
			&cource_data_out.Approved,
			&cource_data_out.Examined,
			&cource_data_out.Continuorse_Assesment,
			&cource_data_out.Completed,
			&cource_data_out.Date,
		)

		if err != nil {
			error_out := fmt.Sprintf("falied to query row: %s", err)
			ErrorPrintOut("acamsfile", "GetFromACAMSCource", error_out)
			log.Fatal(err)
		}

		cource_data_out_list = append(cource_data_out_list, cource_data_out)

	}

	fmt.Println("ACAMS", cource_data_out)

	return cource_data_out_list
}

func CreateABDMS(student_uuid string) bool {
	created := true
	result := GetStudentAllDetails(student_uuid)
	fmt.Println("STUDENT DATA TO ADD TO ACMS: ", result)

	var acms_struct ACMS

	uuid := encription.Generateuudi()
	confirm_creation := true

	dbread := dbcode.SqlRead()
	AddStudent, err := dbread.DB.Begin()

	if err != nil {
		log.Fatal()
	}

	stmt, err := AddStudent.Prepare(`insert into abdms(
		uuid,
		student_uuid,
		program_name,
		first_name,
		last_name,
		email,
		applied,
		approved,
		payment_method,
		paid,
		completed,
		date) values(?,?,?,?,?,?,?,?,?,?,?,?)`)

	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close()
	var program_name = "acms"

	date := fmt.Sprintf("%s", time.Now().Local())

	acms_struct = ACMS{
		UUID:           uuid,
		Student_UUID:   result.UUID,
		Program_Name:   program_name,
		First_Name:     result.First_Name,
		Last_Name:      result.Last_Name,
		Email:          result.Email,
		Applied:        true,
		Approved:       false,
		Payment_Method: "lamp",
		Paid:           "pending",
		Completed:      false,
		Date:           date,
	}

	_, err = stmt.Exec(acms_struct.UUID,
		acms_struct.Student_UUID,
		acms_struct.Program_Name,
		acms_struct.First_Name,
		acms_struct.Last_Name,
		acms_struct.Email,
		acms_struct.Applied,
		acms_struct.Approved,
		acms_struct.Payment_Method,
		acms_struct.Paid,
		acms_struct.Completed,
		acms_struct.Date)

	if err != nil {
		log.Fatal(err)
	}

	err = AddStudent.Commit()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(confirm_creation)

	AddToABDMSCources(student_uuid, date, "lamp")

	return created

}

func GetABDMSAdmin(students_uuid_in, promt string) (bool, ProgramStruct, []ProgramStruct) {
	var confirmacms bool
	var acams_data_out ProgramStruct
	var acams_data_out_list []ProgramStruct

	promtout := promt
	dbread := dbcode.SqlRead().DB

	switch promtout {

	case "one":
		statement, err := dbread.Prepare("select uuid, student_uuid, program_name,first_name, last_name, email, applied, approved, payment_method, paid, completed, date from abdms where student_uuid = ?")

		if err != nil {
			error_out := fmt.Sprintf("%s prepare", err)
			ErrorPrintOut("abdms", "GetABDMS", error_out)
			confirmacms = false
		}

		defer statement.Close()

		err = statement.QueryRow(students_uuid_in).Scan(&acams_data_out.UUID, &acams_data_out.Student_UUID, &acams_data_out.Program_Name, &acams_data_out.First_Name, &acams_data_out.Last_Name, &acams_data_out.Email, &acams_data_out.Applied, &acams_data_out.Approved, &acams_data_out.Payment_Method, &acams_data_out.Paid, &acams_data_out.Completed, &acams_data_out.Date)

		if err != nil {
			error_out := fmt.Sprintf("%s assigning", err)
			ErrorPrintOut("abdms", "GetABDMS", error_out)
		}

		fmt.Println("the approved tag", acams_data_out.Approved)

		if acams_data_out.Applied {
			confirmacms = true
		} else {
			confirmacms = false
		}
		return confirmacms, acams_data_out, acams_data_out_list

	case "multiple":
		rows, err := dbread.Query("select * from abdms")

		if err != nil {
			error_out := fmt.Sprintf("getting multiple acams data: %s", err)
			ErrorPrintOut("abdms", "GetABDMS", error_out)
		}

		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(
				&acams_data_out.UUID,
				&acams_data_out.Student_UUID,
				&acams_data_out.Program_Name,
				&acams_data_out.First_Name,
				&acams_data_out.Last_Name,
				&acams_data_out.Email,
				&acams_data_out.Applied,
				&acams_data_out.Approved,
				&acams_data_out.Payment_Method,
				&acams_data_out.Paid,
				&acams_data_out.Completed,
				&acams_data_out.Date,
			)

			if err != nil {
				error_out := fmt.Sprintf("getting multiple acams data for loop: %s", err)
				ErrorPrintOut("abdms", "GetABDMS", error_out)
			}

			acams_data_out_list = append(acams_data_out_list, acams_data_out)
		}

	}

	return confirmacms, acams_data_out, acams_data_out_list

}

func GetABDMS(students_uuid_in, promt string) (bool, ProgramStruct, []ProgramStruct) {
	var confirmacms bool
	var acams_data_out ProgramStruct
	var acams_data_out_list []ProgramStruct

	promtout := promt
	dbread := dbcode.SqlRead().DB

	switch promtout {

	case "one":
		statement, err := dbread.Prepare("select uuid, student_uuid, program_name,first_name, last_name, email, applied, approved, payment_method, paid, completed, date from abdms where student_uuid = ?")

		if err != nil {
			error_out := fmt.Sprintf("%s prepare", err)
			ErrorPrintOut("abdms 345", "GetABDMS", error_out)
			confirmacms = false
		}

		defer statement.Close()

		err = statement.QueryRow(students_uuid_in).Scan(&acams_data_out.UUID, &acams_data_out.Student_UUID, &acams_data_out.Program_Name, &acams_data_out.First_Name, &acams_data_out.Last_Name, &acams_data_out.Email, &acams_data_out.Applied, &acams_data_out.Approved, &acams_data_out.Payment_Method, &acams_data_out.Paid, &acams_data_out.Completed, &acams_data_out.Date)

		if err != nil {
			error_out := fmt.Sprintf("%s assigning", err)
			ErrorPrintOut("abdms 405", "GetABDMS", error_out)
		}

		fmt.Println("the approved tag", acams_data_out.Approved)

		if acams_data_out.Approved {
			confirmacms = true
		} else {
			confirmacms = false
		}
		return confirmacms, acams_data_out, acams_data_out_list

	case "multiple":
		rows, err := dbread.Query("select * from abdms")

		if err != nil {
			error_out := fmt.Sprintf("getting multiple acams data: %s", err)
			ErrorPrintOut("abdms 422", "GetACBDMS", error_out)
		}

		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(
				&acams_data_out.UUID,
				&acams_data_out.Student_UUID,
				&acams_data_out.Program_Name,
				&acams_data_out.First_Name,
				&acams_data_out.Last_Name,
				&acams_data_out.Email,
				&acams_data_out.Applied,
				&acams_data_out.Approved,
				&acams_data_out.Payment_Method,
				&acams_data_out.Paid,
				&acams_data_out.Completed,
				&acams_data_out.Date,
			)

			if err != nil {
				error_out := fmt.Sprintf("getting multiple acams data for loop: %s", err)
				ErrorPrintOut("abdms 445", "GetABDMS", error_out)
			}

			acams_data_out_list = append(acams_data_out_list, acams_data_out)
		}

	}

	return confirmacms, acams_data_out, acams_data_out_list

}

func LoadABDMS() {

	dbread := dbcode.SqlRead()

	defer dbread.DB.Close()

	create_acams := `
		create table if not exists abdms(
			uuid blob not null,
			student_uuid text,
			program_name text,
			first_name text,
			last_name text,
			email text,
			applied bool,
			approved bool,
			payment_method text,
			paid text,
			completed bool,
			date text
		);`

	_, create_acams_error := dbread.DB.Exec(create_acams)
	if create_acams_error != nil {
		log.Printf("%q: %s\n", create_acams_error, create_acams)
	}

	//CREATE THE COURCE TABLES
	cource_table := []string{
		"creative_writing",
		"understanding_miracles",
		"channeling_skills",
		"enneagram",
		"mythology_on_gods_and_goddess",
		"herbs",
		"meditation_skills",
		"mantras_and_mudras",
		"divinations",
		"archetypes",
		"basics_in_research",
		"understanding_propaganda",
		"great_Spiritual_teachers",
		"reprogramming",
		"shamanism",
		"mystery_schools_in_the_world",
		"law_and_ethics_in_metaphysical_sciences",
		"non_violet_communication"}

	for _, item := range cource_table {
		create_course_table := fmt.Sprintf(`
		create table if not exists %s(uuid blob not null, 
			student_uuid text,
			cource_name text,
			book text,
			module text,
			video text,
			applied bool,
			approved bool,
			continuorse_assesment text,
			examined bool,
			completed bool,
			date text);`, item)

		_, create_course_table_error := dbread.DB.Exec(create_course_table)

		if create_course_table_error != nil {
			log.Printf("%q: %s\n", create_course_table_error, create_course_table)

		}
	}

}
