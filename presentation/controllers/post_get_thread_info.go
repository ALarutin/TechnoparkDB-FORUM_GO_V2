package controllers

import (
	"data_base/database"
	"data_base/presentation/logger"
	"fmt"
	"github.com/gorilla/mux"
	json "github.com/mailru/easyjson"
	"net/http"
	"strconv"
	"strings"
)

func GetThreadInfoPostHandler(w http.ResponseWriter, r *http.Request) {
	//start := time.Now()
	varMap := mux.Vars(r)
	slug, found := varMap["id"]
	if !found {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error.Println("not found")
		return
	}

	id, err := strconv.Atoi(slug)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error.Println(err.Error())
		return
	}

	related := strings.Split(r.URL.Query().Get("related"), ",")

	//if time.Since(start) > time.Millisecond * 10{
	//	logger.Info.Print(time.Since(start))
	//}
	//start = time.Now()
	postInfo, err := database.GetInstance().GetPostInfo(id, related)
	if err != nil {
		if postInfo.Post.ID == 0 {
			myJSON := fmt.Sprintf(`{"%s%s%v"}`, messageCantFind, cantFindPost, id)
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
	//if time.Since(start) > time.Millisecond * 10{
	//	logger.Info.Print(time.Since(start))
	//}
	//start = time.Now()

	data, err := json.Marshal(postInfo)
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
	//if time.Since(start) > time.Millisecond * 10{
	//	logger.Info.Print(time.Since(start))
	//}
}
