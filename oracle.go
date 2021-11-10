package csgo

import (
	"context"
	siam "github.com/m2q/algo-siam"
	"sync"
	"time"
)

// Oracle aggregates and stores all HLTV data
type Oracle struct {
	pastMatches   []Match
	futureMatches []Match
	cfg           *OracleConfig
	buffer        *siam.AlgorandBuffer
}

// OracleConfig defines the oracles behavior
type OracleConfig struct {
	// PrimaryAPI is the primary source of information. Only data proposed from this API
	// is put onto the blockchain.
	PrimaryAPI API

	// VerificationAPIs is an optional list of verification APIs used to verify the correctness
	// of the PrimaryAPI. If the data proposed from the PrimaryAPI does not match data from each
	// of verification APIs, the data is not written to the blockchain.
	VerificationAPIs []API

	// MaxVerifyTime does not need to be set if VerificationAPIs is nil. If the VerificationAPIs
	// do not agree with the PrimaryAPI's proposal within MaxVerifyTime, then the proposed data
	// gets discarded
	MaxVerifyTime time.Duration

	// RefreshInterval is the pause between two API fetch commands.
	// If the API access a rate-limited location, then set RefreshInterval low enough
	// as to not trigger a rate-limit or blacklist event.
	RefreshInterval time.Duration
}

// NewOracle creates and initializes an Oracle struct. Requires an API to fetch and
// collect data, as well as a siam.AlgorandBuffer, in order to publish changes to the
// blockchain.
func NewOracle(b *siam.AlgorandBuffer, cfg *OracleConfig) *Oracle {
	return &Oracle{cfg: cfg, buffer: b}
}

// Serve spawns a cancelable goroutine that manages continuously fetches data, and publishes
// it on the blockchain. The goroutine can be cancelled anytime via the returned context.CancelFunc.
// The returned WaitGroup will be done if the Oracle finishes execution (can be used to make sure
// no goroutine is being leaked).
func (o *Oracle) Serve() (*sync.WaitGroup, context.CancelFunc) {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	wg.Add(1)
	go func() {
		defer wg.Done()
		o.serve(ctx)
	}()
	return &wg, cancel
}

// serve contains the actual implementation of the Oracle behavior.
func (o *Oracle) serve(ctx context.Context) {
}
