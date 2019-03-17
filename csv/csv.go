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

func ReadCSV(file string) ([]account.Transaction, error) {

	var transactions []account.Transaction
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

	return transactions, nil
}
