package main

import (
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
	"net/http"
	s "github.com/auramirea/service"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/auramirea/persistence"
)
type listEpisodesRequest struct {
	S string `json:"showId"`
}

type listEpisodesResponse struct {
	V   []s.Episode `json:"episodes"`
}

type searchRequest struct {
	S string `json:"query"`
}
type searchResponse struct {
	R s.SearchResult `json:"result"`
}

func makeListEndpoint(tvs *s.TvService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listEpisodesRequest)
		episodes := tvs.ListEpisodes(req.S)

		return listEpisodesResponse{V: episodes}, nil
	}
}
func makeSearchEndpoint(tvs *s.TvService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(searchRequest)
		params := &s.SearchParams{Query: req.S}
		result := tvs.Search(params)
		return searchResponse{R: result}, nil
	}
}
func makeUserEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(s.User)
		return s.GetUserService().CreateUser(req), nil
	}
}

func main() {
	c := s.NewClient(nil)
	d := &persistence.DbMigration{}
	d.MigrationsDown()
	d.MigrationsUp()

	tvs := s.NewTvService(c)
	ctx := context.Background()

	listHandler := httptransport.NewServer(ctx, makeListEndpoint(tvs), decodeListRequest, encodeResponse)
	searchHandler := httptransport.NewServer(ctx, makeSearchEndpoint(tvs), decodeSearchRequest, encodeResponse)
	userHandler := httptransport.NewServer(ctx, makeUserEndpoint(), decodeUserRequest, encodeResponse)

	http.Handle("/list", listHandler)
	http.Handle("/search", searchHandler)
	http.Handle("/users", userHandler)
	http.ListenAndServe(":8080", nil)
}

func decodeUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request s.User
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}
func decodeListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return listEpisodesRequest{S: string(r.URL.Query().Get("showId"))}, nil
}
func decodeSearchRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return searchRequest{S: string(r.URL.Query().Get("query"))}, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}