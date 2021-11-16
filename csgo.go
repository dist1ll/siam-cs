package csgo

import (
	"encoding/json"
	"errors"
	"github.com/m2q/siam-cs/model"
	"io/ioutil"
	"log"
)

// API is a provider of CSGO match data. It provides methods for fetching past and
// future matches. Which matches are selected is up to the API implementation.
type API interface {
	// GetPastMatches returns a list of past CSGO pro matches.
	// The reference implementation is provided by HLTV.
	GetPastMatches() ([]model.Match, error)
	// GetFutureMatches returns a list of upcoming CSGO pro matches.
	// The reference implementation is provided by HLTV.
	GetFutureMatches() ([]model.Match, error)
}

// StubAPI is a stub that implements API. You can explicitly set the match data
// that shall be returned by the API functions, by modifying the public fields or
// calling SetMatches. You can also specify if the stub should return an error.
type StubAPI struct {
	Past      []model.Match
	Future    []model.Match
	Logger    *log.Logger
	LogActive bool
}

// SetMatches sets future and past matches to be returned by the API
func (s *StubAPI) SetMatches(past []model.Match, future []model.Match) {
	s.Past = past
	s.Future = future
}

// GetPastMatches returns a static list of past matches, that can be set via SetMatches
func (s *StubAPI) GetPastMatches() ([]model.Match, error) {
	if s.LogActive {
		s.Logger.Println("API request past matches")
	}
	if s.Past == nil {
		return nil, errors.New("API stub has not been assigned")
	}
	return s.Past, nil
}

// GetFutureMatches returns a static list of future matches, that can be set via SetMatches
func (s *StubAPI) GetFutureMatches() ([]model.Match, error) {
	if s.LogActive {
		s.Logger.Println("API request future matches")
	}
	if s.Future == nil {
		return nil, errors.New("API stub has not been assigned")
	}
	return s.Future, nil
}

func CreateData() {
	hltv := &HLTV{}
	hltv.Fetch()
	p, _ := hltv.GetPastMatches()
	f, _ := hltv.GetFutureMatches()
	p = append(p, f...)
	file, _ := json.MarshalIndent(p, "", "\t")
	_ = ioutil.WriteFile("./generator/reference_data_x.json", file, 0644)
}
