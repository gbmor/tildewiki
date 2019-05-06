package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/spf13/viper"
)

// Loads a given wiki page and returns a page object.
// Used for building the initial cache and re-caching.
func buildPage(filename string) (*Page, error) {

	// open the page into *os.File
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("%v\n", err)
		return nil, err
	}

	// the cleanup crew
	defer func() {
		err := file.Close()
		if err != nil {
			log.Printf("%v\n", err)
		}
	}()

	// stat the file to get mod time later
	stat, err := file.Stat()
	if err != nil {
		log.Printf("Couldn't stat %s: %v\n", filename, err)
	}

	// body holds the raw bytes from the file
	var body pagedata
	body, err = ioutil.ReadAll(file)
	if err != nil {
		log.Printf("%v\n", err)
	}

	// extract the file name from the path
	_, shortname := filepath.Split(filename)

	// get meta info on file from the header comment
	title, desc, author := body.getMeta()
	if title == "" {
		title = shortname
	}
	if desc != "" {
		desc = confVars.descSep + " " + desc
	}
	if author != "" {
		author = "`by " + author + "`"
	}

	// longtitle is used in the <title> tags of the output html
	longtitle := title + " " + confVars.titleSep + " " + confVars.wikiName

	// store the raw bytes of the document after parsing
	// from markdown to HTML.
	// keep the unparsed markdown for future use (maybe gopher?)
	bodydata := render(body, longtitle)
	return newPage(filename, shortname, title, author, desc, stat.ModTime(), bodydata, body, false), nil
}

// Scan the page until reaching following fields in the
// header comment:
//		title:
//		author:
//		description:
func (body pagedata) getMeta() (string, string, string) {

	// a bit redundant, but scanner is simpler to use
	bytereader := bytes.NewReader(body)
	metafinder := bufio.NewScanner(bytereader)
	var title, desc, author string

	// scan the file line by line until it finds
	// the comments.
	for metafinder.Scan() {

		splitter := bytes.Split(metafinder.Bytes(), []byte(":"))

		switch string(bytes.ToLower(splitter[0])) {
		case "title":
			title = string(bytes.TrimSpace(splitter[1]))
		case "description":
			desc = string(bytes.TrimSpace(splitter[1]))
		case "author":
			author = string(bytes.TrimSpace(splitter[1]))
		default:
			continue
		}

		if title != "" && desc != "" && author != "" {
			break
		}

	}

	return title, desc, author
}

// Checks the index page's cache. Returns true if the
// index needs to be re-cached.
func (indexCache *indexPage) checkCache() bool {

	interval, err := time.ParseDuration(viper.GetString("IndexRefreshInterval"))
	if err != nil {
		log.Printf("Couldn't parse index refresh interval: %v\n", err)
	}
	// stat the index page to get the mod time for the next block
	stat, err := os.Stat(confVars.assetsDir + "/" + confVars.indexFile)
	if err != nil {
		log.Printf("Couldn't stat index page: %v\n", err)
		return false
	}

	// if the last tally time is zero, or past the
	// interval in the config file, regenerate the index
	if indexCache.LastTally.IsZero() || time.Since(indexCache.LastTally) > interval {
		return true
	}
	// if the modtime is zero (never cached) or stored mod time is different
	// from the file's modtime, (re)generate cache
	if indexCache.Modtime.IsZero() || stat.ModTime() != indexCache.Modtime {
		return true
	}

	return false
}

// Re-caches the index page
func (indexCache *indexPage) cache() {
	body := render(genIndex(), confVars.wikiName+" "+confVars.titleSep+" "+confVars.wikiDesc)
	imutex.Lock()
	indexCache.Body = body
	imutex.Unlock()
}

// Generate the front page of the wiki
func genIndex() []byte {

	var err error
	indexpath := confVars.assetsDir + "/" + confVars.indexFile

	stat, err := os.Stat(indexpath)
	if err != nil {
		log.Printf("Couldn't stat index: %v\n", err)
	}

	if indexCache.Modtime != stat.ModTime() {
		imutex.Lock()
		indexCache.Raw, err = ioutil.ReadFile(indexpath)
		imutex.Unlock()
		if err != nil {
			return []byte("Could not open \"" + indexpath + "\"")
		}

	}

	// body holds the bytes of the generated index page being sent to the client.
	// create the byte array and the buffer used to write to it
	body := make([]byte, 0)
	buf := bytes.NewBuffer(body)

	// scan the file line by line until it finds the anchor
	// comment. replace the anchor comment with a list of
	// wiki pages sorted alphabetically by title.
	imutex.RLock()
	builder := bufio.NewScanner(bytes.NewReader(indexCache.Raw))
	imutex.RUnlock()
	builder.Split(bufio.ScanLines)

	for builder.Scan() {
		if bytes.Equal(builder.Bytes(), []byte("<!--pagelist-->")) {
			buf.Write(tallyPages())
		} else {
			buf.Write(append(builder.Bytes(), byte('\n')))
		}
	}

	imutex.Lock()
	indexCache.LastTally = time.Now()
	imutex.Unlock()

	return buf.Bytes()
}

