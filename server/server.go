package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/sashasych/avito/model"
	"github.com/sashasych/avito/service"
	"log"
	"net/http"
)

func StartServer() {

	//TODO вынести код в дб
	//TODO обработать ошибки
	//TODO вынести все методы в service
	connStr := "user=postgres password=123 dbname=billing sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	//

	http.HandleFunc("/updateBalance", func(w http.ResponseWriter, r *http.Request) {
		balanceChange := model.BalanceChange{}
		err = json.NewDecoder(r.Body).Decode(&balanceChange)
		fmt.Println(balanceChange)
		rows, err := db.Query("SELECT * From user_balance WHERE id=$1", balanceChange.UserID)
		if rows.Next() {
			_, err = db.Exec("UPDATE user_balance SET balance=balance+$2::MONEY WHERE id=$1", balanceChange.UserID, balanceChange.Change)
			if err != nil {
				panic(err)
			}
		} else {
			_, err = db.Exec("INSERT INTO user_balance(balance) VALUES ($1)", balanceChange.Change)
			if err != nil {
				panic(err)
			}
		}
	})

	http.HandleFunc("/transferMoney", func(w http.ResponseWriter, r *http.Request) {
		transferChange := model.TransferChange{}
		err = json.NewDecoder(r.Body).Decode(&transferChange)
		fmt.Println(transferChange)
		tx, err := db.Begin()

		tx.Exec("UPDATE user_balance SET balance=balance+$2::NUMERIC::MONEY WHERE id=$1", transferChange.ToUserID, transferChange.Change)
		if err != nil {
			log.Fatal(err)
			tx.Rollback()
		}
		tx.Exec("UPDATE user_balance SET balance=balance-$2::NUMERIC::MONEY WHERE id=$1", transferChange.FromUserID, transferChange.Change)
		if err != nil {
			log.Fatal(err)
			tx.Rollback()
		}
		tx.Commit()
	})

	//TODO добавить вывод баланса в долларах(по умолчанию) возможно что-то с константами сделать
	http.HandleFunc("/getBalance", func(w http.ResponseWriter, r *http.Request) {
		getBalance := model.GetBalance{}
		err = json.NewDecoder(r.Body).Decode(&getBalance)
		fmt.Println(getBalance)
		balance := model.Balance{}
		rows, err := db.Query("SELECT balance::DECIMAL From user_balance WHERE id=$1", getBalance.ID)
		fmt.Println(rows)
		rows.Next()
		err = rows.Scan(&balance.Balance)
		balance.CurrencyType = "RUB"
		if getBalance.CurrencyType != "RUB" {
			// TODO вызов метода получения нового баланса делаем запрос
			balance.Balance, err = service.FetchDataFromExchangeApi(getBalance.CurrencyType, balance.Balance)
			balance.CurrencyType = getBalance.CurrencyType
		}
		if err != nil {
			panic(err)
		}
		jsonResponse(w, "OK", balance)
	})

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}

func jsonResponse(rw http.ResponseWriter, message string, data interface{}) {
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(struct {
		Message string
		Data    interface{}
	}{
		message,
		data,
	})
}
