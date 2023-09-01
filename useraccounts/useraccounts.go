package useraccounts

import (
	"fmt"
	"github.com/diyor200/go-fintech/helpers"
	"github.com/diyor200/go-fintech/interfaces"
	"github.com/diyor200/go-fintech/transactions"
)

func updateAccount(id uint, amount int) interfaces.ResponseAccount {
	db := helpers.ConnectDB()
	account := interfaces.Account{}
	responseAcc := interfaces.ResponseAccount{}
	db.Where("id=?", id).First(&account)
	account.Balance = uint(amount)
	db.Save(&account)

	responseAcc.ID = account.ID
	responseAcc.Name = account.Name
	responseAcc.Balance = int(account.Balance)

	defer db.Close()
	return responseAcc
}

func getAccount(id uint) *interfaces.Account {
	db := helpers.ConnectDB()
	defer db.Close()
	account := &interfaces.Account{}
	if db.Where("id=?", id).First(&account).RecordNotFound() {
		return nil
	}
	return account
}

func Transaction(userId, from, to uint, amount int, jwtToken string) map[string]interface{} {

	userIdString := fmt.Sprint(userId)
	fmt.Println(userIdString)
	isValid := helpers.ValidToken(userIdString, jwtToken)
	fmt.Println(isValid)
	if isValid {
		fromAccount := getAccount(from)
		toAccount := getAccount(to)

		if fromAccount == nil || toAccount == nil {
			return map[string]interface{}{"message": "Account not found"}
		} else if fromAccount.ID != userId {
			return map[string]interface{}{"message": "You are not owner of thee account"}
		} else if int(fromAccount.Balance) < amount {
			return map[string]interface{}{"message": "Not enough money"}
		}

		updatedAccount := updateAccount(from, int(fromAccount.Balance)-amount)
		updateAccount(to, int(toAccount.Balance)+amount)

		transactions.CreateTransaction(from, to, amount)

		var response = map[string]interface{}{"message": "all is fine"}
		response["data"] = updatedAccount
		return response
	}
	return map[string]interface{}{"message": "not valid token"}
}
