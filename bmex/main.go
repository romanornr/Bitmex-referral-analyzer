package bmex

import (
	"fmt"
	"github.com/btcsuite/btcutil"
	"github.com/romanornr/Bitmex-referral-analyzer/bitcoin"
	"github.com/romanornr/Bitmex-referral-analyzer/client"
	"github.com/romanornr/Bitmex-referral-analyzer/config"
	"github.com/spf13/viper"
	"github.com/zmxv/bitmexgo"
)

var c config.Conf
var apiKey string
var apiSecret string

func main() {

	auth, apiClient := client.GetInstance()

	// Call APIs without parameters by passing the auth context.
	// e.g. getting exchange-wide turnover and volume statistics:
	tx, _, err := apiClient.UserApi.UserGetWalletHistory(auth, nil)
	if err != nil {
		fmt.Println(err)
	}

	//var earned float64

	//referralEarning(tx)
	x := MonthEarned(7)
	fmt.Println(x)
	y := WeeklyEarnings(tx)
	fmt.Println(y)

	//for i := len(tx) - 1; i >= 0; i--{
	//	year := tx[i].Timestamp.Year()
	//	if tx[i].TransactType == "AffiliatePayout" && int64(year) == int64(c.Start_year) {
	//		amount := float64(tx[i].Amount) / 100000000
	//		earned += amount
	//	}
	//
	//}
	//fmt.Println(earned)

}

const BITMEXREFLINK = "https://www.bitmex.com/register/vhT2qm"

// loads wallet history by using bitmex api
// returns Transactions which can be used by other functions
func LoadWalletHistory() (error, []bitmexgo.Transaction){
	auth, apiClient := client.GetInstance()

	// Call APIs without parameters by passing the auth context.
	// e.g. getting exchange-wide turnover and volume statistics:
	tx, _, err := apiClient.UserApi.UserGetWalletHistory(auth, nil)
	if err != nil {
		fmt.Println(err)
	}
	return err, tx
}

var previousMonthEarning float64
var monthlyTransactions [12][]bitmexgo.Transaction

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

	for i := 0; i < len(transactions); i++ {
		if transactions[i].TransactType == "AffiliatePayout" && transactions[i].Timestamp.Year() >= startYear {

			if currentMonth != transactions[i].Timestamp.Month().String() || currentMonth == "" {
				fmt.Printf("currentMonth: %s\n", currentMonth)
				fmt.Printf("tx month: %s\n", transactions[i].Timestamp.Month())
				currentMonth = transactions[i].Timestamp.Month().String()
				result.Btc = monthBTC.String()
				result.Dollar = fmt.Sprintf("$%.2f", monthBTC.ToBTC()*bitcoinPrice)
				stats.AddStat(*result)
				result = new(Stat)
				monthBTC = 0
				monthDollar = 0.0
			}

			btc, _ := btcutil.NewAmount(float64(transactions[i].Amount) / 100000000)
			result.Date = transactions[i].Timestamp.Month().String()
			monthBTC += btc
			monthDollar += btc.ToBTC()*bitcoinPrice

			totalBTC += btc
			totalDollar += btc.ToBTC() * bitcoinPrice
		}
	}

	stats.TotalDollar = fmt.Sprintf("$%.2f", totalDollar)
	stats.TotalBtc = totalBTC.String()

	fmt.Println(stats)
	return stats
}


type Stat struct {
	Date   string
	Btc    string
	Dollar string
	Change string
}

type Stats struct {
	Stat []Stat
	TotalBtc string
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
			result.Btc = btc.String()

			result.Dollar = fmt.Sprintf("$%.2f", btc.ToBTC()*bitcoinPrice)
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

func MonthEarned(month int) float64 {
	var earnedBTC float64
	monthTransactions := monthlyTransactions[month-1]
	for i := 0; i < len(monthTransactions); i++ {
		amount := float64(monthTransactions[i].Amount) / 100000000
		earnedBTC += amount
	}
	return earnedBTC
}
