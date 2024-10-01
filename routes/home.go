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

	data_out := NewsStruct{
		Auther:     auther,
		Title:      title,
		Image_Link: image_link,
		Story:      story,
	}
	uuid := encription.Generateuudi()
	student_create, err := create_news.Begin()
	date := fmt.Sprintf("%s", time.Now().Local())
	if err != nil {
		fmt.Println("not working")
	}

	stmt, err := student_create.Prepare("insert into news (uuid, title,auther, image, story, date) values (?,?,?,?,?,?)")

	if err != nil {
		fmt.Println("failed to insert")
	}

	defer stmt.Close()

	_, err = stmt.Exec(uuid, title, auther, image_link, story, date)
	if err != nil {
		fmt.Println("failed to create")
	}

	err = student_create.Commit()

	if err != nil {
		fmt.Println("failed to commit")
	}

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	err = tpl.ExecuteTemplate(w, "", data_out)

	if err != nil {
		log.Fatal(err)
	}
}

func HomePage(w http.ResponseWriter, r *http.Request) {

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	//debug failure to laod templates

	err := tpl.ExecuteTemplate(w, "index.html", nil)

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
