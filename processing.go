package csgo

import "github.com/m2q/siam-cs/model"

// CreateWinnerMap converts a slice []Match into a map, where the key is the match ID,
// and the value is the Match winner. If there is no winner yet, the value will be empty.
func CreateWinnerMap(m []model.Match) map[string]string {
	return nil
}

// ConstructDesiredState returns the desired state of the AlgorandBuffer that the Oracle wishes
// to achieve. There are two factors for desirability.
//
//  1) Past matches should remain on the buffer until their Date is older than the
//     OracleConfig.PastMatchesTTL.
//  2) Remove past matches that are older than their TTL.
func ConstructDesiredState(past []model.Match, future []model.Match, size int) []model.Match {
	return nil
}

// ReverseMatches reverses the order of a match array.
func ReverseMatches(s []model.Match) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
