package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

type Article struct {
	Id uint16
	Title, Anons, FullText string
}

var posts =[]Article{}

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, "%s", err)
	}

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/golang")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	res, err := db.Query("SELECT *  FROM `article`")
	if err != nil {
		panic(err)
	}

	posts = []Article{}
	for res.Next() {
		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.FullText)
		if err != nil {
			panic(err)
		}

		posts = append(posts, post)

		fmt.Println(fmt.Sprintf("Post: %d %s %s", post.Id, post.Title, post.Anons))
	}
		t.ExecuteTemplate(w, "index", posts)
	}

func create(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/create.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, "%s", err)
	}
	t.ExecuteTemplate(w, "create", nil)
}

func save_article(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	anons := r.FormValue("anons")
	full_text := r.FormValue("full_text")

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/golang")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	insert, err := db.Query(fmt.Sprintf("INSERT INTO `article` (`title`,`anons`,`full_text`) VALUES ('%s','%s','%s')", title, anons, full_text))
	if err != nil {
		panic(err)
	}
	defer insert.Close()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func show_post(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/show.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, "%s", err)
		return
	}

	vars := mux.Vars(r)
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/golang")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	res, err := db.Query("SELECT * FROM `article` WHERE `id` = ?", vars["id"])
	if err != nil {
		panic(err)
	}
	defer res.Close()

	var showPost Article
	for res.Next() {
		err = res.Scan(&showPost.Id, &showPost.Title, &showPost.Anons, &showPost.FullText)
		if err != nil {
			panic(err)
		}
	}

	t.ExecuteTemplate(w, "show", showPost)
}

func handleFunc() {

	rtr := mux.NewRouter()
	rtr.HandleFunc("/", index).Methods("GET")
	rtr.HandleFunc("/create", create).Methods("GET")
	rtr.HandleFunc("/save_article", save_article).Methods("POST")
	rtr.HandleFunc("/post/{id:[0-9]+}", show_post).Methods("GET")


	http.Handle("/",rtr)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.ListenAndServe(":8080", nil)
}

func main() {
	handleFunc()
}
