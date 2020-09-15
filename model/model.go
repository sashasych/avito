package model

import "time"

type BalanceChange struct {
	UserID int
	Change float64
}

type TransferChange struct {
	FromUserID int
	ToUserID   int
	Change     float64
}

type GetBalance struct {
	ID           int
	CurrencyType string
}

type Balance struct {
	Balance      float64
	CurrencyType string
}

type CurrenciesValue struct {
	Currencies   map[string]float64 `json:"rates"`
	BaseCurrency *string            `json:"base"`
	Date         time.Time          `json:"-"`
}
