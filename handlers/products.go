package handlers

import (
	"intro/product-api/data"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	// handle get (standard lib)
	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
		return
	}

	// handle post (standard lib)
	if r.Method == http.MethodPost {
		p.addProduct(rw, r)
		return
	}

	// handle put (standard lib)
	if r.Method == http.MethodPut {
		p.l.Println("Attempt PUT for URI ", r.URL.Path)
		// expect ID in URI
		regex := regexp.MustCompile(`/([0-9]+)`)
		g := regex.FindAllStringSubmatch(r.URL.Path, -1)
		//p.l.Println(g)
		if len(g) != 1 {
			p.l.Println("More than 1 id")
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}
		if len(g[0]) != 2 {
			p.l.Println("More than 1 capture group")
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		idString := g[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			p.l.Println("Unable to convert to number")
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
		}
		p.l.Println("got id", id)

		p.updateProduct(rw, r, id)
		return
	}

	// catch all
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Products")
	lp := data.GetProduct()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Error marshalling product", http.StatusInternalServerError)
	}
}

func (p *Products) addProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Product")
	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Error unmarshalling product", http.StatusBadRequest)
	}

	p.l.Printf("Prod: %#v", prod)
	data.AddProduct(prod)
}

func (p *Products) updateProduct(rw http.ResponseWriter, r *http.Request, id int) {
	p.l.Println("Handle PUT Product")
	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Error unmarshalling product", http.StatusBadRequest)
	}

	err = data.UpdateProduct(id, prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Error updating product", http.StatusInternalServerError)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
	}
}
