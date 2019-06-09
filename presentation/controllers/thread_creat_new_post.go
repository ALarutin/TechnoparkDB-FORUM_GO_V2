package controllers

import (
	"data_base/database"
	"data_base/models"
	"data_base/presentation/logger"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func CreatNewPostHandler(w http.ResponseWriter, r *http.Request) {

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

	var thread models.Thread
	if slug == "" {
		thread, err = database.GetInstance().GetThreadById(id)
	} else {
		thread, err = database.GetInstance().GetThreadBySlug(slug)
	}
	if err != nil {
		if err.Error() == errorPqNoDataFound {
			myJSON := fmt.Sprintf(`{"%s%s%s/%d"}`, messageCantFind, cantFindThread, slug, id)
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

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error.Println(err.Error())
		return
	}

	inputPosts := make([]models.Post, 0)
	err = json.Unmarshal(body, &inputPosts)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error.Println(err.Error())
		return
	}

	created := time.Now()

	outPosts := make([]models.Post, 0)
	for _, post := range inputPosts {
		post, err = database.GetInstance().CreatePost(post, created, thread.ID, thread.Forum)
		if err != nil {
			if err.Error() == errorPqNoDataFound {
				myJSON := fmt.Sprintf(`{"%s%s%s"}`, messageCantFind, cantFindUser, post.Author)
				w.WriteHeader(http.StatusNotFound)
				_, err = w.Write([]byte(myJSON))
				if err != nil {
					logger.Error.Println(err.Error())
				}
				return
			}
			if err.Error() == errorForeignKeyViolation {
				myJSON := fmt.Sprintf(`{"%s%s"}`, messageCantFind, cantFindParentOrUser)
				w.WriteHeader(http.StatusConflict)
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
		outPosts = append(outPosts, post)
	}

	data, err := json.Marshal(outPosts)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error.Println(err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(data)
	if err != nil {
		logger.Error.Println(err.Error())
	}
}
