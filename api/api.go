package api

import (
	"encoding/json"
	"fmt"
	"github.com/diyor200/go-fintech/helpers"
	"github.com/diyor200/go-fintech/transactions"
	"github.com/diyor200/go-fintech/useraccounts"
	"github.com/diyor200/go-fintech/users"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
)

type Login struct {
	Username string
	Password string
}

type Register struct {
	Username, Email, Password string
}

type TransactionBody struct {
	UserID, From, To uint
	Amount           int
}

func readBody(r *http.Request) []byte {
	body, err := io.ReadAll(r.Body)
	helpers.HandleErr(err)
	return body
}

func apiResponse(call map[string]interface{}, w http.ResponseWriter) {
	if call["message"] == "all is fine" {
		resp := call
		json.NewEncoder(w).Encode(&resp)
	} else {
		resp := call
		json.NewEncoder(w).Encode(&resp)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	//	Ready body
	body := readBody(r)

	//	Handle login
	var formattedBody Login
	err := json.Unmarshal(body, &formattedBody)
	helpers.HandleErr(err)
	login := users.Login(formattedBody.Username, formattedBody.Password)

	//	Prepare response
	apiResponse(login, w)
}

func register(w http.ResponseWriter, r *http.Request) {
	//	Ready body
	body := readBody(r)

	//	Handle login
	var formattedBody Register
	err := json.Unmarshal(body, &formattedBody)
	helpers.HandleErr(err)
	register := users.Register(formattedBody.Username, formattedBody.Email, formattedBody.Password)

	//	Prepare response
	apiResponse(register, w)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["id"]
	auth := r.Header.Get("Authorization")

	user := users.GetUser(userId, auth)
	apiResponse(user, w)
}

func getMyTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userID"]
	auth := r.Header.Get("Authorization")

	transacs := transactions.GetMyTransactions(userId, auth)
	apiResponse(transacs, w)
}

func transaction(w http.ResponseWriter, r *http.Request) {
	body := readBody(r)
	auth := r.Header.Get("Authorization")
	var formattedBody TransactionBody
	err := json.Unmarshal(body, &formattedBody)
	helpers.HandleErr(err)

	transaction := useraccounts.Transaction(formattedBody.UserID, formattedBody.From,
		formattedBody.To, formattedBody.Amount, auth)
	apiResponse(transaction, w)
}

func StartAPI() {
	router := mux.NewRouter()
	router.Use(helpers.PanicHandler)
	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/register", register).Methods("POST")
	router.HandleFunc("/transaction", transaction).Methods("POST")
	router.HandleFunc("/transaction/{id}", getMyTransaction).Methods("GET")
	router.HandleFunc("/users/{userID}", getUser).Methods("GET")
	fmt.Println("App is working on port :8888")
	log.Fatal(http.ListenAndServe(":8888", router))
}
