package csgo

import (
	siam "github.com/m2q/algo-siam"
	"github.com/m2q/algo-siam/client"
	"testing"
)

// setupOracleMockedAPI creates an Oracle instance, initialized with a
// AlgorandBuffer (client.AlgorandMock) and StubAPI.
func setupOracleMockedAPI() (*Oracle, *siam.AlgorandBuffer, *StubAPI) {
	c := client.CreateAlgorandClientMock("", "")
	buffer, _ := siam.CreateAlgorandBuffer(c, "")
	api := &StubAPI{}
	cfg := &OracleConfig{
		PrimaryAPI:      nil,
		RefreshInterval: 0,
	}
	o := NewOracle(buffer, cfg)
	return o, buffer, api
}

func getTestingData(path string) []Match {
	return nil
}

// Tests if the Oracle writes the data provided by the API stub to the AlgorandBuffer
func TestOracle_SimpleSmallData(t *testing.T) {
	oracle, _, stub := setupOracleMockedAPI()
	wg, cancel := oracle.Serve()

	// Set match data to stub
	matchData := getTestingData("data_1")
	stub.SetMatches(matchData, 5)

	// Check if data is being written to the AlgorandBuffer
	// contains, _ := buffer.ContainsWithin(MatchesToMap(matchData), time.Second)
	// assert.True(t, contains)

	// Cancel goroutine
	cancel()
	wg.Wait()
}
