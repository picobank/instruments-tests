package main

import (
	"fmt"
	"sync"
	"testing"

	dao "github.com/picobank/instruments-tests/dao-pgx"
)

var (
	symbols   []string
	setupOnce sync.Once
)

func setup(b *testing.B) {
	setupOnce.Do(func() {
		criteria := dao.InstrumentSearchCriteria{}
		instruments, err := dao.SearchInstruments(&criteria)
		if err != nil {
			b.Fatal("Failed to load Instrument symbols", "err", err)
		}

		symbols = make([]string, len(instruments))
		for i, instr := range instruments {
			symbols[i] = instr.Symbol
		}

		fmt.Printf("Symbols: %v\n", symbols)
	})

}

// BenchmarkFindInstrumentBySymbolColumn tests a select of a single row by ID
func BenchmarkFindInstrumentBySymbolColumn(b *testing.B) {
	setup(b)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		id := symbols[i%len(symbols)]
		// fmt.Printf("search for instrument[%s]\n", id)
		criteria := dao.InstrumentSearchCriteria{}
		criteria.Symbol = id
		instrument, err := dao.SearchInstruments(&criteria)
		if err != nil {
			b.Fatal("Failed to select Instrument", "err", err)
		}
		if instrument == nil {
			b.Fatal("Unable to find instrument")
		}
	}
}

// BenchmarkFindInstrumentBySymbolColumn tests a select of a single row by ID
func BenchmarkFindInstrumentBySymbolJSON(b *testing.B) {
	setup(b)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		id := symbols[i%len(symbols)]
		instrument, err := dao.SearchInstrumentsJSON(id)
		if err != nil {
			b.Fatal("Failed to select Instrument", "err", err)
		}
		if instrument == nil {
			b.Fatal("Unable to find instrument")
		}
	}
}
