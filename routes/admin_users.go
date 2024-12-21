package routes

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/kamukwamba/oerisuniversity/dbcode"
	"github.com/kamukwamba/oerisuniversity/encription"
)

type AdminUser struct {
	UUID       string
	First_Name string
	Last_Name  string
	Email      string
	Password   string
	Auth       string
	Date       string
}

func CreatAdminUser(w http.ResponseWriter, r *http.Request) {

	uuid := encription.Generateuudi()
	date := fmt.Sprintf("%s", time.Now().Local())

	r.ParseForm()

	admin_out := AdminUser{
		UUID:       uuid,
		First_Name: r.FormValue("first_name"),
		Last_Name:  r.FormValue("last_name"),
		Email:      r.FormValue("email"),
		Password:   r.FormValue("password"),
		Auth:       r.FormValue("auth"),
		Date:       date,
	}

	// first_name := encription.EncryptData(r.FormValue("first_name"))
	// last_name := encription.EncryptData(r.FormValue("last_name"))
	// email := encription.EncryptData(r.FormValue("email"))
	// password := encription.EncryptData(r.FormValue("password"))
	// auth := encription.EncryptData(r.FormValue("auth"))

	first_name := r.FormValue("first_name")
	last_name := r.FormValue("last_name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	auth := r.FormValue("auth")

	fmt.Println(first_name, last_name, email, password, auth)

	create_admin := dbcode.SqlRead().DB
	stmt, err := create_admin.Prepare("insert into admin (uuid, first_name, last_name, email, password, auth, date) values(?,?,?,?,?,?,?)")

	if err != nil {
		fmt.Println("failed to create admin user", err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(uuid, first_name, last_name, email, password, auth, date)

	if err != nil {
		fmt.Println("failed to insert into admin user")
	}

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	err = tpl.ExecuteTemplate(w, "admin_user_tr", admin_out)

	if err != nil {
		log.Fatal(err)
	}

}

func GetAdminUsers(count, uuid string) (AdminUser, []AdminUser) {
	var admin_user AdminUser
	var admin_user_list []AdminUser
	adminData := dbcode.SqlRead().DB

	switch count {
	case "one":
		stmt, err := adminData.Prepare("select uuid, first_name, last_name, email, password, auth, date from admin where uuid = ?")

		if err != nil {
			log.Fatal(err)
		}

		defer stmt.Close()

		err = stmt.QueryRow(uuid).Scan(&admin_user.UUID, &admin_user.First_Name, &admin_user.Last_Name, &admin_user.Email, &admin_user.Password, &admin_user.Auth, &admin_user.Date)

		if err != nil {
			log.Fatal(err)
		}

	case "many":
		rows, err := adminData.Query("select * from admin")

		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&admin_user.UUID, &admin_user.First_Name, &admin_user.Last_Name, &admin_user.Email, &admin_user.Password, &admin_user.Auth, &admin_user.Date)

			if err != nil {
				log.Fatal(err)
			}
			admin_user_list = append(admin_user_list, admin_user)
		}

	}

	return admin_user, admin_user_list

}

func DeletA(uuid string) bool {

	deleted := true
	deleteuser := dbcode.SqlRead().DB

	stmt, err := deleteuser.Prepare("delete from admin where uuid = ?")

	if err != nil {
		fmt.Println("failed to delete one")
		deleted = false

	}
	defer stmt.Close()

	_, errde := stmt.Exec(uuid)

	if errde != nil {
		fmt.Println("failed to delete two")
		deleted = false
	}

	return deleted

}

func DeleteAdmin(w http.ResponseWriter, r *http.Request) {
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	uuid := r.URL.Query().Get("uuid")
	delete := DeletA(uuid)

	err := tpl.ExecuteTemplate(w, "empty_tr", delete)

	if err != nil {
		log.Fatal(err)
	}

}

func UpdateA(uuid string, update_data AdminUser) AdminUser {

	fmt.Println("Update data", update_data)

	updateadmn := dbcode.SqlRead().DB

	admin_user_out := update_data
	stmt, err := updateadmn.Prepare("UPDATE admin SET first_name = ?, last_name = ?, email = ?, password = ?, auth = ? WHERE uuid = ?")

	if err != nil {
		fmt.Println("failed to update admin")

	}
	defer stmt.Close()

	_, errout := stmt.Exec(update_data.First_Name, update_data.Last_Name, update_data.Email, update_data.Password, update_data.Auth, uuid)

	if errout != nil {
		fmt.Println("failed to update admin")
	}

	return admin_user_out

}

func LoadAdminForm(w http.ResponseWriter, r *http.Request) {
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	err := tpl.ExecuteTemplate(w, "formtemplate", nil)

	if err != nil {
		log.Fatal(err)
	}
}

func GetUpateAdmin(w http.ResponseWriter, r *http.Request) {
	tpl = template.Must(template.ParseGlob("templates/*.html"))
	var set_template string
	var admin_data_out AdminUser

	uuid := r.URL.Query().Get("uuid")
	admin_data_out, _ = GetAdminUsers("one", uuid)
	set_template = "updateform"

	err := tpl.ExecuteTemplate(w, set_template, admin_data_out)

	if err != nil {
		log.Fatal(err)
	}
}

func UpdateAdminUsers(w http.ResponseWriter, r *http.Request) {
	tpl = template.Must(template.ParseGlob("templates/*.html"))
	var set_template string
	var admin_data_out AdminUser

	fmt.Println("It is not working out")

	r.ParseForm()
	uuid := r.URL.Query().Get("uuid")

	// admin_out := AdminUser{
	// 	UUID:       uuid,
	// 	First_Name: encription.EncryptData(r.FormValue("first_name")),
	// 	Last_Name:  encription.EncryptData(r.FormValue("last_name")),
	// 	Email:      encription.EncryptData(r.FormValue("email")),
	// 	Password:   encription.EncryptData(r.FormValue("password")),
	// 	Auth:       encription.EncryptData(r.FormValue("auth")),
	// }

	first_name := r.FormValue("first_name")

	fmt.Println(first_name)

	admin_out := AdminUser{
		UUID:       uuid,
		First_Name: r.FormValue("first_name"),
		Last_Name:  r.FormValue("last_name"),
		Email:      r.FormValue("email"),
		Password:   r.FormValue("password"),
		Auth:       r.FormValue("auth"),
	}

	UpdateA(uuid, admin_out)

	set_template = "admin_user_tr"
	admin_data_out = admin_out

	err := tpl.ExecuteTemplate(w, set_template, admin_data_out)

	if err != nil {
		log.Fatal(err)
	}

}

type AdminUserData struct {
	Admin AdminInfo
	Users []AdminUser
}

func AdminUsers(w http.ResponseWriter, r *http.Request) {

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	admin_id := r.URL.Query().Get("out")
	admin_infor := AdminData(admin_id)

	_, admin_user_data_list := GetAdminUsers("many", "none")

	display_data := AdminUserData{
		Admin: admin_infor,
		Users: admin_user_data_list,
	}

	fmt.Println("The user list", admin_user_data_list)

	err := tpl.ExecuteTemplate(w, "users_admin", display_data)

	if err != nil {
		log.Fatal(err)
	}

}

func LoadAdminUsers() {

	dbread := dbcode.SqlRead().DB

	defer dbread.Close()

	admin_user := `
		create table if not exists admin(
		uuid blob not null,
		first_name text, 
		last_name text,
		email text, 
		password text,
		auth text,
		date text);`

	_, admin_user_error := dbread.Exec(admin_user)
	if admin_user_error != nil {
		log.Println(admin_user_error)
	}

}
