package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kamukwamba/oerisuniversity/dbcode"
	"github.com/kamukwamba/oerisuniversity/routes"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	fs := http.FileServer(http.Dir("assets"))

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

	fmt.Println("::SERVER STARTED::")

	router := http.NewServeMux()

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
	//ADMIN DASHBOARD
	router.HandleFunc("/confirmlogin", routes.ConfirmStudentLogin) //STUDENT LOGIN CODE
	router.HandleFunc("/acamsstudentdata", routes.ACAMSStudentData)
	router.HandleFunc("/programdetails/{id}", routes.ProgramDetails)
	router.HandleFunc("/adminNews", routes.AdminNews)
	router.HandleFunc("/adminMessages", routes.AdminMessagesPage)
	router.HandleFunc("/studentedataadmin/{id}", routes.StudentProfileData)
	router.HandleFunc("/readmessageadmin", routes.ReadMessageAdmin)
	router.HandleFunc("/courceupdateadmin", routes.ApproveCourceUpdate)
	router.HandleFunc("/approve", routes.ApproveProgram)
	router.HandleFunc("/create_cource_code", routes.CreateCourseData)
	router.HandleFunc("/close_cource_div", routes.CloseCreateCourseData)
	router.HandleFunc("/check_working", routes.IsItWorking)
	router.HandleFunc("/add_cource", routes.AddCourceData)
	router.HandleFunc("/getmaterial", routes.GetMaterial)
	router.HandleFunc("/get_admin_users", routes.AdminUsers)
	router.HandleFunc("/createadmin", routes.CreatAdminUser)
	router.HandleFunc("/updateadminuser", routes.UpdateAdminUsers)
	router.HandleFunc("/getupdate", routes.GetUpateAdmin)
	router.HandleFunc("/loadform", routes.LoadAdminForm)

	router.HandleFunc("/deleteadmin", routes.DeleteAdmin)
	router.HandleFunc("/curiculum", routes.Curiculum)

	//EXAM
	router.HandleFunc("/create_page", routes.CreatePage)
	router.HandleFunc("/addexam", routes.AddExam)
	router.HandleFunc("/faq", routes.FAQ)
	router.HandleFunc("/takeexam", routes.TakeExam)

	//STUDENT PORTAL ROUTES
	router.HandleFunc("/approvecource", routes.ApproveCource)
	router.HandleFunc("/messages/{id}", routes.Messages)
	router.HandleFunc("/completed", routes.ProgramCompleted)
	router.HandleFunc("studentsettings/{id}", routes.StudentSettings)
	router.HandleFunc("/studentlogoout/{id}", routes.StudentLogOut)
	router.HandleFunc("/proceed/{id}", routes.StudentProcced)
	router.HandleFunc("/contactinstitution/{id}", routes.ContactInstitution)
	router.HandleFunc("/watchvideo/{id}", routes.WatcVideo)
	router.HandleFunc("/sendmessage", routes.SendMsg)
	router.HandleFunc("/studentprofileportal/{id}", routes.StudentProfilePortal)
	router.HandleFunc("/close_assesment_div", routes.CloseAssesmentDiv)
	router.HandleFunc("/close_admin_div", routes.CloseAdmintDiv)

	router.HandleFunc("/example", routes.Example)

	router.HandleFunc("/login", routes.LoginPage)
	router.HandleFunc("/handinassesment", routes.HandInAssesment)

	//LOAD ASSETS
	router.Handle("/assets/", http.StripPrefix("/assets", fs))

	//RUN SERVER
	log.Fatal(http.ListenAndServe(":3000", router))

}
