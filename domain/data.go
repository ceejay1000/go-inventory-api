package domain

import (
	"encoding/json"
	"log"
	"strings"
	"time"
)

// var Inventories = make(map[string][]Inventory)
var Inventories = map[string][]*Inventory{
	"inventories": {},
}

type Item struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Age      string `json:"age"`
}

type Inventory struct {
	Id            string    `json:"id"`
	Name          string    `json:"name"`
	Owner         string    `json:"owner"`
	DateCreated   time.Time `json:"date-created"`
	DateUpdated   time.Time `json:"date-updated"`
	NumberOfItems int       `json:"number-of-items"`
	Items         []*Item   `json:"items"`
}

type Owner struct {
	Id        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

func GetAllInventories() ([]byte, error) {
	in := &Inventories
	inventory, err := json.Marshal(in)

	if err != nil {
		log.Println("Cannot parse data to JSON")
	}

	return inventory, err
}

func AddInventory(in *Inventory) {
	n := &Inventories
	(*n)["inventories"] = append((*n)["inventories"], in)
	log.Println("Inventories")
	log.Println((*n))
}

func InventoryExists(inventory Inventory) bool {
	for _, in := range Inventories["inventories"] {

		if inventory.Id == in.Id || strings.EqualFold(inventory.Name, in.Name) {
			return true
		}
	}

	return false
}
