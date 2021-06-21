package pantry

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"gorm.io/gorm"
)

// ============================================================================
// ROUTING ====================================================================
// ============================================================================

func addProductViews(router *Router, db *gorm.DB) {
	handler := ProductViewsHandler{db: db}

	views := []View{
		{"products_new", regexp.MustCompile(`/products/new`), []string{"GET", "POST"}, handler.newProductView},
		{"products_edit", regexp.MustCompile(`/products/edit/(?P<id>\d+)`), []string{"GET", "POST"}, handler.editProductView},
		{"products_details", regexp.MustCompile(`/products/details/(?P<id>\d+)`), []string{"GET"}, handler.productDetailsView},
		{"products_index", regexp.MustCompile(`/products$`), []string{"GET"}, handler.productsIndexView},
	}

	for _, view := range views {
		router.addView(view)
	}
}

// ============================================================================
// MODEL ======================================================================
// ============================================================================

type Product struct {
	gorm.Model
	Name   string
	Mesure int8
}

func (p *Product) setName(name string) error {
	if len(name) > 30 {
		return errors.New("name to long")
	}
	p.Name = name
	return nil
}

func (p *Product) setMesure(mesureIdStr string) error {
	unknownMesureMsg := "unknown mesure - choose one from the list."
	// try to convert to int
	mesure64, err := strconv.ParseInt(mesureIdStr, 10, 8)
	if err != nil {
		return errors.New(unknownMesureMsg)
	}
	mesure := int8(mesure64)
	if _, ok := Mesures[mesure]; !ok {
		return errors.New(unknownMesureMsg)
	}
	p.Mesure = mesure
	return nil
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
	db *gorm.DB
}

func (handler *ProductViewsHandler) getProductOr404(w *http.ResponseWriter, r *http.Request) *Product {
	// get id from path
	urlMatches := r.Context().Value(CtxKey("urlMatches"))
	idMatch := urlMatches.(map[string]string)["id"]
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

// FORM ===============================

type ProductForm struct {
	r       *http.Request
	product *Product
	valid   bool
	Errors  map[string]string
	Values  map[string]string // default and submitted values
}

func getProductForm(r *http.Request, product *Product) *ProductForm {
	if product == nil {
		product = &Product{}
	}

	errors := map[string]string{"name": "", "mesure": ""}
	values := map[string]string{
		"name":   product.Name,
		"mesure": strconv.Itoa(int(product.Mesure)),
	}
	return &ProductForm{r, product, false, errors, values}
}

func (form *ProductForm) isValid() (bool, error) {
	form.valid = false
	if form.r.Method == "POST" {
		// initial form parsing
		if err := form.r.ParseForm(); err != nil {
			return false, err
		}
	} else {
		// non POST requests are invalid by default
		return false, nil
	}

	// attach passed strings to display to the user, if validation fails
	form.Values["name"] = form.r.Form["name"][0]
	form.Values["mesure"] = form.r.Form["mesure"][0]

	// check the values themselves
	form.valid = true

	if err := form.product.setName(form.r.Form["name"][0]); err != nil {
		form.Errors["name"] = err.Error()
		form.valid = false
	}
	if err := form.product.setMesure(form.r.Form["mesure"][0]); err != nil {
		form.Errors["mesure"] = err.Error()
		form.valid = false
	}

	return form.valid, nil
}

func (form *ProductForm) save(db *gorm.DB) (*Product, error) {
	// check if form was validated
	if !form.valid {
		return nil, errors.New("form must be valid to save - run isValid()")
	}
	db.Save(form.product)
	return form.product, nil
}

func (form *ProductForm) GetName() string {
	return form.Values["name"]
}

func (form *ProductForm) GetMesure() string {
	return form.Values["mesure"]
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
	// TODO: CSRF
	// get product
	product := handler.getProductOr404(&w, r)
	if product == nil {
		return
	}

	form := getProductForm(r, product)
	// validate form
	isValid, err := form.isValid()
	if err != nil {
		// handle form parsing error (400 error)
		http.Error(w, "400 - Bad request", http.StatusBadRequest)
		return
	}
	if isValid {
		// save model
		product, _ = form.save(handler.db) // TODO do something about the error
		// redirect to details page
		path := "/products/details/" + strconv.Itoa(int(product.ID))
		http.Redirect(w, r, path, http.StatusFound)
		return
	}

	// pass form to template
	ctx := map[string]interface{}{"form": form}
	handler.renderTemplate("edit_product_form", ctx, &w, r)
}

func (handler *ProductViewsHandler) newProductView(w http.ResponseWriter, r *http.Request) {
	form := getProductForm(r, nil)
	// validate form
	isValid, err := form.isValid()
	if err != nil {
		// handle form parsing error (400 error)
		http.Error(w, "400 - Bad request", http.StatusBadRequest)
		return
	}
	if isValid {
		// save model
		product, _ := form.save(handler.db) // TODO do something about the error
		// redirect to details page
		path := "/products/details/" + strconv.Itoa(int(product.ID))
		http.Redirect(w, r, path, http.StatusFound)
		return
	}

	ctx := map[string]interface{}{"form": form}
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
