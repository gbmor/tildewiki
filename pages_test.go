package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"testing"
	"time"
)

// to quiet the function output during
// testing and benchmarks
var hush, _ = os.Open("/dev/null")

var buildPageCases = []struct {
	name     string
	filename string
	want     *Page
	wantErr  bool
}{
	{
		name:     "example.md",
		filename: "pages/example.md",
		want:     &Page{},
		wantErr:  false,
	},
	{
		name:     "fake.md",
		filename: "pages/fake.md",
		want:     &Page{},
		wantErr:  true,
	},
}

func Test_buildPage(t *testing.T) {
	log.SetOutput(hush)
	for _, tt := range buildPageCases {
		t.Run(tt.name, func(t *testing.T) {
			testpage, err := buildPage(tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildPage() error = %v, wantErr %v\n", err, tt.wantErr)
			}
			if testpage == nil && !tt.wantErr {
				t.Errorf("buildPage() returned nil bytes when it wasn't expected.\n")
			}
		})
	}
}
func Benchmark_buildPage(b *testing.B) {
	log.SetOutput(hush)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, c := range buildPageCases {
			_, err := buildPage(c.filename)
			if (err != nil) != c.wantErr {
				b.Errorf("buildPage benchmark failed: %v\n", err)
			}
		}
	}
}

var metaBytes, _ = ioutil.ReadFile("pages/example.md")
var metaTestBytes pagedata = metaBytes
var getMetaCases = []struct {
	name      string
	data      pagedata
	titlewant string
	descwant  string
	authwant  string
}{
	{
		name:      "example",
		data:      metaTestBytes,
		titlewant: "Example Page",
		descwant:  "Example page for the wiki",
		authwant:  "gbmor",
	},
}

func Test_getMeta(t *testing.T) {
	for _, tt := range getMetaCases {
		t.Run(tt.name, func(t *testing.T) {
			if title, desc, auth := tt.data.getMeta(); title != tt.titlewant || desc != tt.descwant || auth != tt.authwant {
				t.Errorf("getMeta() = %v, %v, %v .. want %v, %v, %v", title, desc, auth, tt.titlewant, tt.descwant, tt.authwant)
			}
		})
	}
}
func Benchmark_getMeta(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tt := range getMetaCases {
			tt.data.getMeta()
		}
	}
}

func Test_genIndex(t *testing.T) {
	initConfigParams()
	log.SetOutput(hush)
	genPageCache()
	t.Run("genIndex() test", func(t *testing.T) {
		if got := genIndex(); got == nil {
			t.Errorf("genIndex(), got %v bytes.", got)
		}
	})
}
func Benchmark_genIndex(b *testing.B) {
	initConfigParams()
	log.SetOutput(hush)
	genPageCache()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		indexCache.page.Modtime = time.Time{}
		genIndex()
	}
}

var tallyPagesPagelist = make([]byte, 0, 1)
var tallyPagesBuf = bytes.NewBuffer(tallyPagesPagelist)

// Currently tests for whether the buffer is being written to.
// Also checks if the anchor tag was replaced in the buffer.
func Test_tallyPages(t *testing.T) {
	t.Run("tallyPages test", func(t *testing.T) {
		if tallyPages(tallyPagesBuf); tallyPagesBuf == nil {
			t.Errorf("tallyPages() wrote nil to buffer\n")
		}
		bufscan := bufio.NewScanner(tallyPagesBuf)
		for bufscan.Scan() {
			if bufscan.Text() == "<!--pagelist-->" {
				t.Errorf("tallyPages() - Did not replace anchor tag with page listing.\n")
			}
		}
	})
}
func Benchmark_tallyPages(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// I'm not blanking the *Page values
		// before every run of tallyPages here
		// because the likelihood of
		// tallyPages calling page.cache() for
		// every page is near-zero
		if tallyPages(tallyPagesBuf); tallyPagesBuf == nil {
			b.Errorf("tallyPages() benchmark failed, got nil bytes\n")
		}
	}
}

type fields struct {
	Longname  string
	Shortname string
	Title     string
	Desc      string
	Author    string
	Modtime   time.Time
	Body      []byte
	Raw       []byte
}

type indexFields struct {
	Modtime   time.Time
	LastTally time.Time
}

var IndexCacheCases = []struct {
	name   string
	fields indexFields
	want   bool
}{
	{
		name: "test1",
		fields: indexFields{
			LastTally: time.Now(),
			Modtime:   time.Time{},
		},
		want: false,
	},
	{
		name: "test2",
		fields: indexFields{
			Modtime:   time.Time{},
			LastTally: time.Time{},
		},
		want: true,
	},
}

var testIndex = indexCacheBlk{
	mu: &sync.RWMutex{},
	page: &indexPage{
		Modtime:   time.Time{},
		LastTally: time.Time{},
	},
}

// Check if checkCache() method on indexPage type
// is returning the expected bool
func Test_indexPage_checkCache(t *testing.T) {
	initConfigParams()
	testindexstat, err := os.Stat(confVars.assetsDir + "/" + confVars.indexFile)
	if err != nil {
		t.Errorf("Test_indexPage_checkCache(): Couldn't stat file for first test case: %v\n", err)
	}

	for _, tt := range IndexCacheCases {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "test1" {
				tt.fields.Modtime = testindexstat.ModTime()
			}
			testIndex.page.Modtime = tt.fields.Modtime
			testIndex.page.LastTally = tt.fields.LastTally
			if got := testIndex.checkCache(); got != tt.want {
				t.Errorf("indexPage.checkCache() - got %v, want %v\n", got, tt.want)
			}
		})
	}
}

