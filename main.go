package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync/atomic"
)

var money atomic.Int64

var bank atomic.Int64

func payHandler(w http.ResponseWriter, r *http.Request) {

	httpRequestBody, err := io.ReadAll(r.Body)
	if err != nil {

		fmt.Println("Faild to read http request: ", err)
		return
	}

	httpRequestBodyString := string(httpRequestBody)

	paymentAmount, err := strconv.Atoi(httpRequestBodyString)
	if err != nil {

		fmt.Println(err)

	}

	if money.Load()-int64(paymentAmount) >= 0 {
		money.Add(int64(-paymentAmount))
		fmt.Println("Оплата прошла успешно: ", money.Load())
	} else {

		fmt.Println("Не хватает средств для оплаты!!")
		return
	}

}
func getMoneyHandler(w http.ResponseWriter, r *http.Request) {
	currMoney := money.Load()
	fmt.Println("Возвращаю ваш баланс..")
	w.Write([]byte("Ваш баланс: " + strconv.FormatInt(currMoney, 10)))

}

func getSaveHandler(w http.ResponseWriter, r *http.Request) {

	bankBalance := bank.Load()

	fmt.Println("Возвращаю ваш баланс в банке..")
	w.Write([]byte("Ваш баланс в банке: " + strconv.FormatInt(bankBalance, 10)))

}

func saveHandler(w http.ResponseWriter, r *http.Request) {

	httpRequestBody, err := io.ReadAll(r.Body)

	if err != nil {

		fmt.Println("Произошла ошибка: ", err)
		return
	}

	httpRequestString := string(httpRequestBody)

	saveAmount, err := strconv.Atoi(httpRequestString)

	if err != nil {

		fmt.Println("Произошла ошибка: ", err)
		return
	}

	if money.Load() >= int64(saveAmount) {

		money.Add(int64(-saveAmount))
		bank.Add((int64(saveAmount)))

		fmt.Println("Новое значение переменной money: ", money.Load())
		fmt.Println("Новое значение переменной bank: ", bank.Load())

	} else {

		fmt.Println("Не хватает денег на балансе")
		return

	}

}

func main() {
	money.Add(1000)
	http.HandleFunc("/pay", payHandler)
	http.HandleFunc("/save", saveHandler)
	http.HandleFunc("/myMoney", getMoneyHandler)
	http.HandleFunc("/mySave", getSaveHandler)

	err := http.ListenAndServe(":9091", nil)

	if err != nil {
		fmt.Println(err.Error())
	}

}
