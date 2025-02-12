package main

import (
	"stock_api/api"
	_ "stock_api/docs"
)

func main() {
	api.SetupRoutes()
}
