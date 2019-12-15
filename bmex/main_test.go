package bmex

import "testing"

func TestReferralEarning(t *testing.T) {
	_, txs := LoadWalletHistory()
	ReferralEarning(txs)
}
