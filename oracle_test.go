package csgo

import (
	"fmt"
	siam "github.com/m2q/algo-siam"
	"github.com/m2q/algo-siam/client"
	"github.com/m2q/siam-cs/generator"
	"github.com/m2q/siam-cs/model"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
	"time"
)

// setupOracleMockedAPI creates an Oracle instance, initialized with a
// AlgorandBuffer (client.AlgorandMock) and StubAPI.
func setupOracleMockedAPI(refresh time.Duration) (*Oracle, *siam.AlgorandBuffer, *StubAPI) {
	c := client.CreateAlgorandClientMock("", "")
	buffer, err := siam.CreateAlgorandBuffer(c, client.GeneratePrivateKey64(), nil)
	if err != nil {
		panic(err)
	}
	api := &StubAPI{
		Logger: log.New(os.Stdout, "ORACLE  ", log.LstdFlags),
	}
	cfg := &OracleConfig{
		PrimaryAPI:      api,
		RefreshInterval: refresh,
		SiamCfg:         &siam.ManageConfig{},
	}
	o := NewOracle(buffer, cfg)
	return o, buffer, api
}

// setupOracleWithData creates an Oracle, Stub and AlgorandBuffer instance, and waits until
// the Oracle has filled specified data into the AlgorandBuffer.
func setupOracleWithData(past, future []model.Match, t *testing.T) (*Oracle, *siam.AlgorandBuffer, *StubAPI) {
	oracle, buffer, stub := setupOracleMockedAPI(0)
	oracle.Serve()
	// Set match data to stub
	stub.SetMatches(past, future)
	// Check if desired state is written by Oracle
	desired := ConstructDesiredState(past, future, client.GlobalBytes)
	contains := buffer.ContainsWithin(desired, time.Second*5, 0)
	assert.True(t, contains)
	return oracle, buffer, stub
}

func containsDesiredState(b *siam.AlgorandBuffer, past []model.Match, future []model.Match, t time.Duration) bool {
	desired := ConstructDesiredState(past, future, client.GlobalBytes)
	contains := b.ContainsWithin(desired, t, 0)
	return contains
}

// Tests if the Oracle writes the data provided by the API stub to the AlgorandBuffer
func TestOracle_SimpleFill(t *testing.T) {
	// Generate Data
	past, future := generator.GetData(time.Now())
	// Create Oracle with initial data
	oracle, _, _ := setupOracleWithData(past, future, t)
	defer oracle.Stop()
}

// Check if AFTER data has been inserted, if its updated correctly when matches have concluded
func TestOracle_ProgressTime(t *testing.T) {
	// Generate Data
	past, future := generator.GetData(time.Now())
	// Create Oracle with initial data
	oracle, b, stub := setupOracleWithData(past, future, t)
	defer oracle.Stop()

	// Let one game play out
	past, future = generator.ProgressTime(past, future, 1)
	stub.SetMatches(past, future)

	assert.True(t, containsDesiredState(b, past, future, time.Second*2))
}

func TestOracle_AddNewMatches(t *testing.T) {
	// Generate Data
	past, future := generator.GetData(time.Now())
	// Create Oracle with initial data
	oracle, b, stub := setupOracleWithData(past, future, t)
	defer oracle.Stop()
	// Let one game play out
	g := generator.GenerateFutureData(future[len(future)-1], 20)
	future = append(future, g...)
	stub.SetMatches(past, future)

	assert.True(t, containsDesiredState(b, past, future, time.Second))
}

func TestOracle_RespectFetchTime(t *testing.T) {
	fmt.Print("Expect 4 log prints: \n")
	// Create Oracle with 0.5s refresh time
	oracle, buffer, stub := setupOracleMockedAPI(time.Millisecond * 50)
	buffer.DisableLogger()

	oracle.Serve()
	defer oracle.Stop()

	stub.LogActive = true

	time.Sleep(time.Millisecond * 190)
}
