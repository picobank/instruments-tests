package daopgx

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	// https://github.com/jackc/pgx
	"github.com/jackc/pgx"

	m "github.com/picobank/instruments-tests/models"
)

var pool *pgx.ConnPool

const instrumentCols string = "i.instrument_id, i.symbol, i.name, i.description, i.instrument_class_id, i.currency_id, i.from_date, i.thru_date, i.created_at, i.created_by, i.updated_at, i.updated_by"
const instrumentClassCols string = "ic.instrument_class_id, ic.name"
const getInstrumentByID string = "select " + instrumentCols + ", " + instrumentClassCols + " from instrument i join instrument_class ic on ic.instrument_class_id = i.instrument_class_id where instrument_id = $1"
const getInstrumentClassByID string = "select " + instrumentClassCols + " from instrument_class ic where instrument_class_id = $1"

const searchInstruments string = "select " + instrumentCols + ", " + instrumentClassCols + " from instrument i join instrument_class ic on ic.instrument_class_id = i.instrument_class_id where 1=1"

const listInstrumentClasses string = "select instrument_class_id, name from instrument_class"
const listInstruments string = "select " + instrumentCols + " from instrument i"
const listInstrumentsForClass string = "select " + instrumentCols + " from instrument i where instrument_class_id = $1"

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

	query := buildCriteria(searchInstruments, criteria)
	fmt.Printf("\t[SQL] %s (%v)\n", query, criteria)
	rows, err := connection().Query(query)
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

func buildCriteria(query string, criteria *InstrumentSearchCriteria) string {
	if criteria.InstrumentID != 0 {
		query = query + fmt.Sprintf(" and i.instrument_id = %d", criteria.InstrumentID)
	}
	if criteria.Symbol != "" {
		query = query + fmt.Sprintf(" and i.symbol = '%s'", criteria.Symbol)
	}
	if criteria.Name != "" {
		query = query + fmt.Sprintf(" and i.name like '%%%s%%'", criteria.Name)
	}
	if criteria.ClassName != "" {
		query = query + fmt.Sprintf(" and ic.name = '%s'", criteria.ClassName)
	}
	// TODO Currency,CheckDate
	return query
}

// ListInstrumentsForInstrumentClassID affiche la liste des instruments pour une classe dans la console
func ListInstrumentsForInstrumentClassID(instrumentClassID int32) (int32, error) {
	fmt.Printf("\nListe des instruments pour la classe [%d] ...", instrumentClassID)
	fmt.Println("\n", listInstrumentsForClass)

	rows, err := connection().Query(listInstrumentsForClass, instrumentClassID)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error executing query:", err)
		os.Exit(1)
	}
	defer rows.Close()

	var count int32
	for rows.Next() {
		var instrumentID, instrumentClassID int32
		var currencyID *int32 // nillable value
		var symbol, name, description, createdBy, updatedBy string
		var fromDate, thruDate, createdAt, updatedAt time.Time

		err := rows.Scan(&instrumentID, &symbol, &name, &description, &instrumentClassID, &currencyID, &fromDate, &thruDate, &createdAt, &createdBy, &updatedAt, &updatedBy)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error fetching query result:", err)
			return 0, err
		}
		var ccy string
		if currencyID != nil {
			ccy = strconv.Itoa(int(*currencyID))
		}
		fmt.Printf("\tinstrumentID=%2d symbol='%s' instrumentClassID=%d currencyID=%s fromDate=%s thruDate=%s name='%s' description='%s' \n", instrumentID, symbol, instrumentClassID, ccy, fromDate.Format("02.01.2006"), thruDate.Format("02.01.2006"), name, description)
		count++
	}

	return count, rows.Err()
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
	InstrumentID uint32
	Symbol       string
	Name         string
	ClassName    string
	Currency     string
	CheckDate    time.Time
}