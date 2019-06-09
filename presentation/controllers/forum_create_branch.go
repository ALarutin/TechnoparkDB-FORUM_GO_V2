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
	"strings"
	"time"
)

func CreateBranchHandler(w http.ResponseWriter, r *http.Request) {

	varMap := mux.Vars(r)
	slugUrl, found := varMap["slug"]
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

	var thread models.Thread

	err = json.Unmarshal(body, &thread)
	if err != nil {
		if strings.HasPrefix(err.Error(), `parsing time "{}"`) {
			thread.Created = time.Time{}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error.Println(err.Error())
			return
		}
	}

	logger.Info.Print(thread.Slug)

	thread.Forum = slugUrl

	t, err := database.GetInstance().CreateThread(thread)
	if err != nil {
		if err.Error() == errorPqNoDataFound {
			myJSON := fmt.Sprintf(`{"%s%s%s or %s%s"}`,
				messageCantFind, cantFindUser, thread.Author, cantFindForum, thread.Forum)
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

	data, err := json.Marshal(t)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error.Println(err.Error())
		return
	}

	if !t.IsNew {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	_, err = w.Write(data)
	if err != nil {
		logger.Error.Println(err.Error())
	}
}
