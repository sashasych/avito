# Тестовое задание в Avito на Unit Job
## Установка и запуск БД
```
psql -U postgres
CREATE DATABASE billing;
GRANT ALL PRIVILEGES ON DATABASE billing to postgres;
\c billing
CREATE TABLE user_balance(id SERIAL PRIMARY key, balance MONEY);
CREATE TABLE transaction_history(id SERIAL PRIMARY KEY, UserId INTEGER REFERENCES user_balance (id), info TEXT, amount DECIMAL, date TIMESTAMPTZ);
```
## Установка и запуск сервера
```
git clone https://github.com/sashasych/avito.git
cd avito
go build main.go
./main.exe
```
## Тестовые запросы
Для пополнения и снятия
```
curl -d '{"UserID":1,"Change":50}' -H 'Content-Type: application/json' http://localhost:8080/updateBalance
```
для получения баланса
```
curl -d '{"ID":1, "CurrencyType":"USD"}' -H 'Content-Type: application/json' http://localhost:8080/getBalance
```
для перевода денег
```
curl -d '{"FromUserID":1,"ToUserID":2,"Change":333.33}' -H 'Content-Type: application/json' http://localhost:8080/transferMoney
```
