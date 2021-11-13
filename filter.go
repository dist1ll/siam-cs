package csgo

import (
	"github.com/m2q/siam-cs/model"
	"time"
)

// SplitMatchesAge returns a partition of matches. The first return value contains matches that
// concluded longer than threshold ago, whereas the second return value contains matches that are
// at most threshold old. The distance is measured from a given present time `now`.
func SplitMatchesAge(m []model.Match, threshold time.Duration, now time.Time) ([]model.Match, []model.Match) {
	for i, v := range m {
		if now.Sub(v.Date) <= threshold {
			return m[:i], m[i:]
		}
	}
	return m, []model.Match{}
}
