package etl

import (
	"fmt"
	"os"
	"strconv"
)

var buffSize int

const sqlExtract string = "select company_name,bats_name,isin,currency,mic,reuters_exchange_code,lis_local,live,tick_type,reference_price,bats_prev_close,live_date,bloomberg_primary,bloomberg_bats,mifid_share,asset_class,matching_unit,euroccp_enabled,xclr_enabled,lchl_enabled,reuters_ric_primary,reuters_ric_bats,reference_adt_eur,csd,corporate_action_status,supported_services,trading_segment,printed_name,periodic_auction_max_duration,periodic_auction_min_order_entry_size,periodic_auction_min_order_entry_notional,max_otr_count,max_otr_volume,capped,venue_cap_percentage,venue_uncap_date from etl_bats_instrument"

func init() {
	buffSize, _ = strconv.Atoi(getEnv("ETL_BUFFSIZE", "5"))
}

// Extract data extractor from etl table and send it on the channel
// the method is a channel factory to be piped to the loader
func Extract() chan BatsInstrument {
	// on utilise un channel bufferisé en lecture
	extractCh := make(chan BatsInstrument, buffSize)
	go extractDb(extractCh)
	return extractCh
}

func extractDb(loader chan BatsInstrument) {
	cnx := Connection()
	defer pool.Release(cnx)

	rows, err := cnx.Query(sqlExtract)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error executing query:", err)
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		data := BatsInstrument{}

		err := rows.Scan(&data.CompanyName,
			&data.BatsName,
			&data.Isin,
			&data.Currency,
			&data.Mic,
			&data.ReutersExchangeCode,
			&data.LisLocal,
			&data.Live,
			&data.TickType,
			&data.ReferencePrice,
			&data.BatsPrevClose,
			&data.LiveDate,
			&data.BloombergPrimary,
			&data.BloombergBats,
			&data.MifidShare,
			&data.AssetClass,
			&data.MatchingUnit,
			&data.EuroccpEnabled,
			&data.XclrEnabled,
			&data.LchlEnabled,
			&data.ReutersRicPrimary,
			&data.ReutersRicBats,
			&data.ReferenceAdtEur,
			&data.Csd,
			&data.CorporateActionStatus,
			&data.SupportedServices,
			&data.TradingSegment,
			&data.PrintedName,
			&data.PeriodicAuctionMaxDuration,
			&data.PeriodicAuctionMinOrderEntrySize,
			&data.PeriodicAuctionMinOrderEntryNotional,
			&data.MaxOtrCount,
			&data.MaxOtrVolume,
			&data.Capped,
			&data.VenueCapPercentage,
			&data.VenueUncapDate)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error processing resultset query:", err)
			panic(err)
		}

		loader <- data
	}
	// all lines are extracted and sent to the loader
	close(loader)
}
