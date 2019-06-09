package controllers

import (
	"data_base/database"
	"data_base/presentation/logger"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {

	varMap := mux.Vars(r)
	slug, found := varMap["slug"]
	if !found {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error.Println("not found")
		return
	}

	limit := r.URL.Query().Get("limit")
	var limitInt int
	if limit == "" {
		limitInt = 100
	} else {
		limitInt, _ = strconv.Atoi(limit)
	}

	since := r.URL.Query().Get("since")

	desc := r.URL.Query().Get("desc")
	var descBool bool
	if desc == "true" {
		descBool = true
	} else if desc == "false" {
		descBool = false
	}

	users, err := database.GetInstance().GetUsers(slug, since, descBool, limitInt)
	if err != nil {
		if err.Error() == errorPqNoDataFound {
			myJSON := fmt.Sprintf(`{"%s%s%s"}`, messageCantFind, cantFindForum, slug)
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

	w.WriteHeader(http.StatusOK)
	if len(users) == 0 {
		_, err = w.Write([]byte(`[]`))
		if err != nil {
			logger.Error.Println(err.Error())
		}
	} else {
		data, err := json.Marshal(users)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error.Println(err.Error())
			return
		}
		_, err = w.Write(data)
		if err != nil {
			logger.Error.Println(err.Error())
		}
	}
}
