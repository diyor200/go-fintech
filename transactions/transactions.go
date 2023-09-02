package transactions

import (
	"github.com/diyor200/go-fintech/database"
	"github.com/diyor200/go-fintech/helpers"
	"github.com/diyor200/go-fintech/interfaces"
)

func CreateTransaction(from, to uint, amount int) {
	transaction := &interfaces.Transactions{From: from, To: to, Amount: amount}
	database.DB.Create(&transaction)
}

func GetTransactionByAccount(id uint) []interfaces.ResponseTransaction {
	var transactions []interfaces.ResponseTransaction
	database.DB.Table("transactions").Select("id, transactions.from, transactions.to, amount").
		Where(interfaces.Transactions{From: id}).Or(interfaces.Transactions{To: id}).Scan(&transactions)

	return transactions
}

func GetMyTransactions(id string, jwtToken string) map[string]interface{} {
	isValid := helpers.ValidToken(id, jwtToken)

	if isValid {
		var accounts []interfaces.ResponseAccount
		database.DB.Table("accounts").Select("id, name, balance").
			Where("user_id=?", id).Scan(&accounts)

		var transactions []interfaces.ResponseTransaction
		for i := 0; i < len(accounts); i++ {
			accTransactions := GetTransactionByAccount(accounts[i].ID)
			transactions = append(transactions, accTransactions...)
		}

		var response = map[string]interface{}{"message": "all is fine"}
		response["data"] = transactions
		return response

	} else {
		return map[string]interface{}{"message": "not valid token"}
	}

}
