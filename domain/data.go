package domain

import "time"

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
	Items         []Item    `json:"items"`
}

type Owner struct {
	Id        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}
