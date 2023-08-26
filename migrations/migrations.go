package migrations

import (
	"github.com/diyor200/go-fintech/helpers"
	"github.com/diyor200/go-fintech/interfaces"
)

func createAccount() {
	db := helpers.ConnectDB()

	users := &[2]interfaces.User{
		{Username: "diyorbek", Email: "diyorbek@gmail.com"},
		{Username: "diyorbek1", Email: "diyorbek1@gmail.com"},
	}

	for i := 0; i < 2; i++ {
		generatedPassword := helpers.HashAndSalt([]byte(users[i].Username))
		user := interfaces.User{Username: users[i].Username, Email: users[i].Email, Password: generatedPassword}
		db.Create(&user)

		account := interfaces.Account{Type: "Daily Account", Name: users[i].Username + "'s account",
			Balance: uint(1000 * int(i+1)), UserID: users[i].ID}
		db.Create(&account)
	}
	defer db.Close()
}

func Migrate() {
	db := helpers.ConnectDB()
	db.AutoMigrate(&interfaces.User{}, &interfaces.Account{})
	defer db.Close()

	createAccount()
}
