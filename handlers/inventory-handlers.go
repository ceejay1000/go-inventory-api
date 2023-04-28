package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
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
	}

}

func GetInventory(w http.ResponseWriter, r *http.Request) {

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

}

func DeleteInventory(w http.ResponseWriter, r *http.Request) {

}
