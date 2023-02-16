package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	oj "github.com/multiversx/mx-chain-vm-v1_4-go/scenarios/orderedjson"
)

func main() {
	fmt.Println("Instantiating Executor")
	executor, err := NewExecutor()
	if err != nil {
		panic("Failed to instantiate Executor")
	}

	http.HandleFunc("/DumpAccounts", func(w http.ResponseWriter, r *http.Request) {
		resBody, err := executor.HandleDumpAccounts()
		respond(w, resBody, err)
	})

	http.HandleFunc("/DumpAccount", func(w http.ResponseWriter, r *http.Request) {
		strAddress := r.URL.Query().Get("address")
		resBody, err := executor.HandleDumpAccount(strAddress)
		respond(w, resBody, err)
	})

	http.HandleFunc("/ExecuteTx", func(w http.ResponseWriter, r *http.Request) {
		reqBody, _ := io.ReadAll(r.Body)
		resBody, err := executor.HandleExecuteTx(reqBody)
		respond(w, resBody, err)
	})

	http.HandleFunc("/SetAccount", func(w http.ResponseWriter, r *http.Request) {
		strAddress := r.URL.Query().Get("address")
		reqBody, _ := io.ReadAll(r.Body)
		resBody, err := executor.HandleSetAccount(strAddress, reqBody)
		respond(w, resBody, err)
	})

	http.HandleFunc("/SetCurrentBlockInfo", func(w http.ResponseWriter, r *http.Request) {
		reqBody, _ := io.ReadAll(r.Body)
		resBody, err := executor.HandleSetCurrentBlockInfo(reqBody)
		respond(w, resBody, err)
	})

	http.HandleFunc("/Reset", func(w http.ResponseWriter, r *http.Request) {
		resBody, err := executor.HandleReset()
		respond(w, resBody, err)
	})

	fmt.Println("Starting server at port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func respond(w http.ResponseWriter, resBody oj.OJsonObject, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		fmt.Fprint(w, oj.JSONString(resBody))
	}
}
