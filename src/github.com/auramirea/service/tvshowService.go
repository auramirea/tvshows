package service

import (
	"net/http"
	"github.com/dghubble/sling"
)

const baseURL = "http://api.tvmaze.com/"

type tvServiceInterface interface {
	ListEpisodes(string) []Episode
	Search(*SearchParams) SearchResult
}


type Episode struct {
	Id       int    `json:"id"`
	Url      string `json:"url"`
	Name     string `json:"name"`
	Season   int `json:"season"`
	Number   int `json:"number"`
	Airdate	string `json:"airdate"`
	Runtime int `json:"runtime"`
}

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

func (s *TvService) ListEpisodes(showId string) ([]Episode) {
	result := new([]Episode)
	s.client.ShowsService.sling.New().Get("/shows/" + showId + "/episodes").ReceiveSuccess(result)

	return *result
}
// List returns the authenticated user's issues across repos and orgs.
func (s *TvService) Search(params *SearchParams) (SearchResult) {
	var result SearchResult
	s.client.ShowsService.sling.New().Get("/singlesearch/shows").QueryStruct(params).ReceiveSuccess(&result)
	return result
}


func NewShowsService(httpClient *http.Client) *ShowsService {
	return &ShowsService{
		sling: sling.New().Client(httpClient).Base(baseURL),
	}
}

type TvService struct {
	client *Client
}

// Client to wrap services
type Client struct {
	ShowsService *ShowsService
}

// NewClient returns a new Client
func NewClient(httpClient *http.Client) *Client {
	return &Client{
		ShowsService: NewShowsService(httpClient),
	}
}
func NewTvService(c *Client) *TvService {
	return &TvService{client: c}
}