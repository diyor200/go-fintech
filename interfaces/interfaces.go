package interfaces

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Username string
	Email    string
	Password string
}

type Account struct {
	gorm.Model
	Type    string
	Name    string
	Balance uint
	UserID  uint
}

type ResponseAccount struct {
	ID      uint
	Name    string
	Balance int
}

type ResponseUser struct {
	ID       uint
	Username string
	Email    string
	Accounts []ResponseAccount
}

type ResponseTransaction struct {
	Idm, From, To uint
	Amount        int
}

type Validation struct {
	Value string
	Valid string
}

type ErrResponse struct {
	Message string
}

type Transactions struct {
	gorm.Model
	From, To uint
	Amount   int
}
