package testpgx

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

const instrumentColumns string = "instrument_id, symbol, name, description, instrument_class_id, currency_id, from_date, thru_date, created_at, created_by, updated_at, updated_by"
const getInstrumentsByID string = "select " + instrumentColumns + " from instrument where instrument_id = $1"
const getInstrumentClassByID string = "select instrument_class_id, name from instrument_class where instrument_class_id = $1"

const listInstrumentClasses string = "select instrument_class_id, name from instrument_class"
const listInstruments string = "select " + instrumentColumns + " from instrument"
const listInstrumentsForClass string = "select " + instrumentColumns + " from instrument where instrument_class_id = $1"

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

	fmt.Println("\t\t Connexion a la base établie")
}

func mapInstrument(rows *pgx.Rows) (*m.Instrument, error) {
	var instrumentID, instrumentClassID uint32
	var currencyID *int32 // nillable value
	var symbol, name, description, createdBy, updatedBy string
	var fromDate, thruDate, createdAt, updatedAt time.Time

	err := rows.Scan(&instrumentID, &symbol, &name, &description, &instrumentClassID, &currencyID, &fromDate, &thruDate, &createdAt, &createdBy, &updatedAt, &updatedBy)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error fetching query result:", err)
		return nil, err
	}

	var instrumentClass *m.InstrumentClass
	instrumentClass, err = GetInstrumentClass(instrumentClassID)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error fetching instrumentclass query result:", err)
		return nil, err
	}

	// TODO : mapper currency, class, institutions
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
	rows, err := conn.Query(getInstrumentClassByID, instrumenClassID)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	var instrumentClass *m.InstrumentClass
	for rows.Next() {
		instrumentClass, err = mapInstrumentClass(rows)
	}

	return instrumentClass, err
}

// GetInstrument retourne l'instrument correspondant à un id
func GetInstrument(instrumentID uint32) (*m.Instrument, error) {
	fmt.Printf("\nGetInstrument(%d) ...\n", instrumentID)

	cnx := connection()
	defer pool.Release(cnx)

	conn := connection()
	defer pool.Release(conn)
	rows, err := conn.Query(getInstrumentsByID, instrumentID)
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

// ListInstruments affiche la liste des instruments dans la console
func ListInstruments() ([]m.Instrument, error) {
	fmt.Printf("\nListe des instruments ...")
	fmt.Println("\n", listInstruments)

	rows, err := connection().Query(listInstruments)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error executing query:", err)
		os.Exit(1)
	}
	defer rows.Close()

	var result []m.Instrument
	for rows.Next() {
		var instrumentID, instrumentClassID int32
		var currencyID *int32 // nillable value
		var symbol, name, description, createdBy, updatedBy string
		var fromDate, thruDate, createdAt, updatedAt time.Time

		err := rows.Scan(&instrumentID, &symbol, &name, &description, &instrumentClassID, &currencyID, &fromDate, &thruDate, &createdAt, &createdBy, &updatedAt, &updatedBy)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error fetching query result:", err)
			return nil, err
		}
		var ccy string
		if currencyID != nil {
			ccy = strconv.Itoa(int(*currencyID))
		}
		fmt.Printf("\tinstrumentID=%2d symbol='%s' instrumentClassID=%d currencyID=%s fromDate=%s thruDate=%s name='%s' description='%s' \n", instrumentID, symbol, instrumentClassID, ccy, fromDate.Format("02.01.2006"), thruDate.Format("02.01.2006"), name, description)
		result = append(result, m.Instrument{})
	}

	return result, rows.Err()
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
