package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"encoding/json"

	// https://github.com/valyala/fasthttp
	// https://github.com/julienschmidt/httprouter

	"github.com/julienschmidt/httprouter"

	m "github.com/picobank/instruments-tests/models"
	p "github.com/picobank/instruments-tests/pgx"
)

func listInstrumentClass(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "ListInstrumentClass...")
	instrumentClasses, _ := p.ListInstrumentClass()
	fmt.Fprintf(w, " => %v\n", instrumentClasses)
}

func listInstruments(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	instruments, _ := p.ListInstruments()
	fmt.Fprintf(w, "ListInstruments[%d]: \n%v\n", len(instruments), instruments)
}

func listInstrumentsForInstrumentClassID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	classID, _ := strconv.Atoi(ps.ByName("classId"))
	count, _ := p.ListInstrumentsForInstrumentClassID(int32(classID))
	fmt.Fprintf(w, " => %d rows\n", count)
}

func httpTest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	label := ps.ByName("label")
	fmt.Fprintf(w, "Hello %s\n", label)
}

func getInstrument(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, _ := strconv.Atoi(ps.ByName("instrumentID"))
	instrument, _ := p.GetInstrument(uint32(id))

	// fmt.Fprintf(w, "getInstrument( %d ) => %s\n", id, spew.Sdump(instrument)) // no formatting
	// fmt.Fprintf(w, "%# v", pretty.Formatter(instrument)) // basic generic formatting "github.com/kr/pretty"

	b, _ := json.MarshalIndent(instrument, "", "  ")
	// b, _ := json.Marshal(instrument)
	fmt.Fprintf(w, "%v", string(b)) // json generic formatting "encoding/json"
}

func searchInstruments(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// queryValues := r.URL.Query()
	// fmt.Fprintf(w, "Hello %s\n", label)
}

// test avec httpRouter
func main() {
	router := httprouter.New()
	// url de référence: url sans traitement pour mesure de performance du serveur http seul
	router.GET("/test/:label/", httpTest)
	// liste des classes d'intruments
	router.GET("/instrumentClass/", listInstrumentClass)
	// get d'un instrument par ID
	router.GET("/instrument/:instrumentID/", getInstrument)
	// liste des instruments
	router.GET("/instrument/", listInstruments)
	// liste des instruments d'une classe
	//router.GET("/instrument/:classId/", listInstrumentsForInstrumentClassID)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func checkData() {
	instrumentClasses, _ := p.ListInstrumentClass()
	fmt.Printf("    ==> %d lignes\n", len(instrumentClasses))

	instruments, _ := p.ListInstruments()
	fmt.Printf("    ==> %d lignes\n", len(instruments))

	count, _ := p.ListInstrumentsForInstrumentClassID(m.Equity)
	fmt.Printf("    ==> %d lignes\n", count)

	count, _ = p.ListInstrumentsForInstrumentClassID(m.Future)
	fmt.Printf("    ==> %d lignes\n", count)
}
