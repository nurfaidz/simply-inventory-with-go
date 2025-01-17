package helpers

import "time"

type RegisterInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type LoginInput struct {
	Email    string `json:"email" valid:"required"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type ProductInput struct {
	Name  string `json:"name" valid:"required"`
	Stock uint8  `json:"stock"`
}

type IncomingItemInput struct {
	Qty        string    `json:"qty" valid:"required"`
	IncomingAt time.Time `json:"incoming_at" valid:"required"`
	UserID     uint      `json:"user_id" valid:"required"`
	ProductID  uint      `json:"product_id" valid:"required"`
}

type OutgoingItemInput struct {
	Qty        string    `json:"qty" valid:"required"`
	OutgoingAt time.Time `json:"outgoing_at" valid:"required"`
	UserID     uint      `json:"user_id" valid:"required"`
	ProductID  uint      `json:"product_id" valid:"required"`
}

type DeleteResponse struct {
	Message string `json:"message"`
}

type APIError struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
