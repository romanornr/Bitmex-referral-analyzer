package bitcoin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

func ToDollar() float64 {
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

	valueBTC := Ticker{}
	err = json.Unmarshal(body, &valueBTC)
	if err != nil {
		log.Fatal(err)
	}

	bitcoinPrice, _ := strconv.ParseFloat(valueBTC.Bid, 64)
	return bitcoinPrice
}
