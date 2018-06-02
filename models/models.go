package models

import (
	"time"
)

// enum des classes d'instruments
const (
	_ int32 = iota
	Currency
	Bond
	Equity
	Fund
	Future
	Option
	Entitlement
	Index
	InterestRate
	Commodity
	Miscellaneous
)

// InstrumentClass represents an instrument class
type InstrumentClass struct {
	ID   uint32
	Name string
}

// Instrument represents an instrument
type Instrument struct {
	ID           uint32
	Symbol       string
	Name         string
	Description  string
	Class        *InstrumentClass
	Currency     *Instrument
	Institutions *Institution
	FromDate     time.Time
	ThruDate     time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	CreatedBy    string
	UpdatedBy    string
}

// Institution represents an institution
type Institution struct {
	ID          uint32
	Name        string
	Description string
	Acronym     string
	Instruments []*Instrument
	FromDate    time.Time
	ThruDate    time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CreatedBy   string
	UpdatedBy   string
}
