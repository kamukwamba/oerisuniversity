package routes

import (
	"log"
	"ucmps/dbcode"
)

func LoadExam() {
	dbread := dbcode.SqlRead().DB

	defer dbread.Close()

	//CREATE ADMS
	create_exam := `
		create table if not exists exam_table(
			uuid blod,
			student_uuid blob,
			program_name text,
			cource_name text,
			grade text,
			remark text,
			comment text);`

	_, err := dbread.Exec(create_exam)
	if err != nil {
		log.Fatal(err)
	}
}

func LoadAssesments() {
	dbread := dbcode.SqlRead().DB

	defer dbread.Close()

	//CREATE ADMS
	create_table := `
		create table if not exists assesment_table(
			uuid blod,
			student_uuid blob,
			program_name text,
			cource_name text,
			grade text,
			remark text,
			comment text);`

	_, err := dbread.Exec(create_table)
	if err != nil {
		log.Fatal(err)
	}
}

func LoadCource() {
	dbread := dbcode.SqlRead().DB

	defer dbread.Close()

	//CREATE ADMS
	create_cource := `
		create table if not exists cource_table(
			uuid blod,
			program_name text,
			cource_name text,
			cource_assesment text,
			video_list text,
			module text,
			recomended_book text,
			exam_file text);`

	_, err := dbread.Exec(create_cource)
	if err != nil {
		log.Fatal(err)
	}
}
