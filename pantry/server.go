package pantry

import (
	"html/template"
	"log"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	templateSet, err := template.ParseFiles(
		"pantry/templates/index.gtpl",
		"pantry/templates/base.gtpl")
	if err != nil {
		log.Fatal(err)
	}
	templateSet.Execute(w, nil)
}

func StartServer() {
	http.HandleFunc("/", index)              // set router
	err := http.ListenAndServe(":8080", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
