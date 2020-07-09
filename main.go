package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Account struct {
	Number  string `json:"AccountNumber"`
	Balance string `json:"Balance"`
	Desc    string `json:"AccountDescription"`
}

var Accounts []Account

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to our bank!")
	fmt.Println("Endpoint: /")
}

func returnAllAccounts(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Accounts)
}

func returnAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["number"]
	for _, account := range Accounts {
		if account.Number == key {
			json.NewEncoder(w).Encode(account)
		}
	}
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homePage)
	router.HandleFunc("/accounts", returnAllAccounts)
	router.HandleFunc("/account/{number}", returnAccount).Methods("GET")
	router.HandleFunc("/account", createAccount).Methods("POST")
	router.HandleFunc("/account/{number}", deleteAccount).Methods("DELETE")
	router.HandleFunc("/account/{number}", updateAccount).Methods("PUT")
	log.Fatal(http.ListenAndServe(":10000", router))
}

func createAccount(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var account Account
	json.Unmarshal(reqBody, &account)
	Accounts = append(Accounts, account)
	json.NewEncoder(w).Encode(account)
}

func deleteAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["number"]
	for index, account := range Accounts {
		if account.Number == id {
			Accounts = append(Accounts[:index], Accounts[index+1:]...)
		}
	}
}

func updateAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["number"]
	var index int
	var updates, account Account
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &updates)
	for i, a := range Accounts {
		if a.Number == id {
			account = a
			index = i
			break
		}
	}
	if updates.Number != "" {
		account.Number = updates.Number
	}
	if updates.Balance != "" {
		account.Balance = updates.Balance
	}
	if updates.Desc != "" {
		account.Desc = updates.Desc
	}
	Accounts[index] = account
	json.NewEncoder(w).Encode(account)
}

func main() {
	Accounts = []Account{
		Account{Number: "C983475jhh944385", Balance: "24545.5", Desc: "Checking Account"},
		Account{Number: "S3r4345242342", Balance: "444.4", Desc: "Saving Account"},
	}

	handleRequests()
}
