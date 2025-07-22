package routes

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/kamukwamba/oerisuniversity/dbcode"
	"github.com/kamukwamba/oerisuniversity/encription"

	"github.com/google/uuid"
)

var tpl *template.Template

type NewsStruct struct {
	UUID       string
	Auther     string
	Image_Link string
	Story      string
	Date       string
	Title      string
}

type NewsHomePage struct {
	Present  bool
	NewsMain NewsStruct
	NewsList []NewsStruct
}

type PageName struct {
	Name string
}

func ReadNews(uuid_in, number string) (NewsStruct, []NewsStruct) {

	var news_one NewsStruct
	var news_list []NewsStruct

	get_data := dbcode.SqlRead().DB

	var uuid string
	var title string
	var auther string
	var image_link string
	var story string
	var date string

	switch number {
	case "one":

		stmt, err := get_data.Prepare("select uuid,title,auther, image, story, date  from news where uuid = ?")

		if err != nil {
			fmt.Println("failed to get news")
		}

		defer stmt.Close()

		err = stmt.QueryRow(uuid_in).Scan(&uuid, &title, &auther, &image_link, &story, &date)

		image_link_out := CleanNewsImages(uuid)

		file_link := fmt.Sprintf("/news/%s/%s", image_link_out, image_link)

		news_one = NewsStruct{
			UUID:       uuid,
			Title:      title,
			Auther:     auther,
			Image_Link: file_link,
			Story:      story,
			Date:       date,
		}

		if err != nil {
			fmt.Println("Failed to execute News One")
		}

	case "many":

		rows, err := get_data.Query("select * from news")

		if err != nil {
			fmt.Println("QUERY STATEMENT FAILED: ", err)

		}

		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&uuid, &title, &auther, &image_link, &story, &date)

			image_link_out := CleanNewsImages(uuid)

			file_link := fmt.Sprintf("/news/%s/%s", image_link_out, image_link)

			news_one = NewsStruct{
				UUID:       uuid,
				Title:      title,
				Auther:     auther,
				Image_Link: file_link,
				Story:      story,
				Date:       date,
			}
			if err != nil {
				fmt.Println("FAILED TO LOAD", err)
			}

			news_list = append(news_list, news_one)
		}

	}

	return news_one, news_list
}

func UpdateNews(w http.ResponseWriter, r *http.Request) {
	uuid := r.URL.Query().Get("uuid")
	tpl = template.Must(template.ParseGlob("templates/*.html"))
	data_out, _ := ReadNews("one", uuid)

	err := tpl.ExecuteTemplate(w, "", data_out)

	if err != nil {
		log.Fatal(err)
	}

}

func DeleteNewsRoute(w http.ResponseWriter, r *http.Request) {

	uuid := r.URL.Query().Get("uuid")

	delete_news := dbcode.SqlRead().DB

	delete, err := delete_news.Prepare("delete * from news where uuid = ?")

	if err != nil {
		fmt.Println("failed to delete from news")
	}

	defer delete.Close()

	_, err = delete.Exec(uuid)

	if err != nil {
		fmt.Println("failed to delete news 2")
	}
}

func ReadNewsRoute(w http.ResponseWriter, r *http.Request) {
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	//debug failure to laod templates

	uuid := r.URL.Query().Get("uuid")
	news_out, _ := ReadNews(uuid, "one")

	err := tpl.ExecuteTemplate(w, "newsmain", news_out)

	if err != nil {
		log.Fatal(err)
	}
}

