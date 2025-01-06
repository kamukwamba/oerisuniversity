package main

import (
	"fmt"
	"log"
	"net/http"

	"os"

	"github.com/kamukwamba/oerisuniversity/dbcode"
	"github.com/kamukwamba/oerisuniversity/routes"

	_ "github.com/mattn/go-sqlite3"
)

func Enviroment(w http.ResponseWriter, r *http.Request) {

	result := os.Getenv("Working")
	fmt.Fprintln(w, result)
}

func main() {

	//LOAD DATA TABLES
	dbcode.LoadDB()

	routes.LoadACAMS()
	routes.LoadACMS()
	routes.LoadADMS()
	routes.LoadABDMS()
	routes.LoadAssesments()
	routes.LoadCource()
	routes.LoadExamTable()
	routes.LoadAdminUsers()
	routes.LoadAssesmentTable()

	fmt.Println("::SERVER STARTED::")

	router := http.NewServeMux()

	//ENVIROMENTAL VARIABLES

	router.HandleFunc("/env", Enviroment)

	//MAIN SCREEN
	router.HandleFunc("/", routes.HomePage)
	router.HandleFunc("/aboutus/{id}", routes.AboutUs)
	router.HandleFunc("/programs", routes.Programs)
	router.HandleFunc("/enroll", routes.Enrollment)
	router.HandleFunc("/confirmenrrol", routes.ConfirmEnrollment)
	router.HandleFunc("/adminlogin", routes.AdminLogin)
	router.HandleFunc("/admindashboard", routes.AdminDashboard)
	router.HandleFunc("/programcards", routes.Programcards)
	router.HandleFunc("/news", routes.NewsPage)
	router.HandleFunc("/readnewsstory", routes.ReadNewsRoute)
	router.HandleFunc("/deletemessage", routes.DeleteMessageRouter)
	router.HandleFunc("/getstudentdata", routes.ChangeStudentPassword)
	router.HandleFunc("/updtestudentpassword", routes.ChangePassword)
	router.HandleFunc("/closeupdatedata", routes.CloseUpdateData)

	//ADMIN DASHBOARD
	router.HandleFunc("/acamsstudentdata", routes.ACAMSStudentData)
	router.HandleFunc("/programdetails/{id}", routes.ProgramDetails)
	router.HandleFunc("/adminNews", routes.AdminNews)
	router.HandleFunc("/adminMessages", routes.AdminMessagesPage)
	router.HandleFunc("/studentedataadmin", routes.StudentProfileData)
	router.HandleFunc("/readmessageadmin", routes.ReadMessageAdmin)
	router.HandleFunc("/courceupdateadmin", routes.ApproveCourceUpdate)
	router.HandleFunc("/approve", routes.ApproveProgram)
	router.HandleFunc("/create_cource_data", routes.CreateCourseData)
	router.HandleFunc("/close_cource_div", routes.CloseCreateCourseData)
	router.HandleFunc("/check_working", routes.IsItWorking)
	router.HandleFunc("/add_cource", routes.AddCourceData)
	router.HandleFunc("/getstudymaterial", routes.GetStudyMaterial)
	router.HandleFunc("/get_admin_users", routes.AdminUsers)
	router.HandleFunc("/createadmin", routes.CreatAdminUser)
	router.HandleFunc("/updateadminuser", routes.UpdateAdminUsers)
	router.HandleFunc("/getupdate", routes.GetUpateAdmin)
	router.HandleFunc("/loadform", routes.LoadAdminForm)
	router.HandleFunc("/update_cource_data", routes.UpdateCourceData)
	router.HandleFunc("/saveucdaata", routes.UpdateProgramDetails)
	router.HandleFunc("/gradeca", routes.GradeCA)
	router.HandleFunc("/delete_cource_data", routes.DeleteCource)

	router.HandleFunc("/deleteadmin", routes.DeleteAdmin)
	router.HandleFunc("/createnews", routes.Create_News)
	router.HandleFunc("/curiculum", routes.Curiculum)

	//EXAM
	router.HandleFunc("/create_page", routes.CreatePage)
	router.HandleFunc("/addexam", routes.AddExam)
	router.HandleFunc("/examdetails", routes.AddExamDetails)
	router.HandleFunc("/faq", routes.FAQ)
	router.HandleFunc("/takeexam", routes.TakeExam)
	router.HandleFunc("/submitexam", routes.SubmitExam)
	router.HandleFunc("/grade_exam", routes.GradeExam)
	router.HandleFunc("/gradeexamination", routes.SaveGrades)
	router.HandleFunc("/examddfdea", routes.GetParticularExam)
	router.HandleFunc("/courcecompleted", routes.CourceCompleted)

	//STUDENT PORTAL ROUTES
	router.HandleFunc("/confirmlogin", routes.ConfirmStudentLogin)
	router.HandleFunc("/studentportal", routes.StudentPortal)
	router.HandleFunc("/approvecource", routes.ApproveCource)
	router.HandleFunc("/messages/{id}", routes.Messages)
	router.HandleFunc("/completed", routes.ProgramCompleted)
	router.HandleFunc("/studentsettings", routes.StudentSettings)
	router.HandleFunc("/studentlogoout/{id}", routes.StudentLogOut)
	router.HandleFunc("/proceed/{id}", routes.StudentProcced)
	router.HandleFunc("/contactinstitution/{id}", routes.ContactInstitution)
	router.HandleFunc("/watchvideo", routes.WatcVideo)
	router.HandleFunc("/sendmessage", routes.SendMsg)
	router.HandleFunc("/studentprofileportal/{id}", routes.StudentProfilePortal)
	router.HandleFunc("/close_assesment_div", routes.CloseAssesmentDiv)
	router.HandleFunc("/close_admin_div", routes.CloseAdmintDiv)
	router.HandleFunc("/error", routes.ErrorPage)

	router.HandleFunc("/example", routes.Example)

	router.HandleFunc("/login", routes.LoginPage)
	router.HandleFunc("/handinassesment", routes.HandInAssesment)
	router.HandleFunc("/grade_assesment", routes.GradeAssesment)
	router.HandleFunc("/deleteassesmentresults", routes.DeleteAssesmentAdmin)

	//LOAD ASSETS

	fs := http.FileServer(http.Dir("assets"))
	router.Handle("/assets/", http.StripPrefix("/assets/", fs))

	//RUN SERVER
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))

}
