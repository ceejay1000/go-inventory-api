package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	data "github.com/ceejay1000/inventory_api/domain"
	UUID "github.com/google/uuid"
)

type Item struct{}

func (it *Item) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	httpMethod := r.Method

	switch httpMethod {
	case http.MethodGet:
		GetItem(w, r)
	case http.MethodPost:
		AddItem(w, r)
	case http.MethodPut:
		UpdateItem(w, r)
	case http.MethodDelete:
		DeleteItem(w, r)
	}
}

func (i *Item) GetItem() {
	fmt.Println("Jeans")
}

func GetItem(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)

	urlSegment := strings.Split(r.URL.Path, "/")

	inventoryId := ExtractInventoryId(urlSegment)

	if inventoryId == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Inventory ID not specified"))
		return
	}

	for _, in := range data.Inventories["inventories"] {

		if in.Id == inventoryId {

			itemJSON, err := data.ItemsDataToJson(in.Items)

			if err != nil {
				log.Println("Unable to parse ITEM json")
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(itemJSON)
			return

		}

	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Inventory with ID '" + inventoryId + "' not found"))

	// w.Write([]byte(r.URL.Path))
}

func AddItem(w http.ResponseWriter, r *http.Request) {
	inventoryName := ExtractInventoryName(r.URL.Path)

	inventoryStatus := DoesInventoryExist(inventoryName)

	if inventoryStatus {
		addStatus := AddItemToInventory(r.Body, inventoryName)

		if addStatus {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("Item added to inventory successfully"))
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unable to add item to inventory"))
		return
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Inventory '" + inventoryName + "' does not exist"))
}

func UpdateItem(w http.ResponseWriter, r *http.Request) {
	inventoryName := ExtractInventoryName(r.URL.Path)

	inventoryStatus := DoesInventoryExist(inventoryName)

	if inventoryStatus {

		if itemUpdatedStatus := UpdateItemInInventory(r.Body, inventoryName); itemUpdatedStatus {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Item updated successfully"))
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Cannot update item because inventory '" + inventoryName + "' does not exist"))

}

func DeleteItem(w http.ResponseWriter, r *http.Request) {

	inventoryName := ExtractInventoryName(r.URL.Path)

	inventoryStatus := DoesInventoryExist(inventoryName)

	log.Println(r.URL.Path)

	itemId := ExtractItemId(r.URL.Path)

	if itemId == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Inventory url not specified"))
		return
	}

	if !inventoryStatus {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Cannot delete item because, '" + inventoryName + "' does not exist"))
		return
	}

	for _, in := range data.Inventories["inventories"] {

		for index, item := range in.Items {

			if item.Id == itemId {
				in.Items = append(in.Items[0:index], in.Items[index+1:len(in.Items)]...)
				log.Println("Item deleted successfully")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Item with ID '" + itemId + "' was deleted successfully"))
				return
			}

		}
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Item with ID '" + itemId + "' not found"))

}

func ExtractInventoryName(urlPath string) string {

	log.Println("Url path " + urlPath)
	urlSegments := strings.Split(urlPath, "/")

	inventoryName := strings.Trim(urlSegments[len(urlSegments)-2], " ")

	log.Println("Inventory name " + inventoryName + " with length " + fmt.Sprintf("%d", len(inventoryName)))
	return inventoryName
}

func DoesInventoryExist(inventoryName string) bool {

	for _, in := range data.Inventories["inventories"] {
		if strings.EqualFold(in.Name, inventoryName) {
			return true
		}
	}

	return false

}

// Add Items to inventory
func AddItemToInventory(requestBody io.ReadCloser, inventoryName string) bool {
	item := data.Item{}

	parsedItem, err := item.ItemJsonToStruct(requestBody)
	parsedItem.Id = UUID.NewString()
	parsedItem.Age = time.Hour.String()

	if err != nil {
		return false
	}

	for _, in := range data.Inventories["inventories"] {

		if strings.EqualFold(in.Name, inventoryName) {
			in.Items = append(in.Items, &parsedItem)
			log.Println("Item added to inventory successfully")
			return true
		}

	}

	log.Println("Could not add item to inventory")
	return false
}

func UpdateItemInInventory(itemJson io.ReadCloser, inventoryName string) bool {

	itemData := data.Item{}

	parsedItem, err := itemData.ItemJsonToStruct(itemJson)

	if err != nil {
		log.Println("Cannot parse JSON data")
		return false
	}

	for _, in := range data.Inventories["inventories"] {

		if strings.EqualFold(in.Name, inventoryName) {

			for _, item := range in.Items {

				if strings.EqualFold(item.Name, parsedItem.Name) {
					*item = parsedItem
					log.Println("Item updated successfully")
					return true
				}
			}

		}
	}

	log.Println("Unable to update item")
	return false
}

func ExtractItemId(urlPath string) string {
	urlSegments := strings.Split(urlPath, "/")
	itemId := urlSegments[len(urlSegments)-1]

	regexStr := `[\-\W\0-9\w]+`

	if regexp.MustCompile(regexStr).MatchString(itemId) {
		return itemId

	}

	return ""
}
