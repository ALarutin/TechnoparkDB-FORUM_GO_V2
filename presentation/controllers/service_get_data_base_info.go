package controllers

import (
	"data_base/database"
	"data_base/presentation/logger"
	json "github.com/mailru/easyjson"
	"net/http"
)

func GetDataBaseInfoHandler(w http.ResponseWriter, _ *http.Request) {

	db, err := database.GetInstance().GetDatabase()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error.Println(err.Error())
		return
	}

	db.Post--
	db.User--
	db.Forum--
	db.Thread--

	data, err := json.Marshal(db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error.Println(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		logger.Error.Println(err.Error())
	}
}
