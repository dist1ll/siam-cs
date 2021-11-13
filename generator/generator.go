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
