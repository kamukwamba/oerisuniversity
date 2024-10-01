package routes

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/kamukwamba/oerisuniversity/dbcode"
	"github.com/kamukwamba/oerisuniversity/encription"
)

var tpl *template.Template

type NewsStruct struct {
	UUID       string
	Auther     string
	Image_Link string
	Story      string
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

		stmt, err := get_data.Prepare("select uuid,auther, image, story, date  from news where uuid = ?")

		if err != nil {
			fmt.Println("failed to get news")
		}

		defer stmt.Close()

	case "many":

	}
	return news_one, news_list
}

func ReadNewsRoute(w http.ResponseWriter, r *http.Request) {

}

func Create_News(auther, story, image_link string) {
	create_news := dbcode.SqlRead().DB

	uuid := encription.Generateuudi()
	student_create, err := create_news.Begin()

	if err != nil {
		fmt.Println("not working")
	}

	stmt, err := student_create.Prepare("insert into news (uuid,auther, image, story, date) values (?,?,?,?,?)")

	if err != nil {
		fmt.Println("failed to insert")
	}

	defer stmt.Close()

	_, err = stmt.Exec(uuid, auther, story, image_link)
	if err != nil {
		fmt.Println("failed to create")
	}

	err = student_create.Commit()

	if err != nil {
		fmt.Println("failed to commit")
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

	one, all := ReadNews("o", "many")

	news_main := NewsHomePage{
		NewsMain: one,
		NewsList: all,
	}

	err := tpl.ExecuteTemplate(w, "NewsMainScreen.html", news_main)

	if err != nil {
		log.Fatal(err)
	}
}
