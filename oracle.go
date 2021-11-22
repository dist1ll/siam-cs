package csgo

import (
	"context"
	"log"
	"sync"
	"time"

	siam "github.com/m2q/algo-siam"
	"github.com/m2q/algo-siam/client"
	"github.com/m2q/siam-cs/model"
)

// Oracle aggregates and stores all HLTV data
type Oracle struct {
	pastMatches   []model.Match
	futureMatches []model.Match
	cfg           *OracleConfig
	buffer        *siam.AlgorandBuffer
	cancelOracle  context.CancelFunc
	wgExit        *sync.WaitGroup
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
	// If the API accesses a rate-limited resource, then set RefreshInterval high enough
	// as to not trigger a rate-limit or blacklist event.
	RefreshInterval time.Duration
}

// NewOracle creates and initializes an Oracle struct. Requires an API to fetch and
// collect data, as well as a siam.AlgorandBuffer, in order to publish changes to the
// blockchain.
func NewOracle(b *siam.AlgorandBuffer, cfg *OracleConfig) *Oracle {
	return &Oracle{cfg: cfg, buffer: b}
}

// Serve spawns a cancelable goroutine that aims to keep the AlgorandBuffer
// in a desired state. See ConstructDesiredState.
//
// Any goroutines spawned by the Oracle can be cancelled anytime via Stop.
func (o *Oracle) Serve() {
	// Start ManagingRoutine for AlgorandBuffer
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	wg.Add(1)
	go func() {
		defer wg.Done()

		// Serving Loop
		for ctx.Err() == nil {
			sleep := o.serve(ctx)
			if sleep != 0 {
				select {
				case <-ctx.Done():
					return
				case <-time.After(sleep):
					continue
				}
			}
		}
	}()

	o.cancelOracle = cancel
	o.wgExit = &wg
}

// serve attempts to bring the AlgorandBuffer in a desired state. It returns
// a minimum time that the caller should wait before executing serve again.
func (o *Oracle) serve(ctx context.Context) time.Duration {
	// Update local list of matches
	p, f, err := o.cfg.PrimaryAPI.Fetch()
	if err != nil {
		return o.cfg.RefreshInterval
	}
	o.pastMatches = p
	o.futureMatches = f

	desired := ConstructDesiredState(o.pastMatches, o.futureMatches, client.GlobalBytes)
	err = o.buffer.AchieveDesiredState(ctx, desired)
	if err != nil {
		log.Print(err)
	}
	return o.cfg.RefreshInterval
}

// Stop signals the Oracle to stop its goroutine and stop the siam.AlgorandBuffer
// managing routine. Stop will block until both goroutines have exited.
func (o *Oracle) Stop() {
	if o.cancelOracle != nil {
		o.cancelOracle()
		o.wgExit.Wait()
	} else {
		panic("need to run .Serve() before stopping oracle")
	}
}

// Wait blocks until all goroutines of the Oracle have finished execution.
// To stop execution of an Oracle, call Stop.
func (o *Oracle) Wait() {
	o.wgExit.Wait()
}
