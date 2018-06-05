package httpservice

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	// https://github.com/valyala/fasthttp

	"github.com/julienschmidt/httprouter"
	dao "github.com/picobank/instruments-tests/dao-pgx"
)

func listInstrumentClass(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "ListInstrumentClass...")
	instrumentClasses, _ := dao.ListInstrumentClass()
	fmt.Fprintf(w, " => %v\n", instrumentClasses)
}

func httpTest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	label := ps.ByName("label")
	fmt.Fprintf(w, "Hello %s\n", label)
}

func getInstrument(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, _ := strconv.Atoi(ps.ByName("instrumentID"))
	instrument, _ := dao.GetInstrument(uint32(id))

	b, _ := json.MarshalIndent(instrument, "", "  ")
	fmt.Fprintf(w, "%v", string(b)) // json generic formatting "encoding/json"
}

// url sous la forme:
// http://localhost:8080/instrument/?criteria={"ids":[10001,10002,10005],"symbol":"USD","name":"Dollar","className":"Currency","class":1}
func searchInstruments(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	queryValues := r.URL.Query()
	criteriaP := queryValues.Get("criteria")

	criteria := dao.InstrumentSearchCriteria{}
	criteria.FromJSON([]byte(criteriaP))
	fmt.Printf("---> InstrumentSearchCriteria.FromJSON : %v\n", string(criteria.ToJSON()))

	result, _ := dao.SearchInstruments(&criteria)
	b, _ := json.MarshalIndent(result, "", "  ")
	fmt.Fprintf(w, "%v", string(b)) // json generic formatting "encoding/json"
}

// Start http service
func Start() {
	router := httprouter.New()
	// url de référence: url sans traitement pour mesure de performance du serveur http seul
	router.GET("/test/v1/:label/", httpTest)

	// liste des classes d'intruments
	router.GET("/instrumentClass/", listInstrumentClass)
	// get d'un instrument par ID (raccourcis pour la recherche par critère avec uniquement un id)
	router.GET("/instrument/:instrumentID/", getInstrument)
	// recherche d'instruments par critère
	router.GET("/instrument/", searchInstruments)

	log.Fatal(http.ListenAndServe(":8080", router))
}
