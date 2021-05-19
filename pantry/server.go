package pantry

import (
	"html/template"
	"log"
	"net/http"

	"gorm.io/gorm"
)

type IndexViewHandler struct {
}

func (handler *IndexViewHandler) GetView(r *http.Request) (*View, *http.Request) {
	if r.URL.Path == "/" {
		return &View{Name: "index", Methods: []string{"GET"}, Func: index}, r
	} else {
		return nil, r
	}
}

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
	// prepare router
	router := Router{}

	// add routes
	productViewsHandler := getProductViewsHandler(db)
	router.addViewsHandler(productViewsHandler)
	router.addViewsHandler(&IndexViewHandler{})

	http.HandleFunc("/", router.serve)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
