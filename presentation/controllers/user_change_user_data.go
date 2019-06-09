package controllers

import (
	"data_base/database"
	"data_base/models"
	"data_base/presentation/logger"
	"fmt"
	"github.com/gorilla/mux"
	json "github.com/mailru/easyjson"
	"io/ioutil"
	"net/http"
)

func ChangeUserDataHandler(w http.ResponseWriter, r *http.Request) {
	varMap := mux.Vars(r)
	nickname, found := varMap["nickname"]
	if !found {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error.Println("not found")
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error.Println(err.Error())
		return
	}

	var user models.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error.Println(err.Error())
		return
	}

	user.Nickname = nickname

	u, err := database.GetInstance().UpdateUser(user)
	if err != nil {
		if err.Error() == errorUniqueViolation {
			myJSON := fmt.Sprintf(`{"message": "%s%s"}`, user.Email, emailUsed)
			w.WriteHeader(http.StatusConflict)
			_, err := w.Write([]byte(myJSON))
			if err != nil {
				logger.Error.Println(err.Error())
			}
			return
		}
		if err.Error() == errorPqNoDataFound {
			myJSON := fmt.Sprintf(`{"%s%s%s"}`, messageCantFind, cantFindUser, nickname)
			w.WriteHeader(http.StatusNotFound)
			_, err := w.Write([]byte(myJSON))
			if err != nil {
				logger.Error.Println(err.Error())
			}
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error.Println(err.Error())
		return
	}

	data, err := json.Marshal(u)
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
