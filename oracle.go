package csgo

import (
	"context"
	siam "github.com/m2q/algo-siam"
	"github.com/m2q/siam-cs/model"
	"sync"
	"time"
)

// Oracle aggregates and stores all HLTV data
type Oracle struct {
	pastMatches   []model.Match
	futureMatches []model.Match
	cfg           *OracleConfig
	buffer        *siam.AlgorandBuffer
	cancelOracle  context.CancelFunc
	cancelBuffer  context.CancelFunc
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

	// SiamCfg is the configuration object for the Siam managing routine. It defines
	// how often the Algorand node is being pinged, refreshed, checked for health,
	// and how long to timeout in case of failure.
	SiamCfg *siam.ManageConfig
}

// NewOracle creates and initializes an Oracle struct. Requires an API to fetch and
// collect data, as well as a siam.AlgorandBuffer, in order to publish changes to the
// blockchain.
func NewOracle(b *siam.AlgorandBuffer, cfg *OracleConfig) *Oracle {
	return &Oracle{cfg: cfg, buffer: b}
}

// Serve spawns a cancelable goroutine that continuously fetches data, and publishes
// it on the blockchain. It also spawns a managing routine for the siam.AlgorandBuffer.
// Both goroutines can be cancelled anytime via Stop, which will signal a cancellation
// via context.CancelFunc and block until both goroutines have finished execution.
func (o *Oracle) Serve() {
	// Start ManagingRoutine for AlgorandBuffer
	wg, c := o.buffer.SpawnManagingRoutine(o.cfg.SiamCfg)
	o.cancelBuffer = c

	ctx, cancel := context.WithCancel(context.Background())
	wg.Add(1)
	go func() {
		defer wg.Done()
		o.serve(ctx)
	}()

	o.cancelOracle = cancel
	o.wgExit = wg
}

// serve contains the actual implementation of the Oracle behavior.
func (o *Oracle) serve(ctx context.Context) {
	for ctx.Err() == nil {
		// Update Matches
		err := o.updateLocalMatches()
		if err != nil {
			select {
			case <-ctx.Done():
				return
			case <-time.After(o.cfg.RefreshInterval):
				continue
			}
		}

		// Generate desired state
		// desired := ConstructDesiredState(o.pastMatches, o.futureMatches, client.GlobalBytes)
	}
}

// updateLocalMatches updates the Oracles lists of matches, by querying the API.
// How often this method is called should depend on OracleConfig.RefreshInterval.
func (o *Oracle) updateLocalMatches() error {
	f, err := o.cfg.PrimaryAPI.GetFutureMatches()
	if err != nil {
		return err
	}
	p, err := o.cfg.PrimaryAPI.GetPastMatches()
	if err != nil {
		return err
	}
	o.futureMatches = f
	o.pastMatches = p
	return nil
}

// Stop signals the Oracle to stop its goroutine and stop the siam.AlgorandBuffer
// managing routine. Stop will block until both goroutines have exited.
func (o *Oracle) Stop() {
	if o.cancelOracle != nil {
		o.cancelOracle()
	}
	if o.cancelBuffer != nil {
		o.cancelBuffer()
	}
	o.wgExit.Wait()
}

// Wait blocks until all goroutines of the Oracle have finished execution.
// To stop execution of an Oracle, call Stop.
func (o *Oracle) Wait() {
	o.wgExit.Wait()
}
