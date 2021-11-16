/*
The following file has been adapted from https://github.com/Olament/HLTV-Go
and modified to serve this projects purpose.

MIT License

Copyright (c) 2019 Zixuan Guo

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package csgo

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/m2q/siam-cs/model"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type HLTV struct {
	UpcomingPage *goquery.Document
	ResultsPage  *goquery.Document
}

// Fetch gets the latest version of the HLTV page.
// Note: Do not abuse this function. Exceeding certain rates can be interpreted as
// crawling and result in IP ban.
func (h *HLTV) Fetch() (past, future []model.Match, err error) {
	h.UpcomingPage, err = getDocument("https://www.hltv.org/matches?predefinedFilter=top_tier")
	if err != nil {
		return nil, nil, err
	}
	h.ResultsPage, err = getDocument("https://www.hltv.org/results?stars=1")
	if err != nil {
		return nil, nil, err
	}

	past, err = h.getPastMatches()
	if err != nil {
		return nil, nil, err
	}
	future, err = h.getFutureMatches()
	if err != nil {
		return nil, nil, err
	}
	return past, future, nil
}

func PopSlashSource(selection *goquery.Selection) string {
	res, _ := selection.Attr("src")
	split := strings.Split(res, "/")
	return split[len(split)-1]
}

func (h *HLTV) getPastMatches() ([]model.Match, error) {
	doc := h.ResultsPage

	matches := make([]model.Match, 0, 100)

	doc.Find(".result-con").Each(func(i int, sel *goquery.Selection) {
		selection := sel.Find(".result").First()

		tmp, _ := selection.Parent().Attr("href")
		matchID, _ := strconv.Atoi(strings.Split(strings.TrimPrefix(tmp, "/matches/"), "/")[0])
		event := selection.Find(".event-name").First().Text()
		timeRaw, exists := sel.Attr("data-zonedgrouping-entry-unix")
		date := time.Time{}
		if exists {
			matchTime, _ := strconv.ParseInt(timeRaw[:len(timeRaw)-3], 10, 64)
			date = time.Unix(matchTime, 0)
		} else {
			return
		}
		team1 := selection.Find(".team1").First().Find(".team").First().Text()
		team2 := selection.Find(".team2").First().Find(".team").First().Text()

		winner := selection.Find(".team-won").First().Text()

		scoreWon := selection.Find(".score-won").First().Text()
		scoreLost := selection.Find(".score-lost").First().Text()

		match := model.Match{
			ID: matchID,
			Team1: model.Team{
				Name: team1,
			},
			Team2: model.Team{
				Name: team2,
			},
			Event: model.Event{
				Name: event,
			},
			Date: date,
			Result: model.Result{
				Winner: winner,
				Score:  scoreWon + "-" + scoreLost,
			},
		}

		matches = append(matches, match)
	})
	// past matches are in reverse order on hltv page
	ReverseMatches(matches)
	return matches, nil
}

func (h *HLTV) getFutureMatches() ([]model.Match, error) {
	doc := h.UpcomingPage
	// Get top tier matches
	matches := getMatchesFromMatchesPage(doc, ".liveMatch")
	// Set all live matches to Live=true
	for i, _ := range matches {
		matches[i].Live = true
		matches[i].Date = time.Now()
	}
	matches = append(matches, getMatchesFromMatchesPage(doc, ".upcomingMatch")...)
	return matches, nil
}

// getMatchesFromMatchesPage returns a []Match slice containing matches parsed from the upcoming matches
// page. There are two categories as of now, live and upcoming, specified in the matchType string
func getMatchesFromMatchesPage(doc *goquery.Document, matchType string) []model.Match {
	matches := make([]model.Match, 0)
	doc.Find(matchType).Each(func(i int, selection *goquery.Selection) {
		matchHref, _ := selection.Find("a.match").First().Attr("href")
		matchID, _ := strconv.Atoi(strings.Split(matchHref, "/")[2])
		timeRaw, exists := selection.Find(".matchTime").First().Attr("data-unix")
		date := time.Time{}
		if exists {
			matchTime, _ := strconv.ParseInt(timeRaw[:len(timeRaw)-3], 10, 64)
			date = time.Unix(matchTime, 0)
		}

		event := selection.Find(".matchEventName").First().Text()
		eventID, _ := strconv.Atoi(
			strings.Split(PopSlashSource(selection.Find("img.matchEventLogo")), ".")[0])
		eventLogo, _ := selection.Find(".matchEventLogo").First().Attr("src")

		format := selection.Find(".matchMeta").First().Text()

		team1 := selection.Find(".matchTeamName").First().Text()
		team1IDStr, _ := selection.Attr("team1")
		team1ID, _ := strconv.Atoi(team1IDStr)

		team2 := selection.Find(".matchTeamName").Last().Text()
		team2IDStr, _ := selection.Attr("team2")
		team2ID, _ := strconv.Atoi(team2IDStr)

		match := model.Match{
			ID: matchID,
			Team1: model.Team{
				Name: team1,
				ID:   team1ID,
			},
			Team2: model.Team{
				Name: team2,
				ID:   team2ID,
			},
			Date:   date,
			Format: format,
			Event: model.Event{
				Name:    event,
				ID:      eventID,
				LogoURL: eventLogo,
			},
			Result: model.Result{},
		}

		matches = append(matches, match)
	})
	return matches
}

// getDocument performs a GET-Query to the given URL, and creates a goquery-Document from its response.
func getDocument(url string) (*goquery.Document, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) "+
		"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.71 Mobile Safari/537.36")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, errors.New(res.Status)
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}
