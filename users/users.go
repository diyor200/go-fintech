package users

import (
	"errors"
	"github.com/diyor200/go-fintech/database"
	"github.com/diyor200/go-fintech/helpers"
	"github.com/diyor200/go-fintech/interfaces"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func prepareToken(user *interfaces.User) string {
	//	Sign token
	tokenContent := jwt.MapClaims{
		"user_id": user.ID,
		"expiry":  time.Now().Add(time.Minute ^ 60).Unix(),
	}

	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenContent)
	token, err := jwtToken.SignedString([]byte("TokenPassword"))
	helpers.HandleErr(err)
	return token
}

func prepareResponse(user *interfaces.User, accounts []interfaces.ResponseAccount, withToken bool) map[string]interface{} {
	//Setup response
	responseUser := interfaces.ResponseUser{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
		Accounts: accounts,
	}
	//	Prepare response
	var response = map[string]interface{}{"message": "all is fine"}
	if withToken {
		var token = prepareToken(user)
		response["jwt"] = token
	}
	response["data"] = responseUser
	return response
}

func Login(username, pass string) map[string]interface{} {
	valid := helpers.Validation(
		[]interfaces.Validation{
			{Value: username, Valid: "username"},
			{Value: pass, Valid: "password"},
		})
	if valid {
		user := &interfaces.User{}
		if database.DB.Where("username=?", username).First(&user).RecordNotFound() {
			return map[string]interface{}{"message": "User not found"}
		}

		//	Verify password
		passErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass))
		if errors.Is(passErr, bcrypt.ErrMismatchedHashAndPassword) && passErr != nil {
			return map[string]interface{}{"message": "Wrong password"}
		}

		// Find account to the user
		var accounts []interfaces.ResponseAccount
		database.DB.Table("accounts").Select("id, name, balance").Where("user_id=?", user.ID).Scan(&accounts)

		var response = prepareResponse(user, accounts, true)
		return response
	} else {
		return map[string]interface{}{"message": "not valid values"}
	}

}

func Register(username, email, password string) map[string]interface{} {
	valid := helpers.Validation(
		[]interfaces.Validation{
			{Value: username, Valid: "username"},
			{Value: email, Valid: "email"},
			{Value: password, Valid: "password"},
		})
	if valid {
		generatedPassword := helpers.HashAndSalt([]byte(password))
		user := &interfaces.User{Username: username, Email: email, Password: generatedPassword}
		database.DB.Create(&user)

		account := interfaces.Account{Type: "Daily Account", Name: username + "'s account",
			Balance: uint(1000), UserID: user.ID}
		database.DB.Create(&account)

		accounts := []interfaces.ResponseAccount{}
		respAccount := interfaces.ResponseAccount{ID: account.ID, Name: account.Name, Balance: int(account.Balance)}
		accounts = append(accounts, respAccount)
		var response = prepareResponse(user, accounts, true)

		return response
	} else {
		return map[string]interface{}{"message": "not valid values"}
	}
}

func GetUser(id, jwt string) map[string]interface{} {
	isValid := helpers.ValidToken(id, jwt)

	if isValid {
		user := &interfaces.User{}
		if database.DB.Where("id=?", id).First(&user).RecordNotFound() {
			return map[string]interface{}{"message": "User not found"}
		}
		// Find account to the user
		var accounts []interfaces.ResponseAccount
		database.DB.Table("accounts").Select("id, name, balance").Where("user_id=?", user.ID).Scan(&accounts)
		defer database.DB.Close()

		var response = prepareResponse(user, accounts, false)

		return response
	} else {
		return map[string]interface{}{"message": "not valid token"}
	}
}
