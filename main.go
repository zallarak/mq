package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"
)

type symbols []string

type StockInfo struct {
	price      float64
	changePerc float64
	symbol     string
}

type JsonResp struct {
	List JsonRespList
}

type JsonRespMeta struct {
	Count int
}

type JsonRespList struct {
	Resources []JsonRespResourceCont
	Meta      JsonRespMeta
}

type JsonRespResourceCont struct {
	Resource JsonRespResource
}

type JsonRespResource struct {
	Fields JsonRespFields
}

type JsonRespFields struct {
	Price       float64 `json:",string"`
	Chg_percent float64 `json:",string"`
}

func (s *symbols) String() string {
	return fmt.Sprint(*s)
}

// https://golang.org/src/flag/example_test.go
func (s *symbols) Set(value string) error {
	if len(*s) > 0 {
		return errors.New("symbols value already set")
	}
	for _, sym := range strings.Split(value, ",") {

		sym = strings.ToUpper(sym)

		if sym == "BTC" {
			sym = "BTCUSD=X"
		}

		*s = append(*s, sym)
	}
	return nil
}

var inputSymbols symbols
var inputFile string
var verboseFlag bool

func main() {
	flag.Parse()
	stocks := make(map[string]StockInfo)
	ch := make(chan StockInfo)

	if len(inputFile) > 0 {
		f, err := os.Open(inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "mq: Usage of %s:\n", os.Args[0])
			flag.PrintDefaults()
			return
		}
		appendSymbolsFromFile(f, &inputSymbols)
		f.Close()
	}

	if len(inputSymbols) == 0 {
		fmt.Fprintf(os.Stderr, "mq: Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		return
	}

	for _, sym := range inputSymbols {
		go fetch(sym, ch)
	}

	for range inputSymbols {
		info := <-ch
		stocks[info.symbol] = info
	}
	fmt.Println("")
	printStocks(stocks)
	fmt.Println("")
}

func printStocks(stocks map[string]StockInfo) {
	const headerFmt = "%v\t%v\t%v\t\n"
	const rowFmt = "%s\t%.2f\t%s%.2f%%%s\t\n"
	const colorEnd = "\033[0m"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, headerFmt, "Symbol", "Price ($)", "Change today (%)")
	fmt.Fprintf(tw, headerFmt, "------", "---------", "----------------")
	for sym, info := range stocks {
		colorStart := "\033[32m+"
		if info.changePerc < 0.0 {
			colorStart = "\033[31m"
		}
		fmt.Fprintf(tw, rowFmt, sym, info.price, colorStart, info.changePerc, colorEnd)
	}
	tw.Flush()
}

func getUrl(sym string) string {
	const yahooUrl string = "http://finance.yahoo.com/webservice/v1/symbols/%s/quote?format=json&view=detail"
	return fmt.Sprintf(yahooUrl, sym)
}

func fetch(sym string, ch chan<- StockInfo) {

	var errorResponse StockInfo = StockInfo{price: 0.0, changePerc: 0.0, symbol: sym}

	url := getUrl(sym)
	resp, err := http.Get(url)
	if err != nil {
		if verboseFlag {
			fmt.Println(err)
		}
		ch <- errorResponse
		return
	}
	defer resp.Body.Close()

	// Debug
	// contents, err := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(contents))

	decoder := json.NewDecoder(resp.Body)
	var data JsonResp
	err = decoder.Decode(&data)

	if err != nil {
		ch <- errorResponse
		return
	}

	if data.List.Meta.Count == 0 {
		ch <- errorResponse
		return
	}
	ch <- StockInfo{price: data.List.Resources[0].Resource.Fields.Price,
		changePerc: data.List.Resources[0].Resource.Fields.Chg_percent,
		symbol:     sym}
}

func appendSymbolsFromFile(f *os.File, syms *symbols) {
	input := bufio.NewScanner(f)
	for input.Scan() {
		*syms = append(*syms, input.Text())
	}
}

func init() {
	flag.Var(&inputSymbols, "s", "space separated symbols")
	flag.StringVar(&inputFile, "f", "", "newline separated file of symbols")
	flag.BoolVar(&verboseFlag, "v", false, "verbose")
}