func Benchmark_indexPage_checkCache(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for range IndexCacheCases {
			testIndex.checkCache()
		}
	}
}

// Make sure indexPage.cache() is returning
// non-nil bytes for indexPage.Body field
func Test_indexPage_cache(t *testing.T) {
	for _, tt := range IndexCacheCases {
		t.Run(tt.name, func(t *testing.T) {
			testIndex.page.Modtime = tt.fields.Modtime
			testIndex.page.LastTally = tt.fields.LastTally
			testIndex.cache()
			if testIndex.page.Body == nil {
				t.Errorf("indexPage_cache(): Returning nil for field Body.\n")
			}
		})
	}
}
func Benchmark_indexPage_cache(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for range IndexCacheCases {
			if err := testIndex.cache(); err != nil {
				b.Errorf("testIndex.cache() - %v\n", err)
			}
		}
	}
}

var pageCacheCase2stat, _ = os.Stat("pages/example.md")
var pageCacheCase2bytes, _ = ioutil.ReadFile("pages/example.md")
var pageCacheCase1bytes, _ = ioutil.ReadFile("pages/test1.md")
var PageCacheCases = []struct {
	name      string
	fields    fields
	wantErr   bool
	needCache bool
}{
	{
		name: "test1.md",
		fields: fields{
			Longname:  "pages/test1.md",
			Shortname: "test1.md",
			Modtime:   time.Time{},
			Raw:       pageCacheCase1bytes,
		},
		wantErr:   false,
		needCache: true,
	},
	{
		name: "example.md",
		fields: fields{
			Longname:  "pages/example.md",
			Shortname: "example.md",
			Modtime:   pageCacheCase2stat.ModTime(),
			Raw:       pageCacheCase2bytes,
		},
		wantErr:   false,
		needCache: false,
	},
	{
		name: "doesn't exist",
		fields: fields{
			Longname:  "pages/fake.md",
			Shortname: "fake.md",
			Modtime:   time.Time{},
		},
		wantErr:   true,
		needCache: false,
	},
}

// Tests that the raw field matches
// what's been pulled from disk.
func TestPage_cache(t *testing.T) {
	log.SetOutput(hush)
	for _, tt := range PageCacheCases {
		t.Run(tt.name, func(t *testing.T) {
			page := &Page{
				Longname:  tt.fields.Longname,
				Shortname: tt.fields.Shortname,
				Raw:       tt.fields.Raw,
			}
			if err := page.cache(); !tt.wantErr {
				cachedpage := pageCache.pool[tt.fields.Shortname]
				if !bytes.Equal(cachedpage.Raw, tt.fields.Raw) {
					t.Errorf("page.cache(): byte mismatch for %v: %v\n", page.Shortname, err)
				}
			}
		})
	}
}
func Benchmark_Page_cache(b *testing.B) {
	log.SetOutput(hush)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, tt := range PageCacheCases {
			page := &Page{
				Longname:  tt.fields.Longname,
				Shortname: tt.fields.Shortname,
			}
			if err := page.cache(); err != nil && !tt.wantErr {
				b.Errorf("While benchmarking page.cache, caught: %v\n", err)
			}
		}
	}
}

// Make sure it's returning the appropriate
// bool for zeroed modtime and current modtime
func TestPage_checkCache(t *testing.T) {
	for _, tt := range PageCacheCases {
		t.Run(tt.name, func(t *testing.T) {
			page := &Page{
				Longname:  tt.fields.Longname,
				Shortname: tt.fields.Shortname,
				Modtime:   tt.fields.Modtime,
			}
			got := page.checkCache()
			if got != tt.needCache {
				t.Errorf("Page.checkCache() = %v", got)
			}
		})
	}
}
func Benchmark_Page_checkCache(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tt := range PageCacheCases {
			page := &Page{
				Longname:  tt.fields.Longname,
				Shortname: tt.fields.Shortname,
			}
			page.checkCache()
		}
	}
}

// Check that the fields are filled
// for each page in the cache
func Test_genPageCache(t *testing.T) {
	initConfigParams()
	log.SetOutput(hush)
	genPageCache()
	t.Run("genPageCache", func(t *testing.T) {
		for k, v := range pageCache.pool {
			if v.Body == nil || v.Raw == nil || v.Longname == "" {
				t.Errorf("Test_genPageCache(): %v holds incorrect data or nil bytes\n", k)
			}
		}
	})
}
func Benchmark_genPageCache(b *testing.B) {
	initConfigParams()
	log.SetOutput(hush)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		genPageCache()
	}
}

// Ensure pullFromCache() doesn't return a
// nil page from the cache
func Test_pullFromCache(t *testing.T) {
	initConfigParams()
	log.SetOutput(hush)
	genPageCache()
	t.Run("pullFromCache", func(t *testing.T) {
		for k := range pageCache.pool {
			page, err := pullFromCache(k)
			if page == nil || err != nil {
				t.Errorf("%v returned nil\n", k)
			}
		}
	})
}
func Benchmark_pullFromCache(b *testing.B) {
	initConfigParams()
	log.SetOutput(hush)
	genPageCache()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for k := range pageCache.pool {
			pullFromCache(k)
		}
	}
}

// tests if triggerRecache sets the trip bool
// on all pages in the cache
func Test_triggerRecache(t *testing.T) {
	initConfigParams()
	log.SetOutput(hush)
	genPageCache()
	t.Run("triggerRecache", func(t *testing.T) {
		triggerRecache()
		for k, v := range pageCache.pool {
			if !v.Recache {
				t.Errorf("Recache didn't trip for %v\n", k)
			}
		}
	})
}
