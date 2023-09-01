package migrations

import (
	"github.com/diyor200/go-fintech/database"
	"github.com/diyor200/go-fintech/helpers"
	"github.com/diyor200/go-fintech/interfaces"
)

func createAccount() {
	users := &[2]interfaces.User{
		{Username: "diyorbek", Email: "diyorbek@gmail.com"},
		{Username: "diyorbek1", Email: "diyorbek1@gmail.com"},
	}

	for i := 0; i < 2; i++ {
		generatedPassword := helpers.HashAndSalt([]byte(users[i].Username))
		user := interfaces.User{Username: users[i].Username, Email: users[i].Email, Password: generatedPassword}
		database.DB.Create(&user)

		account := interfaces.Account{Type: "Daily Account", Name: users[i].Username + "'s account",
			Balance: uint(1000 * int(i+1)), UserID: users[i].ID}
		database.DB.Create(&account)
	}
}

func Migrate() {
	User := &interfaces.User{}
	Account := &interfaces.Account{}
	Transactions := &interfaces.Transactions{}
	database.DB.AutoMigrate(&User, &Account, &Transactions)

	createAccount()
}
