package main

import (
	"fmt"
	"github.com/m2q/siam-cs"
)

func main() {

	hltv := csgo.HLTV{}

	err := hltv.Fetch()
	if err != nil {
		fmt.Errorf(err.Error())
	}
	matches, _ := hltv.GetFutureMatches()
	_, _ = hltv.GetPastMatches()

	fmt.Println(matches)
}
