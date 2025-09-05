package routes 



import (
	"log"
	"fmt"
	"html/template"
	"net/http"
	"github.com/kamukwamba/oerisuniversity/dbcode"

)


func UpdateAllProgramDetailsR(w http.ResponseWriter, r *http.Request){

	tpl = template.Must(template.ParseGlob("templates/*.html"))


	program_code := r.URL.Query().Get("programCode")
	program_name := r.URL.Query().Get("programName")

	programData := ProgramDataEntry{
		Name: program_name,
		Code: program_code,
	}

	fmt.Println("Name: ", program_name, "\n Code:", program_code)
	
	err := tpl.ExecuteTemplate(w, "updateprgogramcard", programData)

	if err != nil {
		log.Fatal(err)
	}

}


func UpdateAllProgramDetails(w http.ResponseWriter, r *http.Request){

	r.ParseForm()

	old_program_code := r.URL.Query().Get("oldprogramcode")

	new_program_code := r.FormValue("program_code")
	program_name := r.FormValue("program_name")


	err := ChangeTableName(old_program_code, new_program_code)
	if err != nil {
		fmt.Println("Failed to add to new data base")
	}

	err = UpdateTheCourceNamesTable(old_program_code, new_program_code)
	if err != nil{
		fmt.Println("Calling the fuction returned a failed code")
	}

	err = ChangeProgramDataName(old_program_code, new_program_code, program_name)

	if err != nil{
		fmt.Println("Failed to change the Program Data Names")
	}

	err = ChangeCourceTable(old_program_code, new_program_code)

	if err != nil {
		fmt.Println("Failed to update cource table")
	}




}



func ChangeTableName(old_program_code,  new_program_code string) error{


	dbread := dbcode.SqlRead().DB


	query := fmt.Sprintf("ALTER TABLE %s RENAME TO %s", old_program_code, new_program_code)
    
    _, err := dbread.Exec(query)
    if err != nil {
        return fmt.Errorf("failed to rename table: %v", err)
    }
    
    fmt.Printf("Table renamed from %s to %s\n", old_program_code, new_program_code)
    return nil


}


//
func UpdateTheCourceNamesTable(old_program_code, new_program_code string) error{

	dbread := dbcode.SqlRead().DB

	defer dbread.Close()


	stmt, err:= dbread.Prepare("UPDATE CourseNames SET(programCode = ?) WHERE program_code = ?")
	if err != nil {
		fmt.Println("Failed to Update the program name for the cources")
	}

	_, err = stmt.Exec(new_program_code, old_program_code)

	if err != nil {
		fmt.Println("Failed to update the Cource Names Table")
		return err
	}

	return nil

}

//

func ChangeProgramDataName(old_program_code, new_program_code, new_program_name string) error{



	dbread := dbcode.SqlRead().DB

	defer dbread.Close()


	stmt, err:= dbread.Prepare("UPDATE ProgramData SET(programName = ?, programCode = ?) WHERE program_code = ?")
	if err != nil {
		fmt.Println("Failed to Update the program name for the cources")
		return err
	}

	_, err = stmt.Exec(new_program_name, new_program_code, old_program_code)

	if err != nil {
		fmt.Println("Failed to update the ProgramData Table")
		return err
	}

	return nil
}

//
func ChangeCourceTable(old_program_code, new_program_code string) error{


	dbread := dbcode.SqlRead().DB

	defer dbread.Close()


	stmt, err:= dbread.Prepare("UPDATE cource_table SET(program_name = ?) WHERE program_name = ?")
	if err != nil {
		fmt.Println("Failed to Update the cource_table for the cources")
		return err
	}

	_, err = stmt.Exec(new_program_code, old_program_code)

	if err != nil {
		fmt.Println("Failed to update the cource_table Table")
		return err
	}

	return nil
}

