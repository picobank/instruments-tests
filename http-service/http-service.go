package httpservice

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	// https://github.com/valyala/fasthttp

	"github.com/julienschmidt/httprouter"
	p "github.com/picobank/instruments-tests/dao-pgx"
)

func listInstrumentClass(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "ListInstrumentClass...")
	instrumentClasses, _ := p.ListInstrumentClass()
	fmt.Fprintf(w, " => %v\n", instrumentClasses)
}

func httpTest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	label := ps.ByName("label")
	fmt.Fprintf(w, "Hello %s\n", label)
}

func getInstrument(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, _ := strconv.Atoi(ps.ByName("instrumentID"))
	instrument, _ := p.GetInstrument(uint32(id))

	// fmt.Fprintf(w, "%v", instrument) 						// no formatting
	// fmt.Fprintf(w, "%# v", pretty.Formatter(instrument)) 	// basic generic formatting "github.com/kr/pretty"

	b, _ := json.MarshalIndent(instrument, "", "  ")
	fmt.Fprintf(w, "%v", string(b)) // json generic formatting "encoding/json"
}

func searchInstruments(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// queryValues := r.URL.Query()
	criteria := p.InstrumentSearchCriteria{InstrumentID: []uint32{10001, 10002, 10005}, Name: "Dollar"}
	// criteria := p.InstrumentSearchCriteria{Name: "Dollar"}
	result, _ := p.SearchInstruments(&criteria)
	b, _ := json.MarshalIndent(result, "", "  ")
	fmt.Fprintf(w, "%v", string(b)) // json generic formatting "encoding/json"
}

// Start http service
func Start() {
	router := httprouter.New()
	// url de référence: url sans traitement pour mesure de performance du serveur http seul
	router.GET("/test/:label/", httpTest)
	// liste des classes d'intruments
	router.GET("/instrumentClass/", listInstrumentClass)
	// get d'un instrument par ID
	router.GET("/instrument/:instrumentID/", getInstrument)
	// liste des instruments
	router.GET("/instrument/", searchInstruments)

	log.Fatal(http.ListenAndServe(":8080", router))
}
