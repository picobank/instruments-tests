package daopgx

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	// https://github.com/jackc/pgx
	"github.com/jackc/pgx"

	m "github.com/picobank/instruments-tests/models"
)

var pool *pgx.ConnPool

const instrumentCols = `
	i.instrument_id, 
	i.symbol, 
	i.name, 
	i.description,
	i.instrument_class_id, 
	i.currency_id,
	i.from_date,
	i.thru_date,
	i.created_at,
	i.created_by,
	i.updated_at,
	i.updated_by`

const instrumentClassCols = `
	ic.instrument_class_id, 
		ic.name`

const getInstrumentByID = `
   SELECT ` + instrumentCols + `, ` + instrumentClassCols + `
     FROM instrument i
     JOIN instrument_class ic ON ic.instrument_class_id = i.instrument_class_id 
	 LEFT JOIN instrument i2 ON i2.instrument_id = i.currency_id 
    WHERE i.instrument_id = $1`

const getInstrumentClassByID = `
   SELECT ` + instrumentClassCols + `
     FROM instrument_class ic 
    WHERE instrument_class_id = $1`

const searchInstruments = `
   SELECT ` + instrumentCols + `, ` + instrumentClassCols + `
     FROM instrument i
     JOIN instrument_class ic ON ic.instrument_class_id = i.instrument_class_id 
    WHERE 1 = 1`

const listInstrumentClasses = `
   SELECT instrument_class_id, name 
	 FROM instrument_class`

func init() {
	fmt.Println("\nTest package github.com/jackc/pgx ...")

	config, err := pgx.ParseEnvLibpq()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to parse environment:", err)
		os.Exit(1)
	}
	pool, err = pgx.NewConnPool(pgx.ConnPoolConfig{ConnConfig: config})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to connect to database:", err)
		os.Exit(1)
	}

	fmt.Println("\t Connexion a la base établie")
}

func mapInstrument(rows *pgx.Rows) (*m.Instrument, error) {
	// tester https://marcesher.com/2014/10/13/go-working-effectively-with-database-nulls/
	var instrumentID, instrumentClassID uint32
	var currencyID *int32 // nillable value
	var symbol, name, description, createdBy, updatedBy string
	var fromDate, thruDate, createdAt, updatedAt time.Time
	var classID uint32
	var className string

	err := rows.Scan(&instrumentID, &symbol, &name, &description, &instrumentClassID, &currencyID, &fromDate, &thruDate, &createdAt, &createdBy, &updatedAt, &updatedBy, &classID, &className)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error fetching query result:", err)
		return nil, err
	}

	instrumentClass := &m.InstrumentClass{ID: classID, Name: className}

	// TODO : mapper currency, institutions
	return &m.Instrument{ID: instrumentID, Symbol: symbol, Name: name, Description: description, Class: instrumentClass, Currency: nil, Institutions: nil, FromDate: fromDate, ThruDate: thruDate, CreatedAt: createdAt, UpdatedAt: updatedAt, CreatedBy: createdBy, UpdatedBy: updatedBy}, nil
}

func mapInstrumentClass(rows *pgx.Rows) (*m.InstrumentClass, error) {
	var id uint32
	var name string
	err := rows.Scan(&id, &name)
	if err != nil {
		return nil, err
	}
	return &m.InstrumentClass{ID: id, Name: name}, nil
}

// GetInstrumentClass retourne l'instrumentClass correspondant à un id
func GetInstrumentClass(instrumenClassID uint32) (*m.InstrumentClass, error) {
	fmt.Printf("\nGetInstrumentClass(%d) ...\n", instrumenClassID)

	conn := connection()
	defer pool.Release(conn)
	fmt.Printf("\t[SQL] %s (%v)\n", getInstrumentClassByID, instrumenClassID)
	rows, err := conn.Query(getInstrumentClassByID, instrumenClassID)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	var instrumentClass *m.InstrumentClass
	instrumentClass, err = mapInstrumentClass(rows)

	return instrumentClass, err
}

// GetInstrument retourne l'instrument correspondant à un id
func GetInstrument(instrumentID uint32) (*m.Instrument, error) {
	fmt.Printf("\nGetInstrument(%d) ...\n", instrumentID)

	conn := connection()
	defer pool.Release(conn)
	fmt.Printf("\t[SQL] %s (%v)\n", getInstrumentByID, instrumentID)
	rows, err := conn.Query(getInstrumentByID, instrumentID)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	var instrument *m.Instrument
	for rows.Next() {
		instrument, err = mapInstrument(rows)
		if rows.Next() {
			return nil, errors.New("La requête retourne plus d'une ligne")
		}
	}

	return instrument, err
}

