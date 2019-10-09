package dtos

//LoginDto is a DTO for the login request
type LoginDto struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=16"`
}
