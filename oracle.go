package csgo

import (
	"context"
	"fmt"
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
	current, err := o.buffer.GetBuffer(ctx)
	if err != nil {
		return o.cfg.RefreshInterval
	}
	put, del := computeOverlap(desired, current)

	if len(put)+len(del) == 0 {
		return o.cfg.RefreshInterval
	}

	err = o.buffer.DeleteElements(context.Background(), getKeys(del)...)
	if err != nil {
		log.Print(err)
		return o.cfg.RefreshInterval
	}
	err = o.buffer.PutElements(context.Background(), put)
	if err != nil {
		log.Print(err)
		return o.cfg.RefreshInterval
	}

	return o.cfg.RefreshInterval
}

// waitForFlush waits until the AlgorandBuffer has reached the desired state.
func (o *Oracle) waitForFlush(ctx context.Context, desired map[string]string) {
	for ctx.Err() != nil {
		d, err := o.buffer.GetBuffer(ctx)
		if err != nil && fmt.Sprint(d) == fmt.Sprint(desired) {
			break
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Second):
			continue
		}
	}
}

// computeOverlap returns two maps, m1 and m2. m1 contains the map entries of x, for
// which the keys either don't exist in y, or do exist but with different values than
// in x. m2 contains map entries of y that don't exist in x.
func computeOverlap(x, y map[string]string) (m1, m2 map[string]string) {
	m1 = make(map[string]string)
	m2 = make(map[string]string)

	for k, v := range x {
		yv, ok := y[k]
		// if the key exists in both x and y, and they have the same value, exclude it.
		if !(ok && v == yv) {
			m1[k] = v
		}
	}
	for k, v := range y {
		// keys that exist in y, but not in x
		if _, ok := x[k]; !ok {
			m2[k] = v
		}
	}
	return m1, m2
}

func getKeys(m map[string]string) []string {
	s := make([]string, len(m))
	i := 0
	for k, _ := range m {
		s[i] = k
		i++
	}
	return s
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
