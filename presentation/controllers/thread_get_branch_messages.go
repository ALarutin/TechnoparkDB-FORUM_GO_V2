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
	//start := time.Now()
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
	//if time.Since(start) > time.Millisecond * 10{
	//	logger.Info.Print(time.Since(start))
	//
	//}
	//start = time.Now()
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
	//if time.Since(start) > time.Millisecond * 10{
	//	logger.Info.Print(time.Since(start))
	//}

	limit := r.URL.Query().Get("limit")
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		if err.Error() == `strconv.Atoi: parsing "": invalid syntax` {
			limitInt = 100
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error.Println(err.Error())
			return
		}
	}

	since := r.URL.Query().Get("since")
	sinceInt, err := strconv.Atoi(since)
	if err != nil {
		if err.Error() == `strconv.Atoi: parsing "": invalid syntax` {
			sinceInt = 0
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error.Println(err.Error())
			return
		}
	}

	sort := r.URL.Query().Get("sort")
	desc := r.URL.Query().Get("desc")
	var descBool bool
	if desc == "true" {
		descBool = true
	} else if desc == "false" {
		descBool = false
	}

	//if time.Since(start) > time.Millisecond * 10{
	//	logger.Info.Print(time.Since(start))
	//
	//}
	//start = time.Now()

	posts, err := database.GetInstance().GetPosts(id, limitInt, sinceInt, sort, descBool)
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
