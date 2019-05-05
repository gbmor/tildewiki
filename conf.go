package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// content-type constants
const htmlutf8 = "text/html; charset=utf-8"
const cssutf8 = "text/css; charset=utf-8"

// initialize the basic configuration and
// assign the parsed templates and compiled regex
var validPath = initConfigParams()

// Sets the basic parameters for the default viper (config library) instance
func initConfigParams() *regexp.Regexp {
	conf := viper.GetViper()

	// type of config file to look for
	conf.SetConfigType("yaml")
	// name of config file to look for
	conf.SetConfigName("tildewiki")
	// location of said config file
	conf.AddConfigPath(".")
	conf.AddConfigPath("$HOME/.config/")
	conf.AddConfigPath("/etc/")
	conf.AddConfigPath("/usr/local/etc/")

	err := conf.ReadInConfig()
	if err != nil {
		log.Fatalln("Config file error: ", err)
	}

	// WatchConfig() is a function provided by blackfriday that watches the config
	// file for any changes and automatically reloads it if needed
	conf.WatchConfig()
	conf.OnConfigChange(func(e fsnotify.Event) {
		log.Println("**NOTICE** Config file change detected: ", e.Name)
	})

	// Parse the HTML template file(s) and compile the regex path validation)
	//var Templates = template.Must(template.ParseFiles(viper.GetString("TmplDir")+"/edit.html", viper.GetString("TmplDir")+"/view.html"))
	var validPath = regexp.MustCompile(viper.GetString("ValidPath"))

	//return Templates, ValidPath
	return validPath
}

// this is a custom 500 page using a markdown doc
// in the assets directory.
// if the markdown doc can't be read, default to
// net/http's error handling
func error500(w http.ResponseWriter, _ *http.Request) {
	e500 := viper.GetString("AssetsDir") + "/500.md"
	file, err := ioutil.ReadFile(e500)
	if err != nil {
		log.Printf("Tried to read 500.md: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", htmlutf8)
	_, err = w.Write(render(file, viper.GetString("CSS"), "500: Internal Server Error"))
	if err != nil {
		log.Printf("Failed to write to HTTP stream: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// this is a custom 404 page using a markdown doc
// in the assets directory.
// if the markdown doc can't be read, default to
// net/http's error handling
func error404(w http.ResponseWriter, r *http.Request) {
	e404 := viper.GetString("AssetsDir") + "/404.md"
	file, err := ioutil.ReadFile(e404)
	if err != nil {
		log.Printf("Tried to read 404.md: %v\n", err)
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	w.Header().Set("Content-Type", htmlutf8)
	_, err = w.Write(render(file, viper.GetString("CSS"), "404: File Not Found"))
	if err != nil {
		log.Printf("Failed to write to HTTP stream: %v\n", err)
		error500(w, r)
	}
}
