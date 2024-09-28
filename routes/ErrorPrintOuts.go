package routes

import "fmt"

func ErrorPrintOut(file_name, function, the_error string) {

	error_out_string := fmt.Sprintf("FILE NAME:  %s \nFUNCTION NAME: %s \nTHE ERROR: %s", file_name, function, the_error)
	fmt.Println(error_out_string)
}
