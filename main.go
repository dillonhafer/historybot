package main

import (
	"encoding/json"
	"fmt"
	"github.com/dillonhafer/otd"
	"net/http"
	"os"
)

func main() {
	events := otd.Events()

	println("Running server")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var jsonResp struct {
			Text string `json:"text"`
		}
		jsonResp.Text = otd.RandomEvent(events)
		js, err := json.Marshal(jsonResp)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		w.Write(js)
	})

	err := http.ListenAndServe("127.0.0.1:23000", nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
