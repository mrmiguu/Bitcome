package src

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func init() {
	b, _ := buy()
	fmt.Printf("LTC:buy   $%.2f\n", b)
	s, _ := sell()
	fmt.Printf("LTC:sell  $%.2f\n", s)
	net := s - b
	fmt.Printf("LTC:net   $%.2f\n", net)
	pct := (s/b - 1) * 100
	fmt.Printf("LTC:pct   %.2f%%\n", pct)
}

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
	var r struct {
		Data struct {
			Amount   string
			Currency string
		}
	}
	json.NewDecoder(resp.Body).Decode(&r)
	resp.Body.Close()
	amt, err := strconv.ParseFloat(r.Data.Amount, 64)
	if err != nil {
		return 0, err
	}
	return amt, nil
}
