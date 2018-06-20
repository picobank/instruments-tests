package etl

import "fmt"

func init() {}

func Load(extractCh chan BatsInstrument) chan bool {
	sigend := make(chan bool)
	go loadDb(extractCh, sigend)
	return sigend
}

func loadDb(extractCh chan BatsInstrument, sigend chan bool) {
	for {
		data, open := <-extractCh
		if !open {
			break
		}
		fmt.Println("Loading data: ", data.CompanyName.String)
	}
	sigend <- true
}
