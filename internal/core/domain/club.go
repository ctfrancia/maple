package domain

import "github.com/google/uuid"

type Club struct {
	ID        int
	PublicID  uuid.UUID
	Name      string
	Address   string
	City      string
	State     string
	Zip       string
	Country   string
	Phone     string
	Email     string
	Twitter   string
	Facebook  string
	Instagram string
	Youtube   string
	Tiktok    string
	Discord   string
	Website   string
}
