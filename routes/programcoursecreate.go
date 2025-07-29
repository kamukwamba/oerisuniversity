package routes

import (
	"fmt"
	"log"

	"github.com/kamukwamba/oerisuniversity/dbcode"
	"github.com/kamukwamba/oerisuniversity/encription"
)



type StudentProgramData struct {
	UUID           string
	Student_UUID   string
	Program_Name   string
	First_Name     string
	Last_Name      string
	Email          string
	Payment_Method string
	Paid           string
	Approved       bool
	Applied        bool
	Completed      bool
	Date           string
}

type CourceData struct {
	UUID string
}

// CORRECT THE ACAMS STRUCT ERROR



func ApplyForCource(uuid, cource_name string) bool {
	cource_applied := true

	fmt.Println("Cource name: ", uuid, cource_name)
	dbread := dbcode.SqlRead().DB

	dbupdate_statement := fmt.Sprintf(`UPDATE %s SET applied = ? WHERE uuid = ? `, cource_name)

	statement, err := dbread.Prepare(dbupdate_statement)

	if err != nil {
		error_text := fmt.Sprintf("line 44 error from update prepare:: %s", err)
		ErrorPrintOut("studentportal", "ApplyForCource", error_text)

		cource_applied = false
	}

	defer statement.Close()

	_, errup := statement.Exec(true, uuid)

	if errup != nil {
		error_text := fmt.Sprintf("line 50 error from update prepare:: %s", errup)
		ErrorPrintOut("studentportal", "ApplyForCource", error_text)
		cource_applied = false
	}
	return cource_applied
}


