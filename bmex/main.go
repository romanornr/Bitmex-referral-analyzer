package bmex

import (
	"fmt"
	"github.com/btcsuite/btcutil"
	"github.com/romanornr/Bitmex-referral-analyzer/bitcoin"
	"github.com/romanornr/Bitmex-referral-analyzer/client"
	"github.com/romanornr/Bitmex-referral-analyzer/config"
	"github.com/spf13/viper"
	"github.com/zmxv/bitmexgo"
	"time"
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

	referralEarning(tx)
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

type Month time.Month

const (
	JAN Month = iota + 1
	FEB
	MAR
	APR
	MAY
	JUN
	JUL
	AUG
	SEPT
	OCT
	NOV
	DEC
)

type monthly struct {
	Jan float64
	//Feb  Month
	//Mar  Month
	//Apr  Month
	//May  Month
	//Jun  Month
	//Jul  Month
	//Aug  Month
	//Sept Month
	//Oct  Month
	//Nov  Month
	//Dec  Month
}

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

func referralEarning(transactions []bitmexgo.Transaction) {

	start_year := viper.GetInt("start_year")

	months := [12]Month{JAN, FEB, MAR, APR, MAY, JUN, JUL, SEPT, OCT, NOV, DEC}

	for i := len(transactions) - 1; i >= 0; i-- {

		if transactions[i].TransactType == "AffiliatePayout" && transactions[i].Timestamp.Year() >= start_year {
			for index, _ := range months {
				month := int(transactions[i].Timestamp.Month())

				switch month {
				case index:
					monthlyTransactions[month-1] = append(monthlyTransactions[month-1], transactions[i])
				}
			}

		}
	}
	//fmt.Println(monthlyTransactions[0])
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
