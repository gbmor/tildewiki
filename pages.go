package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// loads a given wiki page and returns a page struct pointer
func loadPage(filename string) (*Page, error) {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println("Couldn't read " + filename)
		return nil, err
	}
	filestat, err := os.Stat(filename)
	if err != nil {
		log.Println("Couldn't stat " + filename)
	}
	var shortname string
	filebyte := []byte(filename)
	for i := len(filebyte) - 1; i > 0; i-- {
		if filebyte[i] == byte('/') {
			shortname = string(filebyte[i+1:])
			break
		}
	}
	title := getTitle(filename)
	author := getAuthor(filename)
	desc := getDesc(filename)
	parsed := render(body, viper.GetString("CSS"), title)
	return &Page{
		Longname:  filename,
		Shortname: shortname,
		Title:     title,
		Author:    author,
		Desc:      desc,
		Modtime:   filestat.ModTime(),
		Body:      parsed,
		Raw:       body}, nil
}

// scan the page for the `title: ` field
// in the header comment. used in the construction
// of the page cache on startup
func getTitle(filename string) string {
	mdfile, err := os.Open(filename)
	if err != nil {
		return filename
	}
	defer func() {
		err := mdfile.Close()
		if err != nil {
			log.Printf("Deferred closing of %s resulted in error: %v\n", filename, err)
		}
	}()
	titlefinder := bufio.NewScanner(mdfile)
	for titlefinder.Scan() {
		splitter := strings.Split(titlefinder.Text(), ":")
		if strings.ToLower(splitter[0]) == "title" {
			return strings.TrimSpace(splitter[1])
		}
	}
	return filename
}

// scan the page for the `description: ` field
// in the header comment. used in the construction
// of the page cache on startup
func getDesc(filename string) string {
	mdfile, err := os.Open(filename)
	if err != nil {
		return ""
	}
	defer func() {
		err := mdfile.Close()
		if err != nil {
			log.Printf("Deferred closing of %s resulted in error: %v\n", filename, err)
		}
	}()
	descfinder := bufio.NewScanner(mdfile)
	for descfinder.Scan() {
		splitter := strings.Split(descfinder.Text(), ":")
		if strings.ToLower(splitter[0]) == "description" {
			return strings.TrimSpace(splitter[1])
		}
	}
	return ""
}

// scan the page for the `author: ` field
// in the header comment. used in the construction
// of the page cache on startup
func getAuthor(filename string) string {
	mdfile, err := os.Open(filename)
	if err != nil {
		return ""
	}
	defer func() {
		err := mdfile.Close()
		if err != nil {
			log.Printf("Deferred closing of %s resulted in error: %v\n", filename, err)
		}
	}()
	authfinder := bufio.NewScanner(mdfile)
	for authfinder.Scan() {
		splitter := strings.Split(authfinder.Text(), ":")
		if strings.ToLower(splitter[0]) == "author" {
			return "`by " + strings.TrimSpace(splitter[1]) + "`"
		}
	}
	return ""
}

// generate the front page of the wiki
func genIndex() []byte {
	body := make([]byte, 0)
	buf := bytes.NewBuffer(body)
	index, err := os.Open(viper.GetString("AssetsDir") + "/" + viper.GetString("Index"))
	if err != nil {
		return []byte("Could not open \"" + viper.GetString("AssetsDir") + "/" + viper.GetString("Index") + "\"")
	}
	defer func() {
		err := index.Close()
		if err != nil {
			log.Printf("Deferred closing of %s resulted in error: %v\n", viper.GetString("Index"), err)
		}
	}()
	builder := bufio.NewScanner(index)
	builder.Split(bufio.ScanLines)
	for builder.Scan() {
		if builder.Text() == "<!--pagelist-->" {
			tmp := tallyPages()
			buf.WriteString(tmp + "\n")
		} else {
			buf.WriteString(builder.Text() + "\n")
		}
	}
	return buf.Bytes()
}

// generate a list of pages for the front page
func tallyPages() string {
	pagelist := make([]byte, 0, 1)
	buf := bytes.NewBuffer(pagelist)
	pagedir := viper.GetString("PageDir")
	viewpath := viper.GetString("ViewPath")
	files, err := ioutil.ReadDir(pagedir)
	if err != nil {
		return "*PageDir can't be read.*"
	}
	var entry string
	if len(files) == 0 {
		return "*No wiki pages! Add some content.*"
	}
	for _, f := range files {
		mutex.RLock()
		page := cachedPages[f.Name()]
		mutex.RUnlock()
		if page.Body == nil {
			page.Shortname = f.Name()
			page.Longname = pagedir + "/" + f.Name()
			err := page.cache()
			if err != nil {
				log.Printf("Couldn't pull new page %s into cache: %v", page.Shortname, err)
			}
		}
		linkname := bytes.TrimSuffix([]byte(page.Shortname), []byte(".md"))
		entry = "* [" + page.Title + "](/" + viewpath + "/" + string(linkname) + ") :: " + page.Desc + " " + page.Author + "\n"
		buf.WriteString(entry)
	}
	return buf.String()
}

// used when refreshing the cached copy
// of a single page
func (page *Page) cache() error {
	page, err := loadPage(page.Longname)
	if err != nil {
		return err
	}
	mutex.Lock()
	cachedPages[page.Shortname] = *page
	mutex.Unlock()
	return nil
}

// compare the recorded modtime of a cached page to the
// modtime of the file. if they're different,
// return `true`, indicating the cache needs
// to be refreshed.
func (page *Page) checkCache() bool {
	newpage, err := os.Stat(page.Longname)
	if err != nil {
		log.Println("Can't stat " + page.Longname + ". Using cached copy...")
		return false
	}
	if newpage.ModTime() != page.Modtime {
		return true
	}
	return false
}

// when tildewiki first starts, pull all available pages
// into cache, saving their modification time as well to
// determine when to re-load the page.
func genPageCache() {
	indexpage, err := os.Stat(viper.GetString("AssetsDir") + "/" + viper.GetString("Index"))
	if err != nil {
		log.Println("Initial Cache Build :: Can't stat index page")
	}
	wikipages, err := ioutil.ReadDir(viper.GetString("PageDir"))
	if err != nil {
		log.Println("Initial Cache Build :: Can't read directory " + viper.GetString("PageDir"))
	}
	wikipages = append(wikipages, indexpage)
	var shortname string
	var longname string
	var page Page
	for _, f := range wikipages {
		shortname = f.Name()
		if shortname == viper.GetString("Index") {
			shortname = viper.GetString("AssetsDir") + "/" + viper.GetString("Index")
			longname = shortname
		} else {
			longname = viper.GetString("PageDir") + "/" + f.Name()
		}
		page.Longname = longname
		page.Shortname = shortname
		err = page.cache()
		if err != nil {
			log.Println("Couldn't cache " + page.Shortname)
		}
		log.Println("Cached page " + page.Shortname)
	}
}
