package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
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
	fmt.Printf("earned ref fees for %s: \t %f BTC\n", time.Month(month), earnedBTC)
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

		switch month {
		case 1:
			monthlyTransactions[month-1] = append(monthlyTransactions[month-1], tx)
		case 2:
			monthlyTransactions[month-1] = append(monthlyTransactions[month-1], tx)
		case 3:
			monthlyTransactions[month-1] = append(monthlyTransactions[month-1], tx)
		case 4:
			monthlyTransactions[month-1] = append(monthlyTransactions[month-1], tx)
		case 5:
			monthlyTransactions[month-1] = append(monthlyTransactions[month-1], tx)
		case 6:
			monthlyTransactions[month-1] = append(monthlyTransactions[month-1], tx)
		case 7:
			monthlyTransactions[month-1] = append(monthlyTransactions[month-1], tx)
		case 8:
			monthlyTransactions[month-1] = append(monthlyTransactions[month-1], tx)
		case 9:
			monthlyTransactions[month-1] = append(monthlyTransactions[month-1], tx)
		case 10:
			monthlyTransactions[month-1] = append(monthlyTransactions[month-1], tx)
		case 11:
			monthlyTransactions[month-1] = append(monthlyTransactions[month-1], tx)
		case 12:
			monthlyTransactions[month-1] = append(monthlyTransactions[month-1], tx)
		}
	}

	monthly.Jan, monthly.Feb, monthly.Mar, monthly.Apr, monthly.May, monthly.Jun, monthly.Jul, monthly.Aug, monthly.Sept, monthly.Oct, monthly.Nov, monthly.Dec = JAN, FEB, MAR, APR, MAY, JUN, JUL, AUG, SEPT, OCT, NOV, DEC

	monthly.Jan.referralEarning(monthlyTransactions[0])
	monthly.Feb.referralEarning(monthlyTransactions[1])
	monthly.Mar.referralEarning(monthlyTransactions[2])
	monthly.Apr.referralEarning(monthlyTransactions[3])
	monthly.May.referralEarning(monthlyTransactions[4])
	monthly.Jun.referralEarning(monthlyTransactions[5])
	monthly.Jul.referralEarning(monthlyTransactions[6])
	monthly.Aug.referralEarning(monthlyTransactions[7])
	monthly.Sept.referralEarning(monthlyTransactions[8])
	monthly.Oct.referralEarning(monthlyTransactions[9])
	monthly.Nov.referralEarning(monthlyTransactions[10])
	monthly.Dec.referralEarning(monthlyTransactions[11])

	fmt.Printf("Total earned ref fees:\t %f BTC\n", earned)
}

func main() {
	file := ScanFiles(".csv")
	transactions, _ := ReadCSV(file)
	calculateTotalReferral(transactions)
}
