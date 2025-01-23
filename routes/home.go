package routes

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/kamukwamba/oerisuniversity/dbcode"
	"github.com/kamukwamba/oerisuniversity/encription"
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
	NewsMain NewsStruct
	NewsList []NewsStruct
}

type PageName struct {
	Name string
}

func ReadNews(uuid, number string) (NewsStruct, []NewsStruct) {

	var news_one NewsStruct
	var news_list []NewsStruct

	get_data := dbcode.SqlRead().DB

	switch number {
	case "one":

		stmt, err := get_data.Prepare("select uuid,title,auther, image, story, date  from news where uuid = ?")

		if err != nil {
			fmt.Println("failed to get news")
		}

		defer stmt.Close()

		err = stmt.QueryRow(uuid).Scan(&news_one.UUID, &news_one.Title, &news_one.Auther, &news_one.Image_Link, &news_one.Story, &news_one.Date)

		if err != nil {
			fmt.Println("Failed to execute News One")
		}

	case "many":

		rows, err := get_data.Query("select * from news")

		if err != nil {
			fmt.Println("failed to laod news list")

		}

		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&news_one.UUID, &news_one.Title, &news_one.Auther, &news_one.Image_Link, &news_one.Story, &news_one.Date)

			if err != nil {
				fmt.Println("failed to get multiple news stories")
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

func Create_News(w http.ResponseWriter, r *http.Request) {

	create_news := dbcode.SqlRead().DB

	r.ParseForm()

	auther := r.FormValue("auther")
	title := r.FormValue("title")
	image_link := r.FormValue("image")
	story := r.FormValue("story")
	
	fmt.Println("Story: ", story)

	uuid := encription.Generateuudi()
	
	date := fmt.Sprintf("%s", time.Now().Local())

	data_out := NewsStruct{
		UUID:       uuid,
		Auther:     auther,
		Title:      title,
		Image_Link: image_link,
		Story:      story,
	}


	stmt, err := create_news.Prepare("insert into news (uuid, title,auther, image, story, date) values (?,?,?,?,?,?)")

	if err != nil {
		fmt.Println("failed to insert", err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(uuid, title, auther, image_link, story, date)
	if err != nil {
		fmt.Println("failed to create")
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
	
	if err !=  nil {
		
		fmt.Println("CheckDate: ", err)
		is_present = false
	}
	
	defer stmt.Close()
	
	var is_date string 

	err  = stmt.QueryRow(date).Scan(&is_date, &current_count)
	
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

func CreateVisitor(date string) bool {
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
		
		if err != nil{
			
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

	_, err := r.Cookie("visited")

	if err != nil {

		http.SetCookie(w, &http.Cookie{
			Name:    "visited",
			Value:   dateVisited,
			Expires: time.Now().Add(100 * 24 * time.Hour),
		})
		is_created := CreateVisitor(dateVisited)
		
		fmt.Println(is_created)
		
		

	} else {
		fmt.Println("Has Aleady Visited Us!!!")
		
	}

	err = tpl.ExecuteTemplate(w, "index.html", nil)

	if err != nil {
		log.Fatal(err)
	}
}

func NewsPage(w http.ResponseWriter, r *http.Request) {

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	_, all := ReadNews("o", "many")

	var latest NewsStruct

	if len(all) >= 1 {
		if len(all) > 1 {
			latest = all[len(all)-1]

		} else {
			latest = all[0]
		}
	}

	news_main := NewsHomePage{
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
