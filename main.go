package main 
import(
	"fmt"
	"net/http"
    "database/sql"
    "log"
	"text/template"
    _ "github.com/lib/pq"
)

type Poem struct{
    Id int
	Author string
	Name string
	Content string
}
func dbConn() (db *sql.DB) {
    dbUrl:="postgres://postgres:postgres@localhost:5432/nydb?sslmode=disable"
    db, err := sql.Open("postgres",dbUrl)
    err = db.Ping()
    if err != nil {
      log.Println(err)
      panic(err)
    }
    log.Println("Database conected")
return db
}
func Init() {
    db:=dbConn()
    execDB,err := db.Query("CREATE TABLE IF NOT EXISTS poems(id SERIAL PRIMARY KEY,Author varchar(50),Name varchar(50),Content text)")
    if(execDB == nil){
        log.Println("doom!")
    }
    if(err != nil){
        panic(err.Error)
        log.Println(err.Error)
    }
    log.Println("init table poems")
}
var tmpl = template.Must(template.ParseGlob("tmpl/*"))
/*func form(w http.ResponseWriter,r *http.Request){
    fmt.Println("r method:",r.Method)
    switch r.Method {
    case "GET":
    	tmpl.ExecuteTemplate(w,"form.html",nil)     
    case "POST":
    	r.ParseForm()
    	poem :=Poem{}
    	poem.Author = r.FormValue("Author")
    	poem.Name = r.FormValue("Name")
    	poem.Content = r.FormValue("Content")
    	//fmt.Println("Content:\n",poem.Content)
    	tmpl.ExecuteTemplate(w,"show.html",poem) 
    }
}*/
func Form(w http.ResponseWriter,r *http.Request){
    fmt.Println("r method:",r.Method)
    if r.Method ==  "GET"{
        tmpl.ExecuteTemplate(w,"form.html",nil)     
    }
}
func Insert(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    if r.Method == "POST" {
        r.ParseForm()
        poem :=Poem{}
        poem.Author = r.FormValue("Author")
        poem.Name = r.FormValue("Name")
        poem.Content = r.FormValue("Content")
        insForm, err := db.Prepare("INSERT INTO poems(Author,Name,Content) VALUES($1,$2,$3)")
        if err != nil {
            panic(err.Error())
        }
        insForm.Exec(poem.Author,poem.Name,poem.Content)
        log.Println("INSERT: Author: " + poem.Author + " | Name: " + poem.Name)
    }
    defer db.Close()
    http.Redirect(w, r, "/show", 301)
}
func Show(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    nId := r.URL.Query().Get("id")
    if(nId == ""){
        nId = "1"
    }
    //nId := 1
    fmt.Println("Show poem id:",nId)
    selDB, err := db.Query("SELECT * FROM poems WHERE id=$1",nId)
    if err != nil {
        panic(err.Error())
    }
    poem := Poem{}
    for selDB.Next() {
        var Id int
        var Author,Name,Content string
        err = selDB.Scan(&Id, &Author, &Name,&Content)
        if err != nil {
            panic(err.Error())
        }
        poem.Id = Id
        poem.Name = Name
        poem.Author = Author
        poem.Content = Content
    }
    tmpl.ExecuteTemplate(w, "show.html", poem)
    defer db.Close()
}
func main() {
    fmt.Println("Sever started on port 8080:")
    Init()
    //http.Handle("/static", http.FileServer(http.Dir("css/")))
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
    http.HandleFunc("/",Show)
    http.HandleFunc("/form",Form)
    http.HandleFunc("/show",Show)
    http.HandleFunc("/insert",Insert)
	err:=http.ListenAndServe(":8080",nil)
	if err!=nil{
		fmt.Println("Sever stopped on port 8080:")
	}
}
/*
postgres://postgres:postgres@localhost:5432/nydb;
*/