func AddToProgramCources(student_uuid, date_in, payment_type, program_code string) bool {
	dbcreate := dbcode.SqlRead().DB

	course_tables, err := GetProgramCourses(program_code)

	added_to_cource_table := true
	

	for _, item := range cource_tables {

		uuid := encription.Generateuudi()

		cource_create, err := dbcreate.Begin()

		if err != nil {
			error_out := fmt.Sprintf("AddToProgramCources: %s", err)
			ErrorPrintOut("programcoursecreate", "AddToProgramCources", error_out)
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
			date) values(?,?,?,?,?,?,?,?,?,?,?,?)`, item.Code)

		statment, err := cource_create.Prepare(insert_String)

		if err != nil {
			error_out := fmt.Sprintf("failed to insert student in: %s,error: %s", item, err)

			ErrorPrintOut("programcoursecreate", "AddToProgramCources", error_out)
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
			ErrorPrintOut("programcoursecreate", "AddToProgramCources", error_out)

		}

		err = cource_create.Commit()

		if err != nil {
			error_out := fmt.Sprintf("failde to commit to cource table: %s, error: %s", item, err)
			ErrorPrintOut("programcoursecreate", "AddToProgramCources", error_out)

		}

	}

	return added_to_cource_table
}

func GetFromProductCources(student_uuid, program_code string) []CourceStruct {
	var cource_data_out CourceStruct

	var cource_data_out_list []CourceStruct

	course_tables, err := GetProgramCourses(program_code)

	

	for _, item := range cource_tables {
		dbread := dbcode.SqlRead().DB

		data_query_string := fmt.Sprintf("select uuid, student_uuid, cource_name, book, module,video, applied, approved, examined, continuorse_assesment,completed, date from %s   where student_uuid = ?", item.Code)

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


	return cource_data_out_list
}

func CreateProgamData(data_in StudentProgramData, payment_type, program_code string) bool {
	created_succesfully := true

	dbcreate := dbcode.SqlRead().DB

	uuid := encription.Generateuudi()

	student_create, err := dbcreate.Begin()

	if err != nil {
		error_out := fmt.Sprintf("%s", err)
		ErrorPrintOut("acams", "CreateACAMS", error_out)
		created_succesfully = false
	}

	prepare_str := Sprintf(`insert into %(
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
		date) values(?,?,?,?,?,?,?,?,?,?,?,?)`, program_code)

	statment, err := student_create.Prepare()

	if err != nil {
		error_out := fmt.Sprintf("the prepare statment: %s", err)
		ErrorPrintOut("programcoursecreate", "CreateProgram", error_out)
		created_succesfully = false

	}

	defer statment.Close()

	

	_, err = statment.Exec(
		uuid,
		data_in.Student_UUID,
		data_in.Program_Name,
		data_in.First_Name,
		data_in.Last_Name,
		data_in.Email,
		data_in.Applied,
		data_in.Approved,
		data_in.Payment_Method,
		data_in.Paid,
		data_in.Completed,
		data_in.Date,
	)

	if err != nil {
		error_out := fmt.Sprintf("execusion statement: %s", err)

		ErrorPrintOut("acams", "CreateACAMS: ", error_out)
		created_succesfully = false

	}

	err = student_create.Commit()

	if err != nil {
		error_out := fmt.Sprintf("commit statement: %s", err)

		ErrorPrintOut("acams", "CreateACAMS: ", error_out)
		created_succesfully = false

	}
	add_to_cources := AddToProgramCources(data_in.Student_UUID, data_in.Date, payment_type)

	fmt.Println("adding to corces was succesful", add_to_cources)

	fmt.Printf("Creating new acams student in acams database complete\n")
	return created_succesfully

}

func GetProgramAdmin(students_uuid_in, promt, program_code string) (bool, StudentProgramData, []StudentProgramData) {
	var confirmacms bool
	var acams_data_out StudentProgramData
	var acams_data_out_list []StudentProgramData


	var prepare_str string
	promtout := promt
	dbread := dbcode.SqlRead().DB

	switch promtout {

	case "one":

		prepare_str = Sprintf("select uuid, student_uuid, program_name,first_name, last_name, email, applied, approved, payment_method, paid, completed, date from %s where student_uuid = ?", program_code)
		statement, err := dbread.Prepare(prepare_str)

		if err != nil {
			error_out := fmt.Sprintf("%s prepare", err)
			ErrorPrintOut("programcoursecreate", "GetPrograms", error_out)
			confirmacms = false
		}

		defer statement.Close()

		err = statement.QueryRow(students_uuid_in).Scan(&acams_data_out.UUID, &acams_data_out.Student_UUID, &acams_data_out.Program_Name, &acams_data_out.First_Name, &acams_data_out.Last_Name, &acams_data_out.Email, &acams_data_out.Applied, &acams_data_out.Approved, &acams_data_out.Payment_Method, &acams_data_out.Paid, &acams_data_out.Completed, &acams_data_out.Date)

		if err != nil {
			error_out := fmt.Sprintf("%s assigning", err)
			ErrorPrintOut("acams", "GetACAMS", error_out)
		}

		fmt.Println("the approved tag", acams_data_out.Approved)

		if acams_data_out.Applied {
			confirmacms = true
		} else {
			confirmacms = false
		}
		return confirmacms, acams_data_out, acams_data_out_list

	case "multiple":
		prepare_str = Sprintf("select * from %s", program_code)
		rows, err := dbread.Query()

		if err != nil {
			error_out := fmt.Sprintf("getting multiple acams data: %s", err)
			ErrorPrintOut("programcoursecreate", "GetPrograms", error_out)
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
				ErrorPrintOut("programcoursecreate", "GetPrograms", error_out)
			}

			acams_data_out_list = append(acams_data_out_list, acams_data_out)
		}

	}

	return confirmacms, acams_data_out, acams_data_out_list

}

func GetProgramsStudents(students_uuid_in, promt, program_name string) (bool, StudentProgramData, []StudentProgramData) {
	var confirmacms bool
	var acams_data_out ACAMS
	var acams_data_out_list []ACAMS

	promtout := promt
	dbread := dbcode.SqlRead().DB

	var prepare_str string

	switch promtout {

	case "one":
		prepare_str = fmt.Sprintf("select uuid, student_uuid, program_name,first_name, last_name, email, applied, approved, payment_method, paid, completed, date from %s where student_uuid = ?", program_name)
		statement, err := dbread.Prepare(prepare_str)

		if err != nil {
			error_out := fmt.Sprintf("%s prepare", err)
			ErrorPrintOut("acams 395", "GetACMS", error_out)
			confirmacms = false
		}

		defer statement.Close()

		err = statement.QueryRow(students_uuid_in).Scan(&acams_data_out.UUID, &acams_data_out.Student_UUID, &acams_data_out.Program_Name, &acams_data_out.First_Name, &acams_data_out.Last_Name, &acams_data_out.Email, &acams_data_out.Applied, &acams_data_out.Approved, &acams_data_out.Payment_Method, &acams_data_out.Paid, &acams_data_out.Completed, &acams_data_out.Date)

		if err != nil {
			error_out := fmt.Sprintf("%s assigning", err)
			ErrorPrintOut("acams 405", "Get Student Program Data", error_out)
		}

		fmt.Println("the approved tag", acams_data_out.Approved)

		if acams_data_out.Approved {
			confirmacms = true
		} else {
			confirmacms = false
		}
		return confirmacms, acams_data_out, acams_data_out_list

	case "multiple":
		prepare_str = Sprintf("select * from %s", program_name)
		rows, err := dbread.Query(prepare_str)

		if err != nil {
			error_out := fmt.Sprintf("getting multiple acams data: %s", err)
			ErrorPrintOut("acams 422", "GetACAMS", error_out)
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
				ErrorPrintOut("acams 445", "GetACAMS", error_out)
			}

			acams_data_out_list = append(acams_data_out_list, acams_data_out)
		}

	}

	return confirmacms, acams_data_out, acams_data_out_list

}


func CreateCourseTable(course_name string) error{

	dbread := dbcode.SqlRead()

	defer dbread.DB.Close()


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

	_, err := dbread.DB.Exec(create_course_table)
	if err != nil {
		log.Printf("%q: %s\n", err, create_course_table)
		return err
	}
	
	return nil

}



func CreateProgramTabel(table_name string) error{


	dbread := dbcode.SqlRead()

	defer dbread.DB.Close()

	create_program := fmt.Sprintf(`
		create table if not exists %s(
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
		);`, table_name)

	_, err := dbread.DB.Exec(create_program)
	if err != nil {
		log.Printf("%q: %s\n", err, create_program)
		return err
	}
	
	return nil
}


