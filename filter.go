package csgo

import (
	"github.com/m2q/siam-cs/model"
	"time"
)

// SplitMatchesAge returns a partition of matches. The first return value contains matches that
// concluded longer than threshold ago, whereas the second return value contains matches that are
// at most threshold old. The distance is measured using time.Now.
func SplitMatchesAge(m []model.Match, threshold time.Duration) ([]model.Match, []model.Match) {
	now := time.Now()
	for i, v := range m {
		if now.Sub(v.Date) <= threshold {
			return m[:i], m[i:]
		}
	}
	return m, []model.Match{}
}