// ListInstrumentClass affiche la liste des classes d'instruments dans la console
func ListInstrumentClass() ([]m.InstrumentClass, error) {
	fmt.Println("\nListe des classes d'instruments ...")
	fmt.Println("\n", listInstrumentClasses)

	rows, err := connection().Query(listInstrumentClasses)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error executing query:", err)
		os.Exit(1)
	}
	defer rows.Close()

	var result []m.InstrumentClass
	for rows.Next() {
		var id uint32
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			return nil, err
		}
		result = append(result, m.InstrumentClass{ID: id, Name: name})
	}

	return result, rows.Err()
}

// SearchInstruments affiche la liste des instruments dans la console
func SearchInstruments(criteria *InstrumentSearchCriteria) ([]m.Instrument, error) {
	fmt.Printf("\nSearchInstruments(%v) ...\n", criteria)

	query, bindings := buildCriteria(searchInstruments, criteria)
	fmt.Printf("\t[SQL] %s\n", query)
	rows, err := connection().Query(query, *bindings...)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error executing query:", err)
		os.Exit(1)
	}
	defer rows.Close()

	var instrument *m.Instrument
	var result []m.Instrument
	for rows.Next() {
		instrument, err = mapInstrument(rows)
		result = append(result, *instrument)
	}

	return result, rows.Err()
}

func buildCriteria(query string, criteria *InstrumentSearchCriteria) (string, *[]interface{}) {
	bindings := make([]interface{}, 0, 0)
	index := 0
	if criteria.InstrumentID != nil && len(criteria.InstrumentID) > 0 {
		// query = query + fmt.Sprintf("\n and i.instrument_id in (%s)", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(criteria.InstrumentID)), ","), "[]"))
		bindings = append(bindings, criteria.InstrumentID)
		index++
		query = query + fmt.Sprintf("\n and i.instrument_id = ANY($%d)", index)
	}
	if criteria.Symbol != "" {
		// query = query + fmt.Sprintf("\n and i.symbol = '%s'", criteria.Symbol)
		bindings = append(bindings, criteria.Symbol)
		index++
		query = query + fmt.Sprintf("\n and i.symbol = $%d", index)
	}
	if criteria.Name != "" {
		// query = query + fmt.Sprintf("\n and i.name like '%%%s%%'", criteria.Name)
		bindings = append(bindings, criteria.Name)
		index++
		query = query + fmt.Sprintf("\n and i.name like '%%' || $%d || '%%'", index)
	}
	if criteria.ClassName != "" {
		// query = query + fmt.Sprintf("\n and ic.name = '%s'", criteria.ClassName)
		bindings = append(bindings, criteria.ClassName)
		index++
		query = query + fmt.Sprintf("\n and ic.name = $%d", index)
	}
	if criteria.ClassID != 0 {
		// query = query + fmt.Sprintf("\n and ic.instrument_class_id = %d", criteria.ClassID)
		bindings = append(bindings, criteria.ClassID)
		index++
		query = query + fmt.Sprintf("\n and ic.instrument_class_id = $%d", index)
	}
	// TODO Currency,CheckDate
	return query, &bindings
}

func connection() *pgx.Conn {
	conn, err := pool.Acquire()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error acquiring connection:", err)
		panic(err)
	}
	// defer release(conn)

	return conn
}

func release(conn *pgx.Conn) {
	pool.Release(conn)
}

// InstrumentSearchCriteria blabla
type InstrumentSearchCriteria struct {
	InstrumentID []uint32 `json:"ids,omitempty"`
	Symbol       string   `json:"symbol,omitempty"`
	Name         string   `json:"name,omitempty"`
	ClassName    string   `json:"className,omitempty"`
	ClassID      uint32   `json:"class,omitempty"`
	// Currency     string
	// CheckDate    time.Time
}

// ToJSON ...
func (isc *InstrumentSearchCriteria) ToJSON() []byte {
	b, _ := json.Marshal(isc)
	return b
}

// FromJSON ...
func (isc *InstrumentSearchCriteria) FromJSON(jsonStr []byte) error {
	return json.Unmarshal(jsonStr, &isc)
}
