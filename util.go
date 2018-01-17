package bitcome

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

var (
	println = log.Println
)

func buy() (float64, error) {
	return buySell("buy")
}

func sell() (float64, error) {
	return buySell("sell")
}

func buySell(action string) (float64, error) {
	resp, err := http.Get("https://api.coinbase.com/v2/prices/LTC-USD/" + action)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	var x struct {
		Data struct {
			Amount   string
			Currency string
		}
	}
	if err := json.NewDecoder(resp.Body).Decode(&x); err != nil {
		return 0, err
	}
	return strconv.ParseFloat(x.Data.Amount, 64)
}

func must(err error) {
	if err == nil {
		return
	}
	panic(err)
}

func warn(err error) error {
	if err != nil {
		println(err)
	}
	return err
}

func ftoa(f float64) string {
	return fmt.Sprintf("%.2f", f)
}

func percent(b, s float64) float64 {
	return (s/b - 1) * 100
}

func appendfile(d Data, file string) (Data, error) {
	dp := make(Data)
	f, err := os.Open(file)
	if err != nil {
		return d, err
	}
	defer f.Close()
	if err = json.NewDecoder(f).Decode(&dp); err != nil {
		return d, err
	}
	for time, usd := range d {
		dp[time] = usd
	}
	return dp, nil
}

func savedata(d Data, file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(d)
}
