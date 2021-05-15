package pantry

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func productsIndex(w http.ResponseWriter, r *http.Request) {
	templateSet, err := template.ParseFiles(
		"pantry/templates/products_index.gtpl",
		"pantry/templates/base.gtpl")
	if err != nil {
		log.Fatal(err)
	}
	templateSet.Execute(w, nil)
}

func newProductView(w http.ResponseWriter, r *http.Request) {
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

func productDetailsView(w http.ResponseWriter, r *http.Request) {
	idMatch := r.URL.Path[len("/products/details/"):]
	id, err := strconv.Atoi(idMatch)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if id != 1 {
		http.NotFound(w, r)
		return
	}

	templateSet, err := template.ParseFiles(
		"pantry/templates/product_details.gtpl",
		"pantry/templates/base.gtpl")
	if err != nil {
		log.Fatal(err)
	}
	templateSet.Execute(w, nil)
}
