package router

import (
	"data_base/presentation/controllers"
	"data_base/presentation/middleware"
	"github.com/gorilla/mux"
	"net/http"
)

type Route struct {
	Info    string
	Name    string
	Path    string
	Method  string
	Handler http.HandlerFunc
}

func GetRouter() (router *mux.Router) {

	router = mux.NewRouter()

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	//forum

	forumSubRouter := router.PathPrefix("/api/forum").Subrouter()

	_forum := []Route{
		{
			Info:    "Handler for creating forum.",
			Name:    "forum_CreatForum",
			Path:    "/create",
			Method:  http.MethodPost,
			Handler: controllers.CreateForumHandler,
		},
		{
			Info:    "Handler for creating branch.",
			Name:    "forum_CreatBranch",
			Path:    "/{slug}/create",
			Method:  http.MethodPost,
			Handler: controllers.CreateBranchHandler,
		},
		{
			Info:    "Handler for obtaining information about the forum.",
			Name:    "forum_GetForumInfo",
			Path:    "/{slug}/details",
			Method:  http.MethodGet,
			Handler: controllers.GetForumInfoHandler,
		},
		{
			Info:    "Handler for getting a list of forum discussion branches.",
			Name:    "forum_GetThreads",
			Path:    "/{slug}/threads",
			Method:  http.MethodGet,
			Handler: controllers.GetThreadsHandler,
		},
		{
			Info:    "Handler for obtaining the users of this forum.",
			Name:    "forum_GetUsers",
			Path:    "/{slug}/users",
			Method:  http.MethodGet,
			Handler: controllers.GetUsersHandler,
		},
	}

	for _, r := range _forum {
		forumSubRouter.
			HandleFunc(r.Path, r.Handler).
			Methods(r.Method).
			Name(r.Name)
	}

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	//post

	postSubRouter := router.PathPrefix("/api/post").Subrouter()

	_post := []Route{
		{
			Info:    "Handler for changing the message.",
			Name:    "post_ChangeMessage",
			Path:    "/{id}/details",
			Method:  http.MethodPost,
			Handler: controllers.ChangeMessageHandler,
		},
		{
			Info:    "Handler for getting information about the discussion thread.",
			Name:    "post_GetThreadInfoPost",
			Path:    "/{id}/details",
			Method:  http.MethodGet,
			Handler: controllers.GetThreadInfoPostHandler,
		},
	}

	for _, r := range _post {
		postSubRouter.
			Methods(r.Method).
			Name(r.Name).
			Path(r.Path).
			HandlerFunc(r.Handler)
	}

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	//service

	serviceSubRouter := router.PathPrefix("/api/service").Subrouter()

	_service := []Route{
		{
			Info:    "Handler for clearing all data in the database.",
			Name:    "service_ClearDataBase",
			Path:    "/clear",
			Method:  http.MethodPost,
			Handler: controllers.ClearDataBaseHandler,
		},
		{
			Info:    "Handler for obtaining information about the database.",
			Name:    "service_GetDataBaseInfo",
			Path:    "/status",
			Method:  http.MethodGet,
			Handler: controllers.GetDataBaseInfoHandler,
		},
	}

	for _, r := range _service {
		serviceSubRouter.
			Methods(r.Method).
			Name(r.Name).
			Path(r.Path).
			HandlerFunc(r.Handler)
	}
	//
	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	//thread

	threadSubRouter := router.PathPrefix("/api/thread").Subrouter()

	_thread := []Route{
		{
			Info:    "Handler for creating new post.",
			Name:    "thread_CreatNewPost",
			Path:    "/{slug_or_id}/create",
			Method:  http.MethodPost,
			Handler: controllers.CreatNewPostHandler,
		},
		{
			Info:    "Handler for updating the branch.",
			Name:    "thread_UpdateBranch",
			Path:    "/{slug_or_id}/details",
			Method:  http.MethodPost,
			Handler: controllers.UpdateBranchHandler,
		},
		{
			Info:    "Handler for voting the discussion thread.",
			Name:    "thread_VoteThread",
			Path:    "/{slug_or_id}/vote",
			Method:  http.MethodPost,
			Handler: controllers.VoteThreadHandler,
		},
		{
			Info:    "Handler for getting information about the discussion thread.",
			Name:    "thread_GetThreadInfoThread",
			Path:    "/{slug_or_id}/details",
			Method:  http.MethodGet,
			Handler: controllers.GetThreadInfoThreadHandler,
		},
		{
			Info:    "Handler for getting messages of this branch of the discussion.",
			Name:    "thread_GetBranchMessages",
			Path:    "/{slug_or_id}/posts",
			Method:  http.MethodGet,
			Handler: controllers.GetBranchMessagesHandler,
		},
	}

	for _, r := range _thread {
		threadSubRouter.
			Methods(r.Method).
			Name(r.Name).
			Path(r.Path).
			HandlerFunc(r.Handler)
	}

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	//user

	userSubRouter := router.PathPrefix("/api/user").Subrouter()

	_user := []Route{
		{
			Info:    "Handler for creating new user.",
			Name:    "user_CreatNewUser",
			Path:    "/{nickname}/create",
			Method:  http.MethodPost,
			Handler: controllers.CreatNewUserHandler,
		},
		{
			Info:    "Handler for changing user data.",
			Name:    "user_ChangUserData",
			Path:    "/{nickname}/profile",
			Method:  http.MethodPost,
			Handler: controllers.ChangeUserDataHandler,
		},
		{
			Info:    "Handler for getting information about user.",
			Name:    "user_GetUserInfo",
			Path:    "/{nickname}/profile",
			Method:  http.MethodGet,
			Handler: controllers.GetUserInfoHandler,
		},
	}

	for _, r := range _user {
		userSubRouter.
			Methods(r.Method).
			Name(r.Name).
			Path(r.Path).
			HandlerFunc(r.Handler)
	}
	router.Use(middleware.ContentType)
	router.Use(middleware.Logger)
	router.Use(middleware.Panic)
	return
}
