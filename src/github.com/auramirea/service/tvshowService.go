package service

import (
	"github.com/dghubble/sling"
	"strings"
	"fmt"
	"github.com/auramirea/client"
	"sync"
	"encoding/json"
)

const baseURL = "http://api.tvmaze.com/"

var instance *TvService
var once sync.Once

const showsKey = "showCache"

type TvService struct {
	sling *sling.Sling
	redisClient *client.RedisClient
}

func NewTvService() *TvService {
	once.Do(func() {
		instance = &TvService{
			sling: sling.New().Client(nil).Base(baseURL),
			redisClient: client.GetRedisClient(),
		}
	})
	return instance
}

type Episode struct {
	Id       int    `json:"id"`
	Url      string `json:"url"`
	Name     string `json:"name"`
	Season   int `json:"season"`
	Number   int `json:"number"`
	Airdate	string `json:"airdate"`
	Runtime int `json:"runtime"`
	Summary string `json:"summary"`
}

type SearchParams struct {
	Query string `url:"q,omitempty"`
}
type Show struct {
	Id       int    `json:"id"`
	Url      string `json:"url"`
	Status   string `json:"status"`
	Rating   `json:"rating"`
	Name     string `json:"name"`
	Language string `json:"language"`
	Summary  string  `json:"summary"`
	Schedule `json:"schedule"`
	Image    `json:"image"`
	_embedded `json:"_embedded"`
	Network `json:"network"`
	Genres []string `json:"genres"`

}
type Network struct {
	Name string `json:"name"`
}
type Schedule struct {
	Time string `json:"time"`
	Days []string `json:"days"`
}
type _embedded struct {
	Episodes []Episode `json:"episodes"`
}
type Rating struct {
	Average float32 `json:"average"`
}
type Image struct {
	Medium   string `json:"medium"`
	Original string `json:"original"`
}

func (s *TvService) ListEpisodes(showId string) ([]Episode) {
	result := new([]Episode)
	s.sling.New().Get("/shows/" + showId + "/episodes").ReceiveSuccess(result)

	return *result
}
// List returns the authenticated user's issues across repos and orgs.
func (s *TvService) Search(params *SearchParams) (Show) {
	var result Show
	s.sling.New().Get("/singlesearch/shows").QueryStruct(params).ReceiveSuccess(&result)
	return result
}

func (s *TvService) GetShow(showId string) *Show {
	result := Show{}
	s.sling.New().Get("/shows/" + showId + "?embed=episodes").ReceiveSuccess(&result)
	return &result
}

func (s *TvService) GetAllShows(page string) []Show {
	result := []Show{}
	path := "/shows"
	if page != "" {
		path += "?page=" + page
	}
	cachedShows := s.redisClient.GetKey(showsKey + page)
	if cachedShows != "" {
		json.Unmarshal([]byte(cachedShows), &result)
		fmt.Println("Unmarshalled cached result")
		return result
	}
	s.sling.New().Get(path).ReceiveSuccess(&result)
	val, _ := json.Marshal(result)
	s.redisClient.SetKey(showsKey, string(val))
	return result
}

func (s *TvService) FilterByGenre(filter string, shows []Show) []Show {
	filteredResult := []Show{}
	for _, show := range(shows) {
		fmt.Println(show)
		for _, genre := range(show.Genres) {
			fmt.Println(genre)
			if strings.Compare(strings.ToLower(genre), strings.ToLower(filter)) == 0 {
				filteredResult = append(filteredResult, show)
			}
		}
	}
	return filteredResult
}

func (s *TvService) FilterByAlphabet(alphabet string, shows []Show) []Show {
	filteredResult := []Show{}
	for _, show := range(shows) {
		if strings.HasPrefix(strings.ToLower(show.Name), strings.ToLower(alphabet)) {
			filteredResult = append(filteredResult, show)
		}
	}
	return filteredResult
}


