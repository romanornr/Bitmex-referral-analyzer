package csv

import (
	"bufio"
	"encoding/csv"
	"github.com/romanornr/Bitmex-referral-analyzer/account"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
)

func ScanFiles(extension string) []string {

	csvFiles := []string{}

	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if !f.IsDir() {
			r, err := regexp.MatchString(extension, f.Name())
			if err == nil && r {
				csvFiles = append(csvFiles, f.Name())
			}
		}
	}
	return csvFiles
}

func ReadCSVFiles(csvfiles []string) ([]account.Transaction, error) {

	var transactions []account.Transaction

	// Range over multiple CSV files
	for _, c := range csvfiles {
		csvFile, _ := os.Open("./" + c)
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

			tx := account.Transaction{
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
	}

	return removeDuplicates(transactions), nil
}

// We combine 2 csv files with some having duplicate transactions
// this function removes duplicate transctions
func removeDuplicates(transactions []account.Transaction) []account.Transaction {
	// Use map to record duplicates as we find them.
	encountered := map[account.Transaction]bool{}
	result := []account.Transaction{}

	for v := range transactions {
		if encountered[transactions[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[transactions[v]] = true
			// Append to result slice.
			result = append(result, transactions[v])
		}
	}
	// Return the new slice.
	return result
}
