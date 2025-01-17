package models

import (
	"github.com/google/uuid"
)

type Client struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}

func NewClient(email string) Client {
	return Client{
		ID:    uuid.New(),
		Email: email,
	}
}
