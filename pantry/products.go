package pantry

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

// ============================================================================
// ROUTING ====================================================================
// ============================================================================

func getProductViewsHandler(db *gorm.DB) *ProductViewsHandler {
	handler := ProductViewsHandler{db: db}

	handler.views = []View{
		{"products_new", regexp.MustCompile(`/products/new`), []string{"GET", "POST"}, handler.newProductView},
		{"products_edit", regexp.MustCompile(`/products/edit/(?P<id>\d+)`), []string{"GET", "POST"}, handler.editProductView},
		{"products_details", regexp.MustCompile(`/products/details/(?P<id>\d+)`), []string{"GET"}, handler.productDetailsView},
		{"products_index", regexp.MustCompile(`/products`), []string{"GET"}, handler.productsIndexView},
	}

	return &handler
}

func (handler *ProductViewsHandler) GetView(r *http.Request) *View {
	var found *View
	found = nil
	for _, view := range handler.views {
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
		// remember match and break
		found = &view
		break
	}
	return found
}

// ============================================================================
// MODEL ======================================================================
// ============================================================================

type Product struct {
	gorm.Model
	Name   string
	Mesure int8
}

var Mesures = map[int8]string{
	1: "szt.",
	2: "kg",
	3: "l",
}

// ============================================================================
// VIEWS ======================================================================
// ============================================================================

// VIEWS HANDLER ======================

type ProductViewsHandler struct {
	db    *gorm.DB
	views []View
}

func (handler *ProductViewsHandler) getProductOr404(w *http.ResponseWriter, r *http.Request) *Product {
	// get id from path
	segments := strings.Split(r.URL.Path, "/")
	idMatch := segments[len(segments)-1]
	id, err := strconv.Atoi(idMatch)
	if err != nil {
		http.NotFound(*w, r)
		return nil
	}

	// get product from DB
	var product Product
	result := handler.db.First(&product, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		http.NotFound(*w, r)
		return nil
	}

	return &product
}

func (handler *ProductViewsHandler) renderTemplate(templateName string, ctx map[string]interface{}, w *http.ResponseWriter, r *http.Request) {
	templateSet, err := template.ParseFiles(
		"pantry/templates/"+templateName+".gtpl",
		"pantry/templates/base.gtpl")
	if err != nil {
		log.Fatal(err)
	}

	type TemplateCtx struct {
		Mesures map[int8]string
		Request *http.Request
		Handler *ProductViewsHandler
		Ctx     map[string]interface{}
	}

	templateSet.Execute(*w, TemplateCtx{
		Mesures: Mesures,
		Request: r,
		Handler: handler,
		Ctx:     ctx,
	})
}

// VIEWS ==============================

func (handler *ProductViewsHandler) productsIndexView(w http.ResponseWriter, r *http.Request) {
	var products []Product
	handler.db.Find(&products)

	ctx := map[string]interface{}{
		"products": products,
	}
	handler.renderTemplate("products_index", ctx, &w, r)
}

func (handler *ProductViewsHandler) editProductView(w http.ResponseWriter, r *http.Request) {
	// get product
	product := handler.getProductOr404(&w, r)
	if product == nil {
		return
	}

	if r.Method == "POST" {
		path := "/products/details/" + strconv.Itoa(int(product.ID))
		http.Redirect(w, r, path, http.StatusFound)
		return
	}

	ctx := map[string]interface{}{
		"product": product,
	}
	handler.renderTemplate("edit_product_form", ctx, &w, r)
}

func (handler *ProductViewsHandler) newProductView(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		productId := "1"
		path := "/products/details/" + productId
		http.Redirect(w, r, path, http.StatusFound)
	}

	ctx := map[string]interface{}{}
	handler.renderTemplate("new_product_form", ctx, &w, r)
}

func (handler *ProductViewsHandler) productDetailsView(w http.ResponseWriter, r *http.Request) {
	product := handler.getProductOr404(&w, r)
	if product == nil {
		return
	}

	ctx := map[string]interface{}{
		"product": product,
	}
	handler.renderTemplate("product_details", ctx, &w, r)
}
