package main // import "github.com/gbmor/tildewiki"

import (
	"log"
	"net/http"
	"os"
	"runtime"
)

// TildeWiki version
const twvers = "0.5.3"

func main() {
	// determine number of parallel processes allowed
	runtime.GOMAXPROCS(runtime.NumCPU())
	// show the logo, repo link, etc
	setUpUsTheWiki()

	// initialize the configuration
	initConfigParams()

	// set up logging if the config file params
	// are set
	if confVars.fileLogging {
		llogfile, err := os.OpenFile(confVars.logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("Couldn't log to file: %v\n", err)
		}
		log.SetOutput(llogfile)
		defer func() {
			err := llogfile.Close()
			if err != nil {
				log.Printf("Couldn't close log file: %v\n", err)
			}
		}()
	}
	// Tell Tildewiki to be quiet,
	// Supersedes file logging
	if confVars.quietLogging {
		llogfile, err := os.Open("/dev/null")
		if err != nil {
			log.Printf("Couldn't quiet logging: %v\n", err)
		}
		log.SetOutput(llogfile)
		defer func() {
			err := llogfile.Close()
			if err != nil {
				log.Printf("Couldn't close log file: %v\n", err)
			}
		}()
	}

	// fill the page cache
	log.Println("**NOTICE** Building initial cache ...")
	genPageCache()

	http.HandleFunc("/", indexHandler)
	http.HandleFunc(confVars.viewPath, validatePath(pageHandler))
	http.HandleFunc("/css", cssHandler)
	http.HandleFunc("/icon", iconHandler)

	log.Println("**NOTICE** Binding to " + confVars.port)

	// let the user know if using reversed page listings
	if confVars.reverseTally {
		log.Printf("**NOTICE** Using reversed page listings on index ... \n")
	}

	log.Fatal(http.ListenAndServe(confVars.port, nil))
}
