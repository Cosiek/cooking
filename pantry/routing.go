package pantry

import (
	"net/http"
	"regexp"
)

type View struct {
	Name    string
	Regex   *regexp.Regexp
	Methods []string
	Func    func(http.ResponseWriter, *http.Request)
}

type ViewHandler interface {
	GetView(r *http.Request) (*View, *http.Request)
}

type Router struct {
	viewHandlers []ViewHandler
}

func (router *Router) serve(w http.ResponseWriter, r *http.Request) {
	for _, vh := range router.viewHandlers {
		view, r := vh.GetView(r)
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
