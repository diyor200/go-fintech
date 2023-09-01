package transactions

import (
	"github.com/diyor200/go-fintech/database"
	"github.com/diyor200/go-fintech/interfaces"
)

func CreateTransaction(from, to uint, amount int) {
	transaction := &interfaces.Transactions{From: from, To: to, Amount: amount}
	database.DB.Create(&transaction)
}
