package server

import (
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/sashasych/avito/db"
	"github.com/sashasych/avito/model"
	"log"
	"net/http"
)

func StartServer() {

	//TODO API методы возвращают читабельные ошибки и коды (ДОП)
	//TODO Разобраться со скриптом sql + скрипт bash под тестовый сценарий
	//TODO Валидация данных и обработка ошибок

	//TODO Исправить показание баланса в евро(частный случай)
	//TODO пагинация + сортировка по сумме и дате в истории операций
	
	//TODO Отрефакторить код
	
	database, err := db.CreateConnection()
	if err != nil {
		log.Fatal(err)
		return
	}

	http.HandleFunc("/updateBalance", func(w http.ResponseWriter, r *http.Request) {
		balanceChangeRequest := model.BalanceChangeRequest{}
		err = json.NewDecoder(r.Body).Decode(&balanceChangeRequest)
		fmt.Println(balanceChangeRequest)
		err = database.UpdateBalance(balanceChangeRequest.UserID, balanceChangeRequest.Change)
		if err != nil {
			fmt.Errorf("%s", err.Error())
		}
	})

	http.HandleFunc("/transferMoney", func(w http.ResponseWriter, r *http.Request) {
		transferChange := model.TransferChange{}
		err = json.NewDecoder(r.Body).Decode(&transferChange)
		fmt.Println(transferChange)
		err = database.TransferMoney(transferChange.ToUserID, transferChange.FromUserID, transferChange.Change)
		if err != nil {
			fmt.Errorf("%s", err.Error())
		}
	})

	http.HandleFunc("/getBalance", func(w http.ResponseWriter, r *http.Request) {
		getBalanceRequest := model.GetBalanceRequest{}
		err = json.NewDecoder(r.Body).Decode(&getBalanceRequest)
		fmt.Println(getBalanceRequest)
		getBalanceResponse := model.GetBalanceResponse{}
		getBalanceResponse.CurrencyType = getBalanceRequest.CurrencyType
		err, getBalanceResponse.Balance  = database.GetBalance(getBalanceRequest.ID, getBalanceRequest.CurrencyType)
		if err != nil {
			fmt.Errorf("%s", err.Error())
		} else {
			jsonResponse(w, "OK", getBalanceResponse)
		}
	})

	http.HandleFunc("/getHistory", func(w http.ResponseWriter, r *http.Request) {
		getHistoryRequest := model.HistoryRequest{}
		err := json.NewDecoder(r.Body).Decode(&getHistoryRequest)
		if err != nil {
			fmt.Errorf("%s", err.Error())
		}
		fmt.Println(getHistoryRequest)
		getHistoryResponse := model.HistoryResponse{}
		err, getHistoryResponse.Transactions = database.GetHistory(getHistoryRequest.UserID)
		if err != nil {
			fmt.Errorf("%s", err.Error())
		} else {
			jsonResponse(w, "OK", getHistoryResponse)
		}
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
