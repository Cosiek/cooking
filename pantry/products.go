package pantry

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

// ============================================================================
// ROUTING ====================================================================
// ============================================================================

type ProductViewsHandler struct {
	db *gorm.DB
}

func addProductHandlers(db *gorm.DB) {
	handler := ProductViewsHandler{db}
	http.HandleFunc("/products", handler.productsIndex)
	http.HandleFunc("/products/new", handler.newProductView)
	http.HandleFunc("/products/details/", handler.productDetailsView)
	http.HandleFunc("/products/edit/", handler.editProductView)
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

func (handler *ProductViewsHandler) productsIndex(w http.ResponseWriter, r *http.Request) {
	var products []Product
	handler.db.Find(&products)

	templateSet, err := template.ParseFiles(
		"pantry/templates/products_index.gtpl",
		"pantry/templates/base.gtpl")
	if err != nil {
		log.Fatal(err)
	}
	templateSet.Execute(w, products)
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

	templateSet, err := template.ParseFiles(
		"pantry/templates/edit_product_form.gtpl",
		"pantry/templates/base.gtpl")
	if err != nil {
		log.Fatal(err)
	}
	templateSet.Execute(w, product)
}

func (handler *ProductViewsHandler) newProductView(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		productId := "1"
		path := "/products/details/" + productId
		http.Redirect(w, r, path, http.StatusFound)
	}
	templateSet, err := template.ParseFiles(
		"pantry/templates/new_product_form.gtpl",
		"pantry/templates/base.gtpl")
	if err != nil {
		log.Fatal(err)
	}
	templateSet.Execute(w, nil)
}

func (handler *ProductViewsHandler) productDetailsView(w http.ResponseWriter, r *http.Request) {
	product := handler.getProductOr404(&w, r)
	if product == nil {
		return
	}
	templateSet, err := template.ParseFiles(
		"pantry/templates/product_details.gtpl",
		"pantry/templates/base.gtpl")
	if err != nil {
		log.Fatal(err)
	}
	templateSet.Execute(w, product)
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
