package models

import "time"

// Customer represent the data of a customer
type Customer struct {
	ID             string    `json:"_id"`
	Name           string    `json:"name"`
	DocumentType   string    `json:"documentType"`
	DocumentNumber string    `json:"documentNumber"`
	CustomerType   string    `json:"customerType"`
	PaymentPeriod  string    `json:"paymentPeriod"`
	City           string    `json:"city"`
	Address        string    `json:"address"`
	WebSite        string    `json:"webSite"`
	PhoneNumber    string    `json:"phoneNumber"`
	Email          string    `json:"email"`
	Genre          string    `json:"genre"`
	Products       string    `json:"products"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
