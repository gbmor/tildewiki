package main

import (
	"sync"
	"time"
)

// The in-memory page cache
var cachedPages = make(map[string]Page)

// Mutex for the page cache
var pmutex = &sync.RWMutex{}

// The in-memory index cache
var indexCache = indexPage{}

// Mutex for the index cache
var imutex = &sync.RWMutex{}

// indexPage and Page types implement
// this interface, currently.
type cacher interface {
	cache()
	checkCache() bool
}

// Page cache object definition
type Page struct {
	Longname  string
	Shortname string
	Title     string
	Desc      string
	Author    string
	Modtime   time.Time
	Body      []byte
	Raw       pagedata
}

// Index cache object definition
type indexPage struct {
	Modtime   time.Time
	LastTally time.Time
	Body      []byte
	Raw       pagedata
}

// Type alias for methods and readability
// in certain situations.
type pagedata []byte

// Creates a filled page object
func newPage(longname, shortname, title, author, desc string, modtime time.Time, body []byte, raw pagedata) *Page {

	return &Page{
		Longname:  longname,
		Shortname: shortname,
		Title:     title,
		Author:    author,
		Desc:      desc,
		Modtime:   modtime,
		Body:      body,
		Raw:       raw}

}

// Creates a page object with the minimal number of fields filled
func newBarePage(longname, shortname string) *Page {
	return &Page{
		Longname:  longname,
		Shortname: shortname,
		Title:     "",
		Author:    "",
		Desc:      "",
		Modtime:   time.Time{},
		Body:      nil,
		Raw:       nil,
	}
}
