package model

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
	ID int
}

type Balance struct {
	Balance float64
}
