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
		if thread.ID == 0 {
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

	if len(inputPosts) == 0 {
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(`[]`))
		if err != nil {
			logger.Error.Println(err.Error())
		}
		return
	}

	if inputPosts[0].Parent != 0 {
		post, err := database.GetInstance().GetPost(inputPosts[0].Parent)
		if err != nil {
			if err.Error() == errorPqNoDataFound {
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
		if post.Thread != thread.ID {
			myJSON := fmt.Sprintf(`{"%s%s"}`, messageCantFind, cantFindParentOrUser)
			w.WriteHeader(http.StatusConflict)
			_, err = w.Write([]byte(myJSON))
			if err != nil {
				logger.Error.Println(err.Error())
			}
			return
		}
	}

	user, err := database.GetInstance().GetUser(inputPosts[0].Author)
	if err != nil {
		if user.ID == 0 {
			myJSON := fmt.Sprintf(`{"%s%s%s"}`, messageCantFind, cantFindUser, inputPosts[0].Author)
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

	outPosts, err := database.GetInstance().CreatePost(inputPosts, thread.ID, thread.Forum)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error.Println(err.Error())
		return
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
