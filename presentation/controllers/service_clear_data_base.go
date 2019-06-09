package controllers

import (
	"data_base/database"
	"data_base/presentation/logger"
	"net/http"
)

func ClearDataBaseHandler(w http.ResponseWriter, r *http.Request) {

	err := database.GetInstance().ClearDatabase()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error.Println(err.Error())
		return
	}
}
