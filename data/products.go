package data

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

// Product defines the structure for an API product
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required,min=3"`
	Description string  `json:"desc" validate:"required,min=3"`
	Price       float32 `json:"price" validate:"required,gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedAt   string  `json:"-"`
	UpdatedAt   string  `json:"-"`
	DeletedAt   string  `json:"-"`
}

// Products is a collection of Product
type Products []*Product

var ErrProductNotFound = fmt.Errorf("Product not found")

func (p *Product) Validate() error {
	validate := validator.New()
	err := validate.RegisterValidation("sku", validateSKU)
	if err != nil {
		return err
	}

	return validate.Struct(p)
}

func validateSKU(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	matches := re.FindAllString(fl.Field().String(), -1)

	return len(matches) == 1
}

func (p *Product) FromJSON(r io.Reader) error {
	dec := json.NewDecoder(r)
	return dec.Decode(p)
}

// ToJSON serializes the contents of the collection to JSON
// NewEncoder provides better performance than json.Unmarshal as it doesn't have to buffer the output into memory; this reduces allocations and
// the overheads of the service
func (p *Products) ToJSON(w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(p)
}

func GetProducts() Products {
	return productList
}

func AddProduct(p *Product) {
	p.ID = nextID()
	productList = append(productList, p)
}

func UpdateProduct(id int, p *Product) error {
	prod, pos, err := findProduct(id)
	if err != nil {
		return err
	}

	p.ID = id

	if p.Name == "" {
		p.Name = prod.Name
	}

	if p.Description == "" {
		p.Description = prod.Description
	}

	if p.SKU == "" {
		p.SKU = prod.SKU
	}

	if p.Price == 0 {
		p.Price = prod.Price
	}

	productList[pos] = p

	return nil
}

func findProduct(id int) (*Product, int, error) {
	for idx, prod := range productList {
		if prod.ID == id {
			return prod, idx, nil
		}
	}

	return nil, -1, ErrProductNotFound
}

func nextID() int {
	lastProduct := productList[len(productList)-1]
	return lastProduct.ID + 1
}

var productList = []*Product{
	{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc323",
		CreatedAt:   time.Now().UTC().String(),
		UpdatedAt:   time.Now().UTC().String(),
	},
	{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and string coffee without milk",
		Price:       1.99,
		SKU:         "fjd34",
		CreatedAt:   time.Now().UTC().String(),
		UpdatedAt:   time.Now().UTC().String()},
}
