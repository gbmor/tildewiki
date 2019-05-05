package main

import (
	"bytes"
	"io/ioutil"
	"reflect"
	"testing"

	bf "github.com/gbmor-forks/blackfriday.v2-patched"
)

var mdTestData1, _ = ioutil.ReadFile("pages/example.md")
var mdTestData2, _ = ioutil.ReadFile("pages/test1.md")
var markdownTests = []struct {
	name       string
	css        string
	title      string
	data       []byte
	renderwant []byte
}{
	{
		name:       "one",
		css:        "assets/wiki.css",
		title:      "Example Page",
		data:       mdTestData1,
		renderwant: bf.Run(mdTestData1, bf.WithRenderer(setupMarkdown("assets/wiki.css", "Example Page"))),
	},
	{
		name:       "two",
		css:        "assets/wiki.css",
		title:      "No Description",
		data:       mdTestData2,
		renderwant: bf.Run(mdTestData2, bf.WithRenderer(setupMarkdown("assets/wiki.css", "No Description"))),
	},
}

func Test_setupMarkdown(t *testing.T) {
	for _, tt := range markdownTests {
		t.Run(string(tt.name), func(t *testing.T) {
			var got interface{} = setupMarkdown(tt.css, tt.title)
			if _, ok := got.(*bf.HTMLRenderer); !ok {
				t.Errorf("setupMarkdown() returned incorrect type: %v", reflect.TypeOf(got))
			}
		})
	}
}
func Benchmark_setupMarkdown(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, c := range markdownTests {
			setupMarkdown(c.css, c.title)
		}
	}
}

func Test_render(t *testing.T) {
	for _, tt := range markdownTests {
		t.Run(string(tt.name), func(t *testing.T) {
			if !bytes.Equal(render(tt.data, tt.css, tt.title), tt.renderwant) {
				t.Errorf("render(): byte mismatch\n")
			}
		})
	}
}
func Benchmark_render(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, c := range markdownTests {
			render(c.data, c.css, c.title)
		}
	}
}
