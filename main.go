// Copyright (c) 2019 Romano (Viacoin developer)
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package main

import (
	"fmt"
	"github.com/bclicn/color"
	"github.com/romanornr/Bitmex-referral-analyzer/account"
	"github.com/romanornr/Bitmex-referral-analyzer/csv"
	"strconv"
	"strings"
	"time"
)

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
		if tx.Type == "AffiliatePayout" {
			earned += tx.Amount / 100000000
		}

		time := strings.Split(tx.Time, ",")
		date := strings.Split(time[0], "/")
		month, _ := strconv.Atoi(date[0])

		for index, _ := range months {
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

	fmt.Printf("\nTotal earned referral fees:\t "+color.Green("%f BTC\n\n"), earned)
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
	changeMessage = fmt.Sprintf("change: " + color.Green("+%.2f%%\n"), change) // change: +347.98%
	if change < 0 {
		changeMessage = fmt.Sprintf("change: " + color.Red("%.2f%%\n"), change) // change: -85.95%
	}

	if earnedBTC <= 0 {
		fmt.Printf("Bitmex referral fees %s \t "+color.Red("%f BTC\n"), time.Month(month), earnedBTC)
	} else {
		fmt.Printf("Bitmex referral fees %s \t "+color.Green("%f BTC \t") + "%s", time.Month(month), earnedBTC, changeMessage)
	}
	previousMonthEarning = earnedBTC
}

func main() {
	file := csv.ScanFiles(".csv")
	transactions, _ := csv.ReadCSV(file)
	calculateTotalReferral(transactions)
	fmt.Printf("%s\n", BITMEXREFLINK)
}