// Generate a list of pages for the front page
func tallyPages() []byte {

	// pagelist and its associated buffer hold the links
	// displayed on the index page
	pagelist := make([]byte, 0, 1)
	buf := bytes.NewBuffer(pagelist)

	// get a list of files in the director specified
	// in the config file parameter "PageDir"
	files, err := ioutil.ReadDir(confVars.pageDir)
	if err != nil {
		return []byte("*PageDir can't be read.*")
	}

	// entry is used in the loop to construct the markdown
	// link to the given page
	if len(files) == 0 {
		return []byte("*No wiki pages! Add some content.*")
	}

	// true if reversing page order, otherwise don't reverse
	switch confVars.reverseTally {
	case true:
		for i := len(files) - 1; i >= 0; i-- {
			writeIndexLinks(files[i], buf)
		}
	default:
		for _, f := range files {
			writeIndexLinks(f, buf)
		}
	}

	buf.WriteByte(byte('\n'))
	return buf.Bytes()
}

// Takes in a file and outputs a markdown link to it
func writeIndexLinks(f os.FileInfo, buf *bytes.Buffer) {

	// pull the page from the cache
	pmutex.RLock()
	page := cachedPages[f.Name()]
	pmutex.RUnlock()

	// if it hasn't been cached, cache it.
	// usually means the page is new.
	if page.Body == nil {
		page.Shortname = f.Name()
		page.Longname = confVars.pageDir + "/" + f.Name()
		page.cache()
	}

	// get the URI path from the file name
	// and write the formatted link to the
	// bytes.Buffer
	linkname := bytes.TrimSuffix([]byte(page.Shortname), []byte(".md"))
	buf.WriteString("* [" + page.Title + "](" + confVars.viewPath + string(linkname) + ") " + page.Desc + " " + page.Author + "\n")
}

// Caches a page
func (page *Page) cache() {

	// buildPage() is defined in this file.
	// it reads the file and builds the Page struct
	newpage, err := buildPage(page.Longname)
	if err != nil {
		log.Printf("Couldn't cache %v: %v", page.Longname, err)
		return
	}

	pmutex.Lock()
	cachedPages[newpage.Shortname] = newpage
	pmutex.Unlock()

}

// Compare the recorded modtime of a cached page to the
// modtime of the file on disk. If they're different,
// return `true`, indicating the cache needs
// to be refreshed.
func (page *Page) checkCache() bool {

	newpage, err := os.Stat(page.Longname)
	if err != nil {
		log.Println("Can't stat " + page.Longname + ". Using cached copy...")
		return false
	}

	if newpage.ModTime() != page.Modtime || page.Recache {
		return true
	}

	return false
}

// When TildeWiki first starts, pull all available pages
// into cache, saving their modification time as well to
// detect changes to a page on disk.
func genPageCache() {

	// build an array of all the (*os.FileInfo)'s
	// needed to build the cache
	wikipages, err := ioutil.ReadDir(confVars.pageDir)
	if err != nil {
		log.Printf("Initial Cache Build :: Can't read directory %s\n", confVars.pageDir)
		panic(err)
	}

	// spawn a new goroutine for each entry, to cache
	// everything as quickly as possible
	var wg sync.WaitGroup
	for _, f := range wikipages {

		wg.Add(1)
		go func(f os.FileInfo) {
			page := newBarePage(confVars.pageDir+"/"+f.Name(), f.Name())
			page.cache()
			log.Println("Cached page " + page.Shortname)
			wg.Done()
		}(f)
	}

	wg.Wait()
}

// Wrapper function to check the cache
// of any cacher type, and if true,
// re-cache the data
func pingCache(c cacher) {

	if c.checkCache() {
		c.cache()
	}
}
