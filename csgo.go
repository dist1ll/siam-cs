package csgo

import "github.com/m2q/siam-cs/model"

// ReverseMatches reverses the order of a match array.
// TODO: Replace this with generics when 1.18 is released.
func ReverseMatches(s []model.Match) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

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
}

// SetMatches sets future and past matches according to the given []Match slice.
// All matches with index < i will be considered past matches, and every other match
// will be considered a future match.
func (s *StubAPI) SetMatches(m []model.Match, i int) {

}
func (s *StubAPI) GetPastMatches() ([]model.Match, error) {
	return nil, nil
}

func (s *StubAPI) GetFutureMatches() ([]model.Match, error) {
	return nil, nil
}
