package etl

import (
	"database/sql"

	"github.com/lib/pq"
)

func init() {}

// EtlBatsInstrument blablabla
type BatsInstrument struct {
	CompanyName                          sql.NullString // 0
	BatsName                             sql.NullString
	Isin                                 sql.NullString
	Currency                             sql.NullString
	Mic                                  sql.NullString
	ReutersExchangeCode                  sql.NullString
	LisLocal                             sql.NullFloat64
	Live                                 sql.NullString
	TickType                             sql.NullString
	ReferencePrice                       sql.NullFloat64
	BatsPrevClose                        sql.NullFloat64 // 10
	LiveDate                             pq.NullTime
	BloombergPrimary                     sql.NullString
	BloombergBats                        sql.NullString
	MifidShare                           sql.NullString
	AssetClass                           sql.NullString
	MatchingUnit                         sql.NullInt64
	EuroccpEnabled                       sql.NullBool
	XclrEnabled                          sql.NullBool
	LchlEnabled                          sql.NullBool
	ReutersRicPrimary                    sql.NullString // 20
	ReutersRicBats                       sql.NullString
	ReferenceAdtEur                      sql.NullFloat64
	Csd                                  sql.NullString
	CorporateActionStatus                sql.NullString
	SupportedServices                    sql.NullString
	TradingSegment                       sql.NullString
	PrintedName                          sql.NullString
	PeriodicAuctionMaxDuration           sql.NullInt64
	PeriodicAuctionMinOrderEntrySize     sql.NullInt64
	PeriodicAuctionMinOrderEntryNotional sql.NullInt64 // 30
	MaxOtrCount                          sql.NullInt64
	MaxOtrVolume                         sql.NullInt64
	Capped                               sql.NullInt64
	VenueCapPercentage                   sql.NullFloat64
	VenueUncapDate                       pq.NullTime
}
