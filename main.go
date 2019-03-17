// Copyright (c) 2019 Romano (Viacoin developer)
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/bclicn/color"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Transaction struct {
	Time          string
	Type          string
	Amount        float64
	Fee           float64
	Address       string
	Status        string
	WalletBalance string
}

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

func (month Month) referralEarning(transactions []Transaction) {
	var earnedBTC float64
	var count int

	for _, tx := range transactions {
		if tx.Type == "AffiliatePayout" {
			earnedBTC += (tx.Amount / 100000000)
			count += 1
		}
	}
	if earnedBTC <= 0{
		fmt.Printf("earned ref fees for %s: \t "+color.Red("%f BTC\n"), time.Month(month), earnedBTC)
	}else {
		fmt.Printf("earned ref fees for %s: \t "+color.Green("%f BTC\n"), time.Month(month), earnedBTC)
	}
}

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

func ScanFiles(extension string) string {
	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if !f.IsDir() {
			r, err := regexp.MatchString(extension, f.Name())
			if err == nil && r {
				//return nil
				return f.Name()
			}
		}
	}
	os.Exit(0)
	return ""
}

func ReadCSV(file string) ([]Transaction, error) {

	var transactions []Transaction
	csvFile, _ := os.Open("./" + file)
	r := csv.NewReader(bufio.NewReader(csvFile))

	r.LazyQuotes = true
	r.Comma = ','

	_, err := r.Read()
	if err != nil && err != io.EOF {
		return nil, err
	}

	for {
		record, error := r.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}

		amount, _ := strconv.ParseFloat(record[2], 64)
		fee, _ := strconv.ParseFloat(record[3], 64)

		tx := Transaction{
			Time:          record[0],
			Type:          record[1],
			Amount:        amount,
			Fee:           fee,
			Address:       record[4],
			Status:        record[5],
			WalletBalance: record[6],
		}

		transactions = append(transactions, tx)
	}

	return transactions, nil
}

func calculateTotalReferral(transactions []Transaction) {

	var earned float64
	var monthly Months

	var monthlyTransactions [12][]Transaction

	for _, tx := range transactions {
		if tx.Type == "AffiliatePayout" {
			earned += tx.Amount / 100000000
		}

		time := strings.Split(tx.Time, ",")
		date := strings.Split(time[0], "/")
		month, _ := strconv.Atoi(date[0])

		months := [12]Month{monthly.Jan, monthly.Feb, monthly.Mar, monthly.Apr, monthly.May, monthly.Jun, monthly.Jul, monthly.Aug, monthly.Sept, monthly.Oct, monthly.Nov, monthly.Dec}
		for index, _ := range months {
			switch month {
			case index:
				monthlyTransactions[month-1] = append(monthlyTransactions[month-1], tx)
			}
		}
	}

	// commit ref earnings
	months := [12]Month{monthly.Jan, monthly.Feb, monthly.Mar, monthly.Apr, monthly.May, monthly.Jun, monthly.Jul, monthly.Aug, monthly.Sept, monthly.Oct, monthly.Nov, monthly.Dec}
	for index, monthly := range months {
		monthly = JAN+Month(index) // assign month an integer to convert later to month name in string form
		monthly.referralEarning(monthlyTransactions[index])
	}

	fmt.Printf("\nTotal earned ref fees:\t " + color.Green("%f BTC\n"), earned)
}

func main() {
	file := ScanFiles(".csv")
	transactions, _ := ReadCSV(file)
	calculateTotalReferral(transactions)
}