func CleanNewsImages(uuid string) string {

	new_uuid := strings.Split(uuid, "-")

	new_file := "nI"

	for _, item := range new_uuid {
		new_file = new_file + item
	}

	return new_file
}
func Create_News(w http.ResponseWriter, r *http.Request) {

	create_news := dbcode.SqlRead().DB

	r.ParseForm()

	auther := r.FormValue("auther")
	title := r.FormValue("title")
	story := r.FormValue("story")

	fmt.Println("Story: ", story)

	uuid := encription.Generateuudi()

	date := fmt.Sprintf("%s", time.Now().Local())

	new_file_dir := CleanNewsImages(uuid)

	file, handler, err := r.FormFile("file")

	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	file_name := handler.Filename

	data_out := NewsStruct{
		UUID:       uuid,
		Auther:     auther,
		Title:      title,
		Image_Link: file_name,
		Story:      story,
	}

	stmt, err := create_news.Prepare("insert into news (uuid, title,auther, image, story, date) values (?,?,?,?,?,?)")

	if err != nil {
		fmt.Println("failed to insert", err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(uuid, title, auther, file_name, story, date)
	if err != nil {
		fmt.Println("failed to create")
	}

	filePath := filepath.Join("./news", new_file_dir, file_name)

	dst, err := os.Create(filePath)

	if err != nil {
		fmt.Println("FAILED TO UPLOAD FILE", err)
	}

	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		fmt.Println("FAILED TO UPLOAD TO SERVER", err)
	}

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	err = tpl.ExecuteTemplate(w, "newssamples", data_out)

	if err != nil {
		log.Fatal(err)
	}
}

func CreateVisitorTable() {

	dbconn := dbcode.SqlRead().DB

	createtable := `CREATE TABLE IF NOT EXISTS visitors (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		visit_time string,
		counter string
	);`

	stmt, err := dbconn.Prepare(createtable)

	if err != nil {
		fmt.Println("Failed to Create Visitor Table")
	}

	defer stmt.Close()

	_, err = stmt.Exec()

	if err != nil {
		fmt.Println("Failed to Execute")
	}
}

func CheckDate(date string) (bool, string) {
	is_present := true
	var current_count string

	dbconn := dbcode.SqlRead().DB

	stmt, err := dbconn.Prepare("select visit_time, counter from visitors where visit_time = ?")

	if err != nil {

		fmt.Println("CheckDate: ", err)
		is_present = false
	}

	defer stmt.Close()

	var is_date string

	err = stmt.QueryRow(date).Scan(&is_date, &current_count)

	if err != nil {
		fmt.Println("There is not date in string ", err)
		is_present = false
	}

	return is_present, current_count
}

func LoadVisited() []Visited {
	dbconn := dbcode.SqlRead().DB
	var data_out Visited
	var data_out_list []Visited
	stmt, err := dbconn.Query("select visit_time, counter from visitors")

	if err != nil {
		fmt.Println("Query Statement failed")
	}

	defer stmt.Close()

	for stmt.Next() {
		err := stmt.Scan(&data_out.Date, &data_out.Count)
		if err != nil {
			fmt.Println("Failed to scan")
			log.Fatal(err)
		}
		data_out_list = append(data_out_list, data_out)

	}

	return data_out_list
}

func CreateVisitor(date, id string) bool {
	year, month, day := time.Now().Date()
	year_out := strconv.Itoa(year)
	month_out := month.String()
	day_out := strconv.Itoa(day)

	createnow := false

	dbconn := dbcode.SqlRead().DB

	data := []string{year_out, month_out, day_out}

	date_out := strings.Join([]string(data), "-")

	is_present, count := CheckDate(date_out)

	if is_present {
		count_out, _ := strconv.Atoi(count)
		counter := count_out + 1
		stmtu, err := dbconn.Prepare("Update visitors SET counter = ?  where visit_time = ?")

		if err != nil {
			fmt.Println("Prepare Failed to Load", err)
		}

		defer stmtu.Close()

		_, err = stmtu.Exec(counter, date_out)
		if err != nil {
			fmt.Println("Failed to Update", err)
		}
	} else {

		stmtc, err := dbconn.Prepare("insert into visitors(counter, visit_time) values(?,?)")

		if err != nil {
			fmt.Println("Prepare statement Failed", err)
		}

		defer stmtc.Close()

		counter := 1
		_, err = stmtc.Exec(counter, date_out)

		if err != nil {

			fmt.Println("Filed to create new visito: ", err)
		}
	}

	return createnow
}

func ClearCookies(w http.ResponseWriter, r *http.Request) {
	dbconn := dbcode.SqlRead().DB

	howMany := r.URL.Query().Get("number")
	date := r.URL.Query().Get("date")

	switch howMany {
	case "all":
		stmtd, err := dbconn.Prepare("delete from visited")

		if err != nil {
			fmt.Println("Pepare Failed")
		}

		defer stmtd.Close()

		_, err = stmtd.Exec()

		if err != nil {
			fmt.Println("Failed to delete all session id")
		}
	case "date":
		stmtd, err := dbconn.Prepare("delete from visitors where visit_time = ?")

		if err != nil {
			fmt.Println("Pepare Failed")
		}

		defer stmtd.Close()

		_, err = stmtd.Exec(date)

		if err != nil {
			fmt.Println("Failed to delete all session id")
		}

	}

}

func HomePage(w http.ResponseWriter, r *http.Request) {

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	year, month, day := time.Now().Date()
	year_out := strconv.Itoa(year)
	month_out := month.String()
	day_out := strconv.Itoa(day)

	data := []string{year_out, month_out, day_out}

	dateVisited := strings.Join([]string(data), "-")

	cookie, err := r.Cookie("visitor_id")

	if err != nil {

		id := uuid.New().String()

		cookie = &http.Cookie{
			Name:     "visitor_id",
			Value:    id,
			Expires:  time.Now().Add(365 * 24 * time.Hour), // Expires in 1 year
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)
		is_created := CreateVisitor(dateVisited, id)

		fmt.Println(is_created)
	} else {
		fmt.Println("Welcome back")
	}

	err = tpl.ExecuteTemplate(w, "index.html", nil)

	if err != nil {
		log.Fatal(err)
	}
}

func NewsPage(w http.ResponseWriter, r *http.Request) {

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	_, all := ReadNews("o", "many")
	var present bool

	var latest NewsStruct

	if len(all) >= 1 {
		if len(all) > 1 {
			latest = all[len(all)-1]

		} else {
			latest = all[0]

		}
		present = true
	} else {
		present = false
	}

	news_main := NewsHomePage{
		Present:  present,
		NewsMain: latest,
		NewsList: all,
	}

	err := tpl.ExecuteTemplate(w, "NewsMainScreen.html", news_main)

	if err != nil {
		log.Fatal(err)
	}
}

func Curiculum(w http.ResponseWriter, r *http.Request) {
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	//debug failure to laod templates

	page_name := PageName{
		Name: "curriculum",
	}

	err := tpl.ExecuteTemplate(w, "curriculum.html", page_name)

	if err != nil {
		log.Fatal(err)
	}
}
