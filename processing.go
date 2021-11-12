package csgo

import (
	"github.com/m2q/siam-cs/model"
	"strconv"
	"time"
)

// PastMatchesTTL is the minimum duration that a past match should live on the blockchain.
// The duration is measured starting from the model.Match Date field. A past match that is
// published on the chain will remain on it for this time.
const PastMatchesTTL = time.Hour * 72 // 3 days

// CreateWinnerMap converts a slice []Match into a map, where the key is the match ID,
// and the value is the Match winner. If there is no winner yet, the value will be empty.
func CreateWinnerMap(m []model.Match) map[string]string {
	result := make(map[string]string, 0)
	for _, v := range m {
		result[strconv.Itoa(v.ID)] = v.Result.Winner
	}
	return result
}

// ConstructDesiredState returns the desired state of a buffer of size l that the Oracle wishes
// to integrate onto the AlgorandBuffer. There are two factors for desirability.
//
//  1) Past matches should remain on the buffer until their Date is older than the PastMatchesTTL
//  2) Remove past matches that are older than their TTL.
func ConstructDesiredState(past []model.Match, future []model.Match, l int) []model.Match {
	// cut off TTL
	pastTTL, desired := SplitMatchesAge(past, PastMatchesTTL)
	// append future matches
	desired = append(desired, future...)
	// truncate if necessary
	if len(desired) > l {
		desired = desired[:l]
	}
	// if the length of desired buffer is STILL not maxed, it means that
	// there are not enough future matches going on. In this case we can
	// fill the rest with old data.
	if len(desired) < l {
		if d := l - len(desired); d <= len(pastTTL) {
			desired = append(pastTTL[len(pastTTL)-d:], desired...)
		}
	}
	return desired
}

// ReverseMatches reverses the order of a match array.
func ReverseMatches(s []model.Match) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
