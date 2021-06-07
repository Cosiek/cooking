package pantry

import (
	"html/template"
	"log"
	"net/http"
	"regexp"

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

func getIndexView() View {
	return View{
		Name:    "index",
		Regex:   regexp.MustCompile(`^/$`),
		Methods: []string{"GET"},
		Func:    index,
	}
}

func StartServer(db *gorm.DB) {
	// prepare router
	router := Router{}

	// add views
	addProductViews(&router, db)
	router.addView(getIndexView())

	http.HandleFunc("/", router.serve)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
