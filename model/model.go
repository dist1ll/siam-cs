package model

import "time"

type Team struct {
	Name string
	// HLTV ID for team
	ID int
}

type Event struct {
	// E.g. "IEM Fall 2021 Europe"
	Name string
	// HLTV ID for event
	ID      int
	LogoURL string
}

type Match struct {
	ID     int
	Team1  Team
	Team2  Team
	Date   time.Time
	Event  Event
	Format string
	Result Result
	Live   bool
}

type Result struct {
	// Winning team's name e.g. "OG", "Astralis"
	Winner string
	// Numbered score (e.g. "1-0", "3-2"). Winner's score is always listed first.
	Score string
}
