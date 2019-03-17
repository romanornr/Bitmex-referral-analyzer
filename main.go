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

type Month int

const (
	JAN Month = iota + 1
	FEB
	MAR
	APR
	MAY
	JUN
	AUG
	SEPT
	OCT
	NOV
	DEC
)

func (month Month) referralEarning(transactions []Transaction) {
	var earnedBTC float64
	var count int
	//t := time.Month(2)
	//fmt.Println(t)

	for _, tx := range transactions {
		if tx.Type == "AffiliatePayout" {
			earnedBTC += (tx.Amount / 100000000)
			count += 1
		}
	}
	fmt.Printf("earned ref fees for %d : %f BTC\n", month, earnedBTC)
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
	var m Months

	var monthlyTransactions [11][]Transaction

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
		}

	}

	m.Jan.referralEarning(monthlyTransactions[0])
	m.Feb.referralEarning(monthlyTransactions[1])
	m.Mar.referralEarning(monthlyTransactions[2])

	fmt.Printf("Total earned ref fees: %f BTC\n", earned)
}

func main() {
	file := ScanFiles(".csv")
	transactions, _ := ReadCSV(file)
	calculateTotalReferral(transactions)
}
