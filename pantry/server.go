package pantry

import (
	"html/template"
	"log"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("pantry/templates/hello.gtpl")
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(w, nil)
}

func StartServer() {
	http.HandleFunc("/", index)              // set router
	err := http.ListenAndServe(":8080", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
