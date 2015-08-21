package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/dillonhafer/otd"
	"net/http"
	"os"
)

const Version = "1.0.0"

var options struct {
	httpAddr string
	version  bool
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage:  %s [options]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.StringVar(&options.httpAddr, "http", "", "HTTP listen address (e.g. 127.0.0.1:3000)")
	flag.BoolVar(&options.version, "version", false, "print version and exit")

	flag.Parse()

	if options.version {
		fmt.Printf("historybot v%v\n", Version)
		os.Exit(0)
	}
	events := otd.Events()

	serveAddress := "127.0.0.1:23000"
	if options.httpAddr != "" {
		serveAddress = options.httpAddr
	}
	fmt.Fprintln(os.Stderr, "Listening on:", serveAddress)
	fmt.Fprintln(os.Stderr, "Use `--httpAddr` flag to change the default address")

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

	err := http.ListenAndServe(serveAddress, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
