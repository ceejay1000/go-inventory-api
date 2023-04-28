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

func GetInventoryById(id string) *Inventory {

	for _, in := range Inventories["inventories"] {

		if in.Id == id {
			return in
		}
	}

	return nil
}

func GetAllInventories() ([]byte, error) {
	in := &Inventories
	inventory, err := json.Marshal(in)

	if err != nil {
		log.Println("Cannot parse data to JSON")
	}

	return inventory, err
}

func InventoryToJSON(in *Inventory) []byte {
	inventoryJSON, err := json.Marshal(in)

	if err != nil {
		log.Panicln("Error parsing Data")
		return nil
	}

	return inventoryJSON
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

func InventoryExistsById(id string) bool {
	for _, in := range Inventories["inventories"] {

		if in.Id == id {
			return true
		}
	}

	return false
}

func GetInventoryJSON(inventoryPayload []byte) (*Inventory, error) {

	parsedInventory := new(Inventory)

	err := json.Unmarshal(inventoryPayload, parsedInventory)

	return parsedInventory, err
}

func UpdateInventory(in Inventory) (bool, string) {

	inventoryMatch := false

	for _, inventory := range Inventories["inventories"] {

		if inventory.Id == in.Id {
			inventory.Name = in.Name
			inventory.Items = in.Items
			inventory.Owner = in.Owner
			inventory.DateUpdated = time.Now()
			inventoryMatch = true
		}
	}

	return inventoryMatch, in.Id
}

func DeleteInventory(id string) {

	for index, in := range Inventories["inventories"] {

		if in.Id == id {
			Inventories["inventories"] = append(
				Inventories["inventories"][0:index],
				Inventories["inventories"][index+1:len(Inventories["inventories"])]...)
		}
	}
}
