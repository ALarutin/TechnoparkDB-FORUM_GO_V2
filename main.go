package main

import (
	"data_base/presentation/logger"
	"data_base/presentation/router"
	"github.com/xlab/closer"
	"net/http"
)

func main() {
	defer closer.Close()
	r := router.GetRouter()
	logger.Info.Printf("Started listening at: 5000")
	logger.Fatal.Println(http.ListenAndServe(":5000", r))
}
