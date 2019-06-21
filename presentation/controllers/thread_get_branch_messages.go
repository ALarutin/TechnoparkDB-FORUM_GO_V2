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

func GetBranchMessagesHandler(w http.ResponseWriter, r *http.Request) {
	varMap := mux.Vars(r)
	slug, found := varMap["slug_or_id"]
	if !found {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error.Println("not found")
		return
	}

	id, err := strconv.Atoi(slug)
	if err != nil {
		id = -1
	} else {
		slug = ""
	}

	id, err = database.GetInstance().GetThreadIdBySlug(slug, id)
	if err != nil {
		if id == 0 {
			myJSON := fmt.Sprintf(`{"%s%s%s"}`, messageCantFind, cantFindThreadSlug, slug)
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



	sort := r.URL.Query().Get("sort")
	desc := r.URL.Query().Get("desc")

	//if time.Since(start) > time.Millisecond * 10{
	//	logger.Info.Print(time.Since(start))
	//
	//}
	//start = time.Now()

	posts, err := database.GetInstance().GetPosts(id, limit, since, sort, desc)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error.Println(err.Error())
		return
	}
	//if time.Since(start) > time.Millisecond * 10{
	//	logger.Info.Print(time.Since(start))
	//
	//}
	//start = time.Now()

	w.WriteHeader(http.StatusOK)
	if len(posts) == 0 {
		_, err = w.Write([]byte(`[]`))
		if err != nil {
			logger.Error.Println(err.Error())
		}
	} else {
		data, err := json.Marshal(posts)
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
	//if time.Since(start) > time.Millisecond * 10{
	//	logger.Info.Print(time.Since(start))
	//
	//}
}
