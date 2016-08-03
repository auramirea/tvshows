package main

import (
	"fmt"
	"net/http"
	"github.com/dghubble/sling"
	"github.com/auramirea/persistence"
)

const baseURL = "http://api.tvmaze.com/"

type SearchParams struct {
	Query string `url:"q,omitempty"`
}
type SearchResult struct {
	Id       int    `json:"id"`
	Url      string `json:"url"`
	Status   string `json:"status"`
	Rating   `json:"rating"`
	Name     string `json:"name"`
	Language string `json:"language"`
	Image    `json:"image"`
}
type Rating struct {
	Average float32 `json:"average"`
}
type Image struct {
	Medium   string `json:"medium"`
	Original string `json:"original"`
}

// Services
type ShowsService struct {
	sling *sling.Sling
}

func NewShowsService(httpClient *http.Client) *ShowsService {
	return &ShowsService{
		sling: sling.New().Client(httpClient).Base(baseURL),
	}
}

// List returns the authenticated user's issues across repos and orgs.
func (s *ShowsService) Search(params *SearchParams) (SearchResult, *http.Response) {
	var result SearchResult
	resp, _ := s.sling.New().Get("/singlesearch/shows").QueryStruct(params).ReceiveSuccess(&result)
	return result, resp
}

// Client to wrap services

// Client is a tiny Github client
type Client struct {
	ShowsService *ShowsService
}

// NewClient returns a new Client
func NewClient(httpClient *http.Client) *Client {
	return &Client{
		ShowsService: NewShowsService(httpClient),
	}
}

func main() {
	client := NewClient(nil)
	params := &SearchParams{Query: "girls"}
	s, _ := client.ShowsService.Search(params)
	fmt.Println(s)
	d := &persistence.DbMigration{}
	d.RunDBMigrations()
}
