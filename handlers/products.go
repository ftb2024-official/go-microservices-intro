package handlers

import (
	"log"
	"microservice/data"
	"net/http"
	"regexp"
	"strconv"
)

// Products is http.Handler
type Products struct {
	l *log.Logger
}

// NewProducts creates a products handler with the given logger
func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// ServeHTTP is the main entry point for the handler and satisfies the http.Handler interface
func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// handle the request for a list of products
	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
		return
	}

	if r.Method == http.MethodPost {
		p.addProduct(rw, r)
		return
	}

	if r.Method == http.MethodPut {
		p.l.Println("PUT:", r.URL.Path)

		reg := regexp.MustCompile(`/([0-9]+)`)
		g := reg.FindAllStringSubmatch(r.URL.Path, -1)

		if len(g) != 1 {
			p.l.Println("Invalid URI, more than one id...")
			http.Error(rw, "Invalid URI...", http.StatusBadRequest)
			return
		}

		if len(g[0]) != 2 {
			p.l.Println("Invalid URI, more than onw capture group...")
			http.Error(rw, "Invalid URI...", http.StatusBadRequest)
			return
		}

		idString := g[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			p.l.Println("Invalid URI, unable to convert to number...", idString)
			http.Error(rw, "Invalid URI...", http.StatusBadRequest)
			return
		}

		p.updateProduct(id, rw, r)
		return
	}

	// catch all if no method is satisfied return an error
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Products...")
	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json...", http.StatusInternalServerError)
	}
}

func (p *Products) addProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Product...")

	newProduct := &data.Product{}

	err := newProduct.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json...", http.StatusBadRequest)
	}

	data.AddProduct(newProduct)
}

func (p *Products) updateProduct(id int, rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT Product")

	prod := &data.Product{}

	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json...", http.StatusBadRequest)
	}

	err = data.UpdateProduct(id, prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found...", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found...", http.StatusInternalServerError)
		return
	}
}
