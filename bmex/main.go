package main

import (
	"fmt"
	"github.com/romanornr/Bitmex-referral-analyzer/config"
	"github.com/spf13/viper"
	"github.com/zmxv/bitmexgo"
	"time"
)

var c config.Conf

func main() {

	config.GetViperConfig()

	// Get your API key/secret pair at https://www.bitmex.com/app/apiKeys
	apiKey := viper.GetString("api_key")
	apiSecret := viper.GetString("api_secret")

	// Create an authentication context
	auth := bitmexgo.NewAPIKeyContext(apiKey, apiSecret)

	// Create a shareable API client instance
	apiClient := bitmexgo.NewAPIClient(bitmexgo.NewConfiguration())

	// Call APIs without parameters by passing the auth context.
	// e.g. getting exchange-wide turnover and volume statistics:
	tx, _, err := apiClient.UserApi.UserGetWalletHistory(auth, nil)
	if err != nil {
		fmt.Println(err)
	}

	//var earned float64

	referralEarning(tx)
	x := MonthEarned(6)
	fmt.Println(x)

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

var previousMonthEarning float64
var monthlyTransactions [12][]bitmexgo.Transaction

func referralEarning(transactions []bitmexgo.Transaction) {

	config.GetViperConfig()
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

func MonthEarned(month int) float64 {
	var earnedBTC float64
	monthTransactions := monthlyTransactions[month-1]
	for i := 0; i < len(monthTransactions); i++ {
		amount := float64(monthTransactions[i].Amount) / 100000000
		earnedBTC += amount
	}
	return earnedBTC
}
