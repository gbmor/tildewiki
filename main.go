package main // import "github.com/gbmor/tildewiki"

import (
	"log"
	"net/http"
)

// TildeWiki version
const twvers = "development version"

// initialize the basic configuration and
// assign the parsed templates and compiled regex
// to these two globals
var templates, validPath = initConfigParams()

// the in-memory page cache
var cachedPages = map[string][]byte{}

func main() {
	// fill the page cache
	genPageCache()

	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/w/", makeHandler(viewHandler))
	//http.HandleFunc("/edit/", makeHandler(editHandler))
	//http.HandleFunc("/save/", makeHandler(saveHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
