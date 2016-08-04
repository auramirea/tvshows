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

func (*Methods) GetAllShows(w rest.ResponseWriter, r *rest.Request) {
	genre := r.URL.Query().Get("genre")
	alphabet := r.URL.Query().Get("alphabet")
	result := tvs.GetAllShows()
	if (genre != "") {
		result = tvs.FilterByGenre(genre, result)
	}
	if (alphabet != "") {
		result = tvs.FilterByAlphabet(alphabet, result)
	}
	w.WriteJson(result)
}

func (*Methods) GetShow(w rest.ResponseWriter, r *rest.Request) {
	showId := r.PathParam("showId")
	result := tvs.GetShow(showId)
	w.WriteJson(result)
}

func main() {
	//dbMigration := p.DbMigration{}
	//dbMigration.MigrationsDown()
	//dbMigration.MigrationsUp()
	methods := Methods{}
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	api.Use(&rest.CorsMiddleware{
		RejectNonCorsRequests: false,
		OriginValidator: func(origin string, request *rest.Request) bool {
			return origin == "http://localhost:8000"
		},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{
			"Accept", "Content-Type", "X-Custom-Header", "Origin", "Accept-Language",
			"Accept-Encoding", "X-Requested-With"},
		AccessControlAllowCredentials: true,
		AccessControlMaxAge:           3600,
	})
	router, err := rest.MakeRouter(
		rest.Get("/shows", methods.GetAllShows),
		rest.Get("/shows/:showId", methods.GetShow),
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


