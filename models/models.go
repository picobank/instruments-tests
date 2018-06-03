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
	ID   uint32 `json:"id"`
	Name string `json:"name"`
}

// Instrument represents an instrument
type Instrument struct {
	ID           uint32           `json:"id"`
	Symbol       string           `json:"symbol"`
	Name         string           `json:"name"`
	Description  string           `json:"desc,omitempty"`
	Class        *InstrumentClass `json:"class"`
	Currency     *Instrument      `json:"currency,omitempty"`
	Institutions *Institution     `json:"institutions,omitempty"`
	FromDate     time.Time        `json:"from"`
	ThruDate     time.Time        `json:"to,omitempty"`
	CreatedAt    time.Time        `json:"create"`
	UpdatedAt    time.Time        `json:"update"`
	CreatedBy    string           `json:"createdBy"`
	UpdatedBy    string           `json:"updatedBy"`
}

// Institution represents an institution
type Institution struct {
	ID          uint32        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"desc,omitempty"`
	Acronym     string        `json:"acronym,omitempty"`
	Instruments []*Instrument `json:"instruments,omitempty"`
	FromDate    time.Time     `json:"from"`
	ThruDate    time.Time     `json:"to"`
	CreatedAt   time.Time     `json:"create"`
	UpdatedAt   time.Time     `json:"update"`
	CreatedBy   string        `json:"createdBy"`
	UpdatedBy   string        `json:"updatedBy"`
}
