package etl

import (
	"fmt"
	"time"
)

const sqlInsert string = "insert into instrument(symbol, name, description, currency_id, instrument_class_id, from_date, thru_date, created_at, created_by, updated_at, updated_by) values ($1, $2, $3, $4,$5, '2001-01-01', '2999-01-01', current_date, 'ETL', current_date, 'ETL' )"

func init() {}

// Load blablabla
func Load(extractCh chan BatsInstrument) chan bool {
	sigend := make(chan bool)
	go loadDb(extractCh, sigend)
	return sigend
}

func loadDb(extractCh chan BatsInstrument, sigend chan bool) {
	start := time.Now()
	count := 0
	for {
		data, open := <-extractCh
		if !open {
			break
		}
		count++
		insertDb(data)
		if count%1000 == 0 {
			fmt.Println("Instruments loaded as by now: ", count, " in ", time.Since(start))
		}
	}
	fmt.Println("Instruments loaded: ", count)
	sigend <- true
}

func insertDb(data BatsInstrument) {
	cnx := Connection()
	defer pool.Release(cnx)

	var currencyID = 3
	_, err := cnx.Exec(sqlInsert, data.Isin, data.BatsName, data.CompanyName, currencyID, 3)
	panicIf(err)
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}
