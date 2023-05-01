package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	data "github.com/ceejay1000/inventory_api/domain"
	UUID "github.com/google/uuid"
)

type ErrorResponse struct {
	Message string    `json:"message"`
	Status  uint      `json:"status"`
	Time    time.Time `json:"time"`
}

type InventoryHandler struct{}

func (in *InventoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		GetInventory(w, r)
	case http.MethodPost:
		PostInventory(w, r)
	case http.MethodPut:
		PutInventory(w, r)
	case http.MethodDelete:
		DeleteInventory(w, r)
	}

}

func GetInventory(w http.ResponseWriter, r *http.Request) {
	urlSegments := strings.Split(r.URL.Path, "/")
	inventoryId := ExtractInventoryId(urlSegments)

	if inventoryId != "" {
		inventoryExists := data.InventoryExistsById(inventoryId)

		if !inventoryExists {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Inventory with ID: " + inventoryId + " does not exist"))
			return
		}

		inventory := data.GetInventoryById(inventoryId)

		if inventory != nil {
			parsedInventory := data.InventoryToJSON(inventory)

			if parsedInventory == nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Could not retrieve data please try again later!"))
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusFound)
			w.Write(parsedInventory)
			return
		}

		if inventory == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("No inventory found"))
			return
		}
	}

	inventoryData, err := data.GetAllInventories()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err != nil {
		errorResponse := new(ErrorResponse)

		errorResponse.Message = "Sorry! An error occurred whiles parsing data"
		errorResponse.Status = http.StatusInternalServerError
		errorResponse.Time = time.Now()

		errMessage, err := json.Marshal(errorResponse)

		if err != nil {
			log.Println("Error parsing data")
			w.Write([]byte("An error occurred"))
		}

		w.Write(errMessage)
		return
	}

	w.Write(inventoryData)
}

func PostInventory(w http.ResponseWriter, r *http.Request) {
	newInventory, err := io.ReadAll(r.Body)
	parsedInventory := new(data.Inventory)

	if err != nil {
		log.Println("Unable to parse JSON data")
		return
	}

	parsedInventory.Id = UUID.NewString()
	err = json.Unmarshal(newInventory, parsedInventory)

	if err != nil {
		log.Println("An error occured")
		w.Write([]byte("Internal server error"))
	}

	inExists := data.InventoryExists(*parsedInventory)

	if inExists {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Sorry!, the added inventory already exists. Start adding items to it"))
		return
	}

	data.AddInventory(parsedInventory)

	inventories, err := data.GetAllInventories()

	if err != nil {
		log.Panicln("Unable to parse data")
	}

	// byteInventories, err := json.Marshal(inventories)

	// if err != nil {
	// 	log.Panicln("Unable to parse data")
	// }

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	log.Println("Added a new inventory")
	w.Write(inventories)
}

func PutInventory(w http.ResponseWriter, r *http.Request) {

	urlSegments := strings.Split(r.URL.Path, "/")
	inventoryId := ExtractInventoryId(urlSegments)

	if inventoryId == "" {
		log.Println("Cannot update inventory. Id is not specified")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Cannot update inventory. Id is not specified"))
		return
	}

	inventoryPayload, err := io.ReadAll(r.Body)

	if err != nil {
		log.Println("Unable to parse request body")
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("An error occured, please try again"))
		return
	}

	parsedInventory, err := data.GetInventoryJSON(inventoryPayload)

	if err != nil {
		log.Println("Unable to parse request body")
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("An error occured, please try again"))
		return
	}

	parsedInventory.Id = inventoryId

	updateStatus, inventoryId := data.UpdateInventory(*parsedInventory)

	if updateStatus {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Inventory with ID: " + inventoryId + " updated successfully"))
		return
	}

	if !updateStatus {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Inventory with ID: " + inventoryId + " does not exist"))
		return
	}

}

func DeleteInventory(w http.ResponseWriter, r *http.Request) {

	urlSegments := strings.Split(r.URL.Path, "/")

	inventoryId := ExtractInventoryId(urlSegments)

	if inventoryId == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Inventory with ID not specified"))
		return
	}

	inventoryStatus := data.InventoryExistsById(inventoryId)

	if !inventoryStatus {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Inventory with ID: " + inventoryId + " does not exist"))
		return
	}

	data.DeleteInventory(inventoryId)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Inventory with ID: " + inventoryId + " deleted successfully"))

}

func ExtractInventoryId(urlSegments []string) string {

	inventoryId := urlSegments[len(urlSegments)-1]

	log.Println("Inventory ID: " + inventoryId)
	// regexStr := `[\-\W\0-9\w]+`
	regexStr := `^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$`

	if regexp.MustCompile(regexStr).MatchString(inventoryId) {
		log.Println("Matched inventory ID")
		return inventoryId
	}

	return ""

}
