package csgo

import (
	siam "github.com/m2q/algo-siam"
)

// Oracle aggregates and stores all HLTV data
type Oracle struct {
	PastMatches   []Match
	FutureMatches []Match
	API           API
	buffer        *siam.AlgorandBuffer
}

// NewOracle creates and initializes an Oracle struct. Requires an API to fetch and
// collect data, as well as a siam.AlgorandBuffer, in order to publish changes to the
// blockchain.
func NewOracle(b *siam.AlgorandBuffer, a API) *Oracle {
	return &Oracle{API: a, buffer: b}
}

func (o *Oracle) Serve() error {
	return nil
}
