// Package generator contains functions for generating []csgo.Match data
// for testing purposes. This data is based on reference data gathered
// from HLTV.org.
package generator

import (
	_ "embed"
	"encoding/json"
	"github.com/m2q/siam-cs/model"
)

// refRaw contains real reference match data as a raw json string
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

func GetData() []model.Match {
	return nil
}
