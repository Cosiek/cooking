package pantry

import (
	"context"
	"net/http"
	"regexp"
)

type View struct {
	Name    string
	Regex   *regexp.Regexp
	Methods []string
	Func    func(http.ResponseWriter, *http.Request)
}

type CtxKey string

type Router struct {
	views []*View
}

func (router *Router) serve(w http.ResponseWriter, r *http.Request) {
	view, r := router.GetView(r)
	if view != nil {
		view.Func(w, r)
		return
	}
	// if no view is willing to manage this, then this is 404
	http.NotFound(w, r)
}

func (router *Router) addView(view View) {
	router.views = append(router.views, &view)
}

func (router *Router) GetView(r *http.Request) (*View, *http.Request) {
	var found *View
	for _, view := range router.views {
		// check for a match
		match := view.Regex.FindStringSubmatch(r.URL.Path)
		if match == nil {
			continue
		}
		// get values of named groups
		matchesMap := make(map[string]string)
		for i, name := range view.Regex.SubexpNames() {
			if i != 0 && name != "" {
				matchesMap[name] = match[i]
			}
		}
		// add these to request context
		ctx := context.WithValue(r.Context(), CtxKey("urlMatches"), matchesMap)
		r = r.WithContext(ctx)
		// remember match and break
		found = view
		break
	}
	return found, r
}

func (router *Router) GetViewByName(viewName string) *View {
	for _, view := range router.views {
		if view.Name == viewName {
			return view
		}
	}
	return nil
}
