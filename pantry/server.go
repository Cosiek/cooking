package pantry

import (
	"html/template"
	"log"
	"net/http"
	"regexp"
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

func StartServer() {
	// prepare router
	router := Router{}

	// add views
	router.addView(getIndexView())

	http.HandleFunc("/", router.serve)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
