package models

type User struct {
	ID             int    `json:"id"`
	Username       string `json:"username"`
	HashedPassword string `json:"-"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}


type CreateUserRequest struct {
	Username       string `json:"username" validate:"required,min=3,max=50"`
	HashedPassword string `json:"hashed_password" validate:"required,min=8,max=100"`
}