package data

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type Product struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       string `json:"price"`
	SKU         string `json:"sku"`
	CreatedOn   string `json:"-"`
	UpdatedOn   string `json:"-"`
	DeletedOn   string `json:"-"`
}

type Products []*Product

func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Product) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}

func GetProduct() Products {
	return productList
}

func AddProduct(p *Product) {
	p.ID = GetNextID()
	productList = append(productList, p)
}

func GetNextID() int {
	lp := productList[len(productList)-1]
	return lp.ID + 1
}

func UpdateProduct(id int, p *Product) error {
	_, idx, err := FindProduct(id)
	if err != nil {
		return err
	}
	p.ID = id
	productList[idx] = p
	return nil
}

func FindProduct(id int) (*Product, int, error) {
	for i, product := range productList {
		if product.ID == id {
			return product, i, nil
		}
	}
	return nil, -1, ErrProductNotFound
}

var ErrProductNotFound = fmt.Errorf("product not found")

var productList = []*Product{
	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       "8.80",
		SKU:         "abc123",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	&Product{
		ID:          2,
		Name:        "Esspresso",
		Description: "Short and strong coffee without milk",
		Price:       "7.80",
		SKU:         "def456",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}
