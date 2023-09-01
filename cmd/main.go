package main

import (
	"github.com/diyor200/go-fintech/api"
	"github.com/diyor200/go-fintech/database"
)

func main() {
	database.InitDatabase()
	api.StartAPI()
}
