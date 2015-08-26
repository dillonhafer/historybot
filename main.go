package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/dillonhafer/otd/on_this_day"
	"github.com/jackc/markovbot/markov"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

const Version = "1.0.0"

var options struct {
	httpAddr      string
	prefixSize    int
	maxOutputSize int
	seed          int64
	version       bool
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage:  %s [options]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.IntVar(&options.prefixSize, "prefix", 2, "prefix size")
	flag.StringVar(&options.httpAddr, "http", "", "HTTP listen address (e.g. 127.0.0.1:3000)")
	flag.IntVar(&options.maxOutputSize, "output", 200, "max output size in words")
	flag.Int64Var(&options.seed, "seed", -1, "seed for random number generator")
	flag.BoolVar(&options.version, "version", false, "print version and exit")

	flag.Parse()

	if options.version {
		fmt.Printf("historybot v%v\n", Version)
		os.Exit(0)
	}

	if options.seed < 0 {
		options.seed = time.Now().UnixNano()
		fmt.Fprintln(os.Stderr, "seed:", options.seed)
	}

	rand.Seed(options.seed)
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

	// Fringe Bot
	http.HandleFunc("/fringe", func(w http.ResponseWriter, r *http.Request) {
		var jsonResp struct {
			Text string `json:"text"`
		}
		fringeEvents, err := readLines("events/fringe")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Create input string
		var nevents []string
		var years []string

		for _, text := range fringeEvents {
			parts := strings.Split(text, "–")
			if len(parts) > 1 {
				nevents = append(nevents, parts[1])
				years = append(years, strings.Replace(parts[0], "On this day in ", "", -1))
			}
		}
		year := RandomYear(years)
		jevents := strings.Join(nevents, "\n")
		var in io.Reader
		in = strings.NewReader(jevents)

		// Build Markov chain
		chain, err := markov.NewChain(in, options.prefixSize)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Write response
		chainText := chain.Generate(options.maxOutputSize)

		jsonResp.Text = fmt.Sprintf("In a parallel universe, on this day in %s- %s", year, chainText)
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

func RandomYear(years []string) string {
	totalYears := len(years) - 1
	rand.Seed(time.Now().UnixNano())
	y := rand.Intn(totalYears)
	return years[y]
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
