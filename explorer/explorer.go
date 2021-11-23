package explorer

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/go_lang_coins/blockchain"
)

var templates *template.Template

const (
	templateDir string = "explorer/templates/"
)



type homeData struct{
	PageTitle string
	Blocks []*blockchain.Block
}

func home(rw http.ResponseWriter, r *http.Request)  {
	//에러를 따로 만들어서 반환해야한다
	// tmpl, err := template.ParseFiles("templates/home.html")
	// if err != nil{
	// 	log.Fatal(err)
	// }

	//template.Must 를 쓰게되면 if문은 과 err변수를 작성하지 않아도된다

	

	data := homeData{"Home", nil} //데이터 전송
	templates.ExecuteTemplate(rw, "home", data)
}

func add(rw http.ResponseWriter, r *http.Request){
	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(rw, "add", nil)
	case "POST":
		blockchain.Blockchain().AddBlock()
		http.Redirect(rw, r, "/", http.StatusPermanentRedirect)
	}
}

func Start(port int){
	handler := http.NewServeMux()
	templates = template.Must(template.ParseGlob(templateDir +"pages/*.gohtml"))
	templates = template.Must(templates.ParseGlob(templateDir + "partials/*gohtml"))
	handler.HandleFunc("/", home)
	handler.HandleFunc("/add", add)
	fmt.Printf("GoGo http://localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d",port),handler)) //에러확인
}

//nil은 기본 dDefaultServeMux 적용