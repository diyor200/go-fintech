package users

import (
	"errors"
	"github.com/diyor200/go-fintech/helpers"
	"github.com/diyor200/go-fintech/interfaces"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func Login(username, pass string) map[string]interface{} {
	//Connect to db
	db := helpers.ConnectDB()
	defer db.Close()
	user := &interfaces.User{}
	if db.Where("username=?", username).First(&user).RecordNotFound() {
		return map[string]interface{}{"message": "User not found"}
	}

	//	Verify password
	passErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass))
	if errors.Is(passErr, bcrypt.ErrMismatchedHashAndPassword) && passErr != nil {
		return map[string]interface{}{"message": "Wrong password"}
	}

	// Find account to the user
	var accounts []interfaces.ResponseAccount{}
	db.Table("accounts").Select("id, name, balance").Where("user_id=?", user.ID).Scan(&accounts)

	//Setup response
	responseUser := interfaces.ResponseUser{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
		Accounts: accounts,
	}

	//	Sign token
	tokenContent := jwt.MapClaims{
		"user_id": user.ID,
		"expiry":  time.Now().Add(time.Minute ^ 60).Unix(),
	}

	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenContent)
	token, err := jwtToken.SignedString([]byte("TokenPassword"))
	helpers.HandleErr(err)

	//	Prepare response
	var response = map[string]interface{}{"message": "all is fine"}
	response["jwt"] = token
	response["data"] = responseUser
	return response
}
