package model

import "time"

type BalanceChangeRequest struct {
	UserID int
	Change float64
}

type TransferChange struct {
	FromUserID int
	ToUserID   int
	Change     float64
}

type GetBalanceRequest struct {
	ID           int
	CurrencyType string
}

type GetBalanceResponse struct {
	Balance      float64
	CurrencyType string
}

type CurrenciesValue struct {
	Currencies   map[string]float64 `json:"rates"`
	BaseCurrency *string            `json:"base"`
	Date         time.Time          `json:"-"`
}

type HistoryRequest struct {
	ID int
}

type HistoryResponse struct {
	Transactions []Transaction
}

type Transaction struct {
	info string
	amount float64
	date time.Time
}