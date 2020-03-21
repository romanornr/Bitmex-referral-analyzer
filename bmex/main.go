package bmex

import (
	"fmt"
	"github.com/btcsuite/btcutil"
	"github.com/romanornr/Bitmex-referral-analyzer/bitcoin"
	"github.com/romanornr/Bitmex-referral-analyzer/client"
	"github.com/spf13/viper"
	"github.com/zmxv/bitmexgo"
	"log"
)

const BITMEXREFLINK = "https://www.bitmex.com/register/vhT2qm"

// loads wallet history by using bitmex api
// returns Transactions which can be used by other functions
func LoadWalletHistory() (error, []bitmexgo.Transaction) {
	auth, apiClient := client.GetInstance()

	// Call APIs without parameters by passing the auth context.
	// e.g. getting exchange-wide turnover and volume statistics:
	tx, _, err := apiClient.UserApi.UserGetWalletHistory(auth, nil)
	if err != nil {
		fmt.Println(err)
	}
	return err, tx
}

func ReferralEarning(transactions []bitmexgo.Transaction) *Stats {

	startYear := viper.GetInt("start_year")

	bitcoinPrice := bitcoin.ToDollar()

	stats := new(Stats)
	var monthBTC btcutil.Amount
	var monthDollar float64
	var totalBTC btcutil.Amount
	var totalDollar float64
	var currentMonth string

	var result *Stat
	result = new(Stat)

	// filter transactions by Type and begin year and put them all bitmexgo.Transaction struct
	var filteredTransactions []bitmexgo.Transaction
	for i := len(transactions) - 1; i > -0; i-- {
		if transactions[i].TransactType == "AffiliatePayout" && transactions[i].Timestamp.Year() >= startYear {
			filteredTransactions = append(filteredTransactions, transactions[i])
		}
	}

	// loop over filteredTransaction. Otherwise it would grab transactions from years before the startYear
	for i := len(filteredTransactions) - 1; i > -0; i-- {
		if filteredTransactions[i].TransactType == "AffiliatePayout" && filteredTransactions[i].Timestamp.Year() >= startYear {

			// first month will be empty, assign the first month here
			if currentMonth == "" {
				currentMonth = filteredTransactions[i].Timestamp.Month().String()
			}

			// detect new month
			if currentMonth != filteredTransactions[i].Timestamp.Month().String() {
				//fmt.Printf("current month: %s  |  filter: %s \n", currentMonth, filteredTransactions[i].Timestamp.Month().String())
				currentMonth = filteredTransactions[i].Timestamp.Month().String()
				//result.Btc = monthBTC.String()
				result.Btc = monthBTC
				//result.Dollar = fmt.Sprintf("$%.2f", monthBTC.ToBTC()*bitcoinPrice)
				result.Dollar = monthBTC.ToBTC() * bitcoinPrice
				stats.AddStat(*result) // commit the previous Stat
				result = new(Stat)     // prepare new Stat for new month. Also reset MonthBTC and MonthDollar for next month
				monthBTC = 0
				monthDollar = 0.0
			}

			btc, _ := btcutil.NewAmount(float64(filteredTransactions[i].Amount) / 100000000)
			result.Date = filteredTransactions[i].Timestamp.Month().String()
			monthBTC += btc
			monthDollar += btc.ToBTC() * bitcoinPrice

			totalBTC += btc
			totalDollar += btc.ToBTC() * bitcoinPrice
		}
	}

	// calculate January
	// get the other months which don't include january and combine that
	// total dollars and btc minus every other earned except january results in earnings January
	// this calculation due some big / wrong logic which puts every month in the stats except January due to the reverse loop
	var dollarJanuary = 0.0
	var bitcoinsJanuary = btcutil.Amount(0)
	for _, st := range stats.Stat {
		dollarJanuary += st.Dollar
		bitcoinsJanuary += st.Btc
	}

	result.Date = "January"
	result.Btc = totalBTC - bitcoinsJanuary
	result.Dollar = totalDollar - dollarJanuary
	stats.AddStat(*result) // push January stats

	stats.TotalDollar = fmt.Sprintf("$%.2f", totalDollar)
	stats.TotalBtc = totalBTC.String()

	return stats
}

type Stat struct {
	Date   string
	Btc    btcutil.Amount
	Dollar float64
	Change string
}

type Stats struct {
	Stat        []Stat
	TotalBtc    string
	TotalDollar string
}

// get earning stats from monday till current day
func WeeklyEarnings(transactions []bitmexgo.Transaction) *Stats {
	startYear := viper.GetInt("start_year")
	bitcoinPrice := bitcoin.ToDollar()

	stats := new(Stats)
	var totalBTC btcutil.Amount
	var totalDollar float64

	for i := 0; i < len(transactions); i++ {
		if transactions[i].TransactType == "AffiliatePayout" && transactions[i].Timestamp.Year() >= startYear {
			result := new(Stat)
			btc, _ := btcutil.NewAmount(float64(transactions[i].Amount) / 100000000)
			result.Date = transactions[i].Timestamp.Weekday().String()
			//result.Btc = btc.String()
			result.Btc = btc

			//result.Dollar = fmt.Sprintf("$%.2f", btc.ToBTC()*bitcoinPrice)
			result.Dollar = btc.ToBTC() * bitcoinPrice
			stats.AddStat(*result)

			totalBTC += btc
			totalDollar += btc.ToBTC() * bitcoinPrice

			if result.Date == "Monday" {
				break
			}
		}
	}

	stats.TotalDollar = fmt.Sprintf("$%.2f", totalDollar)
	stats.TotalBtc = totalBTC.String()
	return stats
}

func (s *Stats) AddStat(item Stat) []Stat {
	s.Stat = append(s.Stat, item)
	return s.Stat
}

type AffiliateStatus struct {
	PrevPayout          string `json:"prevPayout"`
	PrevTurnover        string `json:"prevTurnover"`
	TotalReferrals      int    `json:"totalReferrals"`
	TotalTurnover       string `json:"totalTurnover"`
	TotalComm           string `json:"totalComm"`
	PendingPayout       string `json:"pendingPayout"`
	PendingPayoutDollar string
}

// will send an affiliate status message like this:
// Previous payout: 18.65300541 BTC
// Total turnover: 152550.62037547 BTC
// Total referrals: 468
// Pending payout: 0.00321705 BTC  - $23.93
func Status() (AffiliateStatus, error) {
	auth, apiClient := client.GetInstance()
	status, _, err := apiClient.UserApi.UserGetAffiliateStatus(auth)
	if err != nil {
		log.Printf("Error affiliate status: %s\n", err)
		return AffiliateStatus{}, err
	}

	amountPrevPayout, _ := btcutil.NewAmount(float64(status.PrevPayout) / 100000000)
	amountTotalTurnover, _ := btcutil.NewAmount(float64(status.TotalTurnover) / 100000000)
	amountTotalCommission, _ := btcutil.NewAmount(float64(status.TotalComm) / 100000000)
	amountPendingPayout, _ := btcutil.NewAmount(float64(status.PendingPayout) / 100000000)
	affiliateStatus := AffiliateStatus{
		PrevPayout:          amountPrevPayout.String(),
		TotalReferrals:      status.TotalReferrals,
		TotalTurnover:       amountTotalTurnover.String(),
		TotalComm:           amountTotalCommission.String(),
		PendingPayout:       amountPendingPayout.String(),
		PendingPayoutDollar: fmt.Sprintf("$%.2f", amountPendingPayout.ToBTC()*bitcoin.ToDollar()),
	}
	return affiliateStatus, nil
}
