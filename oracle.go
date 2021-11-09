package csgo

// Oracle aggregates and stores all HLTV data
type Oracle struct {
	PastMatches   []Match
	FutureMatches []Match
}

func (o *Oracle) Serve() error {
	return nil
}
