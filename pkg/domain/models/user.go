// Package models contains the entities of the domain business.
package models

import (
	"time"

	"futuagro.com/pkg/domain/enums"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// User represent the data of a user
type User struct {
	ID              primitive.ObjectID      `json:"_id" bson:"_id"`
	Name            string                  `json:"name" bson:"name"`
	Surname         string                  `json:"surname" bson:"surname"`
	DocumentType    string                  `json:"documentType" bson:"documentType"`
	DocumentNumber  string                  `json:"documentNumber" bson:"documentNumber"`
	CityID          *primitive.ObjectID     `json:"cityId,omitempty" bson:"cityId"`
	City            *City                   `json:"city,omitempty" bson:"city"`
	Email           string                  `json:"email,omitempty" bson:"email"`
	HashedPassword  string                  `json:"hashedPassword,omitempty" bson:"hashedPassword"`
	AddressLine1    string                  `json:"addressLine1,omitempty" bson:"addressLine1"`
	PhoneNumber     string                  `json:"phoneNumber,omitempty" bson:"phoneNumber"`
	IsEmailVerified bool                    `json:"isEmailVerified" bson:"IsEmailVerified"`
	Crops           *[]Crop                 `json:"crops,omitempty" bson:"crops"`
	Role            string                  `json:"role" bson:"role"`
	RecordStatus    *enums.EnumRecordStatus `json:"recordStatus" bson:"recordStatus"`
	CreatedAt       time.Time               `json:"createdAt" bson:"createdAt"`
	UpdatedAt       time.Time               `json:"updatedAt" bson:"updatedAt"`
}

//HashPassword return the hash of a given password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MaxCost)
	return string(bytes), err
}

//CheckPasswordHash compares a hashed password with its plain text password version
func CheckPasswordHash(hashedPwd string, plainPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd))
	return err == nil
}
