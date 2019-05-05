package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	"github.com/spf13/viper"
)

// to quiet the function output during
// testing and benchmarks
var hush, _ = os.Open("/dev/null")

var buildPageCases = []struct {
	name     string
	filename string
	want     *Page
	wantErr  bool
}{{
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
			_, err := buildPage(tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildPage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
func Benchmark_buildPage(b *testing.B) {
	log.SetOutput(hush)
	for i := 0; i < b.N; i++ {
		for _, c := range buildPageCases {
			_, err := buildPage(c.filename)
			if err != nil {
				continue
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
			_, _, _ = tt.data.getMeta()
		}
	}
}

var genIndexCases = []struct {
	name string
}{
	{
		name: "index",
	},
}

func Test_genIndex(t *testing.T) {
	for _, tt := range genIndexCases {
		t.Run(tt.name, func(t *testing.T) {
			if got := genIndex(); len(got) == 0 {
				t.Errorf("genIndex(), got %v bytes.", got)
			}
		})
	}
}
func Benchmark_genIndex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for range genIndexCases {
			out := genIndex()
			if len(out) == 0 {
				continue
			}
		}
	}
}

var tallyPagesCases = []struct {
	name string
}{
	{
		name: "index",
	},
}

func Test_tallyPages(t *testing.T) {
	for _, tt := range tallyPagesCases {
		t.Run(tt.name, func(t *testing.T) {
			if got := tallyPages(); len(got) == 0 {
				t.Errorf("tallyPages() = %v", got)
			}
		})
	}
}
func Benchmark_tallyPages(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tallyPages()
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

var testindexstat, _ = os.Stat(viper.GetString("AssetsDir") + "/" + viper.GetString("Index"))
var IndexCacheCases = []struct {
	name   string
	fields indexFields
	want   bool
}{
	{
		name: "test1",
		fields: indexFields{
			Modtime:   testindexstat.ModTime(),
			LastTally: time.Now(),
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

var testIndex = indexPage{
	Modtime:   time.Time{},
	LastTally: time.Time{},
}

// Check if checkCache() method on indexPage type
// is returning the expected bool
func Test_indexPage_checkCache(t *testing.T) {

	for _, tt := range IndexCacheCases {
		t.Run(tt.name, func(t *testing.T) {
			testIndex.Modtime = tt.fields.Modtime
			testIndex.LastTally = tt.fields.LastTally
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

// No output to test here
func Benchmark_indexPage_cache(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for range IndexCacheCases {
			testIndex.cache()
		}
	}
}

var PageCacheCases = []struct {
	name    string
	fields  fields
	wantErr bool
}{
	{
		name: "test1",
		fields: fields{
			Longname:  "pages/test1.md",
			Shortname: "test1.md",
		},
		wantErr: false,
	},
	{
		name: "example",
		fields: fields{
			Longname:  "pages/example.md",
			Shortname: "example.md",
		},
		wantErr: false,
	},
}

// See if *Page.cache() is returning an error
func TestPage_cache(t *testing.T) {
	for _, tt := range PageCacheCases {
		t.Run(tt.name, func(t *testing.T) {
			page := &Page{
				Longname:  tt.fields.Longname,
				Shortname: tt.fields.Shortname,
			}
			if err := page.cache(); (err != nil) != tt.wantErr {
				t.Errorf("Page.cache() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
func Benchmark_Page_cache(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tt := range PageCacheCases {
			page := &Page{
				Longname:  tt.fields.Longname,
				Shortname: tt.fields.Shortname,
			}
			err := page.cache()
			if err != nil {
				continue
			}
		}
	}
}

// Make sure it's returning a bool. Will modify this
// test later
func TestPage_checkCache(t *testing.T) {
	for _, tt := range PageCacheCases {
		t.Run(tt.name, func(t *testing.T) {
			page := &Page{
				Longname:  tt.fields.Longname,
				Shortname: tt.fields.Shortname,
			}
			var got interface{} = page.checkCache()
			switch got.(type) {
			case bool:
				return
			default:
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
			var maybe interface{} = page.checkCache()
			if maybe.(bool) {
				continue
			}
		}
	}
}

// No output to test for
func Benchmark_genPageCache(b *testing.B) {
	log.SetOutput(hush)
	for i := 0; i < b.N; i++ {
		genPageCache()
	}
}
