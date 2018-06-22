package etl

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"
)

const sqlInsert string = "insert into instrument(symbol, name, description, currency_id, instrument_class_id, from_date, thru_date, created_at, created_by, updated_at, updated_by) values ($1, $2, $3, $4,$5, '2001-01-01', '2999-01-01', current_date, 'ETL', current_date, 'ETL' )"

var batchSize int

func init() {
	batchSize, _ = strconv.Atoi(getEnv("ETL_BATCHSIZE", "10000"))
}

// Load blablabla
func Load(extractCh chan BatsInstrument) chan bool {
	sigend := make(chan bool)

	clearInstruments()

	cnx := Connection()
	defer Release(cnx)
	_, err := cnx.Prepare("pInsert", sqlInsert)
	panicIf(err)

	go loadDb(extractCh, sigend)

	return sigend
}

func loadDb(extractCh chan BatsInstrument, sigend chan bool) {
	start := time.Now()
	count := 0
	cnx := Connection()
	defer pool.Release(cnx)
	batch := cnx.BeginBatch()
	defer batch.Close()
	for {
		data, open := <-extractCh
		if !open {
			break
		}
		count++
		// insertDb(data)
		// insertDbPrepared(data)
		insertDbBatchPrepared(data, batch)
		if count%batchSize == 0 {
			err := batch.Send(context.Background(), nil)
			batch.Close()
			batch = cnx.BeginBatch()
			panicIf(err)
			fmt.Println("Instruments loaded as by now: ", count, " in ", time.Since(start))
		}
	}
	err := batch.Send(context.Background(), nil)
	panicIf(err)
	fmt.Println("Instruments loaded: ", count, " in ", time.Since(start))
	sigend <- true
}

func clearInstruments() {
	start := time.Now()
	cnx := Connection()
	defer pool.Release(cnx)

	cmd, err := cnx.Exec("DELETE FROM instrument WHERE id >= 10000")
	panicIf(err)
	fmt.Println("Cleared ", cmd.RowsAffected(), " instruments in ", time.Since(start))
}

// 20s
func insertDb(data BatsInstrument) {
	cnx := Connection()
	defer pool.Release(cnx)

	var currencyID = 3
	_, err := cnx.Exec(sqlInsert, data.Isin, data.BatsName, data.CompanyName, currencyID, 3)
	panicIf(err)
}

// 14s
func insertDbPrepared(data BatsInstrument) {
	var currencyID = 3
	_, err := pool.Exec("pInsert", data.Isin, data.BatsName, data.CompanyName, currencyID, 3)
	panicIf(err)
}

//  5000 -> 2.5s
// 10000 -> 2.2s
// 20000 -> 2.4s
func insertDbBatchPrepared(data BatsInstrument, batch *pgx.Batch) {
	var currencyID = 3
	batch.Queue("pInsert",
		[]interface{}{data.Isin, data.BatsName, data.CompanyName, currencyID, 3},
		[]pgtype.OID{pgtype.VarcharOID, pgtype.VarcharOID, pgtype.VarcharOID, pgtype.Int4OID, pgtype.Int4OID},
		nil,
	)
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}
