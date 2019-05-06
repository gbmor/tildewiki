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

// Config object initialization
var confVars = &confParams{}

// (Re-)Populates config object
func setConfVars() {
	confVars.port = ":" + viper.GetString("Port")
	confVars.pageDir = viper.GetString("PageDir")
	confVars.assetsDir = viper.GetString("AssetsDir")
	confVars.cssPath = viper.GetString("CSS")
	confVars.viewPath = "/" + viper.GetString("ViewPath") + "/"
	confVars.indexRefreshInterval = viper.GetString("IndexRefreshInterval")
	confVars.wikiName = viper.GetString("Name")
	confVars.wikiDesc = viper.GetString("ShortDesc")
	confVars.descSep = viper.GetString("DescSeparator")
	confVars.titleSep = viper.GetString("TitleSeparator")
	confVars.iconPath = viper.GetString("Icon")
	confVars.indexFile = viper.GetString("Index")
	confVars.reverseTally = viper.GetBool("ReverseTally")
	confVars.validPath = regexp.MustCompile(viper.GetString("ValidPath"))
	confVars.quietLogging = viper.GetBool("QuietLogging")
	confVars.fileLogging = viper.GetBool("FileLogging")
	confVars.logFile = viper.GetString("LogFile")
}

// Sets the basic parameters for the default viper (config library) instance
func initConfigParams() {
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

	// assign the config to the confVars object
	setConfVars()

	// WatchConfig() is a function provided by blackfriday that watches the config
	// file for any changes and automatically reloads it if needed
	conf.WatchConfig()
	conf.OnConfigChange(func(e fsnotify.Event) {
		log.Println("**NOTICE** Config file change detected: ", e.Name)
		setConfVars()
	})

}

// this is a custom 500 page using a markdown doc
// in the assets directory.
// if the markdown doc can't be read, default to
// net/http's error handling
func error500(w http.ResponseWriter, _ *http.Request) {
	e500 := confVars.assetsDir + "/500.md"
	file, err := ioutil.ReadFile(e500)
	if err != nil {
		log.Printf("Tried to read 500.md: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", htmlutf8)
	_, err = w.Write(render(file, "500: Internal Server Error"))
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
	e404 := confVars.assetsDir + "/404.md"
	file, err := ioutil.ReadFile(e404)
	if err != nil {
		log.Printf("Tried to read 404.md: %v\n", err)
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	w.Header().Set("Content-Type", htmlutf8)
	_, err = w.Write(render(file, "404: File Not Found"))
	if err != nil {
		log.Printf("Failed to write to HTTP stream: %v\n", err)
		error500(w, r)
	}
}
