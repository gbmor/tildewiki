package main // import "github.com/gbmor/tildewiki"

import (
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/spf13/viper"
)

// TildeWiki version
const twvers = "0.5.3"

func main() {
	// determine number of parallel processes allowed
	runtime.GOMAXPROCS(runtime.NumCPU())
	// show the logo, repo link, etc
	setUpUsTheWiki()

	// set up logging if the config file params
	// are set
	if viper.GetBool("FileLogging") {
		logfile, err := os.OpenFile(viper.GetString("LogFile"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("Couldn't log to file: %v\n", err)
		}
		log.SetOutput(logfile)
		defer func() {
			err := logfile.Close()
			if err != nil {
				log.Printf("Couldn't close log file: %v\n", err)
			}
		}()
	}
	// Tell Tildewiki to be quiet,
	// Supersedes file logging
	if viper.GetBool("QuietLogging") {
		logfile, err := os.Open("/dev/null")
		if err != nil {
			log.Printf("Couldn't quiet logging: %v\n", err)
		}
		log.SetOutput(logfile)
		defer func() {
			err := logfile.Close()
			if err != nil {
				log.Printf("Couldn't close log file: %v\n", err)
			}
		}()
	}

	// fill the page cache
	log.Println("**NOTICE** Building initial cache ...")
	genPageCache()

	viewpath := "/" + viper.GetString("ViewPath") + "/"

	http.HandleFunc("/", indexHandler)
	http.HandleFunc(viewpath, validatePath(pageHandler))
	http.HandleFunc("/css", cssHandler)
	http.HandleFunc("/icon", iconHandler)

	port := ":" + viper.GetString("Port")
	log.Println("**NOTICE** Binding to " + port)

	// let the user know if using reversed page listings
	if viper.GetBool("ReverseTally") {
		log.Printf("**NOTICE** Using reversed page listings on index ... \n")
	}

	log.Fatal(http.ListenAndServe(port, nil))
}
