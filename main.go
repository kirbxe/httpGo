package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

var money int = 0

var bank int = 0

var mtx sync.RWMutex

func payHandler(w http.ResponseWriter, r *http.Request) {

	httpRequestBody, err := io.ReadAll(r.Body)

	if err != nil {
		// "http status codes"
		w.WriteHeader(http.StatusInternalServerError)
		msg := "Faild to read http request: " + err.Error()
		fmt.Println(msg)
		w.Write([]byte(msg))
		return
	}

	httpRequestBodyString := string(httpRequestBody)

	paymentAmount, err := strconv.Atoi(httpRequestBodyString)

	if err != nil {

		w.WriteHeader(http.StatusBadRequest)
		msg := "Failed to convert request: " + err.Error()
		fmt.Println(err)

		w.Write([]byte(msg))
		return
	}

	mtx.Lock()
	if money-paymentAmount >= 0 {

		money = money - paymentAmount

		mtx.Unlock()

		msg := "Pay has been approved!"
		fmt.Println(msg)
		_, err := w.Write([]byte(msg))

		if err != nil {

			fmt.Println("Failed to write HTTP response: ", err)
		}
		return
	}

	mtx.Unlock()

	msg := "Not enough funds to pay!!"
	fmt.Println(msg)
	_, err = w.Write([]byte(msg))

	if err != nil {
		fmt.Println("Failed to write HTTP response: ", err)
	}
}


func getMoneyHandler(w http.ResponseWriter, r *http.Request) {
	mtx.RLock()
	defer mtx.RUnlock()
	currMoney := money
	fmt.Println("Return your balance in wallet")
	w.Write([]byte("Your balance " + strconv.FormatInt(int64(currMoney), 10)))

}

func getSaveHandler(w http.ResponseWriter, r *http.Request) {
	mtx.RLock()
	defer mtx.RUnlock()

	bankBalance := bank

	fmt.Println("Return your balance in bank...")
	w.Write([]byte("Your balance in bank: " + strconv.FormatInt(int64(bankBalance), 10)))

}

func saveHandler(w http.ResponseWriter, r *http.Request) {

	httpRequestBody, err := io.ReadAll(r.Body)

	if err != nil {
		msg := "Failed to read request: " + err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(msg))
		
		fmt.Println(msg)
		return
	}

	httpRequestString := string(httpRequestBody)

	saveAmount, err := strconv.Atoi(httpRequestString)

	if err != nil {
		msg := "Failed to convert request: " + err.Error()
		
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(msg)

		w.Write([]byte(msg))
		
		return
	}


	
	mtx.Lock()
	if money >= saveAmount {

		money = money - saveAmount
		bank += saveAmount

		mtx.Unlock()

		fmt.Println("New value of the money variable: ", money)
		fmt.Println("New value of the bank variable: ", bank)

		msg := "Money was saved in the bank"
		
		w.Write([]byte(msg))
		
		return

	}

	mtx.Unlock()
	msg := "There is not enough money on the balance to deposit in the bank!"
	fmt.Println(msg)
	w.WriteHeader(http.StatusPaymentRequired)
	w.Write([]byte(msg))

}

func initBalance(input string) {

	input = strings.TrimSpace(input)

	convBalance, err := strconv.ParseInt(input, 10, 64)

	if err != nil {

		fmt.Println("Failed to convert input: ", err)
	}

	money = int(convBalance)

}

func main() {

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Write your balance: ")

	input, err := reader.ReadString('\n')

	if err != nil {

		fmt.Println("Failed to read: ", err)
		return
	}
	initBalance(input)

	fmt.Println("Your balance: ", money)

	http.HandleFunc("/pay", payHandler)
	http.HandleFunc("/save", saveHandler)
	http.HandleFunc("/myMoney", getMoneyHandler)
	http.HandleFunc("/mySave", getSaveHandler)

	errHttp := http.ListenAndServe(":9091", nil)

	if errHttp != nil {
		fmt.Println(errHttp)
	}

}
