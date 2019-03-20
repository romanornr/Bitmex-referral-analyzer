// Copyright (c) 2019 Romano (Viacoin developer)
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package main

import (
	"encoding/json"
	"fmt"
	"github.com/bclicn/color"
	"github.com/romanornr/Bitmex-referral-analyzer/account"
	"github.com/romanornr/Bitmex-referral-analyzer/csv"
	"github.com/romanornr/Bitmex-referral-analyzer/config"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var c config.Conf

type Month time.Month

const BITMEXREFLINK = "https://www.bitmex.com/register/vhT2qm"

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

type Months struct {
	Jan  Month
	Feb  Month
	Mar  Month
	Apr  Month
	May  Month
	Jun  Month
	Jul  Month
	Aug  Month
	Sept Month
	Oct  Month
	Nov  Month
	Dec  Month
}

func calculateTotalReferral(transactions []account.Transaction) {

	var earned float64
	var monthly Months

	var monthlyTransactions [12][]account.Transaction
	months := [12]Month{monthly.Jan, monthly.Feb, monthly.Mar, monthly.Apr, monthly.May, monthly.Jun, monthly.Jul, monthly.Aug, monthly.Sept, monthly.Oct, monthly.Nov, monthly.Dec}

	for _, tx := range transactions {

		time := strings.Split(tx.Time, ",")
		date := strings.Split(time[0], "/")
		year, _ := strconv.Atoi(date[2])
		month, _ := strconv.Atoi(date[0])

		if tx.Type == "AffiliatePayout" && year >= c.Start_year {
			earned += tx.Amount / 100000000
		}

		for index, _ := range months {
			if year < c.Start_year {
				continue
			}
			switch month {
			case index:
				monthlyTransactions[month-1] = append(monthlyTransactions[month-1], tx)
			}
		}
	}

	// commit ref earnings
	for index, monthly := range months {
		monthly = JAN + Month(index) // assign month an integer to convert later to month name in string form
		monthly.referralEarning(monthlyTransactions[index])
	}

	fmt.Printf("\nTotal paid out referral fees:\t "+color.Green("%f BTC\n\n"), earned)
}

var previousMonthEarning float64

func (month Month) referralEarning(transactions []account.Transaction) {
	var earnedBTC float64
	var count int

	for _, tx := range transactions {

			if tx.Type == "AffiliatePayout" {
				earnedBTC += (tx.Amount / 100000000)
				count += 1
		}
	}

	change := ((earnedBTC - previousMonthEarning) / previousMonthEarning) * 100
	if previousMonthEarning == 0 {
		change = 0.00
	}

	var changeMessage string
	changeMessage = fmt.Sprintf("change: "+color.Green("+%.2f%%\n"), change) // change: +347.98%
	if change < 0 {
		changeMessage = fmt.Sprintf("change: "+color.Red("%.2f%%\n"), change) // change: -85.95%
	}

	earnedDollar := fmt.Sprintf(color.Green("$ %.2f"), earnedBTC * bitcoinPrice)

	if earnedBTC <= 0 {
		fmt.Printf("Bitmex referral fees %s \t "+color.Red("%f BTC\n"), time.Month(month), earnedBTC)
	} else {
		fmt.Printf("Bitmex referral fees %s \t "+color.Green("%f BTC \t")+" %s\t %s", time.Month(month), earnedBTC, earnedDollar, changeMessage)
	}
	previousMonthEarning = earnedBTC
}

type bitcoinTicker struct {
	High      string `json:"high"`
	Last      string `json:"last"`
	Timestamp string `json:"timestamp"`
	Bid       string `json:"bid"`
	Vwap      string `json:"vwap"`
	Volume    string `json:"volume"`
	Low       string `json:"low"`
	Ask       string `json:"ask"`
	Open      string `json:"open"`
}

func bitcoinToDollar() {
	url := fmt.Sprintf("https://www.bitstamp.net/api/v2/ticker/btcusd")
	client := http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	valueBTC := bitcoinTicker{}
	err = json.Unmarshal(body, &valueBTC)
	if err != nil {
		log.Fatal(err)
	}

	bitcoinPrice, _ = strconv.ParseFloat(valueBTC.Bid, 64)
}

var bitcoinPrice float64

func init() {
	c.GetConf()
	bitcoinToDollar()
}

func main() {
	file := csv.ScanFiles(".csv")
	transactions, _ := csv.ReadCSV(file)
	calculateTotalReferral(transactions)
	fmt.Printf("%s\n", BITMEXREFLINK)
}
