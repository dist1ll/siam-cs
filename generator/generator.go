// Package generator contains functions for generating []csgo.Match data
// for testing purposes. This data is based on reference data gathered
// from HLTV.org.
package generator

import (
	_ "embed"
	"encoding/json"
	"github.com/m2q/siam-cs/model"
	"time"
)

// refRaw contains real reference match data as a raw json string. The data is sorted, so
// that ALL past matches (i.e. matches that have a non-empty `Winner` field) occur BEFORE
// all future or live matches (i.e. matches where `Winner` is empty).
//go:embed reference_data.json
var refRaw string

// ref is the match data parsed from refRaw.
var ref []model.Match

// init parses the raw match data and initializes ref
func init() {
	if err := json.Unmarshal([]byte(refRaw), &ref); err != nil {
		panic(err)
	}
}

// GetData returns a time-normalized sample of real data. A combination of past and
// future matches. First return value is past, second is future matches.
// Note: The Date fields are normalized w.r.t the reference time. That means that the
// MOST recent past match has a date of refTime, and all other matches are adjusted
// according to this delta
func GetData(refTime time.Time) ([]model.Match, []model.Match) {
	// copy reference data
	result := make([]model.Match, len(ref))
	copy(result, ref)
	// find index of last past match
	lastPast := len(result)
	for i, v := range result {
		if v.Result.Winner == "" {
			lastPast = i
			break
		}
	}
	// normalize time
	diff := refTime.Sub(result[lastPast].Date)
	for i, _ := range result {
		result[i].Date = result[i].Date.Add(diff)
	}
	return result[:lastPast], result[lastPast:]
}

// NormalizeTime shifts the Dates of all matches by the difference between refTime
// and the last past match. Essentially, this means that the time of the data is
// transposed so that the most recent past match happened exactly at refTime.
func NormalizeTime(past []model.Match, future []model.Match, refTime time.Time) {
	// last past match
	diff := refTime.Sub(past[len(past)-1].Date)
	// shift past matches
	for i, _ := range past {
		past[i].Date = past[i].Date.Add(diff)
	}
	// also shift future matches
	for i, _ := range future {
		future[i].Date = future[i].Date.Add(diff)
	}
}

// ProgressTime lets a specified number of future matches conclude, and
// re-normalizes the time. Effectively, this simulates a passing of time.
func ProgressTime(past, future []model.Match, matchCount int) ([]model.Match, []model.Match) {
	// matchCount can't exceed future slice length
	if matchCount > len(future) {
		matchCount = len(future)
	}
	// set result data
	for i := 0; i < matchCount; i++ {
		// just set random stuff
		future[i].Result.Winner = future[i].Team1.Name
		future[i].Result.Score = "16-10"
		// assume that the game took 1 hour, so the Date is set later
		future[i].Date = future[i].Date.Add(time.Hour)
	}
	// re-slice past and future boundaries
	past = append(past, future[:matchCount]...)
	future = future[matchCount:]

	NormalizeTime(past, future, time.Now())
	return past, future
}
