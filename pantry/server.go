package pantry

import (
	"html/template"
	"log"
	"net/http"

	"gorm.io/gorm"
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

func StartServer(db *gorm.DB) {
	// add routes
	http.HandleFunc("/", index)
	addProductHandlers(db)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
