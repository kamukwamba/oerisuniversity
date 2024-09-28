package routes

import (
	"fmt"
	"log"
	"time"

	"github.com/kamukwamba/oerisuniversity/dbcode"
	"github.com/kamukwamba/oerisuniversity/encription"
)

func CheckADMS(student_uuid string) bool {

	confirmadms := true
	dbread := dbcode.SqlRead().DB
	var acms_data_out ProgramStruct

	statement, err := dbread.Prepare("select uuid, student_uuid,program_name ,first_name, last_name, email, applied, approved, payment_method, paid, completed, date from adms where student_uuid = ?")

	if err != nil {
		error_out := fmt.Sprintf("%s", err)
		ErrorPrintOut("adms 1", "GetADMS", error_out)
		confirmadms = false
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
		fmt.Println("Not in the ADMS")
		confirmadms = false
	}

	return confirmadms

}

func CreateADMS(student_uuid string) bool {
	created := true
	result := GetStudentAllDetails(student_uuid)

	var acms_struct ProgramStruct

	uuid := encription.Generateuudi()
	confirm_creation := true

	dbread := dbcode.SqlRead()
	AddStudent, err := dbread.DB.Begin()

	if err != nil {
		log.Fatal()
	}

	stmt, err := AddStudent.Prepare(`insert into adms(
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
	var program_name = "adms"

	date := fmt.Sprintf("%s", time.Now().Local())

	acms_struct = ProgramStruct{
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

	AddToADMSCources(student_uuid, date, "lamp")

	return created
}

func GetADMSAdmin(students_uuid_in, promt string) (bool, ProgramStruct, []ProgramStruct) {
	var confirmacms bool
	var adms_data_out ProgramStruct
	var adms_data_out_list []ProgramStruct

	promtout := promt
	dbread := dbcode.SqlRead().DB

	switch promtout {

	case "one":
		statement, err := dbread.Prepare("select uuid, student_uuid, program_name,first_name, last_name, email, applied, approved, payment_method, paid, completed, date from adms where student_uuid = ?")

		if err != nil {
			error_out := fmt.Sprintf("%s prepare", err)
			ErrorPrintOut("adms", "GetADMS", error_out)
			confirmacms = false
		}

		defer statement.Close()

		err = statement.QueryRow(students_uuid_in).Scan(&adms_data_out.UUID, &adms_data_out.Student_UUID, &adms_data_out.Program_Name, &adms_data_out.First_Name, &adms_data_out.Last_Name, &adms_data_out.Email, &adms_data_out.Applied, &adms_data_out.Approved, &adms_data_out.Payment_Method, &adms_data_out.Paid, &adms_data_out.Completed, &adms_data_out.Date)

		if err != nil {
			error_out := fmt.Sprintf("%s assigning", err)
			ErrorPrintOut("adms", "GetADMS", error_out)
		}

		fmt.Println("the approved tag", adms_data_out.Approved)

		if adms_data_out.Applied {
			confirmacms = true
		} else {
			confirmacms = false
		}
		return confirmacms, adms_data_out, adms_data_out_list

	case "multiple":
		rows, err := dbread.Query("select * from adms")

		if err != nil {
			error_out := fmt.Sprintf("getting multiple adms data: %s", err)
			ErrorPrintOut("adms", "GetADMS", error_out)
		}

		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(
				&adms_data_out.UUID,
				&adms_data_out.Student_UUID,
				&adms_data_out.Program_Name,
				&adms_data_out.First_Name,
				&adms_data_out.Last_Name,
				&adms_data_out.Email,
				&adms_data_out.Applied,
				&adms_data_out.Approved,
				&adms_data_out.Payment_Method,
				&adms_data_out.Paid,
				&adms_data_out.Completed,
				&adms_data_out.Date,
			)

			if err != nil {
				error_out := fmt.Sprintf("getting multiple acams data for loop: %s", err)
				ErrorPrintOut("adms", "GetADMS", error_out)
			}

			adms_data_out_list = append(adms_data_out_list, adms_data_out)
		}

	}

	return confirmacms, adms_data_out, adms_data_out_list
}

func GetADMS(students_uuid_in, promt string) (bool, ProgramStruct, []ProgramStruct) {

	var confirmacms bool
	var adms_data_out ProgramStruct
	var adms_data_out_list []ProgramStruct

	promtout := promt
	dbread := dbcode.SqlRead().DB

	switch promtout {

	case "one":
		statement, err := dbread.Prepare("select uuid, student_uuid, program_name,first_name, last_name, email, applied, approved, payment_method, paid, completed, date from adms where student_uuid = ?")

		if err != nil {
			error_out := fmt.Sprintf("%s prepare", err)
			ErrorPrintOut("adms", "GetADMS", error_out)
			confirmacms = false
		}

		defer statement.Close()

		err = statement.QueryRow(students_uuid_in).Scan(&adms_data_out.UUID, &adms_data_out.Student_UUID, &adms_data_out.Program_Name, &adms_data_out.First_Name, &adms_data_out.Last_Name, &adms_data_out.Email, &adms_data_out.Applied, &adms_data_out.Approved, &adms_data_out.Payment_Method, &adms_data_out.Paid, &adms_data_out.Completed, &adms_data_out.Date)

		if err != nil {
			error_out := fmt.Sprintf("%s assigning", err)
			ErrorPrintOut("adms", "GetADMS", error_out)
		}

		fmt.Println("the approved tag", adms_data_out.Approved)

		if adms_data_out.Approved {
			confirmacms = true
		} else {
			confirmacms = false
		}
		return confirmacms, adms_data_out, adms_data_out_list

	case "multiple":
		rows, err := dbread.Query("select * from adms")

		if err != nil {
			error_out := fmt.Sprintf("getting multiple adms data: %s", err)
			ErrorPrintOut("adms", "GetADMS", error_out)
		}

		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(
				&adms_data_out.UUID,
				&adms_data_out.Student_UUID,
				&adms_data_out.Program_Name,
				&adms_data_out.First_Name,
				&adms_data_out.Last_Name,
				&adms_data_out.Email,
				&adms_data_out.Applied,
				&adms_data_out.Approved,
				&adms_data_out.Payment_Method,
				&adms_data_out.Paid,
				&adms_data_out.Completed,
				&adms_data_out.Date,
			)

			if err != nil {
				error_out := fmt.Sprintf("getting multiple acams data for loop: %s", err)
				ErrorPrintOut("adms", "GetADMS", error_out)
			}

			adms_data_out_list = append(adms_data_out_list, adms_data_out)
		}

	}

	return confirmacms, adms_data_out, adms_data_out_list

}

func AddToADMSCources(student_uuid, date_in, payment_type string) bool {
	applied := true
	dbcreate := dbcode.SqlRead().DB

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
			error_out := fmt.Sprintf("AddToADMSCources: %s", err)
			ErrorPrintOut("acamsfile", "AddToADMSCources", error_out)
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

			ErrorPrintOut("admsfile", "AddToADMSCources", error_out)
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
			ErrorPrintOut("admsfile", "AddToADMSCources", error_out)

		}

		err = cource_create.Commit()

		if err != nil {
			error_out := fmt.Sprintf("failed to commit to cource table: %s, error: %s", item, err)
			ErrorPrintOut("admsfile", "AddToADMSCources", error_out)

		}

	}

	return applied
}

func GetFromADMSCoures(student_uuid string) []CourceStruct {

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
			ErrorPrintOut("acamsfile", "GetFromACAMSCource", error_out)
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

	fmt.Println("ADMS", cource_data_out)

	return cource_data_out_list
}

func LoadADMS() {

	dbread := dbcode.SqlRead()

	defer dbread.DB.Close()

	//CREATE ADMS
	create_acms := `
		create table if not exists adms(
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

	_, create_acms_error := dbread.DB.Exec(create_acms)
	if create_acms_error != nil {
		log.Printf("%q: %s\n", create_acms_error, create_acms)
	}

	//CREATE THE COURSE TABLES
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
