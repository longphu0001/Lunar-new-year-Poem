package main 
import(
	"fmt"
	"net/http"
	"text/template"
)
type Poem struct{
	Author string
	Name string
	Content string
}
var tmpl = template.Must(template.ParseGlob("tmpl/*"))
func form(w http.ResponseWriter,r *http.Request){
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
    	fmt.Println("Content:\n",poem.Content)
    	tmpl.ExecuteTemplate(w,"show.html",poem) 
    }
}
func main() {
	http.HandleFunc("/",form)
    //http.Handle("/static", http.FileServer(http.Dir("css/")))
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	err:=http.ListenAndServe(":8080",nil)
	if err!=nil{
		fmt.Println("Sever started on port 8080:")
	}
}