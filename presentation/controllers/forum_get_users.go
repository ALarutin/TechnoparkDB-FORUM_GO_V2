package controllers

import (
	"data_base/database"
	"data_base/presentation/logger"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {

	varMap := mux.Vars(r)
	slug, found := varMap["slug"]
	if !found {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error.Println("not found")
		return
	}

	forum, err := database.GetInstance().GetForum(slug)
	if err != nil {
		if forum.ID == 0 {
			myJSON := fmt.Sprintf(`{"%s%s%s"}`, messageCantFind, cantFindForum, slug)
			w.WriteHeader(http.StatusNotFound)
			_, err = w.Write([]byte(myJSON))
			if err != nil {
				logger.Error.Println(err.Error())
			}
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error.Println(err.Error())
		return
	}

	limit := r.URL.Query().Get("limit")
	if limit == "" {
		limit = "100"
	}
	since := r.URL.Query().Get("since")
	desc := r.URL.Query().Get("desc")

	users, err := database.GetInstance().GetUsers(slug, since, desc, limit)
	if err != nil {
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
