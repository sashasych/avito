package service

import (
	"encoding/json"
	"fmt"
	"github.com/sashasych/avito/model"
	"io/ioutil"
	"net/http"
)

func FetchDataFromExchangeApi(str string, balance float64) (float64, error) {
	url := "https://api.exchangeratesapi.io/latest"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("error happend", err)
		return 0, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	data := new(model.CurrenciesValue)
	err = json.Unmarshal(respBody, data)
	if err != nil {
		return 0, err
	}
	if value, ok := data.Currencies[str]; ok {
		value = data.Currencies["RUB"] / value
		balance = balance / value
		return balance, nil
	} else {
		return 0, nil
	}
}
