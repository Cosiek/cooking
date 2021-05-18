package pantry

import (
	"html/template"
	"log"
	"net/http"
	"regexp"

	"gorm.io/gorm"
)

type View struct {
	Name    string
	Regex   *regexp.Regexp
	Methods []string
	Func    func(http.ResponseWriter, *http.Request)
}

type ViewHandler interface {
	GetView(r *http.Request) *View
}

type IndexViewHandler struct {
}

func (handler *IndexViewHandler) GetView(r *http.Request) *View {
	if r.URL.Path == "/" {
		return &View{Name: "index", Methods: []string{"GET"}, Func: index}
	} else {
		return nil
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

type Router struct {
	viewHandlers []ViewHandler
}

func (router *Router) serve(w http.ResponseWriter, r *http.Request) {
	for _, vh := range router.viewHandlers {
		view := vh.GetView(r)
		if view != nil {
			view.Func(w, r)
			return
		}
	}
	// if no handler is willing to manage this, then this is 404
	http.NotFound(w, r)
}

func (router *Router) addViewsHandler(handler ViewHandler) {
	router.viewHandlers = append(router.viewHandlers, handler)
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
