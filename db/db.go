package db

import (
	"database/sql"
	"fmt"
	"github.com/sashasych/avito/service"
)

type Database struct {
	connection *sql.DB
}


func CreateConnection() (*Database, error) {
	connection, err := sql.Open("postgres", "user=postgres password=123 dbname=billing sslmode=disable")
	return &Database{connection} , err
}

func (db *Database) logTransaction(userId int, message string, amount float64) error {
	_, err := db.connection.Exec("INSERT INTO transaction_history(userid, info, amount, date) VALUES ($1, $2, $3, now())", userId, message, amount)
	return err
}

func (db *Database) UpdateBalance(userId int, amount float64) error {
	rows, err := db.connection.Query("SELECT * From user_balance WHERE id=$1", userId)
	if rows.Next() {
		_, err = db.connection.Exec("UPDATE user_balance SET balance=balance+$2::NUMERIC::MONEY WHERE id=$1", userId, amount)
	} else {
		rows, err = db.connection.Query("INSERT INTO user_balance(balance) VALUES ($1) RETURNING id", amount)
		rows.Next()
		err = rows.Scan(&userId)
	}
	if err == nil {
		err = db.logTransaction(userId, "", amount)
	}
	return err
}

func (db *Database) TransferMoney(userIdTo, userIdFrom int, amount float64) error {
	tx, err := db.connection.Begin()
	tx.Exec("UPDATE user_balance SET balance=balance+$2::NUMERIC::MONEY WHERE id=$1", userIdTo, amount)
	if err != nil {
		tx.Rollback()
		return err

	}
	tx.Exec("UPDATE user_balance SET balance=balance-$2::NUMERIC::MONEY WHERE id=$1", userIdFrom, amount)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	err = db.logTransaction(userIdTo, "", amount)
	if err == nil {
		err = db.logTransaction(userIdFrom, "", -amount)
	}
	return err
}

func (db *Database) GetBalance(userId int, currencyType string) (err error, balance float64) {

	rows, err := db.connection.Query("SELECT balance::DECIMAL From user_balance WHERE id=$1", userId)
	if err != nil {
		return err, 0
	}
	fmt.Println(rows)
	rows.Next()
	err = rows.Scan(&balance)
	if err != nil {
		return err, 0
	}
	if currencyType != "RUB" {
		balance, err = service.FetchDataFromExchangeApi(currencyType, balance)
	}
	return err, balance
}