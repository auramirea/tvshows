package main

import (
	"net/http"
	s "github.com/auramirea/service"
	"github.com/ant0ine/go-json-rest/rest"
	"fmt"
	"strconv"
	"github.com/auramirea/persistence"
)

var c = s.NewClient(nil)
var tvs = s.NewTvService(c)

type Methods struct {}

func (api *Methods) ListTvShows(w rest.ResponseWriter, r *rest.Request) {
	showId := r.URL.Query().Get("showId")
	if showId == ""{
		rest.Error(w, "'showId' query parameter required", 400)
		return
	}
	if _, err := strconv.Atoi(showId); err != nil {
		rest.Error(w, "'showId' must be an int", 400)
		return
	}
	episodes := tvs.ListEpisodes(showId)
	w.WriteJson(episodes)
}

func (api *Methods) Search(w rest.ResponseWriter, r *rest.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		rest.Error(w, "'q' query parameter required", 400)
		return
	}
	params := &s.SearchParams{query}
	result := tvs.Search(params)
	w.WriteJson(result)
}

func (api *Methods) CreateUser(w rest.ResponseWriter, r *rest.Request) {
	user := s.User{}
	err := r.DecodeJsonPayload(&user)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(s.GetUserService().CreateUser(user))
}

func (api *Methods) DeleteUser(w rest.ResponseWriter, r *rest.Request) {
	userId := r.PathParam("userId")
	if userId == "" {
		rest.Error(w, "'userId' cannot be empty", http.StatusBadRequest)
		return
	}
	s.GetUserService().DeleteUser(userId)
	w.WriteHeader(http.StatusNoContent)
}

func (api *Methods) GetAllUsers(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson(s.GetUserService().FindAllUsers())
}

func (api *Methods) GetUser(w rest.ResponseWriter, r *rest.Request) {
	userId := r.PathParam("userId")
	if userId == "" {
		rest.Error(w, "'userId' cannot be empty", http.StatusBadRequest)
		return
	}
	w.WriteJson(s.GetUserService().FindUser(userId))
}


func main() {
	dbMigration := persistence.DbMigration{}
	//dbMigration.MigrationsDown()
	dbMigration.MigrationsUp()
	methods := Methods{}
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/list", methods.ListTvShows),
		rest.Get("/search", methods.Search),
		rest.Post("/users", methods.CreateUser),
		rest.Get("/users", methods.GetAllUsers),
		rest.Get("/users/:userId", methods.GetUser),
		rest.Delete("/users/:userId", methods.DeleteUser),
	)
	if err != nil {
		fmt.Println(err)
	}
	api.SetApp(router)
	http.ListenAndServe(":8080", api.MakeHandler())
}


