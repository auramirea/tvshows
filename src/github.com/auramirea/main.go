package main

import (
	"net/http"
	"github.com/ant0ine/go-json-rest/rest"
	"fmt"
	"strconv"
	p "github.com/auramirea/persistence"
	"github.com/auramirea/service"
)

var c = service.NewClient(nil)
var tvs = service.NewTvService(c)
var db = p.GetUserRepository()

type Methods struct {}

func (*Methods) ListTvShows(w rest.ResponseWriter, r *rest.Request) {
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

func (*Methods) Search(w rest.ResponseWriter, r *rest.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		rest.Error(w, "'q' query parameter required", 400)
		return
	}
	params := &service.SearchParams{query}
	result := tvs.Search(params)
	w.WriteJson(result)
}

func (*Methods) CreateUser(w rest.ResponseWriter, r *rest.Request) {
	user := p.UserEntity{}
	err := r.DecodeJsonPayload(&user)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(db.CreateUser(user))
}

func (*Methods) DeleteUser(w rest.ResponseWriter, r *rest.Request) {
	userId := r.PathParam("userId")
	if userId == "" {
		rest.Error(w, "'userId' cannot be empty", http.StatusBadRequest)
		return
	}
	db.DeleteUser(userId)
	w.WriteHeader(http.StatusNoContent)
}

func (*Methods) GetAllUsers(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson(db.FindAllUsers())
}

func (*Methods) GetUser(w rest.ResponseWriter, r *rest.Request) {
	userId := r.PathParam("userId")
	if userId == "" {
		rest.Error(w, "'userId' cannot be empty", http.StatusBadRequest)
		return
	}
	user := db.FindUser(userId)
	if user == nil {
		rest.Error(w, "User not found", http.StatusNotFound)
		return
	}
	w.WriteJson(user)
}

func (*Methods) AddShow(w rest.ResponseWriter, r *rest.Request) {
	userId := r.PathParam("userId")
	showId := r.PathParam("showId")
	user := db.AddShow(userId, tvs.GetShow(showId))
	w.WriteJson(user)
}

func (*Methods) DeleteShow(w rest.ResponseWriter, r *rest.Request) {
	userId := r.PathParam("userId")
	showId := r.PathParam("showId")
	user := db.DeleteShow(userId, showId)
	w.WriteJson(user)
}

func main() {
	//dbMigration := p.DbMigration{}
	//dbMigration.MigrationsDown()
	//dbMigration.MigrationsUp()
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
		rest.Put("/users/:userId/shows/:showId", methods.AddShow),
		rest.Delete("/users/:userId/shows/:showId", methods.DeleteShow),
	)
	if err != nil {
		fmt.Println(err)
	}
	api.SetApp(router)
	http.ListenAndServe(":8080", api.MakeHandler())
}